# Dagger in Action

本教程中我们会使用 Dagger 流水线进行构建、测试、打包和发布服务。这些用 Golang 借助 Dagger SDK 编写的流水线，完成了服务的构建和容器镜像封装的任务。还有一条独立的流水线构建并发布了 Helm Chart。

## 前置条件

要在本地运行这些流水线，你需要：

- [安装 Golang](https://go.dev/doc/install)
- [容器运行时（例如 Docker）](https://docs.docker.com/get-docker/)

要想在 Kubernetes 集群里运行流水线，则需要 [Kind](https://kind.sigs.k8s.io/) 或者你找得到的其他 Kubernetes 集群。

## 运行流水线

这些服务都很相似，所以可以使用同样的流水线，用参数化的方式进行构建。你可以克隆代码仓库，在[会议应用代码目录](../../conference-application/)中运行流水线。可以运行任意 `service-pipeline.go` 文件中的任务：

```
go mod tidy
go run service-pipeline.go build <SERVICE DIRECTORY>
```

下面几个任务适用于所有服务：

- `build`：构建服务源码并创建爱你容器镜像。这一任务需要的参数是我们要构建的服务所在的目录。
- `test`：测试服务，但是在启动测试之前，它会启动所有的服务依赖（此处的依赖指的是运行测试所需的容器）。这一任务需要的参数是我们要测试的服务所在的目录。
- `publish`：把构建好的容器镜像发布到镜像仓库。这一步骤需要你登录到镜像仓库，并给容器镜像指定 Tag。该任务的参数包含我们要发布的服务所在的目录，以及容器镜像所需的 Tag。

运行 `go run service-pipeline.go all notifications-service v1.0.0-dagger` 的话，会执行所有任务。在这之前要确保满足所有先决条件，尤其是要登录之后才能将镜像推送到镜像仓。

运行 `go run service-pipeline.go build notifications-service` 则不需要你设置任何凭据。

可以使用环境变量设置容器镜像仓和用户名：

```
CONTAINER_REGISTRY=<YOUR_REGISTRY> CONTAINER_REGISTRY_USER=<YOUR_USER> go run service-pipeline.go publish notifications-service v1.0.0-dagger
```

要完成这一操作，需要你登录到用于发布的镜像仓。

能够用同样的语言来开发应用和流水线，对开发人员来说是很方便的——但是我想你不想在开发人员的笔记本上构建容器镜像吧？

下一节将会展示在 Kubernetes 集群中，远程运行 Dagger 流水线的方法：

## 在 Kubernetes 集群里远程运行 Dagger 流水线

Dagger 流水线引擎能运行在支持容器的任何环境之中。无需复杂配置，就可以在 Kubernetes 中运行你的流水线。

需要一个 Kubernetes 集群来完成后续步骤，可以[使用 Kind 创建一个集群](../../chapter-2/README.zh-cn.md)。

接下来我们会把原本在本地容器运行时中运行的流水线，搬到 Kubernetes 里面去。这种方式还在试验阶段，但是对于我们的体验过程还是有帮助的。

创建一个 Dagger Pod 来运行流水线：

```
kubectl run dagger --image=registry.dagger.io/engine:v0.3.13 --privileged=true
```

另外还可以使用 `kubectl apply -f chapter-3/dagger/k8s/pod.yaml` 命令来启动任务。

查看一下运行情况：

```
kubectl get pods 
```

你应该会看到类似下面的输出：

```
NAME     READY   STATUS    RESTARTS   AGE
dagger   1/1     Running   0          49s
```

**注意**：因为没有为 Dagger 本身设置任何持久化或复制机制，所以在这种情况下，所有存储和缓存机制都是易失的。有关这方面的更多信息，请查阅官方文档。

要在这个远程服务的运行项目的流水线，需要导出以下环境变量：

```
export _EXPERIMENTAL_DAGGER_RUNNER_HOST=kube-pod://<podname>?context=<context>&namespace=<namespace>&container=<container>
```

因为我们手工设置了 Pod 名称，所以 `<podname>` 应该设置为 `dagger`，`<context>` 则应该设置为 Kubernetes 的集群上下文。如果你运行的是一个 Kind 集群，上下文可能是 `kind-dev`。可以用 `kubectl config current-context` 找到当前的上下文。`<namespace>` 应该设置为运行 Dagger Container 的命名空间。`<container>` 的取值设置为 `dagger`。例如我的运行环境中，设置如下：


```
export _EXPERIMENTAL_DAGGER_RUNNER_HOST="kube-pod://dagger?context=kind-dev&namespace=default&container=dagger"
```

注意我的 Kind 集群中不包含任何 Pipeline 相关的东西。

运行构建流水线：

```
go run service-pipeline.go build notifications-service
```

远程测试服务：

```
go run service-pipeline.go test notifications-service
```

可以通过下面的命令来查看 Dagger 日志：

```
kubectl logs -f dagger
```

如果是在远程 Kubernetes 集群（而不是 KinD）上运行，就不需要本地容器运行时来构建服务及其容器镜像了。
