# 第2章 :: クラウドネイティブアプリケーションの課題

---
_🌍 翻訳_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [Português (Portuguese)](README-pt.md) | [日本語](README-ja.md)

> **注意:** このコンテンツは素晴らしいクラウドネイティブコミュニティの [🌟 コントリビューター](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) によってもたらされています！

---

このショートチュートリアルでは、Helmを使用して `Conference Application` をローカルのKubernetes KinD (Kubernetes in Docker) クラスターにインストールします。

> [!NOTE]
> Helmチャートは、HelmチャートリポジトリまたはHelm 3.7以降ではOCI (Open Container Initiative: オープンコンテナイニシアティブ) コンテナーとしてコンテナーレジストリに公開できます。

## KinDでローカルクラスターを作成する

> [!Important]
> すべてのチュートリアルの前提条件を満たしていることを確認してください。前提条件は[こちら](../chapter-1/README-ja.md#チュートリアルの前提条件)で確認できます。

以下のコマンドを使用して、3つのワーカーノードと1つのコントロールプレーンを持つKinDクラスターを作成します。

```shell
cat <<EOF | kind create cluster --name dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
- role: worker
EOF

```

![3つのワーカーノード](imgs/cluster-topology.png)

### アプリケーションとその他のコンポーネントをインストールする前にコンテナーイメージをロードする

`kind-load.sh` スクリプトは、言い換えると、アプリケーションで使用するコンテナーイメージを事前に取得し、KinDクラスターにロードします。

ここでの目的は、クラスターのプロセスを最適化することです。アプリケーションをインストールする際に、必要なすべてのコンテナーイメージを取得するのに10分以上待つ必要がなくなります。すべてのイメージがKinDクラスターに事前にロードされていれば、PostgreSQL、Redis、Kafkaのブートストラップに必要な約1分でアプリケーションが起動するはずです。

それでは、必要なイメージをKinDクラスターに取得しましょう。

> [!Important]
> 次のステップで述べるスクリプトを実行すると、必要なすべてのイメージを取得し、KinDクラスターのすべてのノードにロードします。クラウドプロバイダー上でサンプルを実行しいる場合、クラウドプロバイダーはコンテナーレジストリへのギガビット接続を持っているため、これらのイメージを数秒で取得できるかもしれません。その場合、このステップは必要ないかもしれません。

ターミナルで `chapter-2` ディレクトリにアクセスし、そこからスクリプトを実行します。

```shell
./kind-load.sh
```

> [!Note]
> MacOSでDocker Desktopを実行していて、仮想ディスクのサイズを小さく設定している場合、次のようなエラーが発生することがあります。
>
> ```shell
> $ ./kind-load.sh
> ...
> Command Output: Error response from daemon: write /var/lib/docker/...
> /layer.tar: no space left on device
> ```
>
> 「設定」->「リソース」メニューで仮想ディスク制限の値を変更できます。
>   ![MacOS Docker Desktopの仮想ディスク設定](imgs/macos-docker-desktop-virtual-disk-setting.png)


### NGINX Ingress Controller インストール

ラップトップからクラスター内で実行されているサービスにトラフィックをルーティングするには、NGINX Ingress Controller が必要です。NGINX Ingress Controller は、クラスター内で実行されているルーターとして機能しますが、外部にも公開されています。

```shell
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/release-1.8/deploy/static/provider/kind/deploy.yaml
```

次に進む前に、`ingress-nginx` 内のポッドが正しく起動していることを確認します。
```shell
> kubectl get pods -n ingress-nginx
NAME                                        READY   STATUS      RESTARTS   AGE
ingress-nginx-admission-create-cflcl        0/1     Completed   0          62s
ingress-nginx-admission-patch-sb64q         0/1     Completed   0          62s
ingress-nginx-controller-5bb6b499dc-7chfm   0/1     Running     0          62s
```

これにより、`http://localhost` からクラスター内のサービスにトラフィックをルーティングできるようになります。KinDがこのように動作するためには、クラスターを作成する際にコントロールプレーンノードに追加のパラメーターとラベルを指定しました。
```yaml
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true" #これにより、Ingress Controllerをコントロールプレーンノードにインストールできます
  extraPortMappings:
  - containerPort: 80 # これにより、ローカルホストのポート80をIngress Controllerにバインドできるので、クラスター内で実行されているサービスにトラフィックをルーティングできます
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
```

クラスターとIngress Controllerのインストールと設定が完したら、アプリケーションのインストールに進みます。


## Conference Application のインストール

Helm 3.7以降では、OCI イメージを使用してHelmチャートを公開、ダウンロード、インストールできます。このアプローチでは、Docker HubをHelmチャートレジストリとして使用します。

Conference Applicationをインストールするには、次のコマンドを実行するだけです。

```shell
helm install conference oci://docker.io/salaboy/conference-app --version v1.0.0
```

チャートの詳細を確認するには、次のコマンドを実行できます。

```shell
helm show all oci://docker.io/salaboy/conference-app --version v1.0.0
```

以下のコマンドで、すべてのアプリケーションポッドが起動していることを確認します。

```shell
kubectl get pods
```

> [!Note]
> インターネット接続が遅い場合、アプリケーションの起動に時間がかかることがあります。アプリケーションのサービスは、いくつかのインフラストラクチャーコンポーネント（Redis、Kafka、PostgreSQL）に依存しているため、これらのコンポーネントが起動し、サービスが接続できる状態になる必要があります。
>
> Kafkaのようなコンポーネントは335MB以上、PostgreSQLは88MB以上、Redisは35MB以上と非常に重いです。

最終的には、次のような表示になるはずです。これには数分かかる可能性があります。

```shell
NAME                                                           READY   STATUS    RESTARTS      AGE
conference-agenda-service-deployment-7cc9f58875-k7s2x          1/1     Running   4 (45s ago)   2m2s
conference-c4p-service-deployment-54f754b67c-br9dg             1/1     Running   4 (65s ago)   2m2s
conference-frontend-deployment-74cf86495-jthgr                 1/1     Running   4 (56s ago)   2m2s
conference-kafka-0                                             1/1     Running   0             2m2s
conference-notifications-service-deployment-7cbcb8677b-rz8bf   1/1     Running   4 (47s ago)   2m2s
conference-postgresql-0                                        1/1     Running   0             2m2s
conference-redis-master-0                                      1/1     Running   0             2m2s
```

上記の出力結果から、おそらくKafkaの起動が遅かったため、Kubernetesによってサービスが最初に起動され、Kafkaの準備ができるのを待つために再起動したことが分かります。これは、ポッドの `RESTARTS` 列の値から推測できます。

アプリケーションの起動が確認できたら、ブラウザで [http://localhost](http://localhost) を開いてアプリケーションを確認してみましょう。

![conference app](imgs/conference-app-homepage.png)

------
## [重要] クリーンアップ - _!!! 必ず読んでください!!_

Conference ApplicationはPostgreSQL、Redis、Kafkaをイストールするため、アプリケーションを削除して再インストールしたい場合（ガイドを進めるにつれて行う予定です）、関連するPersistentVolumeClaim（PVC: 永続ボリューム要求）を必ず削除する必要があります。

これらのPVCは、データベースとKafkaのデータを保存するために使用されるボリュームです。インストール間でこれらのPVCを削除しないと、サービスが古い認証情報を使用して新しくプロビジョニングされたデータベースに接続しようとします。

PVCをリストアップして削除するには、次のコマンドを使用します。

```shell
kubectl get pvc
```

次のような出力が表示されるはずです。

```shell
NAME                                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-conference-kafka-0                Bound    pvc-2c3ccdbe-a3a5-4ef1-a69a-2b1022818278   8Gi        RWO            standard       8m13s
data-conference-postgresql-0           Bound    pvc-efd1a785-e363-462d-8447-3e48c768ae33   8Gi        RWO            standard       8m13s
redis-data-conference-redis-master-0   Bound    pvc-5c2a96b1-b545-426d-b800-b8c71d073ca0   8Gi        RWO            standard       8m13s
```

そして、次のコマンドで削除します。
```shell
kubectl delete pvc  data-conference-kafka-0 data-conference-postgresql-0 redis-data-conference-redis-master-0
```

PVCの名前は、チャートをインストールする際に使用したHelmリリース名に基づいて変更されます。

最後に、KinDクラスターを完全に削除したい場合は、次のコマンドを実行できます。

```shell
kind delete clusters dev
```

-------
## 次のステップ

クラウドプロバイダーでホストされている実際のKubernetesクラスターで実践することを強くお勧めします。ほとんどのクラウドプロバイダーは無料トイアルを提供しており、Kubernetesクラスターを作成し、これらすべてのサンプルを実行できます（詳細については[このリポジトリ](https://github.com/learnk8s/free-kubernetes)を確認してください）。

クラウドプロバイダーでクラスターを作成し、アプリケーションを起動して実行できれば、第2章で扱ったすべてのトピックに関する実際の経験を積むことができます。

## まとめと貢献

このショートチュートリアルでは、**Conference Application** のウォーキングスケルトン[^1]のインストールに成功しました。このアプリケーションを残りの章の例として使用します。このアプリケーションがあなたにとって機能することを確認してください。これは、Kubernetesクラスターの使用と対話の基本をカバーしています。

このチュートリアルを改善したいですか？ [issue](https://github.com/salaboy/platforms-on-k8s/issues/new)を作成するか、[Twitter](https://twitter.com/salaboy)でメッセージを送るか、[プルリクエスト](https://github.com/salaboy/platforms-on-k8s/compare)を送信してください。
