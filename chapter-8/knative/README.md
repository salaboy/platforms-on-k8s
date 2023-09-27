# Release Strategies with Knative Serving

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

This tutorial will create a Kubernetes Cluster and install Knative Serving to implement different release strategies. We will use Knative Serving percentage-based traffic splitting and tag and header-based routing.



## Creating a Kubernetes with Knative Serving

You need a Kubernetes Cluster to install [Knative Serving](https://knative.dev). You can create one using Kubernetes KinD, but instead of using the configurations provided in Chapter 2, we will use the following command:

```shell
cat <<EOF | kind create cluster --name dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 31080 # expose port 31380 of the node to port 80 on the host, later to be use by kourier or contour ingress
    listenAddress: 127.0.0.1
    hostPort: 80
EOF
```

When using Knative Serving, there is no need to install an Ingress Controller, as Knative Serving requires a more advanced networking stack to enable features such as traffic routing and splitting. We will be installing [Kourier](https://github.com/knative-extensions/net-kourier) for this, but you can install a fully fledge service mesh like [Istio](https://istio.io/).  


Once you have a Cluster, let's start by installing [Knative Serving](https://knative.dev/docs/install/yaml-install/serving/install-serving-with-yaml/). You can follow the official documentation or copy the installation steps listed here, as the examples had been tested with this version of Knative. 

Install Knative Serving Custom Resource Definitions:

```shell
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.2/serving-crds.yaml

```

Then install the Knative Serving Control Plane: 
```shell
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.2/serving-core.yaml

```

Install a compatible Knative Serving Networking Stack, here you can install Istio, a full-blown service mesh or more simple options like Kourier, that we will use here to limit the resource utilization on your cluster: 

```shell
kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.10.0/kourier.yaml

```

Configure Kourier as the selected networking stack: 

```shell
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress-class":"kourier.ingress.networking.knative.dev"}}'
```

Finally, configure the DNS resolution for your cluster, so it can expose reachable URLs for our services:

```shell
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.2/serving-default-domain.yaml

```

**Only for Knative on KinD** 

For Knative Magic DNS to work in KinD you need to patch the following ConfigMap:

```shell
kubectl patch configmap -n knative-serving config-domain -p "{\"data\": {\"127.0.0.1.sslip.io\": \"\"}}"
```

and if you installed the `kourier` networking layer you need to create an ingress:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: kourier-ingress
  namespace: kourier-system
  labels:
    networking.knative.dev/ingress-provider: kourier
spec:
  type: NodePort
  selector:
    app: 3scale-kourier-gateway
  ports:
    - name: http2
      nodePort: 31080
      port: 80
      targetPort: 8080
EOF
```


One more step is needed to run the examples covered in the following section and this is to enable `tag-header-based-routing` and access to the [Kubernetes Downward API](https://knative.dev/docs/serving/configuration/feature-flags/#kubernetes-downward-api) to fetch information about the pods running in the cluster. We can tune the Knative Serving installation by patching the config-features ConfigMap (https://knative.dev/docs/serving/configuration/feature-flags/#feature-and-extension-flags):

```shell
kubectl patch cm config-features -n knative-serving -p '{"data":{"tag-header-based-routing":"Enabled", "kubernetes.podspec-fieldref": "Enabled"}}'
```

After all these configurations, we should be ready to use Knative Serving on our KinD Cluster. 


# Release Strategies with Knative Serving

Before jumping into implementing different release strategies we need to understand the basics on Knative Serving. For that we will look into defining one Knative Service for the Notification Service. This will give us hands-on experience on using Knative Serving and its capabilities. 

## Knative Services quick intro

Knative Serving simplifies and extends the capabilities offered by Kubernetes by using the concept of a `Knative Service`. A `Knative Service` uses the Knative Serving networking layer to route traffic to our workloads without pushing us to define complex Kubernetes Resources. Because Knative Serving has access to information about how traffic is flowing to our services, it can understand the load that our services are experiencing and make use of a purposefully built autoscaler to upscale or downscale our service instances based on demand. This can be really useful for platform teams looking to implement a function-as-a-service model for their workloads, as Knative Serving can downscale to zero services that are not receiving traffic. 

Knative Services also expose a simplified configuration that resembles a Containers-as-a-Service model like Google Cloud Run, Azure Container Apps, and AWS App Runner, whereby defining which container we want to run, the platform will take care of the rest (no complex configurations for networking, routing traffic, etc). 

Because the Notifications Service uses Kafka for emitting events, we need to install Kafka using Helm:

```shell
helm install kafka oci://registry-1.docker.io/bitnamicharts/kafka --version 22.1.5 --set "provisioning.topics[0].name=events-topic" --set "provisioning.topics[0].partitions=1" --set "persistence.size=1Gi" 
```

Check that Kafka is running before proceeding, as it usually takes a bit of time to fetch the Kafka Container image and start it. 

If you have a Kubernetes Cluster with Knative Serving installed you can apply the following Knative Service resource to run an instance of our Notifications Service:

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: notifications-service 
spec:
  template:
    spec:
      containers:
        - image: salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0
          env:
          - name: KAFKA_URL
            value: kafka.default.svc.cluster.local
          ...<MORE ENVIRONMENT VARIABLES HERE>...  
```

You can apply this resource by running: 

```shell
kubectl apply -f knative/notifications-service.yaml
```

Knative Serving will create an instance of our container and setup all the networking configurations to provide us with a URL to access the service. 

You can list all Knative Services running the following command: 

```shell
> kubectl get ksvc 
NAME                    URL                                                       LATESTCREATED                 LATESTREADY                   READY   REASON
notifications-service   http://notifications-service.default.127.0.0.1.sslip.io   notifications-service-00001   notifications-service-00001   True    

```

You can now curl the `service/info` URL of the service to make sure that it is working, we are using [`jq`](https://jqlang.github.io/jq/download/) a popular json utility to pretty-print the output:

```shell
curl http://notifications-service.default.127.0.0.1.sslip.io/service/info | jq 
```

You should see the following output: 

```json
{
    "name": "NOTIFICATIONS",
    "podIp": "10.244.0.16",
    "podName": "notifications-service-00001-deployment-7ff76b4c77-qkk69",
    "podNamespace": "default",
    "podNodeName": "dev-control-plane",
    "podServiceAccount": "default",
    "source": "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "version": "1.0.0"
}

```

Check that there is a Pod running our container: 

```shell
> kubectl get pods
NAME                                                      READY   STATUS    RESTARTS   AGE
kafka-0                                                   1/1     Running   0          7m54s
notifications-service-00001-deployment-798f8f79f5-jrbr8   2/2     Running   0          4s
```

Notice that there are two containers running, one is our Notification Service and the other is the Knative Serving `queue` proxy which is used to buffer incoming requests and get traffic metrics. 

After 90 seconds (by default), if you are not sending requests to the Notification Service, the service instance will be automatically downscaled. Whenever a new incoming request arrives, Knative Serving will automatically upscale the service and buffer the request until the instance is ready. 


To recap, we get two things out of the box with Knative Serving: 
- A simplified way to run our workloads without creating multiple Kubernetes Resources. This approach resembles a Container-as-a-Service offering, that as a Platform team you might want to offer to your teams.
- Dynamic autoscaling using the Knative Autoscaler, which can be used to downscale your applications to zero when they are not being used. This resembles a Functions-as-a-Service approach, that as a Platform team, you might want to provide your teams.

## Run the Conference application with Knative Services

In this section we will look into implementing different release strategies for our Conference Application, for that we will be deploying all the other application services also using Knative Services. 

Before installing the other services we need to set up PostgreSQL and Redis, as we already installed Kafka before. Before installing PostgreSQL we need to create a ConfigMap containing the SQL statement and create the `Proposals` Table, so the Helm Chart can reference the configMap and execute the statement when the database instance is started.

```shell
kubectl apply -f knative/c4p-sql-init.yaml
```

```shell
helm install postgresql oci://registry-1.docker.io/bitnamicharts/postgresql --version 12.5.7 --set "image.debug=true" --set "primary.initdb.user=postgres" --set "primary.initdb.password=postgres" --set "primary.initdb.scriptsConfigMap=c4p-init-sql" --set "global.postgresql.auth.postgresPassword=postgres" --set "primary.persistence.size=1Gi"

```

and Redis: 

```shell
helm install redis oci://registry-1.docker.io/bitnamicharts/redis --version 17.11.3 --set "architecture=standalone" --set "master.persistence.size=1Gi"
```

Now we can install all the other services (frontend, c4p-service, and agenda-service) by running: 

```shell
kubectl apply -f knative/
```

Check that all the Knative Services are `READY`

```shell
> kubectl get ksvc
NAME                    URL                                                       LATESTCREATED                 LATESTREADY                   READY   REASON
agenda-service          http://agenda-service.default.127.0.0.1.sslip.io          agenda-service-00001          agenda-service-00001          True    
c4p-service             http://c4p-service.default.127.0.0.1.sslip.io             c4p-service-00001             c4p-service-00001             True    
frontend                http://frontend.default.127.0.0.1.sslip.io                frontend-00001                frontend-00001                True    
notifications-service   http://notifications-service.default.127.0.0.1.sslip.io   notifications-service-00001   notifications-service-00001   True    

```

Access the Conference Application Frontend by pointing your browser to the following URL [http://frontend.default.127.0.0.1.sslip.io](http://frontend.default.127.0.0.1.sslip.io)

At this point the application should work as expected, with a small difference, services like the Agenda Service and the C4P Service will be downscaled when they are not being used. If you list the pods after 90 seconds of inactivity you should see the following: 
```shell
> kubectl get pods 
NAME                                                     READY   STATUS    RESTARTS   AGE
frontend-00002-deployment-7fdfb7b8c5-cw67t               2/2     Running   0          60s
kafka-0                                                  1/1     Running   0          20m
notifications-service-00002-deployment-c5787bc49-flcc9   2/2     Running   0          60s
postgresql-0                                             1/1     Running   0          9m23s
redis-master-0                                           1/1     Running   0          8m50s
```

Because the Agenda and C4P service are storing state into persistent storage (Redis and PostgreSQL) data is not lost when the service instance is downscaled. But the Notifications and Frontend services are keeping in-memory data (Notifications and consumed Events), hence we have configured our Knative service to keep at least one instance alive all the time. All the in-memory state kept like this will impact how the application can scale, but remember this is just a Walking Skeleton. 

Now that we have the application up and running let's take a look at some different release strategies. 

# Canary releases

In this section, we will run a simple example showing how to do canary releases by using Knative Services. We will start simply by looking into percentage-based traffic splitting.


Percentage-based traffic-splitting functionalities provided out-of-the-box by Knative Services. We will be updating the Notification Service that we deployed before instead of changing the Frontend as dealing with multiple requests to fetch CCS and JavaScript files can get tricky when using percentage-based traffic-splitting.

To make sure that the service is still up and running you can run the following command: 

```shell
curl http://notifications-service.default.127.0.0.1.sslip.io/service/info | jq
```

You should see the following output: 

```json
{
    "name": "NOTIFICATIONS",
    "podIp": "10.244.0.16",
    "podName": "notifications-service-00001-deployment-7ff76b4c77-qkk69",
    "podNamespace": "default",
    "podNodeName": "dev-control-plane",
    "podServiceAccount": "default",
    "source": "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "version": "1.0.0"
}

```

You can edit the Knative Service (`ksvc`) of the Notification Service and create a new revision by changing the container image that the service is using or changing any other configuration parameter such as environment variables:

```shell
kubectl edit ksvc notifications-service
```

And then change from: 

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: notifications-service 
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/min-scale: "1"
    spec:
      containers:
        - image: salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0  
      ...
```

To `v1.1.0`: 

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: notifications-service 
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/min-scale: "1"
    spec:
      containers:
        - image: salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.1.0  
      ...
```

Before saving this change, that will create a new revision which we can use to split traffic, we need to add the following values into the traffic section:

```yaml
 traffic:
  - latestRevision: false
    percent: 50
    revisionName: notifications-service-00001
  - latestRevision: true
    percent: 50
```

Now if you start hitting the `service/info` endpoint again you will see that half of the traffic is being routed to version `v1.0.0` of our service and the other half to version `v1.1.0`.

```shell
> curl http://notifications-service.default.127.0.0.1.sslip.io/service/info | jq
{
    "name":"NOTIFICATIONS-IMPROVED",
    "version":"1.1.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/v1.1.0/conference-application/notifications-service",
    "podName":"notifications-service-00003-deployment-59fb5bff6c-2gfqt",
    "podNamespace":"default",
    "podNodeName":"dev-control-plane",
    "podIp":"10.244.0.34",
    "podServiceAccount":"default"
}

> curl http://notifications-service.default.127.0.0.1.sslip.io/service/info | jq
{
    "name":"NOTIFICATIONS",
    "version":"1.0.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "podName":"notifications-service-00001-deployment-7ff76b4c77-h6ts4",
    "podNamespace":"default",
    "podNodeName":"dev-control-plane",
    "podIp":"10.244.0.35",
    "podServiceAccount":"default"
}
```

This mechanism is really useful when you need to test a new version but you are not willing to route all the live traffic straightaway to the new version in case that problems arise.

You can modify the traffic rules to have a different percentage split, if you feel confident that the newest version is stable enough to receive more traffic. 

```yaml
 traffic:
  - latestRevision: false
    percent: 10
    revisionName: notifications-service-00001
  - latestRevision: true
    percent: 90
```

The moment that one revision (version) doesn't has any traffic rule point to it, the instance will be downscaled, as no traffic will be routed to it.


# A/B Testing and Blue/Green Deployments

With A/B testing we want to run two or more versions of the same application/service at the same time to enable different groups of users to test changes so we can decide which version works best for them. 

With Knative Serving we have two options: `Header-based` routing and `Tag-based` routing, both use the same mechanisms and configurations behind the covers, but let's see how these mechanisms can be used. 

With Tag/Header-based routing we have more control over where request will go, as we can use an HTTP header or a specific URL to instruct Knative networking mechanisms to route traffic to specific versions of the service. 

This means that for this example we can change the Frontend of our application, as all requests including a Header or a specific URL will be routed to the same version of the service.

Make sure to access the application Frontend by pointing your browser to [http://frontend.default.127.0.0.1.sslip.io](http://frontend.default.127.0.0.1.sslip.io)

![frontend v1.0.0](../imgs/frontend-v1.0.0.png)


Let's now modify the Frontend Knative Service to deploy a new version with the debug feature enabled: 

```shell
kubectl edit ksvc frontend
```

Update the image field to point to `v1.1.0` and add the FEATURE_DEBUG_ENABLED environment variable (remember that we are using the first version of the application that is not using OpenFeature).

```yaml
spec:
      containerConcurrency: 0
      containers:
      - env:
        - name: FEATURE_DEBUG_ENABLED
          value: "true"
       ...
        image: salaboy/frontend-go-1739aa83b5e69d4ccb8a5615830ae66c:v1.1.0
```

Before saving the Knative Service, let's change the traffic rules to match the following:

```yaml
traffic:
  - latestRevision: false
    percent: 100
    revisionName: frontend-00001
    tag: current
  - latestRevision: true
    percent: 0
    tag: version110
```

Notice that no traffic (percent: 0) will be routed to `v1.1.0` unless the tag is specified in the service URL. Users can now point to [http://version110-frontend.default.127.0.0.1.sslip.io](http://version110-frontend.default.127.0.0.1.sslip.io) to access `v1.1.0`


![v1.1.0](../imgs/frontend-v1.1.0.png)

Notice that `v1.1.0` has a different color theme, when you see them side by side you can notice the difference. Check the other sections of the application too. 


If for some reason, you don't want or can't change the URL of the service, you can use HTTP Headers to access `v1.1.0`. Using a Browser plugin like [Chrome ModHeader](https://chrome.google.com/webstore/detail/modheader-modify-http-hea/idgpnmonknjnojddfkpgkljpfnnfcklj) you can modify all the requests that the browser is sending by adding parameters or headers. 

Here we are setting the `Knative-Serving-Tag` header with the value `version110`, which is the name of the tag that we configured in the traffic rules for our frontend Knative Service. 

Now we can access to the normal Knative Service URL (with no changes) to access `v1.1.0`: [http://frontend.default.127.0.0.1.sslip.io](http://frontend.default.127.0.0.1.sslip.io)

![v1.1.0 with header](../imgs/frontend-v1.1.0-with-header.png)


Tag and header-based routing allow us to implement Blue/Green deployments in the same way, as the `green` service (the one we want to test until it is ready for prime time) can be hidden behind a tag with 0% traffic assigned to it.

```yaml
traffic:
    - revisionName: <blue-revision-name>
      percent: 100 # All traffic is still being routed to this revision
    - revisionName: <gree-revision-name>
      percent: 0 # 0% of traffic routed to this version
      tag: green # A named route
```

Whenever we are ready to switch to our `green` service we can change the traffic rules: 

```yaml
traffic:
    - revisionName: <blue-revision-name>
      percent: 0 
      tag: blue 
    - revisionName: <gree-revision-name>
      percent: 100
      
```

To recap, by using Knative Services traffic splitting and header/tag-based routing capabilities we have implemented Canary Releases, A/B testing patterns, and Blue/Green deployments. Check the [Knative Website](https://knative.dev) for more information about the project. 

## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```shell
kind delete clusters dev
```

