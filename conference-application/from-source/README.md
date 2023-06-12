# Building the Conference Application from Source

You can build the application containers using `ko`

For Kubernetes 1.23 you need Knative 1.9 for Kubernetes 1.24 you need Knative 1.10

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

Same with PostgreSQL: 
```
```

and Kafka:

```
```