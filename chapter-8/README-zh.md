# 蓝绿和金丝雀发布

---
_🌍 Available in_: [English](README.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---
本节教程我们将会在 Kubernetes 集群上安装 Knative Serving 和 Argo Rollouts，以实现 Canary 发布、A/B 测试和蓝/绿部署。这些发布策略让团队在发布服务的新版本时拥有更多控制能力。使用不同的技术进行发布，团队可以在一个受控环境中实验和测试他们的新版本，无需在单一时间点将所有实时流量推送到新版本。

- [使用 Knative Serving 实现发布策略](knative/README-zh.md)
- [使用 Argo Rollouts 实现发布策略](argo-rollouts/README-zh.md)

## 清理

可以用如下命令删除 Kind 集群：

```shell
kind delete clusters dev
```

## 下一步

- 如果你对构建 Function As a Service 平台有着浓厚兴趣，可以深入了解一下 [Knative Functions](https://knative.dev/docs/functions/) 项目，这个项目的初衷就是构建能够简化 Function 开发工作的工具。

- 尝试 Argo Rollout 之后，下一步就是创建一个端到端的例子，展示从 Argo CD 到 Argo Rollout 的过程。要完成这一工作，需要创建一个包含 Rollout 定义的仓库。[Argo 项目 FAQ](https://argo-rollouts.readthedocs.io/en/latest/FAQ/) 中包含如何集成这两种工具的进一步说明。

- 可以使用 `AnalysisTemplates` 以及 `AnalysisRuns` 试验更加复杂的场景，这两个能力让开发者在发布新版本时候更有信心。

- 两个项目都能在 [Istio](https://istio.io/) 这样的服务网格上运行。可以更进一步了解一下 Istio 的完整能力。

## 总结和贡献

要改进这些教程，欢迎在 [Twitter](https://twitter.com/salaboy) 上联系我或者提交 PR。
