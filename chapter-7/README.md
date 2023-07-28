# Chapter 7 :: Shared Applications Concerns


# Installation

You need a Kubernetes Cluster to install Dapr. You can create one using Kubernetes KinD as we did in [Chapter 2](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind)

Then you can install Dapr into the cluster by running: 
```
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \
--version=1.11.0 \
--namespace dapr-system \
--create-namespace \
--wait

```

Once Dapr is installed, we can install our Dapr-Enabled version of the application v2.0.0.

# Running v2.0.0

Now you can install v2.0.0 of the application by running: 

```
helm install conference oci://docker.io/salaboy/conference-app --version v2.0.0
```

This version of the Helm Chart installs the same application infrastructure as before (PostgreSQL, Redis, and Kafka). 
Services now interact with Redis and Kafka are now using Dapr APIs. 

If you list the Application pods now, you will see that each service has a Dapr Sidecar (`daprd`) running alongside the service container: 

```
> kubectl get pods
NAME                                                           READY   STATUS    RESTARTS      AGE
conference-agenda-service-deployment-5dd4bf67b-2mnm9           2/2     Running   0             31m
conference-c4p-service-deployment-57b5985757-zz8m6             2/2     Running   1 (31m ago)   31m
conference-frontend-deployment-69d9b479b7-m7r5d                2/2     Running   0             31m
conference-kafka-0                                             1/1     Running   0             31m
conference-notifications-service-deployment-7b6cbf965d-m8w27   2/2     Running   1 (31m ago)   31m
conference-postgresql-0                                        1/1     Running   0             31m
conference-redis-master-0                                      1/1     Running   0             31m
```

The Dapr Sidecar uses Dapr Commponents to enable the Dapr APIs to connect with Redis (Statestore) and Kafka (PubSub). 

You can list Dapr Components by running: 

```
> kubectl get components
NAME                                   AGE
conference-agenda-service-statestore   30m
conference-conference-pubsub           30m
```

You can describe each component to see its configurations:
```
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

```
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

```
> kubectl get subscription
NAME                               AGE
conference-frontend-subscritpion   39m
```

You can also describe this resource to see its configurations: 
```
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

As you can see, this subscription forwards events to the route `/api/new-events/` for the Dapr applications listed in the `Scopes` section, in this case only the `frontend` application. 
