# Chpater 6 :: Let's build a Platform on top of Kubernetes

On this step-by-step tutorial we will create the APIs of our platform by reusing the power of the Kubernetes APIs. The first use case where the platform can assist the development teams is by creating new development environments providing a self-service approach. 

To build this example we will be using Crossplane and `vcluster`, two Open Source projects hosted in the Cloud-Native Computing Foundation. 

## Installation 

To install Crossplane you need to have a Kubernetes Cluster, you can create one using KinD as we did for you [Chapter 2](../chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind). 

Then you can install Crossplane and the Crossplane Helm Provider in your cluster as we did in [Chapter 5](https://github.com/salaboy/platforms-on-k8s/tree/main/chapter-5#installing-crossplane)

We will be using [`vcluster`](https://www.vcluster.com/) in this tutorial, but there is no need to install anything in our cluster for vcluster to work. We just need the `vcluster` CLI to connect to our `vcluster`s you can install it by following the instructions on the official site: [https://www.vcluster.com/docs/getting-started/setup](https://www.vcluster.com/docs/getting-started/setup)


## Defining our Environment API

An environment represent a Kubernetes cluster where the Conference Application will be installed for development purposes. The idea is to provide teams with self-service environments for them to do their work. 

For this tutorial we will define a Environment API and a Crossplane Composition that uses the Helm Provider to create a new instance of `vcluster`. 

Check the Crossplane Composite Resource Definition (XRD) for our [Environments here](resources/env-resource-definition.yaml) and the Crossplane [Composition here](resources/composition-devenv.yaml). This resource configures the provisioning of a new `vcluster` using the Crossplane Helm Provider, [check this configuration here](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-6/resources/composition-devenv.yaml#L24). When a new `vcluster` is created then the composition install our Conference Application into it, once again using the Crossplane Helm Provider, but this time configured [pointing to the just created `vcluster` APIs](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-6/resources/composition-devenv.yaml#L87), you can [check this here](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-6/resources/composition-devenv.yaml#L117).

Let's install both the XRD and the Composition by running: 
```
kubectl apply -f resources/
```

You should see: 

```
composition.apiextensions.crossplane.io/dev.env.salaboy.com created
compositeresourcedefinition.apiextensions.crossplane.io/environments.salaboy.com created
```

With the Environment resource and the Crossplane Composition using `vcluster` our teams can now request their Environments on demand. 


## Requesting a new Environment

To request a new Environment teams can create new environment resources like this one: 

```
apiVersion: salaboy.com/v1alpha1
kind: Environment
metadata:
  name: team-a-dev-env
spec:
  compositionSelector:
    matchLabels:
      type: development
  parameters: 
    installInfra: true
    
```

Once sent to the cluster, the Crossplane Composition will kick in and create a new `vcluster` with an instance of the Conference Application inside. 

```
kubectl apply -f team-a-dev-env.yaml
```
You should see: 

```
environment.salaboy.com/team-a-dev-env created
```

You can always check the state of your Environments by running: 

```
kubect get env
NAME             CONNECT-TO             TYPE          INFRA   DEBUG   SYNCED   READY   CONNECTION-SECRET   AGE
team-a-dev-env   team-a-dev-env-jp7j4   development   true    true    True     False   team-a-dev-env      1s

```

You can check that Crossplane is creating and managing resources related to the composition by running: 

```
kubectl get managed
NAME                            CHART            VERSION          SYNCED   READY   STATE      REVISION   DESCRIPTION        AGE
team-a-dev-env-jp7j4-8lbtj      conference-app   v1.0.0           True     True    deployed   1          Install complete   57s
team-a-dev-env-jp7j4-vcluster   vcluster         0.15.0-alpha.0   True     True    deployed   1          Install complete   57s
```

These managed resources are no other than Helm Releases being created:

```
kubectl get releases
NAME                            CHART            VERSION          SYNCED   READY   STATE      REVISION   DESCRIPTION        AGE
team-a-dev-env-jp7j4-8lbtj      conference-app   v1.0.0           True     True    deployed   1          Install complete   45s
team-a-dev-env-jp7j4-vcluster   vcluster         0.15.0-alpha.0   True     True    deployed   1          Install complete   45s
```


Then we can connect to the provisioned environment by running (use the CONNECT-TO column for the vcluster name): 
```
vcluster connect team-a-dev-env-jp7j4 --server https://localhost:8443 -- zsh
```

Once you are connected to the `vcluster` you are in a different Kubernetes Cluster, so if you list all the available namespaces you should see: 

```
kubectl get ns
NAME              STATUS   AGE
default           Active   64s
kube-system       Active   64s
kube-public       Active   64s
kube-node-lease   Active   64s
```

As you can see, Crossplane is not installed here. But if you list all the pods in this cluster, you should see all the application pods running: 

```
NAME                                                              READY   STATUS    RESTARTS      AGE
conference-app-kafka-0                                            1/1     Running   0             103s
conference-app-postgresql-0                                       1/1     Running   0             103s
conference-app-c4p-service-deployment-57d4ddcd68-45f6h            1/1     Running   2 (99s ago)   104s
conference-app-agenda-service-deployment-9bf7946c9-mmx8h          1/1     Running   2 (98s ago)   104s
conference-app-redis-master-0                                     1/1     Running   0             103s
conference-app-frontend-deployment-c8c64c54d-lntnw                1/1     Running   2 (98s ago)   104s
conference-app-notifications-service-deployment-64ff7bcdf8nbvhl   1/1     Running   3 (80s ago)   104s
```

You can also do port-forwarding to this cluster, to access the application using:
```
kubectl port-forward svc/frontend 8080:80
```
Now your application is available at [http://localhost:8080](http://localhost:8080)


You can exit the vcluster context by typing `exit` in the terminal.


## Simplifying our platform surface

We can go one step further to simplify the interaction with the platform APIs, avoiding teams to connect to the Platform Cluster and remove the need for having access to the Kubernetes APIs. 

In this short section we deploy an Admin User Interface that allow teams to request new environments using a website, or a set of simplified REST APIs. 

Before installing the Admin User Interface you need to make sure that you are not inside a `vcluster` session. (You can exit the vcluster context by typing `exit` in the terminal). Check that you have the `crossplane-system` namespaces in the current cluster were you are connected. 

You can install this Admin User Interface using Helm:

```
helm install admin oci://docker.io/salaboy/conference-admin --version v1.0.0
```



## Next Steps

Can you extend the Admin User Interface to create Databases and Message Brokers like we did in Chapter 5? 



## Sum up and Contribute

Do you want to improve this tutorial? Create an issue, drop me a message on [Twitter](https://twitter.com/salaboy) or send a Pull Request.