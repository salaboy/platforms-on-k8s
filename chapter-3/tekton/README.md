# Tekton for Service Pipelines

In this tutorial we will be installing Tekton into a Kubernetes Cluster, installing our Service Pipeline and creating new instances.

The Service Pipelines clones, build and package and publish the service container image.

## Create or connect to a Kubernetes Cluster

If you don't have a Kubernetes Cluster you can create one using Kubernetes KinD following this tutorial. 


## Install Tekton into your Cluster

You can install Tekton into your Kubernetes Cluster by running: 

```
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.49.0/release.yaml
```

## Installing a generic Service Pipeline

```
kubectl apply -f service-pipeline.yaml
```

This pipeline is reusing a couple of Tasks that we haven't defined such as the `git-clone` task and the `ko` task to build container images. For this reason we need to install these Tasks to our cluster before being able to run our pipeline. 

Install the `git-clone` Tekton Task:
```
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/git-clone/0.9/git-clone.yaml
```

Install the `ko` Tekton Task:
```
kubectl apply -f https://api.hub.tekton.dev/v1/resource/tekton/task/ko/0.1/raw

```
## Creating a secret to publish the container images to a registry

```
kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=DOCKER_USERNAME --docker-password=DOCKER_PASSWORD --docker-email DOCKER_EMAIL
```

Then we need to create a service account to point to that secret: 

```
kubectl apply -f service-account.yaml
```

## Let's run our Service Pipeline


Check the [service-pipeline-run.yaml](service-pipeline-run.yaml) which contains the parameters for our service pipeline to run. 
Here you can configure which service you want to build, which version and other parameters that you might want to change.

```
kubectl apply -f service-pipeline-run.yaml
```

You can check the Service Pipeline Run with: 

```
kubectl get pipelinerun
```

And check the logs by looking at the Pods that are being created for each Pipeline Task: 

```
kubectl get pods
```

## 