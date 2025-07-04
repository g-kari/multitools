# Terraform インフラストラクチャ セットアップ ガイド

このガイドでは、Sakura VPS と Cloudflare を使用して Terraform で OGP 検証サービスを本番環境にデプロイする方法を説明します。

## 📋 前提条件

### 必要なアカウント・サービス
- **Sakura VPS** アカウント（API アクセス有効）
- **Cloudflare** アカウント（ドメイン管理機能付き）
- **GitHub** リポジトリ（CI/CD 用）
- **ドメイン名**（Cloudflare で管理）

### 必要なソフトウェア
- **Terraform** (v1.0+)
- **Git**（バージョン管理用）
- **SSH クライアント**（サーバーアクセス用）

### 必要な認証情報
- Sakura Cloud API キーとシークレット
- Cloudflare API トークン
- サーバーアクセス用 SSH キーペア

## 🚀 クイックスタート

### 1. 認証情報の設定

```bash
# Sakura Cloud の認証情報
export SAKURACLOUD_ACCESS_TOKEN="your-access-token"
export SAKURACLOUD_ACCESS_TOKEN_SECRET="your-access-secret"
export SAKURACLOUD_ZONE="is1b"  # または好みのゾーン

# Cloudflare の認証情報
export CLOUDFLARE_API_TOKEN="your-cloudflare-token"

# サーバーアクセス用 SSH キー
ssh-keygen -t rsa -b 4096 -f ~/.ssh/ogp-service-key
```

### 2. Terraform の初期化

```bash
cd terraform

# Terraform の初期化
terraform init

# 例となる変数をコピー
cp terraform.tfvars.example terraform.tfvars

# 設定の編集
nano terraform.tfvars
```

### 3. 変数の設定

`terraform.tfvars` を編集：

```hcl
# ドメインと DNS
domain_name = "yourdomain.com"
subdomain   = "ogp-api"  # ogp-api.yourdomain.com が作成されます

# サーバー設定
server_name = "ogp-service-prod"
server_plan = "1core-1gb"  # 高負荷時は "2core-2gb"

# SSH アクセス
ssh_key_name = "ogp-service-key"
ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2E..."  # あなたの公開鍵

# アプリケーション設定
app_environment = "production"
cors_origins = "https://yourdomain.com"
rate_limit = "20"  # 1分間あたりのリクエスト数

# Cloudflare 設定
cloudflare_zone_id = "your-zone-id"

# オプション：データベース設定（外部DBを使用する場合）
# database_host = "your-db-host"
# database_name = "ogp_service"
```

### 4. 計画と適用

```bash
# 計画された変更を確認
terraform plan

# インフラストラクチャを適用
terraform apply

# 出力を確認
terraform output
```

## 🛠️ 詳細設定

### Sakura VPS 設定

Terraform 設定により作成されるもの：

#### サーバー仕様
- **OS**: Ubuntu 22.04 LTS
- **メモリ**: 512MB (1core-1gb) または 2GB (2core-2gb)
- **ストレージ**: 20GB SSD
- **ネットワーク**: パブリック IP とファイアウォール ルール

#### セキュリティ設定
```hcl
# ファイアウォール ルール
resource "sakuracloud_simple_monitor" "ogp_service" {
  target = sakuracloud_server.ogp_service.ip_address
  
  health_check {
    protocol    = "http"
    port        = 80
    path        = "/health"
    status      = 200
    timeout     = 10
  }
}
```

#### 自動セットアップ スクリプト
サーバーでは以下を実行する cloud-init スクリプトが動作します：
1. システムパッケージの更新
2. Docker と Docker Compose のインストール
3. ファイアウォール (UFW) の設定
4. SSL 証明書の設定 (Let's Encrypt)
5. OGP サービスのデプロイ
6. 監視の設定

### Cloudflare 設定

#### DNS レコード
- **A レコード**: `ogp-api.yourdomain.com` → サーバー IP
- **CNAME レコード**: `www.ogp-api.yourdomain.com` → `ogp-api.yourdomain.com`

#### セキュリティ設定
- **SSL/TLS**: Full (strict)
- **セキュリティ レベル**: Medium
- **Bot Fight Mode**: 有効
- **レート制限**: 100 リクエスト/分/IP

```hcl
resource "cloudflare_zone_settings_override" "ogp_service" {
  zone_id = var.cloudflare_zone_id
  
  settings {
    ssl = "strict"
    security_level = "medium"
    bot_fight_mode = "on"
    
    security_header {
      enabled = true
    }
  }
}
```

## 🔧 高度な設定

### 変数リファレンス

#### 必須変数
| 変数 | 説明 | 例 |
|------|------|-----|
| `domain_name` | ドメイン名 | `"example.com"` |
| `subdomain` | API サブドメイン | `"ogp-api"` |
| `ssh_public_key` | SSH 公開鍵 | `"ssh-rsa AAAAB3..."` |
| `cloudflare_zone_id` | Cloudflare ゾーン ID | `"abcd1234..."` |

#### オプション変数
| 変数 | 説明 | デフォルト | 例 |
|------|------|-----------|-----|
| `server_plan` | Sakura VPS プラン | `"1core-1gb"` | `"2core-2gb"` |
| `server_name` | サーバー ホスト名 | `"ogp-service"` | `"ogp-prod"` |
| `app_environment` | 環境名 | `"production"` | `"staging"` |
| `cors_origins` | 許可する CORS オリジン | `"*"` | `"https://example.com"` |
| `rate_limit` | 1分間あたりのリクエスト数 | `"10"` | `"20"` |

### サーバー プラン

#### 1core-1gb（開発用推奨）
- **CPU**: 1 vCPU
- **メモリ**: 1GB
- **ストレージ**: 20GB SSD
- **コスト**: 約 ¥680/月

#### 2core-2gb（本番用推奨）
- **CPU**: 2 vCPU
- **メモリ**: 2GB
- **ストレージ**: 20GB SSD
- **コスト**: 約 ¥1,580/月

#### 4core-4gb（高負荷シナリオ）
- **CPU**: 4 vCPU
- **メモリ**: 4GB
- **ストレージ**: 20GB SSD
- **コスト**: 約 ¥3,200/月

### カスタム サーバー設定

カスタム設定用に `terraform/modules/server/user-data.sh` を作成：

```bash
#!/bin/bash
set -e

# カスタム アプリケーション設定
echo "OGP 検証サービスをセットアップしています..."

# 追加パッケージのインストール
apt-get update
apt-get install -y htop nginx-utils

# カスタム監視の設定
cat > /opt/ogp-service/monitoring.sh << 'EOF'
#!/bin/bash
# カスタム監視スクリプト
curl -f http://localhost:8080/health || systemctl restart ogp-service
EOF

chmod +x /opt/ogp-service/monitoring.sh

# ヘルス監視のために crontab に追加
echo "*/5 * * * * /opt/ogp-service/monitoring.sh" | crontab -

echo "カスタム設定が完了しました！"
```

## 🚀 デプロイメント プロセス

### 1. デプロイメント前チェック

```bash
# Terraform 設定の検証
terraform validate

# フォーマットのチェック
terraform fmt -check

# セキュリティスキャン（オプション）
terraform plan -out=tfplan
terraform show -json tfplan | jq > tfplan.json
# tfsec や checkov などのツールを使用してセキュリティ分析
```

### 2. ステージング デプロイメント

```bash
# ステージング ワークスペースの作成
terraform workspace new staging
terraform workspace select staging

# ステージングへのデプロイ
terraform apply -var="app_environment=staging" -var="subdomain=ogp-api-staging"
```

### 3. 本番デプロイメント

```bash
# 本番ワークスペースに切り替え
terraform workspace new production
terraform workspace select production

# 本番へのデプロイ
terraform apply
```

### 4. デプロイメント後の検証

```bash
# サーバー情報の取得
terraform output server_ip
terraform output server_fqdn

# デプロイメントのテスト
FQDN=$(terraform output -raw server_fqdn)
curl https://${FQDN}/health
curl -X POST https://${FQDN}/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'
```

## 🔒 セキュリティ ベストプラクティス

### SSH アクセス
```bash
# サーバーへの接続
SERVER_IP=$(terraform output -raw server_ip)
ssh -i ~/.ssh/ogp-service-key ubuntu@${SERVER_IP}

# パスワード認証の無効化（自動的に実行）
# fail2ban の設定（cloud-init に含まれる）
# 自動セキュリティ更新の設定（含まれる）
```

### ファイアウォール ルール
```hcl
# 必要なポートのみ許可
resource "sakuracloud_simple_monitor" "ogp_service" {
  # HTTP (80) - HTTPS にリダイレクト
  # HTTPS (443) - メインサービス
  # SSH (22) - 管理アクセスのみ
  # その他のポートはデフォルトでブロック
}
```

### SSL/TLS 設定
- **証明書**: Let's Encrypt（自動更新）
- **TLS バージョン**: 1.2+ のみ
- **暗号化**: 強力な暗号のみ
- **HSTS**: 有効（max-age 1年）

### 監視とアラート
```bash
# 内蔵監視チェック：
# - サービス ヘルス（/health エンドポイント）
# - SSL 証明書の期限
# - ディスク使用量
# - メモリ使用量
# - 負荷平均
```

## 📊 監視とメンテナンス

### ヘルス チェック

Cloudflare は以下でサービスを監視します：
- **エンドポイント**: `https://ogp-api.yourdomain.com/health`
- **頻度**: 60秒ごと
- **タイムアウト**: 10秒
- **期待されるステータス**: 200 OK

### ログ管理

```bash
# サーバーにSSH接続
ssh -i ~/.ssh/ogp-service-key ubuntu@${SERVER_IP}

# アプリケーション ログの表示
sudo docker-compose -f /opt/ogp-service/docker-compose.yml logs -f

# システム ログの表示
sudo journalctl -fu ogp-service

# アクセス ログの表示
sudo tail -f /var/log/nginx/access.log
```

### バックアップ戦略

```bash
# 自動日次バックアップ（cloud-init で設定）
# - アプリケーション設定：/opt/ogp-service/
# - Docker イメージとデータ
# - システム設定
# - SSL 証明書

# 手動バックアップ
sudo tar -czf /tmp/ogp-backup-$(date +%Y%m%d).tar.gz \
  /opt/ogp-service/ \
  /etc/letsencrypt/ \
  /etc/nginx/sites-available/
```

### 更新とメンテナンス

```bash
# サーバーにSSH接続
ssh -i ~/.ssh/ogp-service-key ubuntu@${SERVER_IP}

# アプリケーション更新（ゼロダウンタイム）
cd /opt/ogp-service
sudo docker-compose pull
sudo docker-compose up -d

# システム パッケージの更新
sudo apt update && sudo apt upgrade -y
sudo reboot  # カーネル更新の場合
```

## 🐛 トラブルシューティング

### 一般的な問題

#### 1. Terraform Apply の失敗

```bash
# 認証情報の確認
echo $SAKURACLOUD_ACCESS_TOKEN
echo $CLOUDFLARE_API_TOKEN

# ゾーン ID の確認
curl -X GET "https://api.cloudflare.com/client/v4/zones" \
  -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
  -H "Content-Type: application/json"

# terraform 状態の確認
terraform show
terraform refresh
```

#### 2. サーバーにアクセスできない

```bash
# Sakura Cloud コンソールでサーバー状態を確認
# セキュリティ グループ ルールの確認
# サーバー上で cloud-init ログを確認：
ssh ubuntu@${SERVER_IP} sudo cat /var/log/cloud-init-output.log
```

#### 3. SSL 証明書の問題

```bash
# サーバーにSSH接続して証明書を確認
ssh ubuntu@${SERVER_IP}
sudo certbot certificates
sudo nginx -t
sudo systemctl status nginx
```

#### 4. サービスが起動しない

```bash
# Docker 状態の確認
sudo docker ps -a
sudo docker-compose logs

# システム リソースの確認
free -h
df -h
sudo systemctl status ogp-service
```

### デバッグ コマンド

```bash
# Terraform デバッグ
export TF_LOG=DEBUG
terraform apply

# サーバー診断
curl -I https://ogp-api.yourdomain.com
dig ogp-api.yourdomain.com
nslookup ogp-api.yourdomain.com

# ポート接続性
telnet ogp-api.yourdomain.com 443
nc -zv ogp-api.yourdomain.com 80 443
```

## 🔄 CI/CD 統合

### GitHub Actions

`.github/workflows/terraform.yml` に追加：

```yaml
name: Terraform Deploy

on:
  push:
    branches: [main]
    paths: ['terraform/**']

jobs:
  terraform:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Terraform
      uses: hashicorp/setup-terraform@v2
      with:
        terraform_version: 1.5.0
    
    - name: Terraform Init
      run: terraform init
      working-directory: terraform
      env:
        SAKURACLOUD_ACCESS_TOKEN: ${{ secrets.SAKURACLOUD_ACCESS_TOKEN }}
        SAKURACLOUD_ACCESS_TOKEN_SECRET: ${{ secrets.SAKURACLOUD_ACCESS_TOKEN_SECRET }}
        CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
    
    - name: Terraform Plan
      run: terraform plan
      working-directory: terraform
    
    - name: Terraform Apply
      if: github.ref == 'refs/heads/main'
      run: terraform apply -auto-approve
      working-directory: terraform
```

### 必要なシークレット

これらを GitHub リポジトリのシークレットに追加：
- `SAKURACLOUD_ACCESS_TOKEN`
- `SAKURACLOUD_ACCESS_TOKEN_SECRET`
- `CLOUDFLARE_API_TOKEN`

## 🧹 クリーンアップ

### インフラストラクチャの削除

```bash
# すべてのリソースを破棄
terraform destroy

# 削除の確認
terraform state list  # 空であることを確認

# ローカル ファイルのクリーンアップ
rm -rf .terraform/
rm terraform.tfstate*
```

### コスト最適化

```bash
# サーバー停止（データは保持、コンピューティングの課金停止）
terraform apply -var="server_power_state=down"

# サーバー サイズの縮小
terraform apply -var="server_plan=1core-1gb"

# 監視の削除（オプションのコスト削減）
terraform apply -var="enable_monitoring=false"
```

## 📚 追加リソース

### Terraform ドキュメント
- [Sakura Cloud Provider](https://registry.terraform.io/providers/sacloud/sakuracloud/latest/docs)
- [Cloudflare Provider](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs)

### Sakura Cloud ドキュメント
- [API リファレンス](https://manual.sakura.ad.jp/cloud/api/)
- [サーバー プラン](https://cloud.sakura.ad.jp/specification/server-disk/)

### セキュリティ リソース
- [Let's Encrypt](https://letsencrypt.org/)
- [UFW ドキュメント](https://help.ubuntu.com/community/UFW)
- [fail2ban](https://www.fail2ban.org/)

---

**次のステップ**: インフラストラクチャのデプロイメント後は、アプリケーションのデプロイメントについて [DEPLOYMENT.md](DEPLOYMENT.md) を、継続的なメンテナンスについて [OPERATIONS.md](OPERATIONS.md) を参照してください。