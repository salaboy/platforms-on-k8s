# Chapter 7 :: Shared Applications Concerns

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

In this step-by-step tutorial, we will look into using [Dapr](https://dapr.io) to provide Application-level APIs to solve everyday challenges that most Distributed Applications will face. 

Then we will look at [OpenFeature](https://openfeature.dev), a project that aims to standardize Feature Flags so development teams can keep releasing new features and stakeholders can decide when to enable/disable these features for their customers. 

Because both projects are focused on providing developers with new APIs and tools to use inside their service's code, we will deploy a new version of the application (`v2.0.0`). You can find all the changes required for this new version in the `v2.0.0` of this repository. [You can also compare the differences between the branches here](https://github.com/salaboy/platforms-on-k8s/compare/v2.0.0).


# Installation

You need a Kubernetes Cluster to install [Dapr](https://dapr.io) and `flagd` an [OpenFeature](https://openfeature.dev/) Provider. You can create one using Kubernetes KinD as we did in [Chapter 2](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind)

Then you can install Dapr into the cluster by running: 
```shell
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \
--version=1.11.0 \
--namespace dapr-system \
--create-namespace \
--wait

```

Once Dapr is installed, we can install our Dapr-Enabled and FeatureFlag-Enabled versions of the application `v2.0.0`.

# Running v2.0.0

Now you can install v2.0.0 of the application by running: 

```shell
helm install conference oci://docker.io/salaboy/conference-app --version v2.0.0
```

This version of the Helm Chart installs the same application infrastructure as version `v1.0.0` (PostgreSQL, Redis, and Kafka). Services now interact with Redis and Kafka are now using Dapr APIs. This application version also adds OpenFeature Feature Flags using `flagd`.

# Application Level APIs with Dapr

In version `v2.0.0`, if you list the Application pods now, you will see that each service (agenda, c4p, frontend, and notifications) has a Dapr Sidecar (`daprd`) running alongside the service container (READY 2/2): 

```shell
> kubectl get pods
NAME                                                           READY   STATUS    RESTARTS      AGE
conference-agenda-service-deployment-5dd4bf67b-qkctd           2/2     Running   7 (7s ago)    74s
conference-c4p-service-deployment-57b5985757-tdqg4             2/2     Running   6 (19s ago)   74s
conference-frontend-deployment-69d9b479b7-th44h                2/2     Running   2 (68s ago)   74s
conference-kafka-0                                             1/1     Running   0             74s
conference-notifications-service-deployment-7b6cbf965d-2pdkh   2/2     Running   6 (42s ago)   74s
conference-postgresql-0                                        1/1     Running   0             74s
conference-redis-master-0                                      1/1     Running   0             74s
flagd-6bbdc5d999-c42wk                                         1/1     Running   0             74s
```

Notice the `flagd` container also running. We will cover this in the next section.

From the Dapr perspective the application looks like this: 

![conference-app-with-dapr](imgs/conference-app-with-dapr.png)

The Dapr Sidecar exposes the Dapr Components APIs to enable the application to interact with   Statestore (Redis) and PubSub (Kafka) APIs. 

You can list Dapr Components by running: 

```shell
> kubectl get components
NAME                                   AGE
conference-agenda-service-statestore   30m
conference-conference-pubsub           30m
```

You can describe each component to see its configurations:
```shell
> kubectl describe component conference-agenda-service-statestore
Name:         conference-agenda-service-statestore
Namespace:    default
Labels:       app.kubernetes.io/managed-by=Helm
Annotations:  meta.helm.sh/release-name: conference
              meta.helm.sh/release-namespace: default
API Version:  dapr.io/v1alpha1
Auth:
  Secret Store:  kubernetes
Kind:            Component
Metadata:
  Creation Timestamp:  2023-07-28T08:26:55Z
  Generation:          1
  Resource Version:    4076
  UID:                 b4674825-d298-4ee3-8244-a13cdef8d530
Spec:
  Metadata:
    Name:   keyPrefix
    Value:  name
    Name:   redisHost
    Value:  conference-redis-master.default.svc.cluster.local:6379
    Name:   redisPassword
    Secret Key Ref:
      Key:   redis-password
      Name:  conference-redis
  Type:      state.redis
  Version:   v1
Events:      <none>

```

You can see that the Statestore component is connecting to the Redis instance exposed by this service name `conference-redis-master.default.svc.cluster.local` and using the `conference-redis` secret to obtain the password to connect.

Similarly, the PubSub Dapr Component that is connecting to Kafka: 

```shell
kubectl describe component conference-conference-pubsub 
Name:         conference-conference-pubsub
Namespace:    default
Labels:       app.kubernetes.io/managed-by=Helm
Annotations:  meta.helm.sh/release-name: conference
              meta.helm.sh/release-namespace: default
API Version:  dapr.io/v1alpha1
Kind:         Component
Metadata:
  Creation Timestamp:  2023-07-28T08:26:55Z
  Generation:          1
  Resource Version:    4086
  UID:                 e145bc49-18ff-4390-ad15-dcd9a4275479
Spec:
  Metadata:
    Name:   brokers
    Value:  conference-kafka.default.svc.cluster.local:9092
    Name:   authType
    Value:  none
  Type:     pubsub.kafka
  Version:  v1
Events:     <none>
```

The final piece of the puzzle that allows the Frontend service to receive events that are submitted to the PubSub component is the following Dapr Subscription: 

```shell
> kubectl get subscription
NAME                               AGE
conference-frontend-subscritpion   39m
```

You can also describe this resource to see its configurations: 
```shell
> kubectl describe subscription conference-frontend-subscritpion
Name:         conference-frontend-subscritpion
Namespace:    default
Labels:       app.kubernetes.io/managed-by=Helm
Annotations:  meta.helm.sh/release-name: conference
              meta.helm.sh/release-namespace: default
API Version:  dapr.io/v2alpha1
Kind:         Subscription
Metadata:
  Creation Timestamp:  2023-07-28T08:26:55Z
  Generation:          1
  Resource Version:    4102
  UID:                 9f748cb0-125a-4848-bd39-f84e37e41282
Scopes:
  frontend
Spec:
  Bulk Subscribe:
    Enabled:   false
  Pubsubname:  conference-conference-pubsub
  Routes:
    Default:  /api/new-events/
  Topic:      events-topic
Events:       <none>
```

As you can see, this subscription forwards events to the route `/api/new-events/` for the Dapr applications listed in the `Scopes` section, only the `frontend` application. The Frontend Application only needs to expose the `/api/new-events/` endpoint to receive events, in this case the Dapr Sidecar (`daprd`) waits for incoming messages on the PubSub component called `conference-conference-pubsub` and forwards all messages to the Application Endpoint. 

This version of the application removes application dependencies such as the Kafka Client from all services and the Redis Client from the Agenda Service. 


![services without deps](imgs/conference-app-dapr-no-deps.png)

Besides removing the dependencies and making these containers smaller, by consuming the Dapr Components API, we enable the platform team to define how these components will be configured and against which infrastructural components. Configuring the same application to use Google Cloud Platform managed services such as [Google PubSub](https://cloud.google.com/pubsub) or the [MemoryStore databases](https://cloud.google.com/memorystore) will not require any changes in the application code or adding any new dependencies, just new Dapr Component configurations.

![in gcp](imgs/conference-app-dapr-and-gcp.png)


Finally, because this is all about enabling developers with Application Level APIs, let's look at how this looks from the application's service perspective. Because the services are written in Go, I've decided to add the Dapr Go SDK (which is optional). 

When the Agenda Service wants to store or read data from the Dapr Statestore component, it can use the Dapr Client to perform these operations, for example [reading values from the Statestore looks like this](https://github.com/salaboy/platforms-on-k8s/blob/v2.0.0/conference-application/agenda-service/agenda-service.go#L136C2-L136C116): 

```golang
agendaItemsStateItem, err := s.APIClient.GetState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), nil)
```

The `APIClient` reference is just a [Dapr Client instance, which was initialized here](https://github.com/salaboy/platforms-on-k8s/blob/v2.0.0/conference-application/agenda-service/agenda-service.go#L397)

All the application needs to know is the Statestore name (`STATESTORE_NAME`) and the key (`KEY`) to locate the data that wants to be retrieved.

When the application wants to [store state into the Statestore it looks like this](https://github.com/salaboy/platforms-on-k8s/blob/v2.0.0/conference-application/agenda-service/agenda-service.go#L197C2-L199C3):

```golang
if err := s.APIClient.SaveState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), jsonData, nil); err != nil {
		...
}
```  

Finally, if the application code wants to [publish a new event to the PubSub component](https://github.com/salaboy/platforms-on-k8s/blob/v2.0.0/conference-application/agenda-service/agenda-service.go#L225), it will look like this: 

```golang
if err := s.APIClient.PublishEvent(ctx, PUBSUB_NAME, PUBSUB_TOPIC, eventJson); err != nil {
			...
}

```

As we have seen, Dapr provides Application-Level APIs for application developers to use regarding the programming language that they are using. These APIs abstract away the complexities of setting up and managing application infrastructure components, enabling platform teams to have more flexibility to focus on fine-tuning how applications will interact with them without pushing teams to change the application source code. 

Next, let's talk about feature flags. This topic involves not only developers but also enables product managers and roles closer to the business to decide when certain features should be exposed.


## Feature Flags for everyone

The [OpenFeature](https://openfeature.dev/) project aims to standardize how to consume Feature Flags from applications written in different languages. 

In this short tutorial, we will look at how the Conference Application `v2.0.0` uses Open Feature and, more specifically, the `flagd` provider to enable feature flags across all application services. For this example, to keep it simple, I've used the `flagd` provider that allows us to define our feature flag configurations inside a Kubernetes `ConfigMap`.

![openfeature](imgs/conference-app-openfeature.png)

In the same way as Dapr APIs, the idea here is to have a consistent experience no matter which provider we have selected. If the platform team wants to switch providers, for example, LaunchDarkly or Split, there will be no need to change how features are fetched or evaluated. The platform team will be able to swap providers to whichever they think is best. 

`v2.0.0` created a ConfigMap called `flag-configuration` containing the Feature Flag that the application's services will use. 

You can get the flag configuration json file included in the ConfigMap by running: 

```shell
kubectl get cm flag-configuration -o go-template='{{index .data "flag-config.json"}}'
```

You should see the following output: 

```json
{
  "flags": {
    "debugEnabled": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "off"
    },
    "callForProposalsEnabled": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on"  
    },
    "eventsEnabled": {
      "state": "ENABLED",
      "variants": {
        "all": {
          "agenda-service": true,
          "notifications-service": true,
          "c4p-service": true
        },
        "decisions-only": {
          "agenda-service": false,
          "notifications-service": false,
          "c4p-service": true
        },
        "none": {
          "agenda-service": false,
          "notifications-service": false,
          "c4p-service": false
        }
      },
      "defaultVariant": "all"
    }
  }
}
```

There are three feature flags defined for this example: 
- `debugEnabled` is a boolean flag that allows us to turn on and off the Debug tab in the back office of the application. This replaces the need for an environment variable that we used in `v1.0.0`. We can turn the debug section on and off without restarting the application frontend container.
- `callForProposalsEnabled` This boolean flag allows us to disable the **Call for Proposals** section of the application. As conferences have a window to allow potential speakers to submit proposals, when that period is over, this section can be hidden away. Having to release a specific version to just turn off that section would be too complicated to manage, hence having a feature flag for this makes a lot of sense. We can make this change without the need to restart the application frontend container.
- `eventsEnabled` is an object feature flag, this means that it contains a structure and allows teams to define complex configurations. In this case, I've defined different profiles of flags to configure which service can emit events (Events tab in the application back office). By default all services emit events, but by changing the `defaultVariant` value to `none` we can disable events for all services, without the need to restart any container.


You can patch the ConfigMap to turn on the debug feature by following these steps. First, fetch the content of the `flag-config.json` file located inside the ConfigMap and store it locally.
```shell
kubectl get cm flag-configuration -o go-template='{{index .data "flag-config.json"}}' > flag-config.json
```

Modify the content of this file, for example turn on the debug flag: 

```json
{
  "flags": {
    "debugEnabled": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
    **"defaultVariant": "on"**
    },
    ...
```
Then patch the existing `ConfigMap`: 

```shell
kubectl create cm flag-configuration --from-file=flag-config.json=flag-config.json --dry-run=client -o yaml | kubectl patch cm flag-configuration --type merge --patch-file /dev/stdin
```

After 20 seconds or so, you should see the Debug Tab in the application's back office section

![debug feature flag](imgs/feature-flag-debug-tab.png)

You can see that feature flags are now also displayed in this tab. 

Now submit a new proposal and approve it. You will see that in the `Events` tab, Events will be displayed. 

![events for approved proposal](imgs/feature-flag-events-for-proposal.png)


If you repeat the previous process and change the `eventsEnabled` feature flag to `"defaultVariant": "none"`, all services will stop emitting events. Submit a new proposal from the application user interface and approve it, then check the `Events` tab to validate that no event has been emitted. Remember that when changing the `flag-configuration` ConfigMap, `flagd` needs to wait around 10 seconds to refresh the content of the ConfigMap. If you have the Debug tab enabled you can refresh that screen until you see that the value has changed. 

**Notice that this feature flag is being consumed by all services that evaluate the flag before sending any event. **

Finally, if you change the `callForProposalsEnabled` feature flag `"defaultVariant": "off"`, the Call for Proposal menu option will disappear from the application frontend. 

![no call for proposals feature flag](imgs/feature-flag-no-c4p.png)


While we are still using a `ConfigMap` to store the feature flags configurations we have achieved some important improvements that enable teams to go faster. Developers can keep releasing new features to their application services that then product managers (or stakeholders) can decide when to enable/disable. Platform Teams can define where the feature flags will be stored (a managed service or local storage). By using a standard specification driven by a community composed of Feature Flag vendors enables our application development teams to make use of feature flags without defining all the technical aspects required to implement these mechanisms in-house. 

In this example, we haven't used more advanced features to evaluate feature flags like [context-based evaluations](https://openfeature.dev/docs/reference/concepts/evaluation-context#providing-evaluation-context), which can use for example, the geo-location of the user to provide different values for the same feature flag, or [targetting keys](https://openfeature.dev/docs/reference/concepts/evaluation-context#targeting-key). It is up to the reader to go deeper into OpenFeature capabilities as well as which other [Open Feature flag providers](https://openfeature.dev/docs/reference/concepts/provider) are available.  

## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```shell
kind delete clusters dev
```


## Next Steps

A natural next step would be to run `v2.0.0` against infrastructure provisioned by Crossplane, as we did in Chapter 5. This can be managed by our Platform Walking skeleton, which will be in charge of configuring the Conference Application Helm chart to connect to the provisioned infrastructure by wiring Kubernetes resources together. If you are interested in this topic, I've written a blog post about why tools like Crossplane and Dapr are meant to work together: [https://blog.crossplane.io/crossplane-and-dapr/](https://blog.crossplane.io/crossplane-and-dapr/).

Another simple but very useful extension to the application code would be to make sure that the Call for Proposals Service reads the `callForProposalsEnabled` feature flag and returns meaningful errors when this feature is disabled. The current implementation only removes the `Call for Proposals` menu entry, meaning that if you send a `curl` request to the APIs, the functionality should still work. 

## Sum up and Contribute

In this tutorial, we have looked into Application-level APIs using Dapr and Feature Flags using OpenFeature. Application-level APIs like the ones exposed by Dapr Components can be leveraged by application development teams, as most applications will be interested in storing and reading state, emitting, and consuming events and resiliency policies for service-to-service communications. Feature Flags can also help to speed up the development and release process by enabling developers to keep releasing features while other stakeholders can decide when these features are enabled or disabled. 


Do you want to improve this tutorial? Create an issue, drop me a message on [Twitter](https://twitter.com/salaboy), or send a Pull Request.
