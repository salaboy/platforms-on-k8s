# Dagger in Action

In this short tutorial we will be looking at the Dagger Service Pipelines provided to build, test, package and publish each service. 
These pipelines are implemented in Go using the Dagger Go SDK and take care of building each service and creating a container. A separate pipeline is provided to build and publish the application Helm Chart.

## Requirements

For running these pipelines locally you will need: 
- [Go installed](https://go.dev/doc/install)
- [A container runtime (such as Docker running locally)](https://docs.docker.com/get-docker/)

To run the pipelines remotely on a Kubernetes Cluster you can use [KinD](https://kind.sigs.k8s.io/) or any Kubernetes Cluster that you have available. 

## Let's run some pipelines

Because all the services are very similar, the same pipeline definition can be used and parameterized to build each service separately. 

You can clone this repository and from the [Confernece Application directory](../../conference-application/) and run the pipelines locally.  

You can run any defined task inside the `service-pipeline.go` file:

```shell
go mod tidy
go run service-pipeline.go build <SERVICE DIRECTORY>
```

The following tasks are defined for all the services: 
- `build`:  build the service source code and create a container for it. This goal expect as an argument the directory of the service that we want to build.
- `test`: test the service, but first it will start all the service dependencies (containers needed to run the test). This goal expect as an argument the directory of the service that we want to test.
- `publish`: publish the created container image to a container registry. This requires you to be logged in to the container registry and to provide the name of the tag used for the container image. This goal expect as an argument the directory of the service that we want build and publish, as well as the tag used to stamp the container before pushing it.

If you run `go run service-pipeline.go all notifications-service v1.0.0-dagger` all the tasks will be executed. Before being able to run all the tasks you will need to make sure that you have all the pre-requisites set, as for pushing to a Container Registry you will need to provide appropriate credentials. 

You can safely run `go run service-pipeline.go build notifications-service` which doesn't require you to set any credentials. You can set up your container registry and username using envorionment variables, for example: 

```shell
CONTAINER_REGISTRY=<YOUR_REGISTRY> CONTAINER_REGISTRY_USER=<YOUR_USER> go run service-pipeline.go publish notifications-service v1.0.0-dagger
```
This require you to be logged in to the registry where you want to publish your container images to.

Now, for development purposes, this is quite convinient, because you can now build your service code in the same way that your CI (Continuous Integration) system will do. But you don't want to run in production container images that were created in your developer's laptop right? 
The next section shows a simple setup of running Dagger pipelines remotely inside a Kubernetes Cluster. 

## Running your pipelines remotely on Kubernetes

The Dagger Pipeline Engine can be run anywhere where you can run containers, that means that it can runs in Kubernetes without the need of complicated setups. 

For this tutorial you need to have a Kubernetes Cluster, you can create one using [KinD as we did for Chapter 2](../../chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind).

In this short tutorial we will run the pipelines that we were running locally with our local container runtime, now remotely against a Dagger Pipeline Engine that runs inside a Kuberneets Pod. This is an experimental feature, and not a recommended way to run Dagger, but it help us to prove the point. 

Let's run the Dagger Pipeline Engine inside Kubernetes by creating a Pod with Dagger: 

```shell
kubectl run dagger --image=registry.dagger.io/engine:v0.3.13 --privileged=true
```

Alternatively, you can apply the `chapter-3/dagger/k8s/pod.yaml` manifest using `kubectl apply -f chapter-3/dagger/k8s/pod.yaml`.

Check that the `dagger`` pod is running: 
```shell
kubectl get pods 
```
You should see something like this: 
```shell
NAME     READY   STATUS    RESTARTS   AGE
dagger   1/1     Running   0          49s
```

**Note**: this is far from ideal because we are not setting any persistence or replication mechanism for Dagger itself, all the caching mechanism are volatile in this case. Check the official documentation for more about this. 

Now to run the projects pipelines against this remote service you only need to export the following environment variable: 
```shell
export _EXPERIMENTAL_DAGGER_RUNNER_HOST=kube-pod://<podname>?context=<context>&namespace=<namespace>&container=<container>
```

Where `<podname>` is `dagger` (because we created the pod manually), `<context>` is your Kubernetes Cluster context, if you are running against a KinD Cluster this might be `kind-dev`. You can find your current context name by running `kubectl config current-context`. Finally `<namespace>` is the namespace where you run the Dagger Container, and `<container>` is once again `dagger`. For my setup against KinD, this would look like this: 

```shell
export _EXPERIMENTAL_DAGGER_RUNNER_HOST="kube-pod://dagger?context=kind-dev&namespace=default&container=dagger"
```

Notice also that my KinD cluster (named `kind-dev`) didn't had anything related to Pipelines. 

Now if you run in any of the projects: 
```shell
go run service-pipeline.go build notifications-service
```
Or to test your service remotly: 

```shell
go run service-pipeline.go test notifications-service
```

In a separate tab you can tail the logs from the Dagger engine by running: 
```shell
kubectl logs -f dagger
```

The build will happen remotely inside the Cluster. If you were running this against a remote Kubernetes Cluster (not KinD), there will not be need for you to have a local Container Runtime to build your services and their containers. 
