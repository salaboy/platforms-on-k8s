# Conference Application Helm Chart :: K8s

This chart install the basic services required to run the Conference Application. 
This chart also installs the application require infrastructure (PostgreSQL, Redis and Kafka).

This chart doesn't depend on anything to be installed in the target Kubernetes Cluster

## Installation 


```
helm add 
helm install conference-app platforms-on-k8s/conference-app-k8s
```

Chart parameters: 

- install.infrastructure: true 
- install.ingress: true

## Interacting with the application


