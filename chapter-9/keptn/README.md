# Keptn Lifecycle Toolkit, out of the box Deployment Frequency

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---


On this short tutorial we explore the Keptn Lifecycle Toolkit to monitor, observe and react to our cloud native applications lifecycle events. 


## Installation

You need a Kubernetes Cluster to install [Keptn KLT](https://keptn.sh). You can create one using Kubernetes KinD as we did in [Chapter 2](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind)

Then we can install the Keptn Lifecycle Toolkit (KLT). This can be usually done by just installing the Keptn Lifecycle Toolkit Helm chart, but for this tutorial we want to also install Prometheus, Jaeger and Grafana for having dashboards. For that reason, based on the Keptn Lifecycle Toolkit repository, we will use a Makefile to install all the tools that we need for this example. 

Run: 

```shell
make install
```

**Note*: The installation process will take a few minutes to install all the tools needed.

Finally, we need to let KLT know which namespace we want to monitor, and for that we need to annotate the namespaces:

```shell
kubectl annotate ns default keptn.sh/lifecycle-toolkit="enabled"
```

## Keptn Lifecycle toolkit in action

Keptn uses standard Kubernetes annotation to recognized and monitor our workloads. 
The Kubernetes Deployments used by the Conference Application are annotated with the following annotations, for example the Agenda Service: 

```shell
        app.kubernetes.io/name: agenda-service
        app.kubernetes.io/part-of: agenda-service
        app.kubernetes.io/version: {{ .Values.services.tag  }}
```

These annotations allow tools to understand a bit more about our workloads, for example, in this case tools know that the service name is `agenda-service`. We can use the `app.kubernetes.io/part-of` to aggregate multiple services to be part of the same applicaiton. For this example, we wanted to keep each service as a separate entity so we can monitor each individually. 

On this example we will be also using a KeptnTask, that enable us to perform pre and post deployment tasks. You can deploy the following extremely simple example `KeptnTaskDefinition`:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: stdout-notification
spec:
  function:
    inline:
      code: |
        let context = Deno.env.get("CONTEXT");
        console.log("Keptn Task Executed with context: \n");
        console.log(context);

```

As you can see this task is only printing the context from its execution, but here is where you can build any integration with other projects or call external systems. If you look at the Keptn examples, you will find KeptnTaskDefinition to connect for example to Slack, run load tests or to validate that deployments are working as expected after being updated. These tasks uses [Deno](https://deno.land/), a secure JavaScript runtime with Typescript supported out-of-the-box, Python 3 or directly a container image. 

By running: 

```shell
kubectl apply -f keptntask.yaml
```

KeptnTaskDefinitions allow Platform Teams to create reusable tasks that can be hooked into Pre/Post deployment hooks of our applications. By adding the following annotation to our workloads (deployments in this case), Keptn will execute the `stdout-notification` automatically, in this case after performing the deployment (and after any update): 

```shell
  keptn.sh/post-deployment-tasks: stdout-notification
``` 

Let's deploy the Conference application, and lets open Jaeger and Grafana Dashboards. In separate tabs run: 

```shell
make port-forward-jaeger
```

You can point your browser to [http://localhost:16686/](http://localhost:16686/), you should see: 

![jaeger](../imgs/jaeger.png)


and then in a separate terminal: 

```shell
make port-forward-grafana
```

You can point your browser to [http://localhost:3000/](http://localhost:3000/). Use the `admin/admin` credentials and you should see: 

![grafana](../imgs/grafana.png)


Let's now deploy the Conference Application as we did in Chapter 2: 

```shell
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0
```

Check the both Jaeger and Grafana Keptn Dashboards, as by default Keptn Workloads will track the deployment frequency. 

In Grafana go to `Dashboards` -> `Keptn Applications`.  You will see a drop down that allows you to select the different applications services. Check the Notifications Service. Because we have only deployed the first version of the deployment, there is not much to see, but the dashboard will become more interesting after we release new versions of our services.

For example, edit the notifications-service deployment and update the `app.kubernetes.io/version` annotation to have the value `v1.1.0` and update the tag used for the container image to be `v1.1.0`

```shell
kubectl edit deploy conference-notifications-service-deployment
```

After you performed the changes, and the new version is up and running, check the dashboards again. 
In Grafana you will see we are on the second successful deployment, that the average between deployments was 5.83 minutes in my environment, and that `v1.0.0` took 641s while `v1.1.0` took only 40s. There is definitely room for improvement there. 

![grafana](../imgs/grafana-notificatons-service-v1.1.0.png)

If you look at the traces in Jaeger, you will see that the `lifecycle-operator` one of the core components in Keptn is monitoring our deployment resources and performing lifecycle operations, like for example calling pre and post deployments tasks. 

![jager](../imgs/jaeger-notifications-service-v1.1.0.png)

These tasks are executed as Kubernetes Jobs in the same namespace where the workloads are running. You can take a look at the logs from this tasks by tailing the jobs pod's logs. 

```shell
kubectl get jobs
NAME                                   COMPLETIONS   DURATION   AGE
post-stdout-notification-25899-78387   1/1           3s         66m
post-stdout-notification-28367-11337   1/1           4s         61m
post-stdout-notification-54572-93558   1/1           4s         66m
post-stdout-notification-75100-85603   1/1           3s         66m
post-stdout-notification-77674-78421   1/1           3s         66m
post-stdout-notification-93609-30317   1/1           3s         23m
```

The Job with id `post-stdout-notification-93609-30317` was executed after I've performed the update on the Notification Service deployment. 

```shell
> kubectl logs -f post-stdout-notification-93609-30317-vvwp4
Keptn Task Executed with context: 

{"workloadName":"notifications-service-notifications-service","appName":"notifications-service","appVersion":"","workloadVersion":"v1.1.0","taskType":"post","objectType":"Workload"}

```

## Next steps

I strongly recommend you to get more familiar with Keptn Lifecycle Toolkit features and functionalities as what we have seen in this short tutorial are just the basics. Check the concept of [KeptnApplication](https://lifecycle.keptn.sh/docs/concepts/apps/) for more control on how your services are deployed, as Keptn allows you to define fine-grained rules about which services and which versions are allowed to be deployed. 

By grouping multiple services as part of the same Kubernetes Application using the `app.kubernetes.io/part-of` annotation, you can perform pre and post actions on a group of services, allowing you to validate that not only individual services are working as expected but the whole set is. 

