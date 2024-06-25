# 環境パイプライン

---
_🌍 利用可能な言語_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md)
> **注**: これは素晴らしいクラウドネイティブコミュニティの[🌟貢献者](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)によってもたらされました！

---

この短いチュートリアルでは、[ArgoCD](https://argo-cd.readthedocs.io/en/stable/)を使用してステージング環境パイプラインをセットアップします。カンファレンスアプリケーションのインスタンスを含む環境を構成します。

ステージング環境の設定はGitリポジトリを使用して定義します。[`argo-cd/staging`ディレクトリ](argo-cd/staging/)には、複数のKubernetesクラスタに同期できるHelmチャートの定義が含まれています。

## 前提条件とインストール

- Kubernetesクラスタが必要です。このチュートリアルでは[KinD](https://kind.sigs.k8s.io/)を使用します
- クラスタにArgoCDをインストールします。[こちらの手順](https://argo-cd.readthedocs.io/en/stable/getting_started/)に従い、オプションで`argocd` CLIをインストールします
- アプリケーションの設定を変更したい場合は、[このリポジトリ](http://github.com/salaboy/platforms-on-k8s/)をフォーク/コピーする必要があります。リポジトリへの書き込みアクセスが必要です。`chapter-4/argo-cd/staging/`ディレクトリを使用します

[第2章で行ったようにKinDクラスタを作成します](../chapter-2/README-ja.md#kindでローカルクラスターを作成する)。

nginx-ingressコントローラーを使用してクラスタが稼働したら、クラスタにArgo CDをインストールしましょう：

```shell
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

以下のような出力が表示されるはずです：

```shell
> kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
namespace/argocd created
customresourcedefinition.apiextensions.k8s.io/applications.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/applicationsets.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/appprojects.argoproj.io created
serviceaccount/argocd-application-controller created
serviceaccount/argocd-applicationset-controller created
serviceaccount/argocd-dex-server created
serviceaccount/argocd-notifications-controller created
serviceaccount/argocd-redis created
serviceaccount/argocd-repo-server created
serviceaccount/argocd-server created
role.rbac.authorization.k8s.io/argocd-application-controller created
role.rbac.authorization.k8s.io/argocd-applicationset-controller created
role.rbac.authorization.k8s.io/argocd-dex-server created
role.rbac.authorization.k8s.io/argocd-notifications-controller created
role.rbac.authorization.k8s.io/argocd-server created
clusterrole.rbac.authorization.k8s.io/argocd-application-controller created
clusterrole.rbac.authorization.k8s.io/argocd-server created
rolebinding.rbac.authorization.k8s.io/argocd-application-controller created
rolebinding.rbac.authorization.k8s.io/argocd-applicationset-controller created
rolebinding.rbac.authorization.k8s.io/argocd-dex-server created
rolebinding.rbac.authorization.k8s.io/argocd-notifications-controller created
rolebinding.rbac.authorization.k8s.io/argocd-redis created
rolebinding.rbac.authorization.k8s.io/argocd-server created
clusterrolebinding.rbac.authorization.k8s.io/argocd-application-controller created
clusterrolebinding.rbac.authorization.k8s.io/argocd-server created
configmap/argocd-cm created
configmap/argocd-cmd-params-cm created
configmap/argocd-gpg-keys-cm created
configmap/argocd-notifications-cm created
configmap/argocd-rbac-cm created
configmap/argocd-ssh-known-hosts-cm created
configmap/argocd-tls-certs-cm created
secret/argocd-notifications-secret created
secret/argocd-secret created
service/argocd-applicationset-controller created
service/argocd-dex-server created
service/argocd-metrics created
service/argocd-notifications-controller-metrics created
service/argocd-redis created
service/argocd-repo-server created
service/argocd-server created
service/argocd-server-metrics created
deployment.apps/argocd-applicationset-controller created
deployment.apps/argocd-dex-server created
deployment.apps/argocd-notifications-controller created
deployment.apps/argocd-redis created
deployment.apps/argocd-repo-server created
deployment.apps/argocd-server created
statefulset.apps/argocd-application-controller created
networkpolicy.networking.k8s.io/argocd-application-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-applicationset-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-dex-server-network-policy created
networkpolicy.networking.k8s.io/argocd-notifications-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-redis-network-policy created
networkpolicy.networking.k8s.io/argocd-repo-server-network-policy created
networkpolicy.networking.k8s.io/argocd-server-network-policy created
```

`port-forward`を使用してArgoCDユーザーインターフェースにアクセスできます。**新しいターミナル**で以下を実行してください：

```shell
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

**注意**: ArgoCDポッドが起動するまで待つ必要があります。初回実行時は、インターネットからコンテナイメージをフェッチする必要があるため、時間がかかります。

ブラウザで[http://localhost:8080](http://localhost:8080)にアクセスすると、ユーザーインターフェースが表示されます。

<img src="imgs/argocd-warning.png" width="600">

**注意**: デフォルトでは、インストールはHTTPSではなくHTTPを使用して動作します。そのため、警告を受け入れる必要があります（Chromeの「詳細設定」ボタンをクリック）。その後、「localhost（安全ではありません）に進む」を選択してください。

<img src="imgs/argocd-proceed.png" width="600">

これでログインページに移動します：

<img src="imgs/argocd-login.png" width="600">

ユーザー名は`admin`で、ArgoCDダッシュボードのパスワードを取得するには以下を実行します：

```shell
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

ログインすると、空のホーム画面が表示されます：

<img src="imgs/argocd-dashboard.png" width="600">

それでは、ステージング環境をセットアップしましょう。

# ステージング環境用のアプリケーションのセットアップ

このチュートリアルでは、単一の名前空間をステージング環境の現に使用します。Argo CDには制限がなく、ステージング環境を完全に別のKubernetesクラスタにすることもできます。

まず、ステージング環境用の名前空間を作成しましょう：

```shell
kubectl create ns staging
```

以下のような出力が表示されるはずです：

```shell
> kubectl create ns staging
namespace/staging created
```

注意: 代替として、ArgoCDアプリケーション作成時に「Auto Create Namespace」オプションを使用することもできます。

Argo CDをインストールしたら、ユーザーインターフェースにアクセスしてプロジェクトをセットアップできます。

<img src="imgs/argocd-dashboard.png" width="600">

**「+ New App」**ボタンをクリックし、以下の詳細を使用してプロジェクトを設定します：

<img src="imgs/argocd-app-creation.png" width="600">

以下は、私が使用したアプリケーション作成の入力内容です：
- アリケーション名: "staging-environment"
- プロジェクト: "default"
- 同期ポリシー: "Automatic"
- ソースリポジトリ: [https://github.com/salaboy/platforms-on-k8s](https://github.com/salaboy/platforms-on-k8s)（ここであなたのフォークを指定できます）
- リビジョン: "HEAD"
- パス: "chapter-4/argo-cd/staging/"
- クラスタ: "https://kubernetes.default.svc" 
- 名前空間: "staging"

<img src="imgs/argocd-app-creation2.png" width="600">

他の値はデフォルトのままにして、上部の**Create**をクリックします。

アプリケーションが作成されると、**Automatic**モードを選択したため、自動的に変更が同期されます。

<img src="imgs/argocd-syncing.png" width="600">

アプリケーションをクリックして展開すると、作成されているすべてのリソースの完全なビューを確認できます：

<img src="imgs/app-detail.png" width="600">

ローカル環境で実行している場合は`port-forward`を使用してアプリケーションにアクセスできます。**新しいターミナル**で以下を実行してください：

```shell
kubectl port-forward svc/frontend -n staging 8081:80
```

アプリケーションのポッドが起動して実行されるのを待ち、ブラウザで[http://localhost:8081](http://localhost:8081)にアクセスすると、アプリケーションにアクセスできます。

<img src="imgs/app-home.png" width="600">

通常通り、`kubectl`を使用してポッドとサービスのステータスを監視できます。アプリケーションポッドが準備できているかを確認するには、以下を実行します：

```shell
kubectl get pods -n staging
```

以下のような出力が表示されるはずです：

```shell
> kubectl get pods -n staging
NAME                                                              READY   STATUS    RESTARTS        AGE
stating-environment-agenda-service-deployment-6c9cbb9695-xj99z    1/1     Running   5 (6m ago)      8m4s
stating-environment-c4p-service-deployment-69d485ffd8-q96z4       1/1     Running   5 (5m52s ago)   8m4s
stating-environment-frontend-deployment-cd76bdc8c-58vzr           1/1     Running   5 (6m3s ago)    8m4s
stating-environment-kafka-0                                       1/1     Running   0               8m4s
stating-environment-notifications-service-deployment-5c9b5bzb5p   1/1     Running   5 (6m13s ago)   8m4s
stating-environment-postgresql-0                                  1/1     Running   0               8m4s
stating-environment-redis-master-0                                1/1     Running   0               8m4s
```

**注意**: いくつかの再起動（RESTARTSカラム）は問題ありません。一部のサービスは、インフラストラクチャ（Redis、PostgreSQL、Kafka）が起動してから健全になる必要があるためです。

## ステージング環境でのアプリケーション設定の変更

サービスのバージョンや設を更新するには、[staging](staging/)ディレクトリ内の[Chart.yaml](argo-cd/staging/Chart.yaml)ファイルまたは[values.yaml](argo-cd/staging/values.yaml)ファイルを更新できます。

この例では、ArgoCDアプリケーションの詳細とそのパラメータを更新することで、アプリケーション設定を変更できます。

実際のアプリケーションではこのような操作は行いませんが、ここではステージング環境が定義されているGitHubリポジトリの変更をシミュレートしています。

<img src="imgs/argocd-change-parameters.png" width="600">

アプリケーションの詳細/パラメータを編集し、このアプリケーションに使用するvaluesファイルとして`values-debug-enabled.yaml`を選択してください。このファイルはフロントエンドサービスにデバッグフラグを設定し、最初のインストールに使用された元の`values.yaml`ファイルを変更したことをシミュレートしています。

<img src="imgs/argocd-new-values.png" width="600">

port-forwardingを使用していたため、以下のコマンドを再度実行する必要があるかもしれません：

```shell
kubectl port-forward svc/frontend -n staging 8081:80
```

これは、フロントエンドサービスのポッドが新しく設定されたバージョンに置き換えられるため、port-forwardingを再起動して新しいポッドをターゲットにする必要があるからです。

フロントエンドが起動して実行されると、バックオフィスセクションにデバッグタブが表示されるはずです：

![](imgs/app-debug.png)

## クリーンアップ

このチュートリアル用に作成したKinDクラスタを削除したい場合は、以下を実行できます：

```shell
kind delete clusters dev
```

## 次のステップ

Argo CDはGitOpsを実装するための1つのプロジェクトに過ぎません。このチュートリアルをFlux CDで再現できますか？どちらが好みですか？あなたの組織ではすでにGitOpsツールを使用していますか？そのツールを使用してカンファレンスアプリケーションのウォーキングスケルトンをKubernetesクラスタにデプロイするには何が必要でしょうか？

`production-environment`のような別の環境を作成し、`notifications-service`の新しいリリースがステージング環境から本番環境に移行するためのフローを記述できますか？本番環境の設定をどこに保存しますか？

## まとめと貢献

このチュートリアルでは、Argo CDアプリケーションを使用して**ステージング環境**を作成しました。これにより、GitHubリポジトリ内にある設定を、KinDで稼働しているKubernetesクラスタに同期させることができました。GitHubリポジトリの内容を変更してArgoCD アプリケーションを更新すると、ArgoCDは環境が同期されていないことに気づきます。自動同期戦略を使用している場合、ArgoCDは設定に変更があったことに気づくたびに、自動的に同期ステップを実行します。詳細については、[プロジェクトのウェブサイト](https://argo-cd.readthedocs.io/en/stable/)または[私のブログ](https://www.salaboy.com)をチェックしてください。

このチュートリアルを改善したいですか？イシューを作成するか、[Twitter](https://twitter.com/salaboy)でメッセージを送るか、プルリクエストを送信してください。
