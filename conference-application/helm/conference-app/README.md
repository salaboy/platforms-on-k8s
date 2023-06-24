# Conference Application Helm Chart :: K8s

This chart install the basic services required to run the Conference Application. 
This chart also installs the application require infrastructure (PostgreSQL, Redis and Kafka).

This chart doesn't depend on anything to be installed in the target Kubernetes Cluster

## Installation 


```
helm add repo platforms-on-k8s 
helm install conference-app platforms-on-k8s/conference-app-k8s -n chapter02
```

Chart parameters: 

- install.infrastructure: true 
- install.ingress: false

## Interacting with the application

```
kubectl port-forward svc/frontend -n chapter02 8080:80
```

Point your browser to: 

[http://localhost:8080](http://localhost:8080)