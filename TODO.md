# OGP検証サービス開発TODO

READMEの要件に基づいて整理したタスクリストです。

## 高優先度タスク

### 1. プロジェクト構造の作成
- [ ] Go backend用ディレクトリ構成の作成
- [ ] Frontend用ディレクトリ構成の作成  
- [ ] Terraformインフラ設定用ディレクトリ構成の作成

### 2. Goバックエンドアプリケーションの開発
- [ ] `go.mod`の作成とプロジェクト初期化
- [ ] API エンドポイント `/api/v1/ogp/verify` の実装
- [ ] JSON レスポンス形式の実装
- [ ] HTTP サーバーの基本設定

### 3. OGP情報取得機能の実装
- [ ] 指定URLからHTMLを取得する機能
- [ ] HTMLからOGPタグを解析する機能
- [ ] 対象タグ: `og:title`, `og:description`, `og:image`, `og:url`, `og:type`, `og:site_name`
- [ ] その他標準的なOGPタグの取得

### 4. OGP検証機能の実装
- [ ] OGPタグの存在チェック
- [ ] 画像URLの有効性確認
- [ ] 文字数制限チェック（プラットフォーム別）
- [ ] 必須タグの不足警告

### 5. プラットフォーム別プレビュー機能
- [ ] X/Twitter形式対応（タイトル70文字、説明200文字、画像1200x630px）
- [ ] Facebook形式対応（タイトル100文字、説明300文字、画像1200x630px）
- [ ] Discord形式対応（タイトル256文字、説明2048文字）
- [ ] 各プラットフォームでの表示プレビュー生成

## 中優先度タスク

### 6. セキュリティ機能の実装
- [ ] CORS設定の実装
- [ ] レート制限（IP単位: 10req/min）
- [ ] 不正URL検証
- [ ] プライベートIPアドレスへのアクセス制限

### 7. React+Bunフロントエンドアプリケーションの開発
- [ ] Bunプロジェクトの初期化とpackage.json作成
- [ ] TypeScript設定（tsconfig.json）
- [ ] Reactコンポーネントの基本構造
- [ ] URL入力フォームコンポーネントの作成
- [ ] OGP情報表示UIコンポーネントの作成
- [ ] プラットフォーム別プレビュー表示コンポーネント
- [ ] レスポンシブデザインの実装（Tailwind CSS推奨）
- [ ] Bunでのビルド設定

### 8. Terraformインフラ設定
- [ ] `main.tf`の作成（さくらのVPS + Cloudflare構成）
- [ ] `variables.tf`の作成
- [ ] `outputs.tf`の作成
- [ ] `terraform.tfvars.example`の作成

### 9. GitHub Actions CI/CDワークフロー
- [ ] バックエンドのビルド・テスト・デプロイ（Go）
- [ ] フロントエンドのビルド・デプロイ（Bun + React）
- [ ] Bunキャッシュ設定
- [ ] Terraformの実行ワークフロー
- [ ] 自動テスト実行

### 10. Dockerでのローカル開発環境
- [ ] Dockerfileの作成（Go backend用）
- [ ] Dockerfileの作成（Bun + React frontend用）
- [ ] docker-compose.ymlの作成
- [ ] 開発環境用の設定ファイル

### 11. テストの実装
- [ ] ユニットテストの実装
- [ ] 統合テストの実装
- [ ] E2Eテストの実装
- [ ] テストカバレッジ80%以上の達成

## 低優先度タスク

### 12. API仕様書の作成
- [ ] OpenAPI 3.0形式でのAPI仕様書作成
- [ ] Swagger UIの設定

### 13. ドキュメントの作成
- [ ] セットアップガイドの作成
- [ ] Terraformセットアップガイドの作成
- [ ] デプロイ手順書の作成
- [ ] 運用手順書の作成
- [ ] アーキテクチャ図の作成

### 14. モニタリング・ヘルスチェック機能
- [ ] Cloudflare Analytics設定
- [ ] さくらのVPSリソース監視
- [ ] カスタムヘルスチェック実装
- [ ] アラート設定

## 非機能要件

### パフォーマンス
- [ ] レスポンス時間3秒以内の実現
- [ ] 同時リクエスト処理100req/secの対応
- [ ] タイムアウト10秒の設定

### 可用性
- [ ] 99.9%稼働率の実現
- [ ] エラーハンドリングの実装
- [ ] ログ出力機能の実装

## 成果物チェックリスト

### ソースコード
- [ ] Goバックエンドアプリケーション
- [ ] React+Bunフロントエンドアプリケーション（TypeScript）
- [ ] Terraformインフラ設定
- [ ] GitHub Actions ワークフロー

### インフラ設定
- [ ] main.tf（メインTerraform設定）
- [ ] variables.tf（変数定義）
- [ ] outputs.tf（出力値定義）
- [ ] terraform.tfvars.example（設定例）
- [ ] インフラ構築手順書

### ドキュメント
- [ ] README.md
- [ ] API仕様書
- [ ] Terraformセットアップガイド
- [ ] アーキテクチャ図

### テスト
- [ ] ユニットテスト
- [ ] 統合テスト
- [ ] E2Eテスト

## 参考情報
- [Open Graph Protocol仕様](https://ogp.me/)
- [X Cards Documentation](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)
- [Facebook Sharing](https://developers.facebook.com/docs/sharing/webmasters/)
- [Discord Embeds](https://discord.com/developers/docs/resources/channel#embed-object)
- [さくらのクラウド API](https://manual.sakura.ad.jp/cloud/api/)
- [Cloudflare API](https://developers.cloudflare.com/api/)
- [Terraform Sakura Cloud Provider](https://registry.terraform.io/providers/sacloud/sakuracloud/latest/docs)