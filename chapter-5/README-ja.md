# マルチクラウド（アプリケーション）インフラストラクチャ

---
_🌍 利用可能な言語_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md)

> **注意:** これは素晴らしいクラウドネイティブコミュニティの [🌟 コントリビューター](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) によってもたらされました！

---

このステップバイステップのチュートリアルでは、Crossplaneを使用してアプリケーションサービスのためのRedis、PostgreSQL、Kafkaインスタンスをプロビジョニングします。

CrossplaneとCrossplane Compositionsを使用して、これらのコンポーネントのプロビジョニング方法を統一し、エンドユーザー（アプリケーションチーム）からこれらのコンポーネントの場所を隠蔽することを目指します。

アプリケーションチームは、他のKubernetesリソースと同様に、宣言的なアプローチでこれらのリソースを要求できるようにする必要があります。これにより、チームは環境パイプラインを使用して、アプリケーションサービスとアプリケーションが必要とするインフラストラクチャコンポーネントの両方を設定できます。

## Crossplaneのインストール

Crossplaneをインストールするには、Kubernetesクラスターが必要です。[第2章](../chapter-2/README-ja.md#kubernetes-kindを使用してローカルクラスタを作成する)で行ったように、KinDを使用してクラスターを作成できます。

[Crossplane](https://crossplane.io)を独自の名前空間にHelmを使用してインストールしましょう：

```shell
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --version 1.15.0 --wait 
```

次に、Crossplane Helmプロバイダーと、Helmプロバイダーが代わりにチャートをインストールできるようにするための新しい`ClusterRoleBinding`をインストールします。

```shell
kubectl apply -f crossplane/helm-provider.yaml
```

数秒後、設定されたプロバイダーを確認すると、Helmが`INSTALLED`および`HEALTHY`と表示されるはずです：

```shell
❯ kubectl get providers.pkg.crossplane.io
NAME            INSTALLED   HEALTHY   PACKAGE                                                    AGE
provider-helm   True        True      xpkg.upbound.io/crossplane-contrib/provider-helm:v0.17.0   49s
```

次に、Helmプロバイダーにクラスター内でチャートをインストールするためのin-cluster設定を使用するよう指示する`ProviderConfig`を作成します。

```shell
kubectl apply -f crossplane/helm-provider-config.yaml
```

これで、アプリケーションが必要とするすべてのコンーネントを提供するために、データベースとメッセージブローカーのCrossplane compositionsをインストールする準備が整いました。

## Crossplane Compositionsを使用したオンデマンドのアプリインフラストラクチャ

Key-Valueデータベース（Redis）、SQLデータベース（PostgreSQL）、メッセージブローカー（Kafka）のCrossplane Composite Resource Definitions (XRDs)をインストールする必要があります。

```shell
kubectl apply -f resources/definitions
```

次に、対応するCrossplane CompositionsとInitialization Dataをインストールします：

```shell
kubectl apply -f resources/compositions
kubectl apply -f resources/config
```

Crossplane Compositionリソース（`app-database-redis.yaml`）は、どのクラウドリソースを作成し、どのように一緒に設定する必要があるかを定義します。Crossplane Composite Resource Definition (XRD)（`app-database-resource.yaml`）は、アプリケーション開発チームがこのタイプのリソースを作成することで新しいデータベースをすぐに要求できるようにする簡略化されたインターフェースを定義します。

CompositionsとComposite Resource Definitions (XRDs)については、[resources/](resources/)ディレクトリを確認してください。

### アプリケーションインフラストラクチャをプロビジョニングしましょう

次のコマンドを実行して、チームが使用する新しいKey-Valueデータベースをプロビジョニングできます：

```shell
kubectl apply -f my-db-keyvalue.yaml
```

`my-db-keyvalue.yaml`リソースは次のようになっています：

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

`provider: local`、`type: dev`、`kind: keyvalue`いうラベルを使用していることに注意してください。これにより、Crossplaneはラベルに基づいて適切なcompositionを見つけることができます。この場合、HelmプロバイダーはローカルのRedisインスタンスを作成しました。

次のコマンドを使用してデータベースのステータスを確認できます：

```shell
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   97s
```

`default`名前空間に新しいRedisインスタンスが作成されたことを確認できます。

同じ手順に従って、次のコマンドを実行してPostgreSQLデータベースをプロビジョニングできます：

```shell
kubectl apply -f my-db-sql.yaml
```

これで2つの`dbs`が表示されるはずです：

```shell
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   2m
my-db-sql         small   false      sql        True     False   sql.db.local.salaboy.com        5s
```

各データベースに対して1つずつ、2つのPodが実行されていることを確認できます：

```shell
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          3m40s
my-db-sql-postgresql-0           1/1     Running   0          104s
```

4つのKubernetesシークレット（2つのhelmリリース用と2つの新しく作成されたインスタンスに接続するための認証情報を含む）があるはずです：

```shell
> kubectl get secret
NAME                                    TYPE                 DATA   AGE
my-db-keyavalue-redis                   Opaque               1      2m32s
my-db-sql-postgresql                    Opaque               1      36s
sh.helm.release.v1.my-db-keyavalue.v1   helm.sh/release.v1   1      2m32s
sh.helm.release.v1.my-db-sql.v1         helm.sh/release.v1   1      36s
```

同様に、Kafkaメッセージブローカーの新しいインスタンスをプロビジョニングできます：

```shell
kubectl apply -f my-messagebroker-kafka.yaml
```

そして、次のコマンドで一覧表示できます：

```shell
> kubectl get mbs
NAME          SIZE    KIND    SYNCED   READY   COMPOSITION                  AGE
my-mb-kafka   small   kafka   True     True    kafka.mb.local.salaboy.com   2m51s
```

Kafkaはデフォルト設定を使用する場合、シークレットを作成する必要はありません。

3つの実行中のポッド（Kafka用、Redis用、PostgreSQL用）が表示されるはずです。

```shell
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          113s
my-db-sql-postgresql-0           1/1     Running   0          108s
my-mb-kafka-0                    1/1     Running   0          100s
```

**注意**: 同じリソース名を使用してデータベースやメッセージブローカーを削除して再作成する場合は、PersistentVolumeClaimsを削除することを忘れないでください。これらのリソースは、DatabaseやMessageBrokerリソースを削除しても削除されません。

これで、クラスターリソースが処理できる数だけ、データベースやメッセージブローカーのインスタンスを作成できます！

## カンファレンスアプリケーションをデプロイしましょう

さて、2つのデータベースとメッセージブローカーが実行されているので、アプリケーションサービスがこれらのインスタンスに接続されていることを確認する必要があります。まず、Conference ApplicationチャートでHelmの依存関係を無効にし、アプリケーションがインストールされるときにデータベースとメッセージブローカーがインストールされないようにする必要があります。これは、`install.infrastructure`フラグを`false`に設定することで実現できます。

そのために、新しく作成したデータベースに接続するためのサービスの設定を含む`app-values.yaml`ファイルを使用します：

```shell
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0 -f app-values.yaml
```

`app-values.yaml`の内容は次のようになります：
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

`app-values.yaml`ファイルが、例示ファイルで指定したデータベース（`my-db-keyavalue`と`my-db-sql`）とメッセージブローカー（`my-mb-kafka`）の名前に依存していることに注意してください。異なる名前で他のデータベースとメッセージブローカーを要求した場合は、このファイルを新しい名前で適応させる必要があります。

アプリケーションポッドが起動すると、ブラウザで[http://localhost:8080](http://localhost:8080)にアクセスしてアプリケーションにアクセスできるはずです。
ここまでたどり着いた場合、Crossplane Compositionsを使用してマルチクラウドインフラストラクチャをプロビジョニングできるようになりました。[@asarenkansah](https://github.com/asarenkansah)によって貢献され[AWS Crossplane Compositions Tutorial](aws/)をチェックしてください。アプリケーションインフラストラクチャのプロビジョニングをアプリケーションコードから分離することで、クラウドプロバイダー間のポータビリティを可能にするだけでなく、アプリケーションサービスをプラットフォームチームが管理できるインフラストラクチャと接続できるようになります。

## クリーンアップ

このチュートリアル用に作成したKinDクラスターを削除したい場合は、次のコマンドを実行できます：

```shell
kind delete clusters dev
```

## 次のステップ

Google Cloud Platform、Microsoft Azure、Amazon AWSなどのクラウドプロバイダーにアクセスできる場合は、これらのプラットフォーム用の**Crossplaneプロバイダー**を強くお勧めします。これらのプロバイダーをインストールし、Crossplane Helmプロイダーを使用する代わりにクラウドリソースをプロビジョニングすることで、これらのツールがどのように機能するかについての実際の経験を得ることができます。

第5章で言及したように、マネージドサービスとして提供されていないインフラストラクチャコンポーネントを必要とするサービスをどのように扱いますか？ Google Cloud Platformの場合、プロビジョニングできるマネージドKafkaサービスは提供していません。HelmチャートまたはVMを使用してKafkaをインストールするか、KafkaをGoogle PubSubなどのマネージドサービスに切り替えますか？同じサービスの2つのバージョンを維持しますか？

## まとめと貢献

このチュートリアルでは、アプリケーションインフラストラクチャのプロビジョニングをアプリケーションのデプロイメントから分離することに成功しました。これにより、異なるチームがオンデマンドでリソースを要求し（Crossplane compositionsを使用）、独立して進化できるアプリケーションサービスを提供できるようになります。

開発目的でHelmチャートの依存関係を使用し、アプリケーションの完全に機能するインスタンスをすぐに稼働させることは素晴らしいです。より重要な環境では、ここで示したようなアプローチを採用したいかもしれません。このアプローチでは、各サービスが必要とするコンポーネントとアプリケーションを接続する複数の方法があります。

このチュートリアルを改善したいですか？issueを作成するか、[Twitter](https://twitter.com/salaboy)でメッセージを送るか、プルリクエストを送信してください。
