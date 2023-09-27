# Release Strategies with Argo Rollouts

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

In this tutorial, we will be looking at Argo Rollout's built-in mechanisms to implement release strategies. We will also look into the Argo Rollouts Dashboard that allows teams to promote new versions without using the terminal (`kubectl`).

## Installation

You need a Kubernetes Cluster to install [Argo Rollouts](https://argoproj.github.io/rollouts/). You can create one using Kubernetes KinD as we did in [Chapter 2](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind)

Once you have the cluster, we can install Argo Rollouts by running:

```shell
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml

```

or by following the [official documentation that you can find here](https://argoproj.github.io/argo-rollouts/installation/#controller-installation).

You also need to install the [Argo Rollouts `kubectl` plugin](https://argoproj.github.io/argo-rollouts/installation/#kubectl-plugin-installation) 

Once you have the plugin, you can start a local version of the Argo Rollouts Dashboard by running in a new terminal the following command:

```shell
kubectl argo rollouts dashboard
```

Then you can access the dashboard by pointing your browser to [http://localhost:3100/rollouts](http://localhost:3100/rollouts)


![argo rollouts dashboard empty](../imgs/argo-rollouts-dashboard-empty.png)

## Canary Releases

Let's create an Argo Rollout resource to implement a Canary Release on the Notification Service of the Conference Application. You can find the [full definition here](canary-release/rollout.yaml).

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: notifications-service-canary
spec:
  replicas: 3
  strategy:
    canary:
      steps:
      - setWeight: 25
      - pause: {}
      - setWeight: 75
      - pause: {duration: 10}
  revisionHistoryLimit: 2
  selector:
    matchLabels:
      app: notifications-service
  template:
    metadata:
      labels:
        app: notifications-service
    spec:
      containers:
      - name: notifications-service
        image: salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0
        env: 
          - name: KAFKA_URL
            value: kafka.default.svc.cluster.local
          ... 
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        resources:
          limits:
            cpu: "1"
            memory: 256Mi
          requests:
            cpu: "0.1"
            memory: 256Mi

```

The `Rollout` resource replaces our Kubernetes `Deployment` resource. This means we still need to create a Kubernetes Service and an Ingress Resource to route traffic to our Notification Service instance. Notice that we are defining three replicas for the Notification Service.

The previous `Rollout` defines a canary release with two steps: 

```yaml
strategy:
    canary:
      steps:
      - setWeight: 25
      - pause: {}
      - setWeight: 75
      - pause: {duration: 10}
```

First, it will set the traffic split to 25 percent and wait for the team to test the new version (the `pause` step), then after we manually signal that we want to continue the rollout will move to 75 percent to the new version to finally pause for 10 seconds and then move to 100 percent. 

Before applying the Rollout, Service, and Ingress resources located in the `canary-release/` directory, let's install Kafka for the Notification Service to connect. 

```shell
helm install kafka oci://registry-1.docker.io/bitnamicharts/kafka --version 22.1.5 --set "provisioning.topics[0].name=events-topic" --set "provisioning.topics[0].partitions=1" --set "persistence.size=1Gi" 

```

Now when Kafka is running, let's apply all the resources in the `canary-releases/` directory: 

```shell
kubectl apply -f canary-release/
```

Using the argo rollouts plugin you can watch the rollout from the terminal: 

```shell
kubectl argo rollouts get rollout notifications-service-canary --watch
```

You should see something like this: 

```shell
Name:            notifications-service-canary
Namespace:       default
Status:          ✔ Healthy
Strategy:        Canary
  Step:          4/4
  SetWeight:     100
  ActualWeight:  100
Images:          salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0 (stable)
Replicas:
  Desired:       3
  Current:       3
  Updated:       3
  Ready:         3
  Available:     3

NAME                                                      KIND        STATUS     AGE  INFO
⟳ notifications-service-canary                            Rollout     ✔ Healthy  80s  
└──# revision:1                                                                       
   └──⧉ notifications-service-canary-7f6b88b5fb           ReplicaSet  ✔ Healthy  80s  stable
      ├──□ notifications-service-canary-7f6b88b5fb-d86s2  Pod         ✔ Running  80s  ready:1/1
      ├──□ notifications-service-canary-7f6b88b5fb-dss5c  Pod         ✔ Running  80s  ready:1/1
      └──□ notifications-service-canary-7f6b88b5fb-tw8fj  Pod         ✔ Running  80s  ready:1/1
```

As you can see, because we just created the Rollouts, three replicas are created and all the traffic is being routed to this initial `revision:1`, and the Status is set to `Healthy`.

Let's update the Notification Service version to `v1.1.0` by running: 

```shell
kubectl argo rollouts set image notifications-service-canary \
  notifications-service=salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.1.0
```

Now you see the second revision (revision:2) created: 

```shell
Name:            notifications-service-canary
Namespace:       default
Status:          ॥ Paused
Message:         CanaryPauseStep
Strategy:        Canary
  Step:          1/4
  SetWeight:     25
  ActualWeight:  25
Images:          salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0 (stable)
                 salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.1.0 (canary)
Replicas:
  Desired:       3
  Current:       4
  Updated:       1
  Ready:         4
  Available:     4

NAME                                                      KIND        STATUS     AGE    INFO
⟳ notifications-service-canary                            Rollout     ॥ Paused   4m29s  
├──# revision:2                                                                         
│  └──⧉ notifications-service-canary-68fd6b4ff9           ReplicaSet  ✔ Healthy  14s    canary
│     └──□ notifications-service-canary-68fd6b4ff9-jrjxh  Pod         ✔ Running  14s    ready:1/1
└──# revision:1                                                                         
   └──⧉ notifications-service-canary-7f6b88b5fb           ReplicaSet  ✔ Healthy  4m29s  stable
      ├──□ notifications-service-canary-7f6b88b5fb-d86s2  Pod         ✔ Running  4m29s  ready:1/1
      ├──□ notifications-service-canary-7f6b88b5fb-dss5c  Pod         ✔ Running  4m29s  ready:1/1
      └──□ notifications-service-canary-7f6b88b5fb-tw8fj  Pod         ✔ Running  4m29s  ready:1/1
```

Now the Rollout stops at step 1, where only 25 percent of the traffic is routed to `revision:2` and the status is set to `Pause`.

Feel free to hit the `service/info` endpoint to see which version is answering your requests:

```shell
curl localhost/service/info
```
Roughly, one in four requests should be answered by version `v1.1.0`:

```shell
> curl localhost/service/info | jq

{
    "name":"NOTIFICATIONS",
    "version":"1.0.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "podName":"notifications-service-canary-7f6b88b5fb-tw8fj",
    "podNamespace":"default",
    "podNodeName":"dev-worker2",
    "podIp":"10.244.3.3",
    "podServiceAccount":"default"
}

> curl localhost/service/info | jq

{
    "name":"NOTIFICATIONS",
    "version":"1.0.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "podName":"notifications-service-canary-7f6b88b5fb-tw8fj",
    "podNamespace":"default",
    "podNodeName":"dev-worker2",
    "podIp":"10.244.3.3",
    "podServiceAccount":"default"
}

> curl localhost/service/info | jq

{
    "name":"NOTIFICATIONS-IMPROVED",
    "version":"1.1.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/v1.1.0/conference-application/notifications-service",
    "podName":"notifications-service-canary-68fd6b4ff9-jrjxh",
    "podNamespace":"default",
    "podNodeName":"dev-worker",
    "podIp":"10.244.2.4",
    "podServiceAccount":"default"
}

> curl localhost/service/info | jq

{
    "name":"NOTIFICATIONS",
    "version":"1.0.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "podName":"notifications-service-canary-7f6b88b5fb-tw8fj",
    "podNamespace":"default",
    "podNodeName":"dev-worker2",
    "podIp":"10.244.3.3",
    "podServiceAccount":"default"
}

```

Also check the Argo Rollouts Dashboards now, it should show the Canary Release: 

![canary release in dashboard](../imgs/argo-rollouts-dashboard-canary-1.png)

You can move the canary forward by using the promote command or the promote button in the Dashboard. The command looks like this: 

```shell
kubectl argo rollouts promote notifications-service-canary
```

That should move the canary to 75% of the traffic and after 10 more seconds it should be at 100%, as the last puase step is just for 10 seconds. You should see in the terminal: 

```shell
Name:            notifications-service-canary
Namespace:       default
Status:          ✔ Healthy
Strategy:        Canary
  Step:          4/4
  SetWeight:     100
  ActualWeight:  100
Images:          salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.1.0 (stable)
Replicas:
  Desired:       3
  Current:       3
  Updated:       3
  Ready:         3
  Available:     3

NAME                                                      KIND        STATUS        AGE  INFO
⟳ notifications-service-canary                            Rollout     ✔ Healthy     16m  
├──# revision:2                                                                          
│  └──⧉ notifications-service-canary-68fd6b4ff9           ReplicaSet  ✔ Healthy     11m  stable
│     ├──□ notifications-service-canary-68fd6b4ff9-jrjxh  Pod         ✔ Running     11m  ready:1/1
│     ├──□ notifications-service-canary-68fd6b4ff9-q4zgj  Pod         ✔ Running     51s  ready:1/1
│     └──□ notifications-service-canary-68fd6b4ff9-fctjv  Pod         ✔ Running     46s  ready:1/1
└──# revision:1                                                                          
   └──⧉ notifications-service-canary-7f6b88b5fb           ReplicaSet  • ScaledDown  16m  

```
And in the Dashboard: 

![canary promoted](../imgs/argo-rollouts-dashboard-canary-2.png)

Now all requests should be answered by `v1.1.0`:

```shell

> curl localhost/service/info

{
    "name":"NOTIFICATIONS-IMPROVED",
    "version":"1.1.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/v1.1.0/conference-application/notifications-service",
    "podName":"notifications-service-canary-68fd6b4ff9-jrjxh",
    "podNamespace":"default",
    "podNodeName":"dev-worker",
    "podIp":"10.244.2.4",
    "podServiceAccount":"default"
}

> curl localhost/service/info

{
    "name":"NOTIFICATIONS-IMPROVED",
    "version":"1.1.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/v1.1.0/conference-application/notifications-service",
    "podName":"notifications-service-canary-68fd6b4ff9-jrjxh",
    "podNamespace":"default",
    "podNodeName":"dev-worker",
    "podIp":"10.244.2.4",
    "podServiceAccount":"default"
}

> curl localhost/service/info

{
    "name":"NOTIFICATIONS-IMPROVED",
    "version":"1.1.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/v1.1.0/conference-application/notifications-service",
    "podName":"notifications-service-canary-68fd6b4ff9-jrjxh",
    "podNamespace":"default",
    "podNodeName":"dev-worker",
    "podIp":"10.244.2.4",
    "podServiceAccount":"default"
}

```

Before moving forward to Blue/Green deployments let's clean up the canary rollout by running: 

```shell
kubectl delete -f canary-release/
```

## Blue/Green Deployments 

With Blue/Green deployments we want to have two versions of our service running at the same time. The Blue (active) version that all the users will access and the Green (preview) version that internal teams can use to test new features and changes. 

Argo Rollouts provide a blueGreen strategy out-of-the-box: 

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: notifications-service-bluegreen
spec:
  replicas: 2
  revisionHistoryLimit: 2
  selector:
    matchLabels:
      app: notifications-service
  template:
    metadata:
      labels:
        app: notifications-service
    spec:
      containers:
      - name: notifications-service
        image: salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0
        env: 
          - name: KAFKA_URL
            value: kafka.default.svc.cluster.local
          ..
  strategy:
    blueGreen: 
      activeService: notifications-service-blue
      previewService: notifications-service-green
      autoPromotionEnabled: false
```

Once again, we are using our Notifications Service to test the Rollout mechanism. Here we have defined a Blue/Green deployment for the Notification Service which points to two existing Kubernetes Services: `notifications-service-blue` and `notifications-service-green`. Notice that the `autoPromotionEnabled` flag is set to `false`, this stops the promotion from happening automatically when the new version is ready.

Check that you have Kafka already running from the previous section (Canary releases) and apply all the resources located inside the `blue-green/` directory: 

```shell
kubectl apply -f blue-green/
```

This is creating the `Rollout` resource, two Kubernetes Services and two Ingress resources, one of for the Blue Service which forwards traffic from `/` and one for the Green service that forward traffic from `/preview/`

You can monitor the Rollout in the terminal by running: 

```shell
kubectl argo rollouts get rollout notifications-service-bluegreen --watch
```

You should see something like this: 

```
Name:            notifications-service-bluegreen
Namespace:       default
Status:          ✔ Healthy
Strategy:        BlueGreen
Images:          salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0 (stable, active)
Replicas:
  Desired:       2
  Current:       2
  Updated:       2
  Ready:         2
  Available:     2

NAME                                                         KIND        STATUS     AGE    INFO
⟳ notifications-service-bluegreen                            Rollout     ✔ Healthy  3m16s  
└──# revision:1                                                                            
   └──⧉ notifications-service-bluegreen-56bb777689           ReplicaSet  ✔ Healthy  2m56s  stable,active
      ├──□ notifications-service-bluegreen-56bb777689-j5ntk  Pod         ✔ Running  2m56s  ready:1/1
      └──□ notifications-service-bluegreen-56bb777689-qzg9l  Pod         ✔ Running  2m56s  ready:1/1

```

We get two replicas of our Notification Service up and running. If we curl `localhost/service/info` we should get the Notification Service `v1.0.0` information: 

```shell
> curl localhost/service/info | jq

{
    "name":"NOTIFICATIONS",
    "version":"1.0.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "podName":"notifications-service-canary-7f6b88b5fb-tw8fj",
    "podNamespace":"default",
    "podNodeName":"dev-worker2",
    "podIp":"10.244.3.3",
    "podServiceAccount":"default"
}
```


And the Argo Rollouts Dashboard should show us our Blue/Green Rollout: 

![blue green 1](../imgs/argo-rollouts-dashboard-bluegree-1.png)

As we did with the Canary Release, we can update our Rollout configuration, in this case setting the image for version `v1.1.0`.

```shell
kubectl argo rollouts set image notifications-service-bluegreen \
  notifications-service=salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.1.0
```

Now you should see in the terminal both versions of the Notification Service running in parallel: 

```shell
Name:            notifications-service-bluegreen
Namespace:       default
Status:          ॥ Paused
Message:         BlueGreenPause
Strategy:        BlueGreen
Images:          salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0 (stable, active)
                 salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.1.0 (preview)
Replicas:
  Desired:       2
  Current:       4
  Updated:       2
  Ready:         2
  Available:     2

NAME                                                         KIND        STATUS     AGE    INFO
⟳ notifications-service-bluegreen                            Rollout     ॥ Paused   8m54s  
├──# revision:2                                                                            
│  └──⧉ notifications-service-bluegreen-645d484596           ReplicaSet  ✔ Healthy  16s    preview
│     ├──□ notifications-service-bluegreen-645d484596-ffhsm  Pod         ✔ Running  16s    ready:1/1
│     └──□ notifications-service-bluegreen-645d484596-g2zr4  Pod         ✔ Running  16s    ready:1/1
└──# revision:1                                                                            
   └──⧉ notifications-service-bluegreen-56bb777689           ReplicaSet  ✔ Healthy  8m34s  stable,active
      ├──□ notifications-service-bluegreen-56bb777689-j5ntk  Pod         ✔ Running  8m34s  ready:1/1
      └──□ notifications-service-bluegreen-56bb777689-qzg9l  Pod         ✔ Running  8m34s  ready:1/1
```

Both `v1.0.0` and `v1.1.0` are running and Healthy, but the Status of the BlueGreen Rollout is in Pause, as it will keep running both versions until the team responsible for validating the `preview` / `green` version is ready for the prime time. 

Check the Argo Rollouts Dashboard, it should show both versions running too:

![blue green 2](../imgs/argo-rollouts-dashboard-bluegree-2.png)

At this point, you can send requests to both services by using the Ingress routes that we defined. You can curl `localhost/service/info` to hit the Blue service (stable service) and curl `localhost/preview/service/info` to hit the Green service (preview service).

```shell
> curl localhost/service/info

{
    "name":"NOTIFICATIONS",
    "version":"1.0.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service",
    "podName":"notifications-service-canary-7f6b88b5fb-tw8fj",
    "podNamespace":"default",
    "podNodeName":"dev-worker2",
    "podIp":"10.244.3.3",
    "podServiceAccount":"default"
}
```

And now let's check Green Service: 

```shell
> curl localhost/green/service/info

{
    "name":"NOTIFICATIONS-IMPROVED",
    "version":"1.1.0",
    "source":"https://github.com/salaboy/platforms-on-k8s/tree/v1.1.0/conference-application/notifications-service",
    "podName":"notifications-service-canary-68fd6b4ff9-jrjxh",
    "podNamespace":"default",
    "podNodeName":"dev-worker",
    "podIp":"10.244.2.4",
    "podServiceAccount":"default"
}
```

If we are happy with the results we can promote our Green Service to be our new stable service, we do this by hitting the Promote button in the Argo Rollouts Dashboard or by running the following command: 

```shell
kubectl argo rollouts promote notifications-service-bluegreen
```

You should see in the terminal: 

```shell
Name:            notifications-service-bluegreen
Namespace:       default
Status:          ✔ Healthy
Strategy:        BlueGreen
Images:          salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0
                 salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.1.0 (stable, active)
Replicas:
  Desired:       2
  Current:       4
  Updated:       2
  Ready:         2
  Available:     2

NAME                                                         KIND        STATUS     AGE    INFO
⟳ notifications-service-bluegreen                            Rollout     ✔ Healthy  2m44s  
├──# revision:2                                                                            
│  └──⧉ notifications-service-bluegreen-645d484596           ReplicaSet  ✔ Healthy  2m27s  stable,active
│     ├──□ notifications-service-bluegreen-645d484596-fnbg7  Pod         ✔ Running  2m27s  ready:1/1
│     └──□ notifications-service-bluegreen-645d484596-ntcbf  Pod         ✔ Running  2m27s  ready:1/1
└──# revision:1                                                                            
   └──⧉ notifications-service-bluegreen-56bb777689           ReplicaSet  ✔ Healthy  2m44s  delay:9s
      ├──□ notifications-service-bluegreen-56bb777689-k6qxk  Pod         ✔ Running  2m44s  ready:1/1
      └──□ notifications-service-bluegreen-56bb777689-vzsw7  Pod         ✔ Running  2m44s  ready:1/1
```

Now the stable service is `revision:2`. You will see that Argo Rollouts will keep `revision:1` active for a while, just in case that we want to revert back, but after a few seconds, it will be downscaled. 

Check the Dashboard to see that our Rollout is in `revision:2` as well:

![rollout promoted](../imgs/argo-rollouts-dashboard-bluegree-3.png)


If you made it this far, you implemented Canary Releases and Blue/Green Deployments using Argo Rollouts! 


## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```shell
kind delete clusters dev
```

