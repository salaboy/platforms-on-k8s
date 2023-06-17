# Building the Conference Application from Source

You can build the application containers using `ko`

For Kubernetes 1.23 you need Knative 1.8 for Kubernetes 1.24 you need Knative 1.9 or 1.10

If you have [Knative Serving installed](https://knative.dev/docs/install/yaml-install/serving/install-serving-with-yaml/#verify-the-installation) in your cluster you can leverage the `config-knative` directory when running

```
ko apply -f config-knative/
```

You can install Redis using the Bitnami Helm Chart: 

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis --set image.tag=6.2 --set architecture=standalone
```


Before installing PostgreSQL, we need to create a configMap with the tables that we want to create: 
```
kubectl apply -f c4p-service/init-sql-configmap.yaml
```

Same with PostgreSQL: 
```
helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql --set image.debug=true --set global.postgresql.auth.postgresPassword=postgres --set primary.initdb.user=postgres --set primary.initdb.password=postgres --set primary.initdb.scriptsConfigMap=init-sql
```

and Kafka:

```
helm install kafka oci://registry-1.docker.io/bitnamicharts/kafka
```

Optional Ingress Controller: 
```
helm install ingress-controller oci://registry-1.docker.io/bitnamicharts/nginx-ingress-controller
```