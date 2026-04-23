# DevOps Book Club - Week 4

<!-- TOC -->
* [DevOps Book Club - Week 4](#devops-book-club---week-4)
  * [Environment pipelines](#environment-pipelines)
  * [Environment pipelines in action](#environment-pipelines-in-action)
    * [Install Argo CD](#install-argo-cd)
    * [Patch argocd-repo-server to use the CMP plugin](#patch-argocd-repo-server-to-use-the-cmp-plugin)
    * [Create an application](#create-an-application)
    * [Dealing with changes the GitOps way](#dealing-with-changes-the-gitops-way)
  * [Service + environment pipelines](#service--environment-pipelines)
  * [Linking back to platform engineering](#linking-back-to-platform-engineering)
<!-- TOC -->

## Environment pipelines

## Environment pipelines in action

* Deployment to Kubernetes via environment pipelines are standardised via GitOps tools, such as Argo CD and Flux CD.
* We assume that the service artefacts (e.g., container images, Helm charts) are already built and published by service pipelines.
* ArgoCD is both multi-cluster and multi-cloud.
* Nomenclature:

  | Book        | ArgoCD      |
  |-------------|-------------|
  | environment | Application |

### Install Argo CD

```shell
kubectl create namespace argocd
kubectl apply -n argocd --server-side --force-conflicts -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

```shell
kubectl config set-context --current --namespace=argocd
```

```shell
brew install argocd
```

```shell
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'
```

```shell
kubectl get svc argocd-server -n argocd -o=jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

```shell
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

```shell
open http://localhost:8080
```

```shell
argocd admin initial-password -n argocd | head -1 | tee >(pbcopy)
```

```shell
argocd login localhost:8080 --username admin --password $(pbpaste) --insecure
```

This will set the `admin` password to the more memorable `my-password`:
```shell
argocd account update-password \
    --current-password $(pbpaste) \
    --new-password "my-password"
```

### Patch argocd-repo-server to use the CMP plugin

:warning: **Note!**  
In order to override the docker repository for the staging environment,
we will use a [config management plugin](https://argo-cd.readthedocs.io/en/stable/operator-manual/config-management-plugins/).  
This is described in the [CMP postrender README](./argo-cd/cmp-postrender/README.md).  
The custom plugin supports values files overrides via the `Values File` parameter, which we will use in the next step when creating an application.

### Create an application

Follow the [click-ops steps](./README.md#setting-up-our-application-for-the-staging-environment).

(optional)
```shell
kubectl create ns staging
```
### Dealing with changes the GitOps way

* Updating the output of a service pipeline should be as simple as commiting an update to a chart version in the environments repo.
* To roll back a version, it's preferred to revert the commit that updated the chart version, rather than adding a new commit with the previous version. This way, any configuration changes are also rolled back.
* Changes to application parameters and configuration should be treated as versioned features or bug fixes.
* Storing versioned snapshots of service and config versions in environment repos enables the use of preview environments, which can be used for testing and validation before promoting changes to main branch.
* :information_source: Even though ArgoCD uses Helm charts as input, it doesn't use Helm Kubernetes resources. The rendered yaml manifest is applied directly, e.g. via `kubectl`.

## Service + environment pipelines

## Linking back to platform engineering
