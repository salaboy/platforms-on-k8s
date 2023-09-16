# Tekton in Action

This short tutorial covers how to install Tekton and how to create a very simple Task and Pipeline. 

[Tekton](https://tekton.dev) is a non-opinionated Pipeline Engine built for the Cloud (specifically for Kubernetes). You can build any kind of pipelines that you want as the engine doesn't impose any restrictions on the kind of Tasks that it can execute. This makes it perfect for building Service Pipelines where you might need to have special requirements that cannot be met by a managed service.  

After running our first Tekton Pipeline this tutorial also includes a links to more complex Service Pipelines used to build the Conference Application Services. 


## Installing Tekton

Follow the next steps in order to install and setup Tekton in your Kubernetes Cluster. If you don't have a Kubernetes Cluster you can create one with [KinD, as we did for Chapter 2](../../chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind) 

1. **Install Tekton Pipelines**

```
  kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.45.0/release.yaml
```

1. **Install Tekton Dashboard (optional)**

```
kubectl apply -f https://github.com/tektoncd/dashboard/releases/download/v0.33.0/release.yaml
```
You can access the dashboard by port-forwarding using `kubectl`:

```
kubectl port-forward svc/tekton-dashboard  -n tekton-pipelines 9097:9097
```

![Tekton Dashboard](imgs/tekton-dashboard.png)

Then you can access pointing your browser to [http://localhost:9097](http://localhost:9097)


1. **Install Tekton CLI (optional)**:

You can also install [Tekton `tkn` CLI tool](https://github.com/tektoncd/cli). 
If you are in Mac OSX you can run: 

```
brew install tektoncd-cli
```


## Getting started with Tekton Tasks

This section aims to get you started creating Tasks and a Simple Pipeline, so you can then look into the Service Pipelines used to build the artifacts for the Conference Application. 

With Tekton we can define what our tasks do by create Tekton Task definitions. The following is the most simple example of a task: 

```
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: hello-world-task
spec:
  params:
  - name: name
    type: string
    description: who do you want to welcome?
    default: tekton user
  steps:
    - name: echo
      image: ubuntu
      command:
        - echo
      args:
        - "Hello World: $(params.name)" 
```

This Tekton `Task` uses the `ubuntu` image and the `echo` command located inside that image. This `Task` also accept a parameter called `name` that will be used to print a message. Let's apply this `Task` definition to our cluster by running: 

```
kubectl apply -f hello-world/hello-world-task.yaml
```

When we apply this resource to Kubernertes we are not executing the task, we are only making the Task definition available for other to use. This task can now be referenced in multiple pipelines or being executed independently by different users. 

You can now list the available tasks in the cluster by running: 
```
> kubectl get tasks
NAME               AGE
hello-world-task   88s
```

Now let's run our Task. We do this by creating a `TaskRun` resource, which represent a single run for our task. Notice that this concrete run will have a fixed resource name (`hello-world-task-run-1`) and a concrete value for the task parameter called `name`. 

```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  name: hello-world-task-run-1
spec:
  params: 
  - name: name
    value: "Building Platforms on top of Kubernetes reader!"
  taskRef:
    name: hello-world-task
```

Let's apply this `TaskRun` resource to our cluster to create our first Task Run (execution):
```
kubectl apply -f hello-world/task-run.yaml
taskrun.tekton.dev/hello-world-task-run-1 created
```

As soon as the `TaskRun` is created, the Tekton Pipeline Engine is in charge of scheduling the tasks and create the Kubernetes Pod needed to execute it. If you list the pods in the default namespace you should see something like this: 

```
kubectl get pods
NAME                         READY   STATUS     RESTARTS   AGE
hello-world-task-run-1-pod   0/1     Init:0/1   0          2s
```

You can also list `TaskRun`s to check for status: 

```
kubectl get taskrun
NAME                     SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-task-run-1   True        Succeeded   66s         7s
```

Finally, because we were executing a single task you can see the logs of the TaskRun execition by tailing the logs of the pod that was created:

```
kubectl logs -f hello-world-task-run-1-pod 
Defaulted container "step-echo" out of: step-echo, prepare (init)
Hello World: Building Platforms on top of Kubernetes reader!
```

Let's now look into how to sequence multiple tasks together using a Tekton Pipeline.

## Getting Started with Tekton Pipelines

Now we can use Pipelines to coordinate multiple tasks like the one that we defined before. We can also reuse Task definitions created by the Tekton community from the [Tekton Hub](https://hub.tekton.dev/).

![Tekton Hub](imgs/tekton-hub.png)


Before creating the Pipeline we will install the `wget` Tekton task from the Tekton Hub by running: 

```
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/wget/0.1/wget.yaml
```

You should see: 

```
task.tekton.dev/wget created
```

Now let's use our `Hello World` task and the `wget` task that we just installed together into a simple pipeline. 

We will be creating this simple Pipeline Definition, that fetch a file, read its content and then uses the previously defined `Hello World` Task.

![Hello World Pipeline](imgs/hello-world-pipeline.png)

Let's create the following pipeline definition:

```yaml
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: hello-world-pipeline
  annotations:
    description: |
      Fetch resource from internet, cat content and then say hello
spec:
  results: 
  - name: message
    type: string
    value: $(tasks.cat.results.messageFromFile)
  params:
  - name: url
    description: resource that we want to fetch
    type: string
    default: ""
  workspaces:
  - name: files
  tasks:
  - name: wget
    taskRef:
      name: wget
    params:
    - name: url
      value: "$(params.url)"
    - name: diroptions
      value:
        - "-P"  
    workspaces:
    - name: wget-workspace
      workspace: files
  - name: cat
    runAfter: [wget]
    workspaces:
    - name: wget-workspace
      workspace: files
    taskSpec: 
      workspaces:
      - name: wget-workspace
      results: 
        - name: messageFromFile
          description: the message obtained from the file
      steps:
      - name: cat
        image: bash:latest
        script: |
          #!/usr/bin/env bash
          cat $(workspaces.wget-workspace.path)/welcome.md | tee /tekton/results/messageFromFile
  - name: hello-world
    runAfter: [cat]
    taskRef:
      name: hello-world-task
    params:
      - name: name
        value: "$(tasks.cat.results.messageFromFile)"
```

It end up not being that easy to fetch a file, read its content and then use our previously defined hello-world taks to print the content of the file that we have fetched. 
With pipelines we have the flexiblity to add new tasks if needed to do transformations or further processing of the inputs and outputs of each individual tasks. 

For this example, we are using the `wget` Task that we installed from the Tekton Hub, a task that is defined inline called `cat` that basically fetch the content of the downloaded file and store it into a Tekton Result that can be referenced later into our `hello-world-task`. 

Go ahead and install this pipeline definition by running: 

```
kubectl apply -f hello-world/hello-world-pipeline.yaml
```

Then we can create a new `PipelineRun` everytime that we want to execute this pipeline:

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: hello-world-pipeline-run-1
spec:
  workspaces:
    - name: files
      volumeClaimTemplate: 
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 1M 
  params:
  - name: url
    value: "https://raw.githubusercontent.com/salaboy/salaboy/main/welcome.md"
  pipelineRef:
    name: hello-world-pipeline
  
```

Because our tasks needs to download and store files in the filesystem, we are using Tekton workspaces as abstractions to provide storage for our `PipelineRun`s. As we did before with our `TaskRun` we can also provide parameters for the `PipelineRun` allowing us to parameterize each run to use different configurations, or in this case different files. 

Both with `PipelineRuns` and `TaskRuns` you will need to generate a new resource name for each run. As if you try to reapply the same resource twice, the Kubernetes API server will not allow you to mutate the existing resource with the same name. 

Run this pipeline by running: 

```
kubectl apply -f hello-world/pipeline-run.yaml
```

Check the pods that are created: 

```
> kubectl get pods
NAME                                         READY   STATUS        RESTARTS   AGE
affinity-assistant-ca1de9eb35-0              1/1     Terminating   0          19s
hello-world-pipeline-run-1-cat-pod           0/1     Completed     0          11s
hello-world-pipeline-run-1-hello-world-pod   0/1     Completed     0          5s
hello-world-pipeline-run-1-wget-pod          0/1     Completed     0          19s
```

Notice that there is one Pod per Task and a pod called `affinity-assistant-ca1de9eb35-0` which is making sure that the Pods are created in the correct node (where the volume was bound).

Check the TaskRuns too: 

```
> kubectl get taskrun
NAME                                     SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-pipeline-run-1-cat           True        Succeeded   109s        104s
hello-world-pipeline-run-1-hello-world   True        Succeeded   103s        98s
hello-world-pipeline-run-1-wget          True        Succeeded   117s        109s

```

And of course, if all the tasks are successful, the PipelineRun will be too: 

```
kubectl get pipelinerun
NAME                         SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-pipeline-run-1   True        Succeeded   2m13s       114s
```


Make sure you check the pipeline and task executions in the Tekton Dashboard if you installed it.
![Tekton Dashboard](imgs/tekton-dashboard-hello-world-pipeline.png)


## Tekton for Service Pipelines

Service Pipelines in real life are much more complex that the previous simple examples. This is mostly because the pipeline tasks will need to have special configurations and credentials to access external systems. 

An example Service Pipeline definition can be found in this directory in a file called [service-pipeline.yaml](service-pipeline.yaml). 


![Service Pipeline](imgs/service-pipeline.png)


The example Service Pipelines uses [`ko`] to build and publish the container image for our Service. This pipeline is very specify to our Go Services, as if we were building services using a different programming language we will need to use other tools. The example service pipeline can be parameterized to build different services.

To be able to run this Service Pipeline you need to set up credentials to a Container Registry, this means allowing the pipelines to push containers to a container registry such as Docker Hub. To authenticate with a container registry from a Tekton Task/Pipeline [check the official documentation](https://tekton.dev/docs/how-to-guides/kaniko-build-push/#container-registry-authentication).

For this example we will create a Kubernetes Secret with our Docker Hub credentials:

```
kubectl create secret docker-registry docker-credentials --docker-server=https://index.docker.io/v1/ --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
```


Then we will install the `Git Clone` and the `ko` Tekton Tasks: 

```
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/git-clone/0.9/git-clone.yaml
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/ko/0.1/ko.yaml
```

Let's install our Service Pipeline defintion to our cluster:

```
kubectl apply -f service-pipeline.yaml
```

Now we can create new pipeline instances to build and publish our services container images. The folllowing `PipelineRun` resource configure our Service Pipeline to build the Notifications Service.

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: service-pipeline-run-1
  annotations:
    kubernetes.io/ssh-auth: kubernetes.io/dockerconfigjson
spec:
  params:
  - name: target-registry
    value: docker.io/salaboy
  - name: target-service
    value: notifications-service
  - name: target-version 
    value: 1.0.0-from-pipeline-run
  workspaces:
    - name: sources
      volumeClaimTemplate: 
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 100Mi 
    - name: docker-credentials
      secret:  
        secretName: docker-credentials
  pipelineRef:
    name: service-pipeline
```

Apply this `PipelineRun` definition to the cluster to create a new instance of the Service Pipeline: 

```
kubectl apply -f service-pipeline-run.yaml
```

Notice the `spec.params` section, that you will need to modify so the pipeline pushes the resulting container image to your own registry. In other words, replace `docker.io/salaboy` with your registry + username. The `target-service` parameter allows you to choose from which service from the conference application you want to build (from the available services: `notifications-service`, `agenda-service`, `c4p-service`, `frontend`).


There is a separate pipeline ([app-helm-chart-pipeline.yaml](app-helm-chart-pipeline.yaml)) that package and publish the Helm Chart which includes all the application services. 
When the team decides the combination of services and version that they want to bundle inside a helm chart they can run another pipeline to package and publish the chart to the same container registry where the services container images are published. 

![Helm Chart Application Pipeline](imgs/app-helm-pipeline.png)

You can install the Application Helm Chart Pipeline by running: 

```
kubectl apply -f app-helm-chart-pipeline.yaml
```

Then you can create new instances by creating new `PipelineRun` resources:

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: app-helm-chart-pipeline-run-1
  annotations:
    kubernetes.io/ssh-auth: kubernetes.io/dockerconfigjson
spec:
  params:
  - name: target-registry
    value: docker.io/salaboy
  - name: target-version
    value: v0.9.9
  workspaces:
    - name: sources
      volumeClaimTemplate: 
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 100Mi 
    - name: dockerconfig
      secret:  
        secretName: docker-credentials
  pipelineRef:
    name: app-helm-chart-pipeline
```

Apply this `PipelineRun` definition to the cluster to create a new instance of the Application Helm Chart Pipeline: 

```
kubectl apply -f app-helm-chart-pipeline-run.yaml
```

Notice that the Application Helm Chart pipeline also uses the same `docker-credentials` to push the Helm Chart as an OCI container image. The pipeline accepts the `target-version` parameter, which is used to patch the `Chart.yaml` file before packaging and pushing the helm chart to the OCI container registry. Notice that this pipeline doesn't patch the versions of the containers referenced by the chart, which means that is it up to the reader to adapt the pipeline to accept as parameters the versions of each service and validate that the referenced container images exists in the referenced container registry.


**Note**: These pipelines are just examples to illustrate the work required to configure Tekton to build containers and charts. For example, the Application Helm Chart Pipeline doesn't change the version of the chart or the version of the container images referenced inside the chart. If we really want to automate all the process we can get the image versions and the chart version from a Git tag that represent the version that we want to release.


## Clean up

If you want to get rid of the KinD Cluster created for these tutorials, you can run:

```
kind delete clusters dev
```

