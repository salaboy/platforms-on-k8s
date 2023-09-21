# AWS Cloud Provider

在本教程中，我们将使用 Crossplane 在 AWS 中部署 Redis、PostgreSQL 和 Kafka。

## 安装 Crossplane

按照第二章的说明部署一个 Kind 集群。然后使用 Helm 把 [Crossplane](https://crossplane.io) 部署到它自己的命名空间中：

```
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --wait
```

安装 `kubectl crossplane` 插件：

```
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh
sudo mv kubectl-crossplane /usr/local/bin
```

Then install the Crossplane AWS provider: 
```
kubectl crossplane install provider crossplane/provider-aws:v0.21.2
```

几秒钟之后，用下面的命令查询 Provider，会看到 Helm Provider 进入了 `INSTALLED`、`HEALTHY` 状态。

```
> kubectl get providers.pkg.crossplane.io
NAME                             INSTALLED   HEALTHY   PACKAGE                               AGE
crossplane-provider-aws         True        True      crossplane/provider-aws:v0.21.2       49s
```

现在，我们准备安装应用程序运行所需的数据库和消息中间件。

## 使用 Crossplane Compositions 实现基础设施的按需供给

我们需要为我们的 Redis、PostgreSQL 和 Kafka 安装 Crossplane Composition。

```
kubectl apply -f resources/
```

Crossplane Composition 资源（`app-database-redis.yaml`）定义了需要创建哪些云资源以及如何将它们配置在一起。跨平面复合资源定义 (XRD) (`app-database-resource.yml`)定义了一个简化接口，使应用开发团队能够通过创建这种类型的资源快速申请新数据库。

检视 [resources/](resources/) 目录，看看其中的 Compositions 以及 Composite Resource Definition (XRDs)。

创建一个包含 AWS 账户 `aws_access_key_id` 和 `aws_secret_access_key` 的文本文件。

```ini
[default]
aws_access_key_id = 
aws_secret_access_key = 
```

用 AWS 凭据创建一个 Kubernetes Secret：

```
kubectl create secret \
generic aws-secret \
-n crossplane-system \
--from-file=creds=./aws-credentials.txt
```

创建一个 `ProviderConfig` 资源：

```yaml
cat <<EOF | kubectl apply -f -
apiVersion: aws.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: aws-secret
      key: creds
EOF
```

### 部署应用所需的基础设施

下面我们要执行命令来创建应用所需的基础设施：

```
kubectl apply -f my-db-keyvalue.yaml
kubectl apply -f my-db-sql.yaml
kubectl apply -f aws-messagebroker-kafka.yaml
```

## 部署应用程序

现在我们已经有了数据库和消息中间件，接下来要确保我们的应用能够连接到这些实例。第一件要做的事情就是禁用会议应用的 Chart 中定义的依赖关系，这样我们安装应用的时候就不用安装数据库和消息中间件了。简单的把 `install.infrastructure` 设置为 `false` 就可以达成这种目的。

我们要使用 `app-values.yaml` 文件，其中包含了应用连接数据库所需的凭据：

```
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0 -f app-values.yaml
```

确保根据新创建的 AWS 资源的值改写 yaml 文件中注释掉的部分。

`app-values.yaml` 内容如下： 

```yaml
install:
  infrastructure: false
frontend:
  kafka:
    url: #aws-kafka-endpoint
agenda:
  kafka:
    url: #aws-kafka-endpoint
  redis: 
    host: #aws-redis-endpoint
    secretName: #aws-redis-password
c4p: 
  kafka:
    url: #aws-kafka-endpoint
  postgresql:
    host: #aws-psql-endpoint
    secretName: #aws-psql-secret

notifications: 
  kafka:
    url: #aws-kafka-endpoint
```
