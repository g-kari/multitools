# OGP検証サービス開発依頼書

## 概要
WebサービスのOGP（Open Graph Protocol）情報を検証し、主要なSNSプラットフォームでの表示プレビューを確認できるサービスの開発を依頼します。

## 推奨技術スタック（2025年版）

### アーキテクチャ概要
```
Cloudflare Pages (フロントエンド) → さくらのVPS (Golang API)
```

### フロントエンド
- **ホスティング**: Cloudflare Pages
- **技術**: HTML/CSS/JavaScript (React/Vue.js可)
- **主な機能**: 
  - URL入力フォーム
  - OGP情報の表示
  - プラットフォーム別プレビュー表示
  - レスポンシブデザイン

### バックエンド
- **言語**: Go (Golang)
- **実行環境**: さくらのVPS (512MB以上推奨)
- **主な機能**: 
  - 指定されたURLからOGP情報を取得・解析
  - OGPタグの検証とエラーチェック
  - レスポンス形式はJSON
  - CORS対応
  - 固定IPアドレス対応

## 機能要件

### 1. OGP情報取得機能
- URLを受け取り、HTMLを取得してOGPタグを解析
- 以下のOGPタグを取得：
  - `og:title`
  - `og:description`  
  - `og:image`
  - `og:url`
  - `og:type`
  - `og:site_name`
  - その他の標準的なOGPタグ

### 2. 検証機能
- OGPタグの存在チェック
- 画像URLの有効性確認
- 文字数制限チェック（プラットフォーム別）
- 必須タグの不足警告

### 3. プレビュー機能
対象プラットフォームでの表示プレビューを生成：

#### X (旧Twitter)
- タイトル: 最大70文字
- 説明: 最大200文字
- 画像: 1200x630px推奨
- Summary Card / Summary Card with Large Image形式

#### Facebook
- タイトル: 最大100文字
- 説明: 最大300文字
- 画像: 1200x630px推奨
- リンクプレビュー形式

#### Discord
- タイトル: 最大256文字
- 説明: 最大2048文字
- 画像: 制限なし（ただし表示最適化考慮）
- Embed形式

## API仕様

### エンドポイント
```
POST /api/v1/ogp/verify
```

### リクエスト
```json
{
  "url": "https://example.com/page"
}
```

### レスポンス
```json
{
  "url": "https://example.com/page",
  "ogp": {
    "title": "ページタイトル",
    "description": "ページの説明",
    "image": "https://example.com/image.jpg",
    "url": "https://example.com/page",
    "type": "website",
    "site_name": "サイト名"
  },
  "validation": {
    "errors": [],
    "warnings": ["タイトルがX用に長すぎます"]
  },
  "previews": {
    "twitter": {
      "card_type": "summary_large_image",
      "display_title": "切り詰められたタイトル",
      "display_description": "切り詰められた説明"
    },
    "facebook": {
      "display_title": "切り詰められたタイトル",
      "display_description": "切り詰められた説明"
    },
    "discord": {
      "display_title": "表示タイトル",
      "display_description": "表示説明"
    }
  }
}
```

## 非機能要件

### パフォーマンス
- レスポンス時間: 3秒以内
- 同時リクエスト処理: 100req/sec
- タイムアウト: 10秒

### セキュリティ
- CORS設定
- レート制限（IP単位: 10req/min）
- 不正URLの検証
- プライベートIPアドレスへのアクセス制限

### 可用性
- 99.9%の稼働率
- エラーハンドリング
- ログ出力機能

## インフラ要件

### さくらのVPS
- **プラン**: VPS 512MB（月額¥685）
- **OS**: Ubuntu 22.04 LTS
- **固定IP**: 標準で付属
- **SSL**: Let's Encrypt使用
- **リージョン**: 日本（東京・大阪）

### Cloudflare（無料プラン）
- **Pages**: 静的ホスティング
- **CDN**: グローバル配信
- **SSL**: 自動証明書
- **Analytics**: アクセス解析
- **セキュリティ**: DDoS保護

### Infrastructure as Code
- **Terraform**: インフラ構成管理
- **対象リソース**: 
  - さくらのVPS作成・設定
  - Cloudflare DNS/Pages設定
  - Nginx・SSL証明書設定
  - systemdサービス設定

## 開発・運用要件

### 開発環境
- Dockerでのローカル開発環境
- GitHub ActionsでのCI/CD
- テストカバレッジ80%以上

### Infrastructure as Code
- **Terraform**: 全インフラをコードで管理
- **プロバイダー**: 
  - Sakura Cloud Provider
  - Cloudflare Provider
- **リソース管理**: 
  - VPS作成・設定
  - DNS設定
  - SSL証明書管理
  - サービス設定

### モニタリング
- Cloudflare Analytics（無料）
- さくらのVPSリソース監視
- カスタムヘルスチェック
- アラート設定

### ドキュメント
- API仕様書（OpenAPI 3.0）
- Terraform設定ガイド
- デプロイ手順書
- 運用手順書

## 成果物

1. **ソースコード**
   - Goバックエンドアプリケーション
   - フロントエンドアプリケーション（HTML/CSS/JS）
   - Terraformインフラ設定
   - GitHub Actions ワークフロー

2. **Infrastructure as Code**
   - `main.tf`: メインのTerraform設定
   - `variables.tf`: 変数定義
   - `outputs.tf`: 出力値定義
   - `terraform.tfvars.example`: 設定例
   - インフラ構築手順書

3. **ドキュメント**
   - README.md
   - API仕様書
   - Terraformセットアップガイド
   - アーキテクチャ図

4. **テスト**
   - ユニットテスト
   - 統合テスト
   - E2Eテスト

## 工期・予算
- 開発期間: 3-4週間
- 月額運用費: 約¥685（約$4.50）+ ドメイン費用
- 予算: 要相談

## Terraformセットアップ例

### 必要な設定ファイル

**main.tf**
```hcl
terraform {
  required_providers {
    sakuracloud = {
      source  = "sacloud/sakuracloud"
      version = "~> 2.25"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

# さくらのVPS作成
resource "sakuracloud_server" "ogp_server" {
  name        = "ogp-verification-server"
  core        = 1
  memory      = 1
  disks       = [sakuracloud_disk.ogp_disk.id]
  description = "OGP検証サービス用サーバー"
  
  network_interface {
    upstream = "shared"
  }
  
  user_data = file("cloud-init.yaml")
  
  tags = ["ogp", "production"]
}

resource "sakuracloud_disk" "ogp_disk" {
  name              = "ogp-verification-disk"
  source_archive_id = data.sakuracloud_archive.ubuntu.id
  size              = 20
}

# Cloudflare DNS設定
resource "cloudflare_record" "ogp_api" {
  zone_id = var.cloudflare_zone_id
  name    = "api"
  value   = sakuracloud_server.ogp_server.ip_address
  type    = "A"
  ttl     = 300
}

# Cloudflare Pages設定
resource "cloudflare_pages_project" "ogp_frontend" {
  account_id        = var.cloudflare_account_id
  name              = "ogp-verification-frontend"
  production_branch = "main"
  
  source {
    type = "github"
    config {
      owner                         = var.github_owner
      repo_name                     = "ogp-verification-frontend"
      production_branch             = "main"
      pr_comments_enabled           = true
      deployments_enabled          = true
      production_deployment_enabled = true
      preview_deployment_setting    = "custom"
      preview_branch_includes       = ["dev", "staging"]
    }
  }
  
  build_config {
    build_command       = "npm run build"
    destination_dir     = "dist"
    root_dir           = ""
    web_analytics_tag  = var.cloudflare_web_analytics_tag
    web_analytics_token = var.cloudflare_web_analytics_token
  }
}
```

**variables.tf**
```hcl
variable "sakuracloud_token" {
  description = "さくらのクラウド APIトークン"
  type        = string
  sensitive   = true
}

variable "sakuracloud_secret" {
  description = "さくらのクラウド APIシークレット"
  type        = string
  sensitive   = true
}

variable "cloudflare_api_token" {
  description = "Cloudflare APIトークン"
  type        = string
  sensitive   = true
}

variable "cloudflare_zone_id" {
  description = "Cloudflare Zone ID"
  type        = string
}

variable "cloudflare_account_id" {
  description = "Cloudflare Account ID"
  type        = string
}

variable "domain_name" {
  description = "ドメイン名"
  type        = string
  default     = "example.com"
}
```

**outputs.tf**
```hcl
output "server_ip" {
  description = "サーバーのIPアドレス"
  value       = sakuracloud_server.ogp_server.ip_address
}

output "api_endpoint" {
  description = "API エンドポイント"
  value       = "https://api.${var.domain_name}"
}

output "frontend_url" {
  description = "フロントエンド URL"
  value       = cloudflare_pages_project.ogp_frontend.subdomain
}
```

## デプロイ手順

### 1. 事前準備
```bash
# 必要なツールのインストール
# Terraform
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install terraform

# 認証情報の設定
export SAKURACLOUD_ACCESS_TOKEN="your-token"
export SAKURACLOUD_ACCESS_TOKEN_SECRET="your-secret"
export CLOUDFLARE_API_TOKEN="your-api-token"
```

### 2. インフラ構築
```bash
# リポジトリクローン
git clone <repository-url>
cd ogp-verification-service

# Terraform初期化
cd terraform
terraform init

# 設定ファイルコピー
cp terraform.tfvars.example terraform.tfvars
# terraform.tfvarsを編集

# インフラ作成
terraform plan
terraform apply
```

### 3. アプリケーションデプロイ
```bash
# バックエンドデプロイ
cd ../backend
GOOS=linux GOARCH=amd64 go build -o ogp-service
scp ogp-service user@server-ip:/home/user/
ssh user@server-ip 'sudo systemctl restart ogp-service'

# フロントエンドデプロイ（GitHub連携で自動）
cd ../frontend
git push origin main  # Cloudflare Pagesが自動デプロイ
```

## 料金概算

### 月額運用費
- **さくらのVPS 512MB**: ¥685
- **Cloudflare Pages**: ¥0（無料）
- **ドメイン**: ¥100-500/月
- **合計**: 約¥785-1,185/月（$5-8）

### 初期費用
- **開発費**: 要相談
- **さくらのVPS初期費用**: ¥0
- **ドメイン取得**: ¥1,000-3,000/年

## アーキテクチャのメリット

### コスト効率
- 月額$10以下での運用
- 固定IPアドレス標準装備
- SSL証明書無料

### パフォーマンス
- Cloudflare CDNによる高速配信
- 日本国内データセンター（低レイテンシー）
- 静的フロントエンドによる高速表示

### 運用性
- Infrastructure as Codeによる再現可能な構築
- GitHub Actionsによる自動デプロイ
- 日本語サポート（さくらインターネット）

### セキュリティ
- Cloudflare DDoS保護
- 自動SSL証明書更新
- CORS設定による適切なAPI保護

## 参考情報
- [Open Graph Protocol仕様](https://ogp.me/)
- [X Cards Documentation](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)
- [Facebook Sharing](https://developers.facebook.com/docs/sharing/webmasters/)
- [Discord Embeds](https://discord.com/developers/docs/resources/channel#embed-object)
- [さくらのクラウド API](https://manual.sakura.ad.jp/cloud/api/)
- [Cloudflare API](https://developers.cloudflare.com/api/)
- [Terraform Sakura Cloud Provider](https://registry.terraform.io/providers/sacloud/sakuracloud/latest/docs)

## 連絡先
プロジェクト責任者: [担当者名]
Email: [メールアドレス]
Slack: [チャンネル名]

---

**追記**: 本提案では、Oracle Cloudを除外し、コスト効率と運用性を重視した「Cloudflare Pages + さくらのVPS」構成を推奨しています。Terraformによる Infrastructure as Code により、環境の再現性と運用効率を確保できます。
