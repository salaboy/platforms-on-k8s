# DORA Metrics + CloudEvents & CDEvents for Kubernetes

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

This tutorial install a set of components that consumes [CloudEvents](https://cloudevents.io) from multiple sources and allows you to track the DORA metrics, using  a Kubernetes-native architecture (cloud-agnostic).

This demo focus on observing different event sources to then map these events to meaninful events related to our software delivery practices that can be aggregated to calculate DORA metrics.

The Events transformation flow goes like this:
- The inputs are [CloudEvents](https://cloudevents.io) coming from diffeerent sources
- These CloudEvents can be mapped and transformed to [CDEvents](https://cdevents.dev) for further processing
- Aggregation functions can be defined to calculate DORA (or other) metrics
- Metrics can be exposed for consumption (in this example via REST endpoints)

## Installation

We will be using a Kubernetes Cluster with Knative Serving for running our transformation functions. You can follow the instructions from [Chapter 8 to create a KinD cluster with Knative Serving installed in it](https://github.com/salaboy/platforms-on-k8s/tree/main/chapter-8/knative#creating-a-kubernetes-with-knative-serving). 

Then we will install Knative Eventing, this is optional, as we will use Knative Eventing to install the Kubernetes API Event Source, which takes internal Kubernetes Events and transform them in CloudEvents.


1. Install [Knative Eventing](https://knative.dev/docs/install/yaml-install/eventing/install-eventing-with-yaml/)
```shell
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.11.0/eventing-crds.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.11.0/eventing-core.yaml
```

1. Create your "dora-cloudevents" namespace: 
```shell
kubectl create ns dora-cloudevents
```

1. Install PostgreSQL and Create Tables
```shell
kubectl apply -f resources/dora-sql-init.yaml
helm install postgresql oci://registry-1.docker.io/bitnamicharts/postgresql --version 12.5.7 --namespace dora-cloudevents --set "image.debug=true" --set "primary.initdb.user=postgres" --set "primary.initdb.password=postgres" --set "primary.initdb.scriptsConfigMap=dora-init-sql" --set "global.postgresql.auth.postgresPassword=postgres" --set "primary.persistence.size=1Gi"
```


1. Install Sockeye, a simple CloudEvents monitor, it requires Knative Serving to be installed: 

```shell
kubectl apply -f https://github.com/n3wscott/sockeye/releases/download/v0.7.0/release.yaml
```

1. Install the [Kubernetes API Server CloudEvent Event Source](https://knative.dev/docs/eventing/sources/apiserversource/getting-started/#create-an-apiserversource-object): 
```shell
kubectl apply -f api-serversource-deployments.yaml
```


## Components

This demo deploys the following components to transform CloudEvents into CDEvents and then aggregate the available data. 

- **CloudEvents Endpoint**: Endpoint to send all CloudEvents to; these CloudEvents will be stored in the SQL database to the `cloudevents-raw` table. 

- **CloudEvents Router**: Router, with a routing table, which routes events to be transformed to `CDEvents`. This mechanism allows the same event type to be transformed into multiple `CDEvents`, if needed. This component reads from the `cloudevents-raw` table and processes events. This component is triggered via configurable fixed period of time. 

- **CDEvents Transformers**: These functions receive events from the `CloudEvents Router`  and transforms the CloudEvents to CDEvents. The result is sent to the `CDEvents Endpoint`. 


- **CDEvents Endpoint**: Endpoint to send `CDEvents`, these CloudEvents will be stored in the SQL database to the `cdevents-raw` table, as they do not need any transformation. This endpoint validates that the CloudEvent received is a CD CloudEvent. 

- **Metrics Functions**: These functions are in charge of calculating different metrics and storing them into special tables, probably one per table. To calculate said metrics, these functions read from `cdevents-raw`. An example on how to calculate the **Deployment Frequency** metric is explained below. 

- **Metrics Endpoint**: Endpoint that allows you to query the metrics by name and add some filters. This is an optional component, as you can build a dashboard from the metrics tables without using these endpoints.


![dora-cloudevents-architecture](../imgs/dora-cloudevents-architecture.png)


## Deploying Components and generating data

First deploy the components and transformation functions by running: 

```shell
kubectl apply -f resources/components.yaml
```

Open Sockeye to monitor CloudEvents by pointing your browser to [http://sockeye.default.127.0.0.1.sslip.io/](http://sockeye.default.127.0.0.1.sslip.io/).


Then, create a new Deployment in the `default` namespace to test that your configuration is working.

```shell
kubectl apply -f test/example-deployment.yaml
```

At this point you should see tons of events in Sockeye:

![sockeye](../imgs/sockeye.png)

If the Deployment Frequency functions (transformation and calculation) are installed you should be able to query the deployment frequency endpoint and see the metric. Notice that this can take up to a couple of minutes, as Cron jobs are being used to aggregate the data periodically:  

```shell
curl http://dora-frequency-endpoint.dora-cloudevents.127.0.0.1.sslip.io/deploy-frequency/day | jq
```
And see something like this, depending on which deployments you created (I've created two deployments: `nginx-deployment` and `nginx-deployment-3`): 

```shell
[
  {
    "DeployName": "nginx-deployment",
    "Deployments": 1,
    "Time": "2023-08-05T00:00:00Z"
  },
  {
    "DeployName": "nginx-deployment-3",
    "Deployments": 1,
    "Time": "2023-08-05T00:00:00Z"
  }
]

```

Try modifying the deployments or creating new ones, the components are configured to monitor all Deployments in the `default` Namespace.

Notice that all the components were installed in the `dora-cloudevents` namespace. You check the pods and the URL for the Knative Services by running the following commands: 

Check the URL for the Knative Services in the `dora-cloudevents` namespace:
```shell
kubectl get ksvc -n dora-cloudevents
```

Check which pods are running, I find this interesting as because we are using Knative Serving, all transformation functions that are not being used doesn't need to be running all the time: 

```shell
kubectl get pods -n dora-cloudevents
```

Finally you can check all the CronJob executions that aggregates data by running: 

```shell
kubectl get cronjobs -n dora-cloudevents
```


## Development 

Deploy the `dora-cloudevents` components using `ko` for development:

```shell
ko apply -f config/
```

# Metrics

From [https://github.com/GoogleCloudPlatform/fourkeys/blob/main/METRICS.md](https://github.com/GoogleCloudPlatform/fourkeys/blob/main/METRICS.md)

## Deployment Frequency

![deployment frequency](imgs/deployment-frequency-metric.png)

We look for new or updated deployment resources. This is done by using the `APIServerSource` that we configured earlier. 

The flow should look like: 
```mermaid
graph TD
    A[API Server Source] --> |writes to `cloudevents_raw` table| B[CloudEvent Endpoint]
    B --> |read from `cloudevents_raw` table| C[CloudEvents Router]
    C --> D(CDEvent Transformation Function)
    D --> |writes to `cdevents_raw` table| E[CDEvents Endpoint]
    E --> F(Deployment Frequency Function)
    F --> |writes to `deployments` table| G[Deployments Table]
    G --> |read from `deployments` table| H[Metrics Endpoint]
```

Calculate buckets: Daily, Weekly, Monthly, Yearly.


This counts the number of deployments per day: 

```sql
SELECT
distinct deploy_name AS NAME,
DATE_TRUNC('day', time_created) AS day,
COUNT(distinct deploy_id) AS deployments
FROM
deployments
GROUP BY deploy_name, day;
```


## TODOs and Extensions

- Add processed events mechanism for `cloudevents_raw` and `cdevents_raw` tables. This should avoid the `CloudEvents Router` and the `Metrics Calculation Functions` to recalculate already processed events. This can be achieved by having a table that keeps track of which was the last processed event and then making sure that the `CloudEvents Router` and the `Metrics Calculation Functions` join against the new tables. 
- Add queries to calculate buckets for Deployment Frequency Metric: Weekly, Monthly, Yearly to the `deployment-frequency-endpoint.go`. Check blog post to calculate frequency and not volume: https://codefresh.io/learn/software-deployment/dora-metrics-4-key-metrics-for-improving-devops-performance/
- Create Helm Chart for generic components (CloudEvents Endpoint, CDEvents Endpoint, CloudEvents Router)
- Automate table creation for PostgreSQL helm chart (https://stackoverflow.com/questions/66333474/postgresql-helm-chart-with-initdbscripts)
- Create functions for **Lead Time for Change**


## Other Sources and Extensions

- [Install Tekton](https://github.com/cdfoundation/sig-events/tree/main/poc/tekton)
  - Tekton dashboard: `k port-forward svc/tekton-dashboard 9097:9097 -n tekton-pipelines`
  - Cloud Events Controller: `kubectl apply -f https://storage.cloud.google.com/tekton-releases-nightly/cloudevents/latest/release.yaml`
  - ConfigMap: `config-defaults` for <SINK URL>
- https://github.com/GoogleCloudPlatform/fourkeys
- https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance
- Continuously Delivery Events aka [CDEvents](https://cdevents.dev)
- CloudEvents aka [CEs](https://cloudevents.io/)  
- GitHub Source: https://github.com/knative/docs/tree/main/code-samples/eventing/github-source
