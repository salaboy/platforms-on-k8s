# 第1章 :: Kubernetesの上に構築されるプラットフォーム（の台頭）

---
_🌍 翻訳_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [Português (Portuguese)](README-pt.md) | [Español](README-es.md) | [日本語](README-ja.md)
> **注意:** このコンテンツは素晴らしいクラウドネイティブコミュニティの [🌟 コントリビューター](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) によってもたらされています！

---

## チュートリアルの前提条件

本書にリンクされているステップバイステップのチュートリアルに従うには、以下のツールが必要です：
- [Docker](https://docs.docker.com/engine/install/), v24.0.2
- [kubectl](https://kubernetes.io/docs/tasks/tools/), Client v1.27.3
- [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/), v0.20.0
- [Helm](https://helm.sh/docs/intro/install/), v3.12.3

これらは、チュートリアルをテストする際に使用されるテクノロジーとバージョンです。

> [!Warning]
> Dockerの代わりに[Podman](https://podman.io/)などの他のテクノロジーを使用したい場合は、以下のコマンドを使用してrootfulコンテナ実行をオンにすることで可能です。
```shell
podman machine set --rootful
```

## カンファレンスアプリケーションのシナリオ

本書の各章で変更・使用するアプリケーションは、シンプルな "ウォーキングスケルトン" を表しています。つまり、仮定やツール、フレームワークをテストするのに十分な複雑さを持ちつつ、顧客が最終的に使用する製品ではありません。

"カンファレンスアプリケーション" のウォーキングスケルトンは、潜在的な _スピーカー_ が提案を提出し、カンファレンスの _主催者_ がそれを価するという単純なユースケースを実装しています。以下はアプリのホームページす：

![ホーム](imgs/homepage.png)

このアプリケーションがどのように一般的に使用されるかは以下をご確認ください：
1. **C4P:** 潜在的な _スピーカー_ は、アプリケーションの **Call for Proposals** (C4P) セクションから新しい提案を提出できます。
   ![提案](imgs/proposals.png)

2. **レビューと承認**: 提案が提出されると、カンファレンスの _主催者_ はアプリケーションの **バックオフィス** セクションを使用してそれらをレビュー（承認または拒否）できます。
   ![バックオフィス](imgs/backoffice.png)

3. **発表**: _主催者_ に承認された場合、その提案はカンファレンスの **アジェンダ** ページに自動的に公開されます。
   ![アジェンダ](imgs/agenda.png)

4. **スピーカーへの通知**: **バックフィス** の **通知** タブで、_スピーカー_ は自分に送信された全ての通知（メール）を確認できます。スピーカーはこのタブで承認と拒否の両方のメールを見ることができます。
   ![通知](imgs/notifications-backoffice.png)

### イベント駆動型アプリケーション

**アプリケーション内のすべてのアクションによって、新しいイベントが発行されます。** 例えば、以下のタイミングでイベントが発行されます：
- 新しい提案が提出されたとき
- 提案が承認または拒否されたとき
- 通知が送信されたとき

これらのイベントは送信され、フロントエンドアプリケーションによってキャプチャされます。幸いなことに、読者の皆さんは **バックオフィス** セクションの **イベント** タブにアクセスすることで、アプリ内でこれらの詳細を確認できます。

![イベント](imgs/events-backoffice.png)

## まとめと貢献

このチュートリアルを改善したいですか？ [issue](https://github.com/salaboy/platforms-on-k8s/issues/new) を作成するか、[Twitter](https://twitter.com/salaboy) でメッセージを送るか、[プルリクエスト](https://github.com/salaboy/platforms-on-k8s/compare) を送信してください。
