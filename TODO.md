# OGP検証サービス開発TODO

READMEの要件に基づいて整理したタスクリストです。

## ✅ 完了済みタスク

### 1. プロジェクト構造の作成
- [x] Go backend用ディレクトリ構成の作成
- [x] Frontend用ディレクトリ構成の作成  
- [x] Terraformインフラ設定用ディレクトリ構成の作成

### 2. Goバックエンドアプリケーションの開発
- [x] `go.mod`の作成とプロジェクト初期化
- [x] API エンドポイント `/api/v1/ogp/verify` の実装
- [x] JSON レスポンス形式の実装
- [x] HTTP サーバーの基本設定

### 3. OGP情報取得機能の実装
- [x] 指定URLからHTMLを取得する機能
- [x] HTMLからOGPタグを解析する機能
- [x] 対象タグ: `og:title`, `og:description`, `og:image`, `og:url`, `og:type`, `og:site_name`
- [x] その他標準的なOGPタグの取得

### 4. OGP検証機能の実装
- [x] OGPタグの存在チェック
- [x] 画像URLの有効性確認
- [x] 文字数制限チェック（プラットフォーム別）
- [x] 必須タグの不足警告

### 5. プラットフォーム別プレビュー機能
- [x] X/Twitter形式対応（タイトル70文字、説明200文字、画像1200x630px）
- [x] Facebook形式対応（タイトル100文字、説明300文字、画像1200x630px）
- [x] Discord形式対応（タイトル256文字、説明2048文字）
- [x] 各プラットフォームでの表示プレビュー生成

### 6. セキュリティ機能の実装
- [x] CORS設定の実装
- [x] レート制限（IP単位: 10req/min）
- [x] 不正URL検証
- [x] プライベートIPアドレスへのアクセス制限

### 7. React+Bunフロントエンドアプリケーションの開発
- [x] Bunプロジェクトの初期化とpackage.json作成
- [x] TypeScript設定（tsconfig.json）
- [x] Reactコンポーネントの基本構造
- [x] URL入力フォームコンポーネントの作成
- [x] OGP情報表示UIコンポーネントの作成
- [x] プラットフォーム別プレビュー表示コンポーネント
- [x] レスポンシブデザインの実装（Tailwind CSS推奨）
- [x] Bunでのビルド設定

### 8. Terraformインフラ設定
- [x] `main.tf`の作成（さくらのVPS + Cloudflare構成）
- [x] `variables.tf`の作成
- [x] `outputs.tf`の作成
- [x] `terraform.tfvars.example`の作成

### 9. GitHub Actions CI/CDワークフロー
- [x] バックエンドのビルド・テスト・デプロイ（Go）
- [x] フロントエンドのビルド・デプロイ（Bun + React）
- [x] Bunキャッシュ設定
- [x] Terraformの実行ワークフロー
- [x] 自動テスト実行
- [x] ワークフローの手動実行設定（自動実行は無効化済み）

### 10. Dockerでのローカル開発環境
- [x] Dockerfileの作成（Go backend用）
- [x] Dockerfileの作成（Bun + React frontend用）
- [x] docker-compose.ymlの作成
- [x] 開発環境用の設定ファイル

### 11. 基本テストの実装
- [x] バックエンドユニットテストの基本実装
- [x] フロントエンドコンポーネントテストの基本実装

### 12. 環境設定
- [x] .env.example ファイルの作成
- [x] 開発・本番環境用設定の準備

## 🔄 進行中・今後のタスク

### 🚨 高優先度：動作確認と品質保証
- [ ] **サービス動作確認**
  - [ ] Go バックエンドサービスの起動確認
  - [ ] React フロントエンドの起動確認
  - [ ] Docker Compose環境での統合動作確認
  - [ ] API エンドポイントの動作テスト
  - [ ] エラーハンドリングの動作確認
  - [ ] レート制限機能の動作確認
  - [ ] CORS設定の動作確認
  - [ ] プラットフォーム別プレビュー機能の確認

- [ ] **問題修正（実装ベース）**
  - [ ] 発見された問題の実装修正
  - [ ] パフォーマンス問題の改善
  - [ ] セキュリティ問題の修正
  - [ ] UI/UX問題の改善
  - [ ] ※ テスト仕様を変更して問題を回避することは禁止

- [ ] **品質保証**
  - [ ] 実際のWebサイトでのOGP検証テスト
  - [ ] 異常系テストの実施
  - [ ] 負荷テストの実施
  - [ ] セキュリティテストの実施

### テストの拡張
- [ ] 統合テストの実装
- [ ] E2Eテストの実装
- [ ] テストカバレッジ80%以上の達成

### API仕様書の作成
- [ ] OpenAPI 3.0形式でのAPI仕様書作成
- [ ] Swagger UIの設定

### ドキュメントの作成
- [ ] セットアップガイドの作成
- [ ] Terraformセットアップガイドの作成
- [ ] デプロイ手順書の作成
- [ ] 運用手順書の作成
- [ ] アーキテクチャ図の作成

### モニタリング・ヘルスチェック機能
- [ ] Cloudflare Analytics設定
- [ ] さくらのVPSリソース監視
- [ ] カスタムヘルスチェック実装
- [ ] アラート設定

## 📋 非機能要件

### パフォーマンス
- [x] タイムアウト10秒の設定
- [ ] レスポンス時間3秒以内の実現
- [ ] 同時リクエスト処理100req/secの対応

### 可用性
- [x] エラーハンドリングの実装
- [ ] 99.9%稼働率の実現
- [ ] ログ出力機能の実装

## ✅ 成果物チェックリスト

### ソースコード
- [x] Goバックエンドアプリケーション
- [x] React+Bunフロントエンドアプリケーション（TypeScript）
- [x] Terraformインフラ設定
- [x] GitHub Actions ワークフロー

### インフラ設定
- [x] main.tf（メインTerraform設定）
- [x] variables.tf（変数定義）
- [x] outputs.tf（出力値定義）
- [x] terraform.tfvars.example（設定例）
- [ ] インフラ構築手順書

### ドキュメント
- [x] README.md
- [ ] API仕様書
- [ ] Terraformセットアップガイド
- [ ] アーキテクチャ図

### テスト
- [x] 基本ユニットテスト
- [ ] 統合テスト
- [ ] E2Eテスト

## 🚀 デプロイ状況

- **開発環境**: Docker Compose設定完了
- **CI/CDパイプライン**: GitHub Actions設定完了（手動実行のみ）
- **インフラ**: Terraform設定完了（未デプロイ）
- **本番環境**: 準備完了（デプロイ待ち）

## 📝 次のステップ

### 🎯 最優先（即座に実施）
1. **サービス動作確認**: Go + React の実際の動作テスト
2. **問題発見と修正**: 実装レベルでの問題解決
3. **品質保証**: 実際のWebサイトでの検証

### 📈 中期目標
4. **API仕様書の作成**: OpenAPI/Swagger設定
5. **テストカバレッジの向上**: 統合・E2Eテスト追加
6. **ドキュメント整備**: セットアップ・運用ガイド作成
7. **本番デプロイ**: Terraformでのインフラ構築

### ⚠️ 重要な制約
- **テスト仕様変更による回避禁止**: 発見された問題は実装修正で解決すること
- **実装ベース解決**: モックやテスト仕様変更ではなく、実際のコード修正で対応
- **品質優先**: 動作しない機能は完成とみなさない

## 📚 参考情報
- [Open Graph Protocol仕様](https://ogp.me/)
- [X Cards Documentation](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)
- [Facebook Sharing](https://developers.facebook.com/docs/sharing/webmasters/)
- [Discord Embeds](https://discord.com/developers/docs/resources/channel#embed-object)
- [さくらのクラウド API](https://manual.sakura.ad.jp/cloud/api/)
- [Cloudflare API](https://developers.cloudflare.com/api/)
- [Terraform Sakura Cloud Provider](https://registry.terraform.io/providers/sacloud/sakuracloud/latest/docs)