# 平台工程的度量

---
_🌍 Available in_: [English](README.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

本文包含两个不同的教程，讲述如何使用 DORA 指标来度量平台的绩效。

- [DORA metrics and CloudEvents](dora-cloudevents/README-zh.md)
- [Keptn Lifecycle Toolkit](keptn/README-zh.md)

## 总结

本章内容展示了两种不同又互补的度量方法，我们可以使用这两种方法来观察和监控我们的应用程序。第一个教程侧重于 CloudEvents 和 CDEvents，展示平台工程团队如何利用不同的事件源来得出 DORA 指标；第二个教程侧重于 Keptn 生命周期工具包，该工具包通过扩展 Kubernetes 调度器，收集有关应用程序的信息，提供了开箱可用的部署频率指标。

平台团队理应考虑采用此类工具，这不仅是为了计算指标，更是为了证明其平台投资的合理性。如果平台工程提高了团队的部署频率，缩短了变更的准备时间，同时缩短了从故障中恢复的时间，就证明了平台工程的正确性。反过来看，如果团队的部署频率降低，变更需要更长的时间才能提交给客户，那么可能应该重新评估。

## 清理

可以用如下命令清理 Kind 创建的集群：

```
kind delete clusters dev
```
