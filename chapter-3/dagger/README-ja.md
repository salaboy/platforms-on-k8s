# Dagger in Action

このショートチュートリアルでは、各サービスのビルド、テスト、パッケージ化、公開のために提供されているDaggerサービスパイプラインを見ていきます。
これらのパイプラインは、Dagger Go SDKを使用してGoで実装されており、各サービスのビルドとコンテナの作成を行います。アプリケーションのHelmチャートをビルドおよび公開するための別のパイプラインも提供されています。

## 要件

これらのパイプラインをローカルで実行するには、以下が必要です。
- [Goのインストール](https://go.dev/doc/install)
- [コンテナランタイム（ローカルで実行されているDockerなど）](https://docs.docker.com/get-docker/)

パイプラインをKubernetesクラスタ上でリモートで実行するには、[KinD](https://kind.sigs.k8s.io/)または利用可能な任意のKubernetesクラスタを使用できます。

## パイプラインを実行してみましょう

すべてのサービスが非常によく似ているため、同じパイプライン定義を使用して、各サービスを個別にビルドするようにパラメータ化できます。

このリポジトリをクローンし、[Conference Applicationディレクトリ](../../conference-application/)からパイプラインをローカルで実行できます。

`service-pipeline.go`ファイル内で定義されているタスクを実行できます。

```shell
go mod tidy
go run service-pipeline.go build <SERVICE DIRECTORY>
```

すべてのサービスに対して、以下のタスクが定義されています。
- `build`: サービスのソースコードをビルドし、そのコンテナを作成します。このゴールは、ビルドしたいサービスのディレクトリを引数として期待します。
- `test`: サービスをテストします。ただし、最初にすべてのサービスの依存関係（テストの実行に必要なコンテナ）を開始します。このゴールは、テストしたいサービスのディレクトリを引数として期待します。
- `publish`: 作成されたコンテナイメージをコンテナレジストリに公開します。これには、コンテナレジストリにログインし、コンテナイメージに使用するタグ名を指定する必要があります。このゴールは、ビルドおよび公開したいサービスのディレクトリと、コンテナをプッシュする前にスタンプするために使用されるタグを引数として期待します。

`go run service-pipeline.go all notifications-service v1.0.0-dagger`を実行すると、すべてのタスクが実行されます。すべてのタスクを実行する前に、コンテナレジストリにプッシュするには適切な認証情報を提供する必要があるため、すべての前提条件が設定されていることを確認する必要があります

`go run service-pipeline.go build notifications-service`を安全に実行できます。これには認証情報を設定する必要はありません。環境変数を使用してコンテナレジストリとユーザー名を設定できます。例えば、次のようになります。

```shell
CONTAINER_REGISTRY=<YOUR_REGISTRY> CONTAINER_REGISTRY_USER=<YOUR_USER> go run service-pipeline.go publish notifications-service v1.0.0-dagger
```

これには、コンテナイメージを公開したいレジストリにログインしておく必要があります。

さて、開発目的では、これは非常に便利です。なぜなら、CI（継続的インテグレーション）システムが行うのと同じ方法でサービスコードをビルドできるからです。しかし、開発者のラップトップで作成されたコンテナイメージを本番環境で実行したくないですよね？
次のセクションでは、DaggerパイプラインをKubernetesクラスタ内リモートで実行するシンプルなセットアップを紹介します。

## KubernetesでパイプラインをリモートでRunning する

Dagger Pipeline Engineは、コンテナを実行できる場所ならどこでも実行できます。つまり、複雑なセットアップを必要とせずに、Kubernetesで実行できるということです。

このチュートリアルでは、Kubernetesクラスタが必要です。[第2章でKinDを使用して作成した](../../chapter-2/README-ja.md#kindでローカルクラスタを作成する)ように、KinDを使用して作成できます。

このショートチュートリアルでは、ローカルのコンテナランタイムでローカルに実行していたパイプラインを、今度はKubernetes Pod内で実行されているDagger Pipeline Engineに対してリモートで実行します。これは実験的な機能であり、Daggerを実行するための推奨される方法ではありませんが、ポイントを証明するのに役立ちます。

Daggerを使用してKubernetes内でDagger Pipeline Engineを実行してみましょう。

```shell
kubectl run dagger --image=registry.dagger.io/engine:v0.3.13 --privileged=true
```

または、`kubectl apply -f chapter-3/dagger/k8s/pod.yaml`を使用して`chapter-3/dagger/k8s/pod.yaml`マニフェストを適用することもできます。

`dagger` ポッドが実行されていることを確認します。
```shell
kubectl get pods
```

以下のような出力が表示されるはずです。
```shell
NAME     READY   STATUS    RESTARTS   AGE
dagger   1/1     Running   0          49s
```

**注意**: これは理想からは程遠いです。なぜなら、Dagger自体に対して永続化やレプリケーションのメカニズムを設定していないため、この場合、すべてのキャッシュメカニズムは揮発性だからです。詳細については、公式ドキュメントを確認してください。

さて、このリモートービスに対してプロジェクトのパイプラインを実行するには、次の環境変数をエクスポートするだけです。
```shell
export _EXPERIMENTAL_DAGGER_RUNNER_HOST=kube-pod://<podname>?context=<context>&namespace=<namespace>&container=<container>
```

ここで、`<podname>`は`dagger`（ポッドを手動で作成したため）、`<context>`はKubernetesクラスタのコンテキストです。KinDクラスタに対して実行している場合、これは`kind-dev`になる可能性があります。現在のコンテキスト名は、`kubectl config current-context`を実行すると確認できます。最後に、`<namespace>`はDaggerコンテナを実行する名前空間で、`<container>`は再び`dagger`です。KinDに対する私のセットアップでは、次のようになります。

```shell
export _EXPERIMENTAL_DAGGER_RUNNER_HOST="kube-pod://dagger?context=kind-dev&namespace=default&container=dagger"
```

また、私のKinDクラスタ（`kind-dev`とう名前）には、パイプラインに関連するものは何もなかったことにも注意してください。

さて、いずれかのプロジェクトで次のコマンドを実行すると、
```shell
go run service-pipeline.go build notifications-service
```

またはサービスをリモートでテストするには、

```shell
go run service-pipeline.go test notifications-service
```

別のタブで、次のコマンドを実行してDaggerエンジンのログを追跡できます。
```shell
kubectl logs -f dagger
```

ビルドはクラスタ内でリモートで行われます。リモートのKubernetesクラスタ（KinDではない）に対して実行していた場合、サービスとそのコンテナをビルドするためにローカルのコンテナランタイムを用意する必要はありません。
