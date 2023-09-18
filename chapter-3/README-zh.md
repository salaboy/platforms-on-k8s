# 使用流水线构建云原生应用

这里的教程包含了 Tekton 和 Dagger 的流水线。Tekton 扩展了 Kubernetes API，提供了用于定义流水线和任务的能力；而 Dagger 让我们可以用可编程的方式定义流水线，这种流水线可以在 Kubernetes 集群上远程执行，也能在本地的开发环境中执行。最后还使用 Github Action 来对比不同的流水线方法。

These short tutorials cover both Tekton and Dagger for Service Pipelines. With Tekton we extend the Kubernetes APIs to define our Pipelines and Tasks. With Dagger we programatically define Pipelines that can be executed remotely in Kubernentes or locally in our development laptops. Finally, a pointer to a set of GitHub Actions is provided to be able to compare between these different approaches.

- [Tekton 流水线教程](tekton/README-zh.md)
- [Dagger 流水线教程](dagger/README-zh.md)
- [GitHub Actions](github-actions/README-zh.md)

## 清理

要删掉教程中使用的 Kind 集群，请执行如下命令：

```
kind delete clusters dev
```

## 下一步

强烈建议按照教程在本地运行一下 Tekton 和 Dagger。有开发背景的读者还可以用自定义步骤扩展 Dagger 流水线。

如果你不是 Golang 开发者，你会使用 Tekton 和 Dagger 为自己的技术栈构建流水线么？

如果你有一个容器镜像仓（例如 Docker Hub）的账号，那么可以试着在流水线中配置凭据，把镜像推送到镜像仓中。接下来在 Helm Chart 的 `values.yaml` 修改，使用自己的镜像仓地址替换原有的 `docker.io/salaboy`。

最后，你可以把镜像仓 `salaboy/platforms-on-k8s` Fork 到你自己的 Github 账号下，试试位于 [../../.github/workflows/](../../.github/workflows/) 目录的 Github Action。

## 总结和贡献

本教程中，我们试用了两种截然不同的流水线。起初用的是 [Tekton](https://tekton.dev) ，这是一个中立的面向 Kubernetes 的产品，其中使用了 Kubernetes 风格的声明式 API 以及 Kubernetes 的调谐机制。接下来我们又尝试了 [Dagger](https://dagger.io)，它是一个可以使用用户自有技术栈进行扩展的面向容器编排的流水线引擎。

可以肯定的是，无论选择了哪种管道引擎，因为无需定义每个步骤、凭据和配置细节，所以开发人员能从流水线中获益。如果平台团队能为您的服务创建共享步骤/任务甚至默认管道，开发人员就可以专注于编写应用代码。

要改进这些教程，欢迎在 [Twitter](https://twitter.com/salaboy) 上联系我或者提交 PR。
