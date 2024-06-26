# 第6章 :: Kubernetes上にプラットフォームを構築しよう

---
_🌍 利用可能な言語_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md)

> **注意:** これは素晴らしいクラウドネイティブコミュニティの [🌟 コントリビューター](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) によってもたらされました！

---

このステップバイステップのチュートリアルでは、Kubernetes APIの力を再利用してプラットフォームのAPIを作成します。プラットフォームが開発チームを支援できる最初のユースケースは、新しい開発環境を作成し、セルフサービスアプローチを提供することです。

この例を構築するために、Cloud Native Computing Foundationでホストされている2つのオープンソースプロジェクト、Crossplaneとvclusterを使用します。

## インストール

Crossplaneをインストールするには、Kubernetesクラスターが必要です。[第2章](../chapter-2/README-ja.md#kindでローカルクラスターを作成する)で行ったように、KinDを使用してクラスターを作成できます。

その後、[第5章](https://github.com/salaboy/platforms-on-k8s/tree/main/chapter-5#crossplaneのインストール)で行ったように、クラスターにCrossplaneとCrossplane Helmプロバイダーをインストールできます。

このチュートリアルでは[`vcluster`](https://www.vcluster.com/)を使用しますが、vclusterを動作させるためにクラスターに何かをインストールする必要はありません。`vcluster`に接続するために`vcluster` CLIが必要です。公式サイトの指示に従ってインストールできます: [https://www.vcluster.com/docs/getting-started/setup](https://www.vcluster.com/docs/getting-started/setup)

## 環境APIの定義

環境は、Conference Applicationが開発用インストールされるKubernetesクラスターを表します。アイデアは、チームが作業を行うためのセルフサービス環境を提供することです。

このチュートリアルでは、環境APIと、Helmプロバイダーを使用して新しい`vcluster`インスタンスを作成するCrossplane Compositionを定義します。

[ここにある環境のCrossplane Composite Resource Definition (XRD)](resources/env-resource-definition.yaml)と[ここにあるCrossplane Composition](resources/composition-devenv.yaml)を確認してください。このリソースはCrossplane Helmプロバイダーを使用して新しい`vcluster`のプロビジョニングを設定します。[この設定はここで確認できます](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-6/resources/composition-devenv.yaml#L24)。新しい`vcluster`が作成されると、compositionは再びCrossplane Helmプロバイダーを使用してConference Applicationをその中にインストールしますが、今回は[作成されたばかりの`vcluster` APIを指すように設定されています](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-6/resources/composition-devenv.yaml#L87)。[ここで確認できます](https://github.com/salaboy/platforms-on-k8s/blob/main/chapter-6/resources/composition-devenv.yaml#L117)。

両方のXRDをインストールしましょう：

```shell
kubectl apply -f resources/definitions
```

XRDが定義されたので、Compositionをインストールしましょう：

```shell
kubectl apply -f resources/compositions
```

以下が表示されるはずです：

```shell
composition.apiextensions.crossplane.io/dev.env.salaboy.com created
compositeresourcedefinition.apiextensions.crossplane.io/environments.salaboy.com created
```

`vcluster`を使用する環境リソースとCrossplane Compositionにより、チームはオンデマンドで環境をリクエストできるようになりました。

## 新しい環境のクエスト

新しい環境をリクエストするために、チームは次のような環境リソースを作成できます：

```yaml
apiVersion: salaboy.com/v1alpha1
kind: Environment
metadata:
  name: team-a-dev-env
spec:
  compositionSelector:
    matchLabels:
      type: development
  parameters: 
    installInfra: true
```

クラスターに送信すると、Crossplane Compositionが起動し、Conference Applicationのインスタンスを含む新しい`vcluster`を作成します。

```shell
kubectl apply -f team-a-dev-env.yaml
```

以下が表示されるはずです：

```shell
environment.salaboy.com/team-a-dev-env created
```

環境の状態は常に次のコマンドで確認できます：

```shell
> kubectl get env
NAME             CONNECT-TO             TYPE          INFRA   DEBUG   SYNCED   READY   CONNECTION-SECRET   AGE
team-a-dev-env   team-a-dev-env-jp7j4   development   true    true    True     False   team-a-dev-env      1s
```

Crossplaneがcompositionに関連すリソースを作成および管理していることを確認するには、次のコマンドを実行します：

```shell
> kubectl get managed
NAME                            CHART            VERSION          SYNCED   READY   STATE      REVISION   DESCRIPTION        AGE
team-a-dev-env-jp7j4-8lbtj      conference-app   v1.0.0           True     True    deployed   1          Install complete   57s
team-a-dev-env-jp7j4-vcluster   vcluster         0.15.0-alpha.0   True     True    deployed   1          Install complete   57s
```

これらの管理されたリソースは、作成されているHelmリリースに他なりません：

```shell
kubectl get releases
NAME                            CHART            VERSION          SYNCED   READY   STATE      REVISION   DESCRIPTION        AGE
team-a-dev-env-jp7j4-8lbtj      conference-app   v1.0.0           True     True    deployed   1          Install complete   45s
team-a-dev-env-jp7j4-vcluster   vcluster         0.15.0-alpha.0   True     True    deployed   1          Install complete   45s
```

次に、以下のコマンドを実行してプロビジョニングされた環境に接続できます（CONNECT-TOカラムのvcluster名を使用します）：

```shell
vcluster connect team-a-dev-env-jp7j4 --server https://localhost:8443 -- zsh
```

`vcluster`に接続すると、異なるKubernetesクラスターにいます。利用可能なすべての名前空間をリストすると、次のように表示されるはずです：

```shell
kubectl get ns
NAME              STATUS   AGE
default           Active   64s
kube-system       Active   64s
kube-public       Active   64s
kube-node-lease   Active   64s
```

ご覧の通り、ここにはCrossplaneはインストールされていません。しかし、このクラスター内のすべてのポッドをリストすると、すべてのアプリケーションポッドが実行されているのが分かります：

```shell
NAME                                                              READY   STATUS    RESTARTS      AGE
conference-app-kafka-0                                            1/1     Running   0             103s
conference-app-postgresql-0                                       1/1     Running   0             103s
conference-app-c4p-service-deployment-57d4ddcd68-45f6h            1/1     Running   2 (99s ago)   104s
conference-app-agenda-service-deployment-9bf7946c9-mmx8h          1/1     Running   2 (98s ago)   104s
conference-app-redis-master-0                                     1/1     Running   0             103s
conference-app-frontend-deployment-c8c64c54d-lntnw                1/1     Running   2 (98s ago)   104s
conference-app-notifications-service-deployment-64ff7bcdf8nbvhl   1/1     Running   3 (80s ago)   104s
```

このクラスターへのポートフォワーディングを行って、次のコマンドを使用してアプリケーションにアクセスすることもできます：

```shell
kubectl port-forward svc/frontend 8080:80
```

これで、アプリケーシンは[http://localhost:8080](http://localhost:8080)で利用可能になります。

ターミナルで`exit`と入力することで、`vcluster`コンテキストを終了できます。

## プラットフォームのインターフェースの簡素化

プラットフォームAPIとの対話を簡素化し、チームがプラットフォームクラスターに接続する必要をなくし、Kubernetes APIへのアクセスの必要性を排除するためにさらに一歩進めることができます。

このセクションでは、チームがウェブサイトまたは簡素化されたREST APIを使用して新しい環境をリクエストできるようにする管理者用ユーザーインターフェースをデプロイします。

管理者用ユーザーインターフェースをインストールする前に、`vcluster`セッション内にいないことを確認する必要があります（ターミナルで`exit`と入力して`vcluster`コンテキストを終了できます）。現在接続しているクラスターに`crossplane-system`名前空間があることを確認してください。

この管理者用ユーザーインターフェースはHelmを使用してインストールできます：

```shell
helm install admin oci://docker.io/salaboy/conference-admin --version v1.0.0
```

インストールが完了したら、次のコマンドを実行して管理者UIにポートフォワーディングできます：

```shell
kubectl port-forward svc/admin 8081:80
```

これで、[http://localhost:8081](http://localhost:8081)にあるシンプルなインターフェースを使用して環境を作成および確認できます。環境の準備が整うまで待つと、接続に使用する`vcluster`コマンドが表示されます。

![imgs/admin-ui.png](imgs/admin-ui.png)

このシンプルなインターフェースを使用することで、開発チームはすべてのプラットフォームツール（CrossplaneやArgo CDなど）を持つクラターのKubernetes APIに直接アクセスする必要がなくなります。

ユーザーインターフェースの他に、プラットフォーム管理アプリケーションは、Kubernetes Resource Modelに従う代わりにリソースがどのように見えるかを柔軟に定義できる簡素化されたRESTエンドポイントのセットを提供します。例えば、Kubernetes APIが必要とするすべてのメタデータを持つKubernetesリソースを持つ代わりに、新しい環境を作成するために次のようなJSONペイロードを使用できます：

```json
{
    "name": "team-curl-dev-env",
    "parameters":{
        "type": "development",
        "installInfra": true,
        "frontend":{
            "debug": true
        }
    }
}
```

次のコマンドを実行してこの環境を作成できます：

```shell
curl -X POST -H "Content-Type: application/json" -d @team-a-dev-env-simple.json http://localhost:8081/api/environments/
```

そして、すての環境を次のコマンドでリストできます：

```shell
curl localhost:8081/api/environments/
```

または、次のコマンドで1つの環境を削除できます：

```shell
curl -X DELETE http://localhost:8081/api/environments/team-curl-dev-env
```

このアプリケーションは、KubernetesとKubernetesの外部の世界との間のファサードとして機能します。組織のニーズに応じて、これらの抽象化（API）を早い段階で持つことで、プラットフォームチームは裏側でツールとワークフローの決定を柔軟に変更できるようになります。

## クリーンアップ

これらのチュートリアル用に作成したKinDクラスターを削除したい場合は、次のコマンドを実行できます：

```shell
kind delete clusters dev
```

## 次のステップ

第5章で行ったように、管理者用ユーザーインターフェースを拡張してデータベースとメッセージブローカーを成できますか？何が必要ですか？変更が必要な場所を理解することで、Kubernetes APIと対話し、消費者向けに簡素化されたインターフェースを提供するコンポーネントの開発に関する実践的な経験が得られます。

`vcluster`の代わりに実際のクラスターを使用するための独自のcompositionを作成できますか？どのようなシナリオで実際のクラスターを使用し、どのような場合に`vcluster`を使用しますか？

Kubernetes KinDで実行する代わりに、実際のKubernetesクラスターでこれを実行するには、どのような追加のステップが必要ですか？

## まとめと貢献

このチュートリアルでは、オンデマンドの開発環境をプロビジョニングするためにKubernetesリソースモデルを再利用して新しいプラットフォームAPIを構築しました。さらに、プラットフォーム管理アプリケーションを使用して、チームにKubernetesの仕組みや、プラットフォームの構築に使用した基盤となる詳細、プロジェクト、技術について学ぶ必要を押し付けることなく、同じ機能を公開する簡素化された層を作成しました。

契約（この例では環境リソース定義）に依存することで、プラットフォームチームは要件と利用可能なツールに応じて環境をプロビジョニングするメカニズムを柔軟に変更できます。

このチュートリアルを改善したいですか？issueを作成するか、[Twitter](https://twitter.com/salaboy)でメッセージを送るか、プルリクエトを送信してください。
