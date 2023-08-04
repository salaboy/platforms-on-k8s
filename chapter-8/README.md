# Chapter 8 :: Enabling Teams to Experiment

In these tutorials you will install Knative Serving and Argo Rollouts on a Kubernetes cluster to implement Canary Releases, A/B testing and Blue/Green Deployments. The release strategies discussed here aim to enable teams to have more control when releasing new versions of their services. By applying different techniques when releasing software, teams can experiment and test their new versions in a controlled setup, without pushing all the live traffic to a new version at a single point in time. 

- [Release Strategies using Knative Serving](knative/README.md)
- [Release Strategies using Argo Rollouts](argo-rollouts/README.md)


## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```
kind delete clusters dev
```

## Next Steps



## Sum up and Contribute