# 第8章 :: チームに実験の機会を与える

---
_🌍 利用可能な言語_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md)

> **注意:** これは素晴らしいクラウドネイティブコミュニティの [🌟 コントリビューター](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) によってもたらされました！

---

これらのチュートリアルでは、Kubernetesクラスターに Knative Serving と Argo Rollouts をインストールして、カナリアリリース、A/Bテスト、ブルー/グリーンデプロイメントを実装します。ここで説明するリリース戦略は、チームがサービスの新バージョンをリリースする際により多くの制御を持てるようにすることを目的としています。ソフトウェアをリリースする際に異なる技術を適用することで、チームは新バージョンを制御された環境で実験してストすることができ、一度にすべてのライブトラフィックを新バージョンにプッシュする必要がありません。

- [Knative Servingを使用したリリース戦略](knative/README-ja.md)
- [Argo Rolloutsを使用したリリース戦略](argo-rollouts/README-ja.md)

## クリーンアップ

このチュートリアル用に作成したKinDクラスターを削除したい場合は、次のコマンドを実行できます：

```shell
kind delete clusters dev
```

## 次のステップ

- Function-as-a-Serviceプラットフォームの構築に興味がある場合は、[Knative Functions](https://knative.dev/docs/functions/)プロジェクトをチェックしてください。このイニシアチブは、関数開発者の作業を容易にするツールの開発に取り組んでいます。

- Argo Rolloutsを試した後の次のステップは、Argo CDからArgo Rolloutsへのフローを示すエンドツーエンドの例を作成することです。こには、Rollouts定義を含むリポジトリの作成が必要です。詳細については、[Argoプロジェクトのよくある質問セクション](https://argo-rollouts.readthedocs.io/en/latest/FAQ/)で統合に関する情報を確認してください。

- `AnalysisTemplates`と`AnalysisRuns`を使用したより複雑な例を実験してください。この機能は、チームがより自信を持って新バージョンをデプロイするのに役立ちます。

- 両プロジェクトは[Istio](https://istio.io/)のようなサービスメッシュと連携できるため、Istioがあなたにとって何ができるかを熟知してください。

## まとめと貢献

このチュートリアルを改善したいですか？issueを作成するか、[Twitter](https://twitter.com/salaboy)でメッセージを送るか、プルリクエストを送信してください。
