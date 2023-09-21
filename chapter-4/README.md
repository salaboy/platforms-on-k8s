# Environment Pipelines

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md)
> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

In this short tutorial, we will set up our Staging Environment Pipeline using [ArgoCD](https://argo-cd.readthedocs.io/en/stable/). We will configure the environment to contain an instance of the Conference Application.

We will define the configuration of the Staging environment using a Git repository. The [`argo-cd/staging` directory](argo-cd/staging/) contains the definition of a Helm chart that can be synced to multiple Kubernetes Clusters. 

## Prerequisites and installation

- We need a Kubernetes Cluster. We will use Kubernetes [KinD](https://kind.sigs.k8s.io/) in this tutorial
- Install ArgoCD in your cluster, [follow this instructions](https://argo-cd.readthedocs.io/en/stable/getting_started/) and optionally install the `argocd` CLI 
- You can fork/copy [this repository](http://github.com/salaboy/platforms-on-k8s/) as if you want to change the configuration for the application, you will need to have write access to the repository. We will be using the directory `chapter-4/argo-cd/staging/`

[Create a KinD Cluster as we did in Chapter 2](../chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind).


Once you have the cluster up and running with the nginx-ingress controller, let's install Argo CD in the cluster: 

```shell
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

You should see something like this: 

```shell
> kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
namespace/argocd created
customresourcedefinition.apiextensions.k8s.io/applications.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/applicationsets.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/appprojects.argoproj.io created
serviceaccount/argocd-application-controller created
serviceaccount/argocd-applicationset-controller created
serviceaccount/argocd-dex-server created
serviceaccount/argocd-notifications-controller created
serviceaccount/argocd-redis created
serviceaccount/argocd-repo-server created
serviceaccount/argocd-server created
role.rbac.authorization.k8s.io/argocd-application-controller created
role.rbac.authorization.k8s.io/argocd-applicationset-controller created
role.rbac.authorization.k8s.io/argocd-dex-server created
role.rbac.authorization.k8s.io/argocd-notifications-controller created
role.rbac.authorization.k8s.io/argocd-server created
clusterrole.rbac.authorization.k8s.io/argocd-application-controller created
clusterrole.rbac.authorization.k8s.io/argocd-server created
rolebinding.rbac.authorization.k8s.io/argocd-application-controller created
rolebinding.rbac.authorization.k8s.io/argocd-applicationset-controller created
rolebinding.rbac.authorization.k8s.io/argocd-dex-server created
rolebinding.rbac.authorization.k8s.io/argocd-notifications-controller created
rolebinding.rbac.authorization.k8s.io/argocd-redis created
rolebinding.rbac.authorization.k8s.io/argocd-server created
clusterrolebinding.rbac.authorization.k8s.io/argocd-application-controller created
clusterrolebinding.rbac.authorization.k8s.io/argocd-server created
configmap/argocd-cm created
configmap/argocd-cmd-params-cm created
configmap/argocd-gpg-keys-cm created
configmap/argocd-notifications-cm created
configmap/argocd-rbac-cm created
configmap/argocd-ssh-known-hosts-cm created
configmap/argocd-tls-certs-cm created
secret/argocd-notifications-secret created
secret/argocd-secret created
service/argocd-applicationset-controller created
service/argocd-dex-server created
service/argocd-metrics created
service/argocd-notifications-controller-metrics created
service/argocd-redis created
service/argocd-repo-server created
service/argocd-server created
service/argocd-server-metrics created
deployment.apps/argocd-applicationset-controller created
deployment.apps/argocd-dex-server created
deployment.apps/argocd-notifications-controller created
deployment.apps/argocd-redis created
deployment.apps/argocd-repo-server created
deployment.apps/argocd-server created
statefulset.apps/argocd-application-controller created
networkpolicy.networking.k8s.io/argocd-application-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-applicationset-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-dex-server-network-policy created
networkpolicy.networking.k8s.io/argocd-notifications-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-redis-network-policy created
networkpolicy.networking.k8s.io/argocd-repo-server-network-policy created
networkpolicy.networking.k8s.io/argocd-server-network-policy created
```

You can access the ArgoCD User Interface by using `port-forward`, in a **new terminal** run:

```shell
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

**Note**: you need to wait for the ArgoCD pods to be started. The first time you do this, it will take more time because it needs to fetch the container images from the internet.

You can access the user interface by pointing your browser to [http://localhost:8080](http://localhost:8080)

<img src="imgs/argocd-warning.png" width="600">

**Note**: by default the installation works using HTTP and not HTTPS, hence you need to accept the warning (hit the "Advanced Button" on Chrome) and proceed (**Process to localhost unsafe**). 

<img src="imgs/argocd-proceed.png" width="600">

That should take you to the Login Page:

<img src="imgs/argocd-login.png" width="600">

The user is `admin`, and to get the password for the ArgoCD Dashboard by running: 

```shell
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

Once in, you should see the empty home screen: 

<img src="imgs/argocd-dashboard.png" width="600">

Let's now set up our Staging Environment.


# Setting up our application for the Staging Environment

For this tutorial, we will use a single namespace to represent our Staging Environment. With Argo CD there are no limits, and our Staging environment could be a completely different Kubernetes Cluster. 

First, let's create a namespace for our Staging Environment:

```shell
kubectl create ns staging
```

You should see something like this: 

```shell
> kubectl create ns staging
namespace/staging created
```

Note: Alternatively, you can use the "Auto Create Namespace" option in the ArgoCD application creation. 

Once you have Argo CD installed, you can access the user interface to set up the project. 

<img src="imgs/argocd-dashboard.png" width="600">

Hit the **"+ New App"** button and use the following details to configure your project: 

<img src="imgs/argocd-app-creation.png" width="600">

Here are the Create Application inputs that I've used: 
- Application Name: "staging-environment"
- Project: "default"
- Sync Policy: "Automatic"
- Source Repository: [https://github.com/salaboy/platforms-on-k8s](https://github.com/salaboy/platforms-on-k8s) (here you can point to your fork)
- Revision: "HEAD"
- Path: "chapter-4/argo-cd/staging/"
- Cluster: "https://kubernetes.default.svc" 
- Namespace: "staging"

<img src="imgs/argocd-app-creation2.png" width="600">

And left the other values to their default ones, hit **Create** on the top 

Once the App is created, it will automatically synchronize the changes, as we selected the **Automatic** mode.

<img src="imgs/argocd-syncing.png" width="600">

You can expand the app by clicking on it to see the full view of all the resources that are being created: 

<img src="imgs/app-detail.png" width="600">

If you are running in a local environment, you can always access the application using `port-forward`, in a **new terminal** run:

```shell
kubectl port-forward svc/frontend -n staging 8081:80
```

Wait for the applications pod to be up and running and then you can access the application pointing your browser to [http://localhost:8081](http://localhost:8081).

<img src="imgs/app-home.png" width="600">


As usual, you can monitor the status of the pods and services using `kubectl`. To check if the application pods are ready, you can run: 

```shell
kubectl get pods -n staging
```

You should see something like this: 

```shell
> kubectl get pods -n staging
NAME                                                              READY   STATUS    RESTARTS        AGE
stating-environment-agenda-service-deployment-6c9cbb9695-xj99z    1/1     Running   5 (6m ago)      8m4s
stating-environment-c4p-service-deployment-69d485ffd8-q96z4       1/1     Running   5 (5m52s ago)   8m4s
stating-environment-frontend-deployment-cd76bdc8c-58vzr           1/1     Running   5 (6m3s ago)    8m4s
stating-environment-kafka-0                                       1/1     Running   0               8m4s
stating-environment-notifications-service-deployment-5c9b5bzb5p   1/1     Running   5 (6m13s ago)   8m4s
stating-environment-postgresql-0                                  1/1     Running   0               8m4s
stating-environment-redis-master-0                                1/1     Running   0               8m4s
```

**Note**: a few restarts are OK (RESTARTS column), as some services need to wait for the infrastructure (Redis, PostgreSQL, Kafka) to be up before them being healthy.

## Changing the Application's configuration in the Staging Environment

To update the version of configurations of your services, you can update the files located in the [Chart.yaml](argo-cd/staging/Chart.yaml) file or [values.yaml](argo-cd/staging/values.yaml) file located inside the [staging](staging/) directory.

For the sake of this example, you can change the application configuration by updating the ArgoCD application details and its parameters. 

While you will not do this for your applications, here we are simulating a change in the GitHub repository where the staging environment is defined. 

<img src="imgs/argocd-change-parameters.png" width="600">

Go ahead and edit the Application Details / Parameters and select `values-debug-enabled.yaml` for the values file that we want to use for this application. This file sets the debug flag into the frontend service and it simulates us changing the original `values.yaml` file that was used for the first installation.

<img src="imgs/argocd-new-values.png" width="600">

Because we were using port-forwarding, you might need to run this command again: 

```shell
kubectl port-forward svc/frontend -n staging 8081:80
```

This is due, the Frontend Service Pod is going to be replaced by the newly configured version, hence the port-forwarding needs to be restarted to target the new pod. 

Once the Frontend is up and running you should see the Debug tab in the Back Office section:

![](imgs/app-debug.png)

## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```shell
kind delete clusters dev
```

## Next Steps

Argo CD is just one project to implement GitOps, can you replicate this tutorial with Flux CD? Which one do you prefer? Is your organization already using a GitOps tool? What would it take to deploy the Conference Application walking skeleton to a Kubernetes Cluster using that tool? 

Can you create another environment, let's say a `production-environment` and describe the flow that will need to follow a new release of the `notifications-service` from the staging environment to the production environment? Where would you store the production environment configuration? 

## Sum up and Contribute

In this tutorial we created our **Staging Environment** using an Argo CD application. This allowed us to sync the configuration located inside a GitHub repository to our running Kubernetes Cluster in KinD. If you make changes to the content of the GitHub repository and refresh the ArgoCD application, ArgoCD will notice that our Environment is out of sync. If we use an automated sync strategy, ArgoCD will run the sync step automatically for us every time it notices that there have been changes in the configuration. For more information check the [project website](https://argo-cd.readthedocs.io/en/stable/) or [my blog](https://www.salaboy.com). 

Do you want to improve this tutorial? Create an issue, drop me a message on [Twitter](https://twitter.com/salaboy), or send a Pull Request.
