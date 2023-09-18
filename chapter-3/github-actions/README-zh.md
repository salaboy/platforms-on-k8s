# GitHub Actions in Action

Github Action 让我们能够无需任何基础设施，自动完成我们的流水线。所有的流水线/工作流都运行在 Github.com 的基础设施里，但是一旦规模达到了收费门槛，费用可能就非常昂贵了。

可以在下面的链接中找到会议应用所需的流水线：

- [Agenda Service GH Service Pipeline](../../.github/workflows/agenda-service-service-pipeline.yaml)
- [C4P Service GH Service Pipeline](../../.github/workflows/c4p-service-service-pipeline.yaml)
- [Notifications Service GH Service Pipeline](../../.github/workflows/notifications-service-service-pipeline.yaml)

这些流水线是靠事件触发的，可以使用过滤器来确保只有特定[源码变更事件](https://github.com/dorny/paths-filter)才会触发流水线。这里使用 [ko-build](https://github.com/ko-build/setup-ko) 构建和发布容器镜像，他能为我们的 Golang 应用生成多平台镜像。`docker/login-action@v2` 需要两个环境 Secret（`secrets.DOCKERHUB_USERNAME`、`secrets.DOCKERHUB_TOKEN`）进行配置，这两个 Secret 让流水线能够把容器镜像推送到 Docker Hub。
