# Service Pipelines

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)
> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

These short tutorials cover both Tekton and Dagger for Service Pipelines. With Tekton we extend the Kubernetes APIs to define our Pipelines and Tasks. With Dagger we programatically define Pipelines that can be executed remotely in Kubernentes or locally in our development laptops. Finally, a pointer to a set of GitHub Actions is provided to be able to compare between these different approaches.

- [Tekton for Service Pipelines Tutorial](tekton/README.md)
- [Dagger for Service Pipelines Tutorial](dagger/README.md)
- [GitHub Actions](github-actions/README.md)


## Clean up

If you want to get rid of the KinD Cluster created for these tutorials, you can run:

```shell
kind delete clusters dev
```

## Next Steps

I strongly recommend following the tutorials listed for Tekton and Dagger in your local environments. If you have a development background you can extend the Dagger pipelines with your custom steps. 

If you are not a Go developer, would you dare to build a pipeline for your tech stack using Tekton and Dagger? 

If you have a container registry account such as a Docker Hub account you can try configuring the credentails for the pipelines to push the container images to the registry. You can then use the Helm Chart `values.yaml` to consume the images from your account instaed of the official ones hosted in `docker.io/salaboy`.

Finally, you can fork this repository `salaboy/platforms-on-k8s` in your own GitHub user (organization) to experiment with running the GitHub actions located in this directory [../../.github/workflows/](../../.github/workflows/).


## Sum up and Contribute

In this tutorials we experimented with two completely different approaches for Service Pipelines. We started with [Tekton](https://tekton.dev) a very non-opinionated Pipeline Engine that was designed to be Kubernetes-native, leveraging the declarative power of the Kubernetes APIs and Kubernetes reconcilation loops. Then we tried [Dagger](https://dagger.io) a pipeline engine designed to orchestrate containers and that can be configured using your favourite tech stack using their SDKs. 

One thing is for sure, no matter which Pipeline Engine your organization has chosen development teams will greatly benefit from being able to consume Service Pipelines without the need of defining every single step, credentials and configuration details. If the platform team can create shared steps/tasks or even default pipelines for your services, developers can focus on writing application code. 

Do you want to improve this tutorial? Create an issue, drop me a message on [Twitter](https://twitter.com/salaboy) or send a Pull Request.
