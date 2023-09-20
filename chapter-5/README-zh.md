# 多云基础设施

在这个教程里，我们会使用 Crossplan 来为我们的会议应用交付 Redis、PostgreSQL 以及 Kafka 实例。

借用 Crossplane，我们希望统一这些组件的交付方式，对应用团队隐藏组件的地址。

应用团队应该能够用一种和其他 Kubernetes 资源一样的声明式的方法来申请这些资源。这样一来，流水线不仅能够配置应用，也能够配置应用所需的基础设施组件。

## 部署 Crossplane

首先要有个 Kubernetes 集群才能部署 Crossplane，可以按照第二章的说明部署一个 Kind 集群。然后使用 Helm 把 [Crossplane](https://crossplane.io) 部署到它自己的命名空间中：

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

然后安装 Crossplane 的 Helm Provider：

```
kubectl crossplane install provider crossplane/provider-helm:v0.10.0
```

我们需要一个合适的 `ServiceAccount`，用于创建新的 `ClusterRoleBinding`，Helm Provider 就能用这个权限来安装 Chart：

```
SA=$(kubectl -n crossplane-system get sa -o name | grep provider-helm | sed -e 's|serviceaccount\/|crossplane-system:|g')
kubectl create clusterrolebinding provider-helm-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
```

```
kubectl apply -f crossplane/helm-provider-config.yaml
```

几秒钟之后，用下面的命令查询 Provider，会看到 Helm Provider 进入了 `INSTALLED`、`HEALTHY` 状态。

```
> kubectl get providers.pkg.crossplane.io
NAME                             INSTALLED   HEALTHY   PACKAGE                               AGE
crossplane-provider-helm         True        True      crossplane/provider-helm:v0.10.0      49s
```

现在我们准备好了安装我们的数据库和消息中间间，Crossplane 已经有能力供给应用所需的所有组件：

## 使用 Crossplane Compositions 实现基础设施的按需供给

接下来安装 Redis、PostgreSQL 以及 Kafka 的 Composition 对象：

```
kubectl apply -f resources/
```

`app-database-redis.yaml` 中的 Crossplane Composition 资源，其中定义了将要创建的资源及其相关配置。`app-database-resource.yaml` 中的 Crossplane 的资源定义（XRD）中定义了一个简化的界面，让应用开发团队能够方便的通过定义资源的方式来申请新的数据库。

[resources/](resources/) 目录中包含了 Composite Resource 的定义（XRD）以及 Composition 的资源。

### 创建应用所需的基础设施

使用下面的命令就能创建新的 Redis 实例：

```
kubectl apply -f my-db-keyvalue.yaml
```

 `my-db-keyvalue.yaml` 中的资源如下：

```yaml
apiVersion: salaboy.com/v1alpha1
kind: Database
metadata:
  name: my-db-keyvalue
spec:
  compositionSelector:
    matchLabels:
      provider: local
      type: dev
      kind: keyvalue
  parameters: 
    size: small
```

注意，我们用到了标签 `provider: local`、`type: dev` 以及 `kind: keyvalue`，这些标签让 Crossplane 能够找到正确的 Composition。在本例中，使用 Helm Provider 创建了一个本地的 Redis 实例。

用下面的命令查询数据库状态：

```
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   97s
```

新的 Redis 实例位于 `default` 命名空间中。可以用类似的方法创建 PostgreSQL 数据库：

```
kubectl apply -f my-db-sql.yaml
```

查询这两个 `dbs` 实例：

```
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   2m
my-db-sql         small   false      sql        True     False   sql.db.local.salaboy.com        5s
```

再看看运行中的 Pod，每个 Pod 都是一个数据库：

```
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          3m40s
my-db-sql-postgresql-0           1/1     Running   0          104s
```

上面的步骤应该创建了 4 个 Kubernetes Secret（其中两个是 Helm Release 使用的，另外两个是用于连接新建的数据库实例）

```
> kubectl get secret
NAME                                    TYPE                 DATA   AGE
my-db-keyavalue-redis                   Opaque               1      2m32s
my-db-sql-postgresql                    Opaque               1      36s
sh.helm.release.v1.my-db-keyavalue.v1   helm.sh/release.v1   1      2m32s
sh.helm.release.v1.my-db-sql.v1         helm.sh/release.v1   1      36s
```

再用同样的方法创建新的 Kafka 实例：

```
kubectl apply -f my-messagebroker-kafka.yaml
```

查询一下 `mbs` 实例：

```
> kubectl get mbs
NAME          SIZE    KIND    SYNCED   READY   COMPOSITION                  AGE
my-mb-kafka   small   kafka   True     True    kafka.mb.local.salaboy.com   2m51s
```

Kafka 的默认配置中不需要任何 Secret。

现在你能看到三个 Pod 正在运行，分别是 Kafka、PosgreSQL 以及 Redis：

```
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          113s
my-db-sql-postgresql-0           1/1     Running   0          108s
my-mb-kafka-0                    1/1     Running   0          100s
```

**注意**：如果要使用同样的资源名称删除并新建这些数据库或者消息中间件，因为这些实例不会在删除过程中删除 PVC，所以需要手动删除 PVC。

如上所述，你可以在集群资源承受范围之内，创建很多的数据库或者消息中间件了。

## 部署我们的会议应用

现在我们已经有了数据库和消息中间件，接下来要确保我们的应用能够连接到这些实例。第一件要做的事情就是禁用会议应用的 Chart 中定义的依赖关系，这样我们安装应用的时候就不用安装数据库和消息中间件了。简单的把 `install.infrastructure` 设置为 `false` 就可以达成这种目的。

`app-values.yaml` 文件中包含了应用连接数据所需的信息：

```
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0 -f app-values.yaml
```

`app-values.yaml` 的内容大概这样子：

```yaml
install:
  infrastructure: false
frontend:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
agenda:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
  redis: 
    host: my-db-keyavalue-redis-master.default.svc.cluster.local
    secretName: my-db-keyavalue-redis
c4p: 
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
  postgresql:
    host: my-db-sql-postgresql.default.svc.cluster.local
    secretName: my-db-sql-postgresql
notifications: 
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
```

注意 `app-values.yaml` 中需要指定我们的数据库以及消息中间件的名称（`my-db-keyavalue`、`my-db-sql` 和 `my-mb-kafka`）。如果你需要用其它的数据库和消息中间件，就要修改文件的内容。

应用 Pod 启动之后，就可以用浏览器访问 [http://localhost](http://localhost)，来浏览应用页面了。

走到这一步，你就可以通过 Crossplane Composition 来使用多云了。[@asarenkansah](https://github.com/asarenkansah) 贡献的 [AWS Crossplane Composition 教程](aws/) 中很好的展示了这种能力。把应用和基础资源分开之后，不仅脱离了云供应商的锁定，同样也让团队能够连接到受平台团队管理的基础设施。

## 清理

可以用如下命令清理 Kind 创建的集群：

```
kind delete clusters dev
```

## 下一步

如果你能够使用 GCP、Azure 或 AWS 等云提供商，强烈建议你浏览这些云厂商的 Crossplane Provider。使用这些 Provider 代替 Crossplane Helm Provider，会让你更好的理解这些工具的运作机制。

如第 5 章所述，如果服务需要一些基础设施组件，但是所需组件又不是托管服务，那应该怎么办？GCP 没有提供 Kafka 服务，你会使用 Helm Charts 或虚拟机安装 Kafka，还是将 Kafka 换成 Google PubSub 这样的托管服务？你会维护同一服务的两个版本吗？

## 总结和贡献

在本教程中，我们设法将基础设施与应用部署分离开来。这样，不同的团队就可以按需请求资源（使用 Crossplane 组Composite），应用服务也可以独立发展。

Helm Chart 依赖能够快速启动和运行一个功能完备的应用程序实例，对于开发过程是非常友好。但是对于更正式的环境，您可能希望采用这种分离的方法，把应用程序和所需资源连接起来。

要改进这些教程，欢迎在 [Twitter](https://twitter.com/salaboy) 上联系我或者提交 PR。
