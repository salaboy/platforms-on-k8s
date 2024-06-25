# Tekton 実践ガイド

この短いチュートリアルでは、Tektonのインストール方法と、非常にシンプルなTaskとPipelineの作成方法について説明します。

[Tekton](https://tekton.dev)は、クラウド（特にKubernetes）向けに構築された、柔軟性の高いパイプラインエンジンです。実行できるTaskの種類に制限を設けていないため、どのようなパイプラインでも構築できます。これは、管理サービスでは対応できない特別な要件が必要な場合があるサービスパイプラインの構築に最適です。

最初のTektonパイプラインを実行した後、このチュートリアルには、カンファレンスアプリケーションサービスの構築に使用されるより複雑なサービスパイプラインへのリンクも含まれています。

## Tektonのインストール

Kubernetes クラスタにTektonをインストールし、セットアップするには、下の手順に従ってください。Kubernetesクラスタをお持ちでない場合は、[第2章で行ったようにKinDを使用してクラスタを作成](../../chapter-2/README.md#creating-a-local-cluster-with-kubernetes-kind)できます。

1. **Tekton Pipelinesのインストール**

```
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.45.0/release.yaml
```

2. **Tekton Dashboard のインストール（オプション）**

```
kubectl apply -f https://github.com/tektoncd/dashboard/releases/download/v0.33.0/release.yaml
```

ダッシュボードには、`kubectl`を使用してポートフォワーディングすることでアクセスできます：

```
kubectl port-forward svc/tekton-dashboard  -n tekton-pipelines 9097:9097
```

![Tekton Dashboard](imgs/tekton-dashboard.png)

その後、ブラウザで[http://localhost:9097](http://localhost:9097)にアクセスできます。

3. **Tekton CLI のインストール（オプション）**：

[Tekton `tkn` CLIツール](https://github.com/tektoncd/cli)もインストールできます。
Mac OSXを使用している場合は、次のコマンドを実行できます：

```
brew install tektoncd-cli
```

## Tekton Tasksの基本

このセクションでは、TaskとシンプルなPipelineを作成する方法を説明し、その後、カンファレンスアプリケーションのアーティファクトを構築するために使用されるサービスパイプラインを見ていくことができるようになります。

Tektonでは、Tekton Taskの定義を作成することで、タスクの動作を定義できます。以下は、最もシンプルなタスクの例です：

```yaml
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: hello-world-task
spec:
  params:
  - name: name
    type: string
    description: who do you want to welcome?
    default: tekton user
  steps:
    - name: echo
      image: ubuntu
      command:
        - echo
      args:
        - "Hello World: $(params.name)" 
```

このTekton `Task`は、`ubuntu`イメージと、そのイメージ内にある`echo`コマンドを使用します。このタスクは`name`というパラメータも受け付け、メッセージの出力に使用します。このタスク定義をクラスタに適用するには、次のコマンドを実行します：

```
kubectl apply -f hello-world/hello-world-task.yaml
```

このリソースをKubernetesに適用する際、タスクを実行しているわけではなく、他の人が使用できるようにタスク定義を利用可能にしているだけです。このタスクは、複数のパイプラインで参照したり、異なるユーザーが独立して実行したりできるようになりました。

クラスタで利用可能なタスクを一覧表示するには、次のコマンドを実行します：

```
> kubectl get tasks
NAME               AGE
hello-world-task   88s
```

では、タスクを実行してみましょう。これは`TaskRun`リソースを作成することで行います。`TaskRun`は、タスクの単一の実行を表します。この具体的な実行には、固定のリソース名（`hello-world-task-run-1`）と、`name`というタスクパラメータの具体的な値があることに注意してください。

```yaml
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  name: hello-world-task-run-1
spec:
  params: 
  - name: name
    value: "Building Platforms on top of Kubernetes reader!"
  taskRef:
    name: hello-world-task
```

この`TaskRun`リソースをクラスタに適用して、最初のタスク実行（実行）を作成しましょう：

```
kubectl apply -f hello-world/task-run.yaml
taskrun.tekton.dev/hello-world-task-run-1 created
```

`TaskRun`が作成されるとすぐに、Tektonパイプラインエンジンがタスクのスケジューリングと、実行に必要なKubernetes Podの作成を担当します。デフォルトの名前空間でPodを一覧表示すると、次のような示が見られるはずです：

```
kubectl get pods
NAME                         READY   STATUS     RESTARTS   AGE
hello-world-task-run-1-pod   0/1     Init:0/1   0          2s
```

`TaskRun`を一覧表示してステータスを確認することもできます：

```
kubectl get taskrun
NAME                     SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-task-run-1   True        Succeeded   66s         7s
```

最後に、単一のタスクを実行していたため、作成されたPodのログをtailすることで、TaskRunの実行ログを見ることができます：

```
kubectl logs -f hello-world-task-run-1-pod 
Defaulted container "step-echo" out of: step-echo, prepare (init)
Hello World: Building Platforms on top of Kubernetes reader!
```

では、Tekton Pipelineを使用して複数のタスクをシーケンス化する方法を見ていきましょう。

## Tekton Pipelinesの基本

これで、先ほど定義したようなタスクを複数調整するためパイプラインを使用できます。また、[Tekton Hub](https://hub.tekton.dev/)からTektonコミュニティによって作成されたタスク定義を再利用することもできます。

![Tekton Hub](imgs/tekton-hub.png)

パイプラインを作成する前に、Tekton Hubから`wget`Tektonタスクをインストールします：

```
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/wget/0.1/wget.yaml
```

次のように表示されるはずです：

```
task.tekton.dev/wget created
```

では、先ほどインストールした`Hello World`タスクと`wget`タスクを一緒に使用して、シンプルなパイプラインを作成しましょう。

ファイルをフェッチし、その内容を読み取り、そして前に定義した`Hello World`タスクを使用するこのシンプルなパイプライン定義を作成します。

![Hello World Pipeline](imgs/hello-world-pipeline.png)

以下のようなパイプライン定義を作成ましょう：

```yaml
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: hello-world-pipeline
  annotations:
    description: |
      Fetch resource from internet, cat content and then say hello
spec:
  results: 
  - name: message
    type: string
    value: $(tasks.cat.results.messageFromFile)
  params:
  - name: url
    description: resource that we want to fetch
    type: string
    default: ""
  workspaces:
  - name: files
  tasks:
  - name: wget
    taskRef:
      name: wget
    params:
    - name: url
      value: "$(params.url)"
    - name: diroptions
      value:
        - "-P"  
    workspaces:
    - name: wget-workspace
      workspace: files
  - name: cat
    runAfter: [wget]
    workspaces:
    - name: wget-workspace
      workspace: files
    taskSpec: 
      workspaces:
      - name: wget-workspace
      results: 
        - name: messageFromFile
          description: the message obtained from the file
      steps:
      - name: cat
        image: bash:latest
        script: |
          #!/usr/bin/env bash
          cat $(workspaces.wget-workspace.path)/welcome.md | tee /tekton/results/messageFromFile
  - name: hello-world
    runAfter: [cat]
    taskRef:
      name: hello-world-task
    params:
      - name: name
        value: "$(tasks.cat.results.messageFromFile)"
```

ファイルをフェッチし、その内容を読み取り、そして前に定義した hello-world タスクを使用してファイルの内容を出力するのは、思ったほど簡単ではありませんでした。
パイプラインを使用すると、必要に応じて新しいタスクを追加して、各個別タスクの入力と出力のさらなる処理や変換を行う柔軟性があります。

この例では、Tekton Hubからインストールした`wget`タスク、`cat`と呼ばれるインラインで定義されたタスク（基本的にダウンロードしたファイルの内容を取得し、後で`hello-world-task`で参照できるTekton Resultに格納します）を使用しています。

このパイプライン定義をインストールするには、次のコマンドを実行します：

```
kubectl apply -f hello-world/hello-world-pipeline.yaml
```

そして、このパイプラインを実行したい時はいつでも、新しい`PipelineRun`を作成できます：

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: hello-world-pipeline-run-1
spec:
  workspaces:
    - name: files
      volumeClaimTemplate: 
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 1M 
  params:
  - name: url
    value: "https://raw.githubusercontent.com/salaboy/salaboy/main/welcome.md"
  pipelineRef:
    name: hello-world-pipeline
```

タスクがファイルシステムにファイルをダウンロードして保存する必要があるため、`PipelineRun`にストレージを提供するための抽象化としてTekton workspacesを使用しています。前に`TaskRun`で行ったように、`PipelineRun`にもパラメータを提供でき、各実行で異なる設定を使用したり、この場合は異なるファイルを使用したりできます。

`PipelineRun`と`TaskRun`の両方で、各実行に新しいリソース名を生成する必要があります。同じリソースを2回再適用しようとすると、Kubernetes APIサーバーは同じ名前の既存リソースの変更を許可しないからです。

このパイプラインを実行するには、次のコマンドを実行します：

```
kubectl apply -f hello-world/pipeline-run.yaml
```

作成されたポッドを確認してください：

```
> kubectl get pods
NAME                                         READY   STATUS        RESTARTS   AGE
affinity-assistant-ca1de9eb35-0              1/1     Terminating   0          19s
hello-world-pipeline-run-1-cat-pod           0/1     Completed     0          11s
hello-world-pipeline-run-1-hello-world-pod   0/1     Completed     0          5s
hello-world-pipeline-run-1-wget-pod          0/1     Completed     0          19s
```

タスクごとに1つのPodがあり、`affinity-assistant-ca1de9eb35-0`というPodがあることに注目してください。これは、Podが正しいノード（ボリュームがバインドされた場所）に作成されることを確認しています。

TaskRunも確認してください：

```
> kubectl get taskrun
NAME                                     SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-pipeline-run-1-cat           True        Succeeded   109s        104s
hello-world-pipeline-run-1-hello-world   True        Succeeded   103s        98s
hello-world-pipeline-run-1-wget          True        Succeeded   117s        109s
```

もちろん、すべてのタスクが成功すれば、PipelineRunも成功します：

```
kubectl get pipelinerun
NAME                         SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-pipeline-run-1   True        Succeeded   2m13s       114s
```

インストールした場合は、Tekton Dashboardでパイプラインとタスクの実行を確認してください。
![Tekton Dashboard](imgs/tekton-dashboard-hello-world-pipeline.png)

## サービスパイプラインのためのTekton

実際のサービスパイプラインは、前述の簡単な例よりもはるかに複雑です。これは主に、パイプラインタスクが外部システムにアクセスするための特別な設定と認証情報を必要とするためです。

サービスパイプライン定義の例は、このディレクトリの[service-pipeline.yaml](service-pipeline.yaml)ファイルにあります。

![Service Pipeline](imgs/service-pipeline.png)

この例のサービスパイプラインは、[`ko`]を使用してサービスのコンテナイメージをビルドおよび公開します。このパイプラインは、Goサービス専用のものです。異なるプログラミング言語でサービスを構築する場合は、他のツールを使用する必要があります。例のサービスパイプラインは、異なるサービスをビルドするようにパラメータ化できます。

このサービスパイプラインを実行するには、コンテナレジストリの認証情報を設定する必要があります。これは、パイプラインがDocker Hubなどのコンテナレジストリにコンテナをプッシュできるようにすることを意味します。Tekton Task/PipelineからコンテナレジストリでAuthenticate（認証）する方法については、[公式ドキュメント](https://tekton.dev/docs/how-to-guides/kaniko-build-push/#container-registry-authentication)を確認してください。

この例では、Docker Hubの認証情報を含むKubernetesシークレットを作成します：

```
kubectl create secret docker-registry docker-credentials --docker-server=https://index.docker.io/v1/ --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
```

次に、`Git Clone`と`ko` Tektonタスクをインストールします：

```
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/git-clone/0.9/git-clone.yaml
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/ko/0.1/ko.yaml
```

サービスパイプライン定義をクラスタにインストールしましょう：

```
kubectl apply -f service-pipeline.yaml
```

これで、サービスのコンテナイメージをビルドおよび公開するための新しいパイプラインインスタンスを作成できます。以下の`PipelineRun`リソースは、Notifications Serviceをビルドするようにサービスパイプラインを設定します。

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: service-pipeline-run-1
  annotations:
    kubernetes.io/ssh-auth: kubernetes.io/dockerconfigjson
spec:
  params:
  - name: target-registry
    value: docker.io/salaboy
  - name: target-service
    value: notifications-service
  - name: target-version 
    value: 1.0.0-from-pipeline-run
  workspaces:
    - name: sources
      volumeClaimTemplate: 
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 100Mi 
    - name: docker-credentials
      secret:  
        secretName: docker-credentials
  pipelineRef:
    name: service-pipeline
```

このサービスパイプラインの新しいインスタンスを作成するために、この`PipelineRun`定義をクラスタに適用します：

```
kubectl apply -f service-pipeline-run.yaml
```

`spec.params`セクションに注目してください。パイプラインが結果のコンテナイメージを自分のレジストリにプッシュするように、この部分を変更する必要があります。つまり、`docker.io/salaboy`を自分のレジストリ + ユーザー名に置き換えてください。`target-service`パラメータを使用して、カンファレンスアプリケーションのどのサービスをビルドするかを選択できます（利用可能なサービス：`notifications-service`、`agenda-service`、`c4p-service`、`frontend`）。

アプリケーションサービスをすべて含むHelmチャートをパッケージ化して公開する別のパイプライン（[app-helm-chart-pipeline.yaml](app-helm-chart-pipeline.yaml)）もあります。
チームがHelmチャート内にバンドルしたいサービスとバージョンの組み合わせを決定したら、別のパイプラインを実行してチャートをパッケージ化し、サービスのコンテナイメージが公開されているのと同じコンテナレジストリに公開できます。

![Helm Chart Application Pipeline](imgs/app-helm-pipeline.png)

アプリケーションHelmチャートパイプラインをインストールするには、次のコマンドを実行します：

```
kubectl apply -f app-helm-chart-pipeline.yaml
```

その後、新しい`PipelineRun`リソースを作成して新しいインスタンスを作成できます：

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: app-helm-chart-pipeline-run-1
  annotations:
    kubernetes.io/ssh-auth: kubernetes.io/dockerconfigjson
spec:
  params:
  - name: target-registry
    value: docker.io/salaboy
  - name: target-version
    value: v0.9.9
  workspaces:
    - name: sources
      volumeClaimTemplate: 
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 100Mi 
    - name: dockerconfig
      secret:  
        secretName: docker-credentials
  pipelineRef:
    name: app-helm-chart-pipeline
```

アプリケーションHelmチャートパイプラインの新しいインスタンスを作成するために、この`PipelineRun`定義をクラスタに適用します：

```
kubectl apply -f app-helm-chart-pipeline-run.yaml
```

アプリケーションHelmチャートパイプラインも、HelmチャートをOCIコンテナイメージとしてプッシュするために同じ`docker-credentials`を使用していることに注意してください。パイプラインは`target-version`パラータを受け取り、これを使用してHelmチャートをパッケージ化してOCIコンテナレジストリにプッシュする前に`Chart.yaml`ファイルにパッチを当てます。このパイプラインはチャートが参照するコンテナのバージョンをパッチしないことに注意してください。つまり、読者がパイプラインを適応させて各サービスのバージョンをパラメータとして受け取り、参照されるコンテナイメージが参照されるコンテナレジストリに存在することを検証することは読者次第です。

**注意**: これらのパイプラインは、コンテナとチャートを構築するためにTektonを設定するのに必要な作業を説明するための例に過ぎません。例えば、アプリケーションHelmチャートパイプラインは、チャートのバージョンやチャート内で参照されるコンテナイメージのバージョンを変更しません。プロセス全体を本当に自動化したい場合は、リリースしたいバージョンを表すGitタグからイメージバージョンとチャートバージョンを取得することができます。

## クリーンアップ

これらのチュートリアル用に作成したKinDクラスタを削除したい場合は、次のコマンドを実行できます：

```
kind delete clusters dev
```
