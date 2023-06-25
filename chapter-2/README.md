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


## Installing the Conference Application

From Helm 3.7+, we can use OCI images to publish, download and install Helm Charts. This approach uses Docker Hub as a Helm Chart registry and to install the Conference Application you only need to run the following command:

```
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0
```

You can also run the following command to see the details of the chart: 

```
helm show all oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0
```



## Important !!! READ!!

Because the Conference Application is installing PostgreSQL, Redis and Kafka, if you want to remove and install the application again you need to make sure to delete the associated PersistenceVolumeClaims (PVCs). These PVCs are the volumes uses to store the data from the databases and Kafka, failing to delete these PVCs in between installations will cause the services to use old credentials to connect to the new provisioned databases. 

You can delete all PVCs by listing them with:

```
kubectl get pvc
```

You should see:

```
NAME                                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-conference-kafka-0                Bound    pvc-2c3ccdbe-a3a5-4ef1-a69a-2b1022818278   8Gi        RWO            standard       8m13s
data-conference-postgresql-0           Bound    pvc-efd1a785-e363-462d-8447-3e48c768ae33   8Gi        RWO            standard       8m13s
redis-data-conference-redis-master-0   Bound    pvc-5c2a96b1-b545-426d-b800-b8c71d073ca0   8Gi        RWO            standard       8m13s
```

And then delete with: 
```
kubectl delete pvc  data-conference-kafka-0 data-conference-postgresql-0 redis-data-conference-redis-master-0
```

The name of the PVCs will change based on the Helm Release name that you used when installing the chart.



