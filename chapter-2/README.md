# Cloud-Native Application Challenges

In this short tutorial we will be installing the `Conference Application` using Helm into a local KinD Kubernetes Cluster. 

Helm Charts can be published to Helm Chart repositories or also, since Helm 3.7 as OCI containers to container registries. 

## Pre Requisites
- [Install Docker](https://docs.docker.com/get-docker/): **check docker configurations for CPU and RAM allowences**
- [Install KinD](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Install Helm](https://helm.sh/docs/intro/install/)

## Creating a local cluster with Kubernetes KinD

Create a KinD Cluster with 3 worker nodes and 1 Control Plane

```
cat <<EOF | kind create cluster --name dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
- role: worker
EOF

```

### Installing NGINX Ingress Controller

We need NGINGX Ingress Controller to route traffic from our laptop to the services that are running inside the cluster. NGINX Ingress Controller act as a router that is running inside the cluster, but exposed to the outside world. 

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
```

This allows you to route traffic from http://localhost to services running inside the cluster. Notice that for KinD to work in this way, when we created the cluster we provided extra parameters and labels for the control plane node:
```
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true" #This allow the ingress controller to be installed in the control plane node
  extraPortMappings:
  - containerPort: 80 # This allows us to bind port 80 in local host to the ingress controller, so it can route traffic to services running inside the cluster.
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
```

Once we have our cluster and our Ingress Controller installed and configured we can move ahead and install our application.

## Registering new Helm Chart repositories

For this example, I've published the Chart into my own Helm Chart repository that you need to add to your local Helm setup. To do this you need to run: 

```
helm repo add platforms-on-k8s https://salaboy.github.io/helm/
```

Once the repository is added locally, you need to run the following command to fetch the description of all the avialable Helm Charts in that repository: 

```
helm repo update
```

Running `helm search repo conference-app` should now return the chart for the conference application. 

```
helm repo search conference-app
```
You should see something like this: 
```
```

## Installing the Conference Application

To install a new instance of the application in your Kubernetes Cluster you can run: 

```
helm install conference platforms-on-k8s/conference-app
```



