# Conference Application Helm Chart :: K8s

This chart install the basic services required to run the Conference Application. 
This chart also installs the application require infrastructure (PostgreSQL, Redis and Kafka).

This chart doesn't depend on anything to be installed in the target Kubernetes Cluster

## Installation 


```shell
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0 -n chapter02
```

Chart parameters: 

- install.infrastructure: true 
- install.ingress: false


## Packaging and distributing the chart

```shell
helm dependency build
helm package .
helm push conference-app-v1.0.0.tgz oci://registry-1.docker.io/salaboy/
```

## Interacting with the application

```shell
kubectl port-forward svc/frontend -n chapter02 8080:80
```

Point your browser to: 

[http://localhost:8080](http://localhost:8080)