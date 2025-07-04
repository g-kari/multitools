# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリでコードを操作する際のガイダンスを提供します。

## プロジェクト概要

これは **OGP (Open Graph Protocol) 検証サービス** プロジェクトです。ウェブサイトの OGP メタデータを分析し、Twitter/X、Facebook、Discord 向けのプラットフォーム固有のプレビューと検証結果を提供します。

## 技術スタック

### バックエンド
- **言語**: Go (Golang)
- **フレームワーク**: JSON API 付き標準 Go HTTP サーバー
- **デプロイ**: Sakura VPS (512MB) with Ubuntu 22.04 LTS

### フロントエンド
- **フレームワーク**: React with TypeScript
- **ランタイム**: Bun (パッケージマネージャー兼ビルドツール)
- **スタイリング**: Tailwind CSS (推奨)
- **デプロイ**: Cloudflare Pages

### インフラストラクチャ
- **IaC**: Terraform
- **CI/CD**: GitHub Actions
- **監視**: Cloudflare Analytics + カスタムヘルスチェック

## 開発コマンド

### バックエンド (Go)
```bash
# Go モジュールを初期化
go mod init ogp-verification-service

# 開発サーバーを実行
go run main.go

# 本番用にビルド
go build -o ogp-service

# テストを実行
go test ./...

# カバレッジ付きでテストを実行
go test -cover ./...
```

### フロントエンド (React + Bun)
```bash
# Bun プロジェクトを初期化
bun create react-app frontend --template typescript

# 依存関係をインストール
bun install

# 開発サーバー
bun dev

# 本番用にビルド
bun run build

# テストを実行
bun test

# 型チェック
bun run type-check
```

### インフラストラクチャ (Terraform)
```bash
# Terraform を初期化
terraform init

# インフラ変更を計画
terraform plan

# インフラ変更を適用
terraform apply

# インフラを削除
terraform destroy
```

### Docker 開発
```bash
# 開発環境を起動
docker-compose up -d

# サービスをビルドして起動
docker-compose up --build

# サービスを停止
docker-compose down
```

## プロジェクト構造

```
├── backend/              # Go バックエンドアプリケーション
│   ├── cmd/             # アプリケーションエントリーポイント
│   ├── internal/        # プライベートアプリケーションコード
│   │   ├── handlers/    # HTTP ハンドラー
│   │   ├── models/      # データモデル
│   │   ├── services/    # ビジネスロジック
│   │   └── validators/  # 入力検証
│   ├── pkg/            # パブリックライブラリコード
│   └── go.mod          # Go モジュール定義
├── frontend/            # React + Bun フロントエンド
│   ├── src/
│   │   ├── components/  # React コンポーネント
│   │   ├── hooks/       # カスタム React フック
│   │   ├── services/    # API サービス
│   │   └── types/       # TypeScript 型定義
│   ├── package.json
│   └── bun.lockb
├── terraform/           # Infrastructure as Code
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── docker-compose.yml   # ローカル開発環境
└── .github/workflows/   # GitHub Actions CI/CD
```

## API 仕様

### メインエンドポイント
- **POST** `/api/v1/ogp/verify`
- **リクエスト**: `{"url": "https://example.com"}`
- **レスポンス**: OGP データ、検証結果、プラットフォームプレビューを含む JSON

### プラットフォームサポート
- **Twitter/X**: タイトル (70文字)、説明 (200文字)、画像 (1200x630px)
- **Facebook**: タイトル (100文字)、説明 (300文字)、画像 (1200x630px)
- **Discord**: タイトル (256文字)、説明 (2048文字)、画像 (柔軟)

## セキュリティ要件

- フロントエンドドメイン用の CORS 設定
- レート制限: IP あたり 10 リクエスト/分
- プライベート IP アドレスのブロック
- 入力検証とサニタイゼーション
- ログやレスポンスに機密データを含めない

## パフォーマンス要件

- レスポンス時間: < 3 秒
- 同時リクエスト: 100 req/sec
- リクエストタイムアウト: 10 秒
- テストカバレッジ: 80%+

## 開発ワークフロー

1. **TODO.md でのタスク管理**: すべての作業状況、コンテンツ、進捗は TODO.md で管理する必要があります
   - **必須**: 作業を開始する前に必ず TODO.md を更新
   - **必須**: 作業開始時にタスクを進行中としてマーク
   - **必須**: 完了時にタスクを完了としてマーク
   - **必須**: 発見された新しいタスクを追加
   - **必須**: TODO.md でブロッカーや問題を文書化
   - **必須**: TODO.md をプロジェクトの状況の信頼できる情報源として使用
   - **絶対必要**: どんな小さな作業でも TODO.md に記載せずに作業することは禁止

2. **完了した作業を常にコミット**: タスクの完了や大幅な進捗時に、説明的なメッセージでコミット
3. **機能ブランチを使用**: 新機能や大きな変更には専用ブランチを作成
4. **テストを書く**: 新機能には単体テストを実装
5. **変更を文書化**: 変更時は関連ドキュメントを更新
6. **デプロイ前にテスト**: プッシュ前にすべてのテストと型チェックを実行

## コミットガイドライン

- 慣例的なコミットメッセージを使用
- タスク完了時に常にコミット
- AI 生成コミットには 🤖 絵文字を含める
- 例: `feat: implement OGP validation service 🤖`

## 環境変数

### バックエンド
- `PORT`: サーバーポート (デフォルト: 8080)
- `CORS_ORIGINS`: 許可された CORS オリジン
- `RATE_LIMIT`: IP あたりの分間リクエスト数

### フロントエンド
- `REACT_APP_API_URL`: バックエンド API URL
- `REACT_APP_ENV`: 環境 (development/production)

## テスト

### バックエンドテスト
```bash
# すべてのテストを実行
go test ./...

# カバレッジ付きで実行
go test -cover ./...

# 特定のテストを実行
go test ./internal/handlers -v
```

### フロントエンドテスト
```bash
# すべてのテストを実行
bun test

# ウォッチモードで実行
bun test --watch

# カバレッジ付きで実行
bun test --coverage
```

## デプロイ

### 本番デプロイ
1. バックエンド: Go バイナリをビルドして Sakura VPS にデプロイ
2. フロントエンド: GitHub にプッシュ (Cloudflare Pages に自動デプロイ)
3. インフラストラクチャ: Terraform 変更を適用

### 開発環境
ホットリロードを有効にしたローカル開発には Docker Compose を使用。

### テストと検証
- **Docker 専用テスト**: すべてのサービステスト、検証、開発作業は Docker を使用する必要があります
- **ローカルインストール禁止**: Go、Bun、その他のランタイム依存関係をホストシステムに直接インストールしないでください
- **コンテナ化された検証**: すべてのサービス検証とテストに `docker-compose up -d` を使用
- **クリーンな環境**: このアプローチによりローカル環境の汚染を防止し、一貫性を保証

## 監視

- **ヘルスチェック**: `/health` エンドポイント
- **メトリクス**: レスポンス時間、エラー率、リクエスト数
- **ログ**: 構造化 JSON ログ
- **アラート**: 高エラー率やダウンタイムに対する設定

## 開発制限

### 禁止される行為
- **モックサーバーを作成しない**: このプロジェクトは実際の Go と React 実装のみを使用
- **Python/Flask 代替を作成しない**: 既存の Go バックエンドサービスを使用
- **公式技術スタックを迂回しない**: Go + React + TypeScript を使用

### 必要なアプローチ
- **公式実装を使用**: 実際の Go バックエンドと React フロントエンドで常に作業
- **実際のサービスでテスト**: テストと開発には本番対応コードを使用
- **技術の一貫性を保持**: 確立された Go + React + Bun + Terraform スタックに従う
- **必須の TODO.md 更新**: すべての作業は TODO.md でリアルタイムに追跡・更新する必要があります
- **Docker ファーストテスト**: すべてのサービステストと検証は Docker を使用してローカル環境の汚染を回避する必要があります

## アーキテクチャ注記

これは以下の特徴を持つ分散システムです：
- **ステートレスバックエンド**: 水平スケーリング対応
- **静的フロントエンド**: CDN 配信用
- **Terraform IaC**: 再現可能なインフラストラクチャ
- **GitHub Actions**: 自動化された CI/CD

このサービスは、パフォーマンス、セキュリティ、信頼性を重視し、ソーシャルメディアプラットフォーム向けの OGP 検証とプレビュー生成に焦点を当てています。