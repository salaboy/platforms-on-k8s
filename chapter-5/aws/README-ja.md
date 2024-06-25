# AWSクラウドプロバイダー

---
_🌍 利用可能な言語_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md)

> **注意:** これは素晴らしいクラウドネイティブコミュニティの [🌟 コントリビューター](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) によってもたらされました！

---

このステップバイステップのチュートリアルでは、Crossplaneを使用してAWSにRedis、PostgreSQL、Kafkaをプロビジョニングします。

## Crossplaneのインストール

Crossplaneをインストールするには、Kubernetesクラスターが必要です。[第2章](../../chapter-2/README-ja.md#kubernetes-kindを使用してローカルクラスタを作成する)で行ったように、KinDを使用してクラスターを作成できます。

[Crossplane](https://crossplane.io)を独自の名前空間にHelmを使用してインストールしましょう：

```shell
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --wait
```

`kubectl crossplane`プラグインをインストールします：

```shell
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh
sudo mv kubectl-crossplane /usr/local/bin
```

次に、Crossplane AWSプロバイダーをインストールします：
```shell
kubectl crossplane install provider crossplane/provider-aws:v0.21.2
```

数秒後、設定されたプロバイダーを確認すると、Helmが`INSTALLED`および`HEALTHY`と表示されるはずです：

```shell
> kubectl get providers.pkg.crossplane.io
NAME                             INSTALLED   HEALTHY   PACKAGE                               AGE
crossplane-provider-aws         True        True      crossplane/provider-aws:v0.21.2       49s
```

これで、アプリケーションが動作するために必要すべてのコンポーネントをプロビジョニングするために、データベースとメッセージブローカーのCrossplane compositionsをインストールする準備が整いました。

## Crossplane Compositionsを使用したオンデマンドのアプリインフラストラクチャ

Key-Valueデータベース（Redis）、SQLデータベース（PostgreSQL）、メッセージブローカー（Kafka）のCrossplane Compositionsをインストールする必要があります。

```shell
kubectl apply -f resources/
```

Crossplane Compositionリソース（`app-database-redis.yaml`）は、どのクラウドリソースを作成し、どのように一緒に設定する必要があるかを定義します。Crossplane Composite Resource Definition (XRD)（`app-database-resource.yaml`）は、アプリケーション開発チームがこのタイプのリソースを作成することで新しいデータベースをすぐに要求できるようにする簡略化されたイターフェースを定義します。

CompositionsとComposite Resource Definitions (XRDs)については、[resources/](resources/)ディレクトリを確認してください。

AWSアカウントのaws_access_key_idとaws_secret_access_keyを含むテキストファイルを作成します。

```text
[default]
aws_access_key_id = 
aws_secret_access_key = 
```

AWS認証情報を含むKubernetesシークレットを作成します。

```shell
kubectl create secret \
generic aws-secret \
-n crossplane-system \
--from-file=creds=./aws-credentials.txt
```

ProviderConfigを作成します。

```shell
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

### アプリケーションインフラストラクチャをプロビジョニングしましょう

チームが使用する新しいKey-Valueデータベースをプロビジョニングするために、必要なすべてのインフラストラクチャを作成する次のコマンドを実行できます：

```shell
kubectl apply -f my-db-keyvalue.yaml
kubectl apply -f my-db-sql.yaml
kubectl apply -f aws-messagebroker-kafka.yaml
```

## カンファレンスアプリケーションをデプロイしましょう

さて、2つのデータベースとメッセージブローカーが実行されているので、アプリケーションサービスがこれらのインスタンスに接続されていることを確認する必要があります。まず、Conference ApplicationチャートでHelmの依存関係を無効にし、アプリケーションがインストールされるときにデータベースとメッセージブローカーがインストールされないようにする必要があります。これは、`install.infrastructure`フラグを`false`に設定することで実現できます。

そのために、新しく作成したータベースに接続するためのサービスの設定を含む`app-values.yaml`ファイルを使用します：

```shell
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0 -f app-values.yaml
```

新しく作成されたAWSインフラストラクチャの値に基づいて、yamlファイルのコメントアウトされた部分を必ず入力してください。

`app-values.yaml`の内容は次のようになります：
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
