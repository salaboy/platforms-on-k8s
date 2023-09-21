# Multi-Cloud (App) Infrastructure

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

This step-by-step tutorial use Crossplane to Provision the Redis, PostgreSQL, and Kafka instances for our application services. 

Using Crossplane and Crossplane Compositions, we aim to unify how these components are provisioned, hiding away where these components are for the end users (application teams).

Application teams should be able to request these resources using a declarative approach as with any other Kubernetes Resource. This enables teams to use Environment Pipelines to configure both the application services and the application infrastructure components needed by the application.


## Installing Crossplane

To install Crossplane, you need to have a Kubernetes Cluster; you can create one using KinD as we did for you [Chapter 2](../chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind). 

Let's install [Crossplane](https://crossplane.io) into its own namespace using Helm: 

```shell
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --wait
```

Install the `kubectl crossplane` plugin: 

```shell
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh
sudo mv kubectl-crossplane /usr/local/bin
```

Then install the Crossplane Helm provider: 
```shell
kubectl crossplane install provider crossplane/provider-helm:v0.10.0
```

We need the correct `ServiceAccount` to create a new `ClusterRoleBinding` so the Helm Provider can install Charts on our behalf. 

```shell
SA=$(kubectl -n crossplane-system get sa -o name | grep provider-helm | sed -e 's|serviceaccount\/|crossplane-system:|g')
kubectl create clusterrolebinding provider-helm-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
```

```shell
kubectl apply -f crossplane/helm-provider-config.yaml
```


After a few seconds, if you check the configured providers, you should see the Helm `INSTALLED` and `HEALTHY`: 

```shell
> kubectl get providers.pkg.crossplane.io
NAME                             INSTALLED   HEALTHY   PACKAGE                               AGE
crossplane-provider-helm         True        True      crossplane/provider-helm:v0.10.0      49s
```

Now we are ready to install our Databases and Message Brokers Crossplane compositions to provide all the components our application needs.


## App Infrastructure on demand using Crossplane Compositions

We need to install our Crossplane Compositions for our Key-Value Database (Redis), our SQL Database (PostgreSQL), and our Message Broker(Kafka). 

```shell
kubectl apply -f resources/
```

The Crossplane Composition resource (`app-database-redis.yaml`) defines which cloud resources need to be created and how they need to be configured together. The Crossplane Composite Resource Definition (XRD) (`app-database-resource.yaml`) defines a simplified interface that enables application development teams to quickly request new databases by creating resources of this type.

Check the [resources/](resources/) directory for the Compositions and the Composite Resource Definitions (XRDs). 


### Let's provision Application Infrastructure

We can provision a new Key-Value Database for our team to use by executing the following command: 

```shell
kubectl apply -f my-db-keyvalue.yaml
```

The `my-db-keyvalue.yaml` resource looks like this: 

```yaml
apiVersion: salaboy.com/v1alpha1
kind: Database
metadata:
  name: my-db-keyvalue
spec:
  compositionSelector:
    matchLabels:
      provider: local
      type: dev
      kind: keyvalue
  parameters: 
    size: small
```

Notice that we are using the labels `provider: local`, `type: dev`, and `kind: keyvalue`. This allows Crossplane to find the right composition based on the labels. In this case, the Helm Provider created a local Redis instance.

You can check the database status using:

```shell
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   97s
```

You can check that a new Redis instance was created in the `default` namespace. 

You can follow the same steps to provision a PostgreSQL database by running: 

```shell
kubectl apply -f my-db-sql.yaml
```

You should see now two `dbs`

```shell
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   2m
my-db-sql         small   false      sql        True     False   sql.db.local.salaboy.com        5s
```


You can now check that there are two Pods running, one for each database:

```shell
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          3m40s
my-db-sql-postgresql-0           1/1     Running   0          104s
```

There should be 4 Kubernetes Secrets (two for our two helm releases and two containing the credentials to connect to the newly created instances):

```shell
> kubectl get secret
NAME                                    TYPE                 DATA   AGE
my-db-keyavalue-redis                   Opaque               1      2m32s
my-db-sql-postgresql                    Opaque               1      36s
sh.helm.release.v1.my-db-keyavalue.v1   helm.sh/release.v1   1      2m32s
sh.helm.release.v1.my-db-sql.v1         helm.sh/release.v1   1      36s
```

We can do the same to provision a new instance of our Kafka Message Broker: 

```shell
kubectl apply -f my-messagebroker-kafka.yaml
```

And then list with: 

```shell
> kubectl get mbs
NAME          SIZE    KIND    SYNCED   READY   COMPOSITION                  AGE
my-mb-kafka   small   kafka   True     True    kafka.mb.local.salaboy.com   2m51s
```

Kafka doesn't require creating any secret when using its default configuration. 

You should see three running pods (one for Kafka, one for Redis, and one for PostgreSQL).

```shell
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          113s
my-db-sql-postgresql-0           1/1     Running   0          108s
my-mb-kafka-0                    1/1     Running   0          100s
```

**Note**: if you are deleting and recreating databases or message brokers using the same resource name, remember to delete the PersistentVolumeClaims, as these resources don't get removed when you delete the Database or MessageBroker resources. 

You can now create as many database or message broker instances as your cluster resources can handle! 

## Let's deploy our Conference Application

Ok, now that we have our two databases and our message broker running, we need to make sure that our application services connect to these instances. The first thing that we need to do is to disable the Helm dependencies defined in the Conference Application chart so that when the application gets installed, don't install the databases and the message broker. We can do this by setting the `install.infrastructure` flag to `false`.

For that, we will use the `app-values.yaml` file containing the configurations for the services to connect to our newly created databases:

```shell
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0 -f app-values.yaml
```

The `app-values.yaml` content looks like this: 
```yaml
install:
  infrastructure: false
frontend:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
agenda:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
  redis: 
    host: my-db-keyavalue-redis-master.default.svc.cluster.local
    secretName: my-db-keyavalue-redis
c4p: 
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
  postgresql:
    host: my-db-sql-postgresql.default.svc.cluster.local
    secretName: my-db-sql-postgresql
notifications: 
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
```

Notice that the `app-values.yaml` file relies on the names that we specified for our databases (`my-db-keyavalue` and `my-db-sql`) and our message brokers (`my-mb-kafka`) in the example files. If you request other databases and message brokers with different names you will need to adapt this file with the new names.

Once the application pods start you should have access to the application by pointing your browser to [http://localhost](http://localhost). 
If you made it this far, you can now provision multi-cloud infrastructure using Crossplane Compositions. Check the [AWS Crossplane Compositions Tutorial](aws/) which was contributed by [@asarenkansah](https://github.com/asarenkansah). By separating the application infrastructure provision from the application code you not only enable cross-cloud provider portability but also enable teams to connect the application's services with infrastructure that can be managed by the platform team.


## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```shell
kind delete clusters dev
```


## Next Steps

If you have access to a Cloud Provider such as Google Cloud Platform, Microsoft Azure, or Amazon AWS, I strongly recommend you check the **Crossplane Providers** for these platforms. Installing these providers and provisioning Cloud Resources, instead of using the Crossplane Helm Provider will give you real-life experience on how these tools work. 

As mentioned in Chapter 5, how would you deal with services that need infrastructure components that are not offered as managed services? In the case of Google Cloud Platform, they don't offer a Managed Kafka Service that you can provision. Would you install Kafka using Helm Charts or VMs or would you switch Kafka for a managed service such as Google PubSub? Would you maintain two versions of the same service? 


## Sum up and Contribute

In this tutorial, we have managed to separate the provisioning for the application infrastructure from the application deployment. This enables different teams to request resources on-demand (using Crossplane compositions) and application services that can evolve independently. 

Using Helm Chart dependencies for development purposes and quickly getting a fully functional instance of the application up and running is great. For more sensitive environments, you might want to follow an approach like the one shown here, where you have multiple ways to connect your application with the components required by each service. 

Do you want to improve this tutorial? Create an issue, drop me a message on [Twitter](https://twitter.com/salaboy), or send a Pull Request.
