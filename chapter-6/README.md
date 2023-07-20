# Chpater 6 :: Let's build a Platform on top of Kubernetes

On this step-by-step tutorial we will create the APIs of our platform by reusing the power of the Kubernetes APIs. The first use case where the platform can assist the development teams is by creating new development environments providing a self-service approach. 

To build this example we will be using Crossplane and `vcluster`, two Open Source projects hosted in the Cloud-Native Computing Foundation. 

## Installation 

To install Crossplane you need to have a Kubernetes Cluster, you can create one using KinD as we did for you [Chapter 2](../chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind). 

Then you can install Crossplane and the Crossplane Helm Provider in your cluster as we did in [Chapter 5]()

We will be using vcluster in this tutorial, but there is no need to install anything in our cluster for vcluster to work. We just need the `vcluster` CLI to connect to our `vcluster`s you can install it by following the instructions on the official site: [https://www.vcluster.com/docs/getting-started/setup](https://www.vcluster.com/docs/getting-started/setup)


## Defining our Environment API

An environment represent a Kubernetes cluster where the Conference Application will be installed for development purposes. The idea is to provide teams with self-service environments for them to do their work. 

For this tutorial we will define a Environment API and a Crossplane Composition that uses the Helm Provider to create a new instance of `vcluster`. 

Check the Crossplane Composite Resource Definition (XRD) for our [Environments here](resources/env-resource-definition.yaml) and the Crossplane [Composition here](resources/composition-devenv.yaml)

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
    infraInstall: true
    
```

Once sent to the cluster, the Crossplane Composition will kick in and create a new `vcluster` with an instance of the Conference Application inside. 

```
kubectl apply -f team-a-dev-env.yaml
```
Then we can connect to our environment by running: 

```
vcluster connect team-a-dev-env --server https://localhost:8443 -- zsh
```
You can exit the vcluster context by typing `exit` in the terminal.



## Next Steps

## Sum up and Contribute

Do you want to improve this tutorial? Create an issue, drop me a message on [Twitter](https://twitter.com/salaboy) or send a Pull Request.