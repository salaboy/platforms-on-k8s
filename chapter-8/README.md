# Chapter 8 :: Enabling Teams to Experiment

In these tutorials you will install Knative Serving and Argo Rollouts on a Kubernetes cluster to implement different release strategies. The release strategies discussed here aim to enable teams to have more control when releasing new versions of their services. 

## Installation

You need a Kubernetes Cluster to install [Knative Serving](https://knative.dev) and [Argo Rollouts](). You can create one using Kubernetes KinD, but instead of using the configurations provided in Chapter 2, we will use the following command:

```
cat <<EOF | kind create cluster --name dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 31080 # expose port 31380 of the node to port 8080 on the host, later to be use by kourier or contour ingress
    listenAddress: 127.0.0.1
    hostPort: 8080
EOF
```

Notice that we are using the 8080 port so you need to make sure that this port is available on your local setup. 

When using Knative Serving there is no need to install an Ingress Controller, as Knative Serving requires a more advanced networking stack to enable features such as traffic routing and splitting. 

If you are limited on resources you can create a cluster for the Knative Serving section and then another cluster for Argo Rollouts. 

Once you have a Cluster let's start by installing [Knative Serving](https://knative.dev/docs/install/yaml-install/serving/install-serving-with-yaml/) you can follow the official documentation or copy the installation steps listed here, as the examples had been tested with this version of Knative. 

Install Knative Serving Custom Resource Definitions:

```
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.2/serving-crds.yaml

```

Then install the Knative Serving Control Plane: 
```
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.2/serving-core.yaml

```

Install a compatible Knative Serving Networking Stack, here you can install Istio, a full-blown service mesh or more simple options like Kourier, that we will use here to limit the resource utilization on your cluster: 

```
kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.10.0/kourier.yaml

```

Configure Kourier as the selected networking stack: 

```
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress-class":"kourier.ingress.networking.knative.dev"}}'
```

Finally, configure the DNS resolution for your cluster, so it can expose reachable URLs for our services:

```
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.2/serving-default-domain.yaml

```

**Only for Knative on KinD** 

For Knative Magic DNS to work in KinD you need to patch the following ConfigMap:

```
kubectl patch configmap -n knative-serving config-domain -p "{\"data\": {\"127.0.0.1.sslip.io\": \"\"}}"
```

and if you installed the `kourier` networking layer you need to create an ingress:

```
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: kourier-ingress
  namespace: kourier-system
  labels:
    networking.knative.dev/ingress-provider: kourier
spec:
  type: NodePort
  selector:
    app: 3scale-kourier-gateway
  ports:
    - name: http2
      nodePort: 31080
      port: 8080
      targetPort: 8080
EOF
```

Unfortunately, we will be using port 8080 in our host environment to access the Knative Services, this will require us to append the port ``:8080`` to our service URLs, something that when using real life domains hosted in a cloud or public IP address is not required. 

One more step is needed to run the examples covered in the following section and this is to enable `tag-header-based-routing` and access to the [Kubernetes Downward API](https://knative.dev/docs/serving/configuration/feature-flags/#kubernetes-downward-api) to fetch information about the pods running in the cluster. We can tune the Knative Serving installation by patching the config-features ConfigMap (https://knative.dev/docs/serving/configuration/feature-flags/#feature-and-extension-flags):

```
kubectl patch cm config-features -n knative-serving -p '{"data":{"tag-header-based-routing":"Enabled", "kubernetes.podspec-fieldref": "Enabled"}}'
```

After all these configurations, we should be ready to use Knative Serving on our KinD Cluster. 


## Release Strategies with Knative Serving

Before jumping into implementing different release strategies we need to understand the basics on Knative Serving. For that we will look into defining one Knative Service for the Notification Service. This will give us hands-on experience on using Knative Serving and its capabilities. 

### Knative Services

Knative Serving simplifies and extend the capabilities offered by Kubernetes by using the concept of a `Knative Service``. A `Knative Service`` uses the Knative Serving networking layer to route traffic to our workloads without pushing us to define complex Kubernetes Resources. Because Knative Serving has access to information about how traffic is flowing to our services it can understand the load that our services are experiencing and make use of a purposefully build autoscaler to upscale or downscale our service instances based on demand. This can be really useful for platform teams looking to implement a function-as-a-service model for their workloads, as Knative Serving can downscale to zero services that are not receiving traffic. 

Knative Services also expose a simplified configuration that resembles Containers as a Service models like Google Cloud Run, Azure Container Apps and AWS App Runner, where by defining which container we want to run, the platform will take care of the rest (no complex configurations for networking, routing traffic, etc). 

Because the Notifications Service uses Kafka for emitting events, we need to install Kafka using Helm:

```
helm install kafka bitnami/kafka --version 22.1.5 --set "provisioning.topics[0].name=events-topic" --set "provisioning.topics[0].partitions=1" 
```

Check that Kafka is running before proceeding, as it usually takes a bit of time to fetch the Kafka Container image and start it. 

If you have a Kubernetes Cluster with Knative Serving installed you can apply the following Knative Service resource to run an instance of our Notifications Service:

```
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: notifications-service 
spec:
  template:
    spec:
      containers:
        - image: salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0
          env:
          - name: KAFKA_URL
            value: kafka.default.svc.cluster.local
          ...<MORE ENVIRONMENT VARIABLES HERE>...  
```

You can apply this resource by running: 

```
kubectl apply -f knative/notifications-service.yaml
```

Knative Serving will create an instance of our container and setup all the networking configurations to provide us with an URL to access the service. 

You can list all Knative Services running the following command: 

```
> kubectl get ksvc 
NAME                    URL                                                       LATESTCREATED                 LATESTREADY                   READY   REASON
notifications-service   http://notifications-service.default.127.0.0.1.sslip.io   notifications-service-00001   notifications-service-00001   True    

```

You can now curl the `service/info` URL of the service to make sure that it is working. Notice the extra `:8080` on the service URL:

```
curl http://notifications-service.default.127.0.0.1.sslip.io:8080/service/info
```

You should see the following output: 

```
{"name":"NOTIFICATIONS","version":"1.0.0","source":"https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/notifications-service","podName":"notifications-service-00002-deployment-76ffd7fbcc-ljkdt","podNamespace":"default","podNodeName":"dev-control-plane","podIp":"10.244.0.22","podServiceAccount":"default"}

```

Check that there is a Pod running our container: 

```
> kubectl get pods
NAME                                                      READY   STATUS    RESTARTS   AGE
kafka-0                                                   1/1     Running   0          7m54s
notifications-service-00001-deployment-798f8f79f5-jrbr8   2/2     Running   0          4s
```

Notice that there are two containers running, one is our Notification Service and the other is the Knative Serving `queue` proxy which is used to buffer incoming requests and get traffic metrics. 

After 90 seconds (by default), if you are not sending requests to the Notification Service, the service instance will be automatically downscaled. Whenever a new incoming request arrives, Knative Serving will automatically upscale the service and buffer the request until the instance is ready. 


To recap, we get two things out of the box with Knative Serving: 
- A simplified way to run our workloads without creating multiple Kubernetes Resources. This approach resembles a Container-as-a-Service offering, that as a Platform team you might want to offer to your teams.
- Dynamic autoscaling using the Knative Autoscaler, which can be used to downscaled your applications to zero when they are not being used. This resembles a Functions-as-a-Service approach, that as a Platform team you might want to provide your teams.

### Canary Releases

In this section we will look into implementing different release strategies for our Notification Service, the same can be applied for all the other services of the Conference Application.


### A/B Testing

### Blue/Green Deployments


## Release Strategies with Argo Rollouts


## Clean up

If you want to get rid of the KinD Cluster created for this tutorial, you can run:

```
kind delete clusters dev
```

## Next Steps



## Sum up and Contribute