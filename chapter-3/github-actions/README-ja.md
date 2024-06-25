# GitHub Actions in Action

GitHub Actionsを使用すると、インフラストラクチャを実行することなく、サービスパイプラインを自動化できます。すべてのパイプライン/ワークフローは、GitHub.comのインフラストラクチャ上で実行されます。大規模になると（課金が始まると）、これは非常に高価になります。

Conference Applicationのサービスパイプラインは、以下の場所で定義されています。
- [Agenda Service GH Service Pipeline](../../.github/workflows/agenda-service-service-pipeline.yaml)
- [C4P Service GH Service Pipeline](../../.github/workflows/c4p-service-service-pipeline.yaml)
- [Notifications Service GH Service Pipeline](../../.github/workflows/notifications-service-service-pipeline.yaml)

これらのパイプラインは、サービスのソースコードが変更された場合にのみトリガーされるように、フィルターを使用しています（#https://github.com/dorny/paths-filter）。各サービスのコンテナをビルドおよび公開するために、これらのパイプラインは`ko-build`（https://github.com/ko-build/setup-ko）を使用して、Goアプリケーション用のマルチプラットフォームコンテナを生成します。

Docker Hubに公開するには、`docker/login-action@v2`アクションが使用され、2つの環境シークレット（`secrets.DOCKERHUB_USERNAME`、`secrets.DOCKERHUB_TOKEN`）の設定が必要です。これにより、パイプラインの実行でコンテナを私のDocker Hubアカウントにプッシュできます。
