# Chapter 8 :: Enabling Teams to Experiment

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)

> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

In these tutorials you will install Knative Serving and Argo Rollouts on a Kubernetes cluster to implement Canary Releases, A/B testing and Blue/Green Deployments. The release strategies discussed here aim to enable teams to have more control when releasing new versions of their services. By applying different techniques when releasing software, teams can experiment and test their new versions in a controlled setup, without pushing all the live traffic to a new version at a single point in time. 

- [Release Strategies using Knative Serving](knative/README.md)
- [Release Strategies using Argo Rollouts](argo-rollouts/README.md)


## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```shell
kind delete clusters dev
```

## Next Steps

- Check the [Knative Functions](https://knative.dev/docs/functions/) project if you are interested in building a Function-as-a-Service platform, as this initiative is working on tooling to make Function developers life easier.

- After trying out Argo Rollouts, the next step is to create an example end-to-end showing the flow of Argo CD to Argo Rollouts. This requires to create a repository that contains your Rollouts definitions. Check the [FAQ section on the Argo Projects](https://argo-rollouts.readthedocs.io/en/latest/FAQ/) for more details about their integration.

- Experiment with more complex examples using `AnalysisTemplates` and `AnalysisRuns` as this feature helps teams to deploy new versions with more confident. 

- As both projects can work with a Service Mesh like [Istio](https://istio.io/), familiarize yourself with what Istio can do for you. 

## Sum up and Contribute


Do you want to improve this tutorial? Create an issue, drop me a message on [Twitter](https://twitter.com/salaboy) or send a Pull Request.
