# OGP 検証サービス - 運用マニュアル

このマニュアルでは、OGP 検証サービスを本番環境で運用・保守するための包括的なガイダンスを提供します。

## 📋 概要

### サービス アーキテクチャ
- **バックエンド**: Go API サーバー（ポート 8080）
- **フロントエンド**: Nginx で配信される React SPA
- **プロキシ**: Nginx（SSL 終端、負荷分散）
- **インフラストラクチャ**: Sakura VPS + Cloudflare CDN

### 主要コンポーネント
- API エンドポイント: `/api/v1/ogp/verify`
- ヘルスチェック: `/health`
- レート制限: IP あたり 10 リクエスト/分
- SSL: Let's Encrypt 証明書

## 🔧 日常運用

### サービス ヘルス監視

#### 自動ヘルスチェック
```bash
#!/bin/bash
# /opt/scripts/health_check.sh

API_URL="https://api.yourdomain.com"
FRONTEND_URL="https://yourdomain.com"
LOG_FILE="/var/log/ogp-service/health.log"

# ログディレクトリの作成
mkdir -p /var/log/ogp-service

# タイムスタンプ付きログ記録関数
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> $LOG_FILE
}

# API ヘルスチェック
if curl -f -s "$API_URL/health" > /dev/null; then
    log "API health check: OK"
else
    log "API health check: FAILED"
    # アラート送信
    echo "API health check failed" | mail -s "OGP Service Alert" admin@yourdomain.com
fi

# フロントエンドチェック
if curl -f -s "$FRONTEND_URL" > /dev/null; then
    log "Frontend health check: OK"
else
    log "Frontend health check: FAILED"
    echo "Frontend health check failed" | mail -s "OGP Service Alert" admin@yourdomain.com
fi

# API 機能チェック
TEST_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/ogp/verify" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}')

if echo "$TEST_RESPONSE" | grep -q '"is_valid":true'; then
    log "API functionality check: OK"
else
    log "API functionality check: FAILED"
    echo "API functionality test failed" | mail -s "OGP Service Alert" admin@yourdomain.com
fi
```

#### 手動ヘルス検証
```bash
# クイック ヘルスチェック
curl https://api.yourdomain.com/health

# 詳細 API テスト
curl -X POST https://api.yourdomain.com/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}' \
  | jq '.validation.is_valid'

# レスポンス時間の確認
time curl -s https://api.yourdomain.com/health

# SSL 証明書の確認
openssl s_client -connect api.yourdomain.com:443 -servername api.yourdomain.com < /dev/null 2>/dev/null | openssl x509 -noout -dates
```

### ログ管理

#### ログの確認
```bash
# アプリケーション ログ
sudo docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs -f

# Nginx アクセス ログ
sudo tail -f /var/log/nginx/access.log

# Nginx エラー ログ
sudo tail -f /var/log/nginx/error.log

# システム ログ
sudo journalctl -u docker -f

# SSL 証明書更新ログ
sudo tail -f /var/log/letsencrypt/letsencrypt.log
```

#### ログローテーション設定
```bash
# アプリケーション ログ用 logrotate 設定
sudo cat > /etc/logrotate.d/ogp-service << EOF
/var/log/ogp-service/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
    postrotate
        /usr/bin/docker-compose -f /opt/ogp-service/docker-compose.prod.yml restart nginx
    endscript
}
EOF

# logrotate のテスト
sudo logrotate -d /etc/logrotate.d/ogp-service
```

## 📊 パフォーマンス監視

### リソース監視

#### システム リソース
```bash
#!/bin/bash
# /opt/scripts/system_monitor.sh

# CPU 使用率
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
echo "CPU Usage: ${CPU_USAGE}%"

# メモリ使用率
MEMORY_USAGE=$(free | grep Mem | awk '{printf("%.1f", $3/$2 * 100.0)}')
echo "Memory Usage: ${MEMORY_USAGE}%"

# ディスク使用率
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | cut -d'%' -f1)
echo "Disk Usage: ${DISK_USAGE}%"

# 負荷平均
LOAD_AVG=$(uptime | awk -F'load average:' '{print $2}')
echo "Load Average:${LOAD_AVG}"

# しきい値を超えた場合のアラート
if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
    echo "High CPU usage: ${CPU_USAGE}%" | mail -s "System Alert" admin@yourdomain.com
fi

if (( $(echo "$MEMORY_USAGE > 85" | bc -l) )); then
    echo "High memory usage: ${MEMORY_USAGE}%" | mail -s "System Alert" admin@yourdomain.com
fi

if (( DISK_USAGE > 90 )); then
    echo "High disk usage: ${DISK_USAGE}%" | mail -s "System Alert" admin@yourdomain.com
fi
```

#### Docker コンテナー監視
```bash
# コンテナー リソース使用量
docker stats --no-stream

# コンテナー ヘルス状態
docker-compose -f /opt/ogp-service/docker-compose.prod.yml ps

# 詳細コンテナー検査
docker inspect ogp-service_backend_1 | jq '.State.Health'
```

#### アプリケーション パフォーマンス指標
```bash
#!/bin/bash
# /opt/scripts/performance_monitor.sh

API_URL="https://api.yourdomain.com"

# レスポンス時間の測定
RESPONSE_TIME=$(curl -o /dev/null -s -w '%{time_total}' "$API_URL/health")
echo "API Response Time: ${RESPONSE_TIME}s"

# レート制限のテスト
echo "Testing rate limiting..."
for i in {1..15}; do
    STATUS=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API_URL/api/v1/ogp/verify" \
      -H "Content-Type: application/json" \
      -d '{"url":"https://example.com"}')
    echo "Request $i: HTTP $STATUS"
    sleep 1
done

# レスポンス時間が遅い場合のアラート
if (( $(echo "$RESPONSE_TIME > 5.0" | bc -l) )); then
    echo "Slow API response: ${RESPONSE_TIME}s" | mail -s "Performance Alert" admin@yourdomain.com
fi
```

## 🔄 メンテナンス手順

### 定期メンテナンス タスク

#### 週次タスク
```bash
#!/bin/bash
# /opt/scripts/weekly_maintenance.sh

echo "週次メンテナンスを開始しています..."

# システム パッケージの更新
sudo apt update && sudo apt upgrade -y

# Docker のクリーンアップ
docker system prune -f
docker image prune -f

# 必要に応じてログの手動ローテーション
sudo logrotate -f /etc/logrotate.d/ogp-service

# SSL 証明書の期限確認
openssl x509 -in /etc/letsencrypt/live/yourdomain.com/cert.pem -noout -dates

# 設定のバックアップ
tar -czf /backup/ogp-config-$(date +%Y%m%d).tar.gz /opt/ogp-service/

# ディスク容量の確認
df -h

echo "週次メンテナンスが完了しました。"
```

#### 月次タスク
```bash
#!/bin/bash
# /opt/scripts/monthly_maintenance.sh

echo "月次メンテナンスを開始しています..."

# 完全システム バックアップ
rsync -av /opt/ogp-service/ /backup/monthly/ogp-service-$(date +%Y%m)/

# セキュリティ更新
sudo unattended-upgrades

# エラー ログのレビュー
grep -i error /var/log/nginx/error.log | tail -20
sudo docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs --since="30d" | grep -i error

# パフォーマンス分析
# 月次レポートの生成
echo "Monthly Performance Report - $(date)" > /tmp/monthly_report.txt
echo "==================================" >> /tmp/monthly_report.txt
echo "" >> /tmp/monthly_report.txt

# 月間平均レスポンス時間
grep "health" /var/log/nginx/access.log | awk '{print $NF}' | awk '{sum+=$1; count++} END {print "Average response time: " sum/count "s"}' >> /tmp/monthly_report.txt

# エラー率
TOTAL_REQUESTS=$(grep -c "api/v1/ogp/verify" /var/log/nginx/access.log)
ERROR_REQUESTS=$(grep "api/v1/ogp/verify" /var/log/nginx/access.log | grep -c " 5[0-9][0-9] ")
ERROR_RATE=$(echo "scale=2; $ERROR_REQUESTS * 100 / $TOTAL_REQUESTS" | bc)
echo "Error rate: ${ERROR_RATE}%" >> /tmp/monthly_report.txt

# レポート送信
mail -s "Monthly Performance Report" admin@yourdomain.com < /tmp/monthly_report.txt

echo "月次メンテナンスが完了しました。"
```

### アプリケーション更新

#### バックエンド更新プロセス
```bash
#!/bin/bash
# /opt/scripts/update_backend.sh

cd /opt/ogp-service

echo "バックエンド更新を開始しています..."

# 最新コードの取得
git pull origin main

# 現在の状態をバックアップ
docker-compose -f docker-compose.prod.yml stop backend
docker tag ogp-service_backend:latest ogp-service_backend:backup-$(date +%Y%m%d)

# 新しいイメージのビルド
docker-compose -f docker-compose.prod.yml build backend

# 新しいイメージで開始
docker-compose -f docker-compose.prod.yml up -d backend

# ヘルスチェックの待機
sleep 30

# ヘルスの確認
if curl -f https://api.yourdomain.com/health; then
    echo "バックエンド更新が成功しました"
    # 古いバックアップのクリーンアップ
    docker rmi ogp-service_backend:backup-$(date --date="7 days ago" +%Y%m%d) 2>/dev/null || true
else
    echo "バックエンド更新が失敗しました。ロールバックしています..."
    docker-compose -f docker-compose.prod.yml stop backend
    docker tag ogp-service_backend:backup-$(date +%Y%m%d) ogp-service_backend:latest
    docker-compose -f docker-compose.prod.yml up -d backend
    exit 1
fi
```

#### フロントエンド更新プロセス
```bash
#!/bin/bash
# /opt/scripts/update_frontend.sh

cd /opt/ogp-service

echo "フロントエンド更新を開始しています..."

# 最新コードの取得
git pull origin main

# 新しいフロントエンドのビルド
docker-compose -f docker-compose.prod.yml build frontend

# ローリング更新
docker-compose -f docker-compose.prod.yml up -d frontend

echo "フロントエンド更新が完了しました。"
```

## 🚨 インシデント対応

### アラート カテゴリ

#### クリティカル アラート（即座の対応）
- API の完全停止（ヘルスチェック失敗）
- SSL 証明書の期限切れ
- サーバーが応答しない
- ディスク容量 >95%

#### 警告アラート（1時間以内の対応）
- 高いレスポンス時間（>5秒）
- エラー率 >10%
- メモリ使用率 >85%
- レート制限が機能しない

#### 情報アラート（24時間以内の対応）
- SSL 証明書が7日以内に期限切れ
- 高い CPU 使用率（>80%）
- ログ エラー

### インシデント対応手順

#### API 停止インシデント
```bash
# 1. サービス状態の確認
docker-compose -f /opt/ogp-service/docker-compose.prod.yml ps

# 2. エラー ログの確認
docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs --tail=50

# 3. バックエンド サービスの再起動
docker-compose -f /opt/ogp-service/docker-compose.prod.yml restart backend

# 4. 再起動が失敗した場合、システム リソースを確認
free -h
df -h
docker system df

# 5. 緊急再起動
sudo systemctl restart docker
docker-compose -f /opt/ogp-service/docker-compose.prod.yml up -d

# 6. 復旧の確認
curl https://api.yourdomain.com/health
```

#### 高エラー率インシデント
```bash
# 1. ログの最近のエラーを確認
tail -100 /var/log/nginx/error.log

# 2. アプリケーション エラーの確認
docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs backend --tail=100 | grep -i error

# 3. トラフィック パターンの分析
tail -100 /var/log/nginx/access.log | awk '{print $1}' | sort | uniq -c | sort -nr | head -10

# 4. DDoS や悪用の確認
fail2ban-client status nginx-req-limit

# 5. 必要に応じて一時的にレート制限を強化
# nginx 設定を編集してリロード
sudo nginx -s reload
```

#### SSL 証明書期限切れ
```bash
# 1. 証明書状態の確認
sudo certbot certificates

# 2. 更新の試行
sudo certbot renew --force-renewal

# 3. 更新が失敗した場合、DNS を確認
dig yourdomain.com
dig api.yourdomain.com

# 4. 手動証明書インストール（緊急時）
sudo certbot certonly --manual -d yourdomain.com -d api.yourdomain.com

# 5. nginx の再起動
sudo systemctl restart nginx
```

### 復旧手順

#### データベース復旧（該当する場合）
```bash
# 1. アプリケーションの停止
docker-compose -f /opt/ogp-service/docker-compose.prod.yml stop

# 2. バックアップからの復旧
sudo tar -xzf /backup/database-backup-YYYYMMDD.tar.gz -C /

# 3. アプリケーションの開始
docker-compose -f /opt/ogp-service/docker-compose.prod.yml up -d

# 4. データ整合性の確認
# アプリケーション固有のチェックを実行
```

#### 完全システム復旧
```bash
# 1. レスキューメディアからの起動（必要に応じて）
# 2. ファイルシステムのマウント
# 3. バックアップからの復旧
sudo tar -xzf /backup/full-system-backup.tar.gz -C /

# 4. 必要に応じて Docker の再インストール
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 5. サービスの開始
cd /opt/ogp-service
docker-compose -f docker-compose.prod.yml up -d

# 6. すべてのサービスの確認
./scripts/health_check.sh
```

## 📈 容量計画

### トラフィック分析
```bash
#!/bin/bash
# /opt/scripts/traffic_analysis.sh

# 日次リクエスト数
TODAY=$(date +%Y-%m-%d)
REQUESTS_TODAY=$(grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | wc -l)
echo "Requests today: $REQUESTS_TODAY"

# ピーク時間分析
grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | awk '{print $4}' | cut -d: -f2 | sort | uniq -c | sort -nr | head -5

# レスポンス時間分析
grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | awk '{print $NF}' | sort -n | awk '{all[NR] = $0} END{print "Min: " all[1] "s, Max: " all[NR] "s, Median: " all[int(NR/2)] "s"}'

# エラー率
TOTAL=$(grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | wc -l)
ERRORS=$(grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | grep -c " 5[0-9][0-9] ")
ERROR_RATE=$(echo "scale=2; $ERRORS * 100 / $TOTAL" | bc)
echo "Error rate: ${ERROR_RATE}%"
```

### スケーリング推奨事項

#### 水平スケーリング指標
- CPU 使用率が継続的に >70%
- メモリ使用率が継続的に >80%
- レスポンス時間 >3秒
- リクエスト量 >1000/時間

#### 垂直スケーリング手順
```bash
# 1. 現在の使用状況を7日間監視
# 2. アップグレードが必要な場合、メンテナンス ウィンドウをスケジュール
# 3. サーバー スナップショットの作成（サポートされている場合）
# 4. Sakura Cloud コンソール経由でサーバー プランをアップグレード
# 5. アップグレード後にサービスを再起動
# 6. 24時間パフォーマンスを監視
```

## 🔒 セキュリティ運用

### セキュリティ監視
```bash
#!/bin/bash
# /opt/scripts/security_monitor.sh

# ログイン失敗試行の確認
sudo grep "Failed password" /var/log/auth.log | tail -10

# fail2ban 状態の確認
sudo fail2ban-client status
sudo fail2ban-client status sshd

# 異常なトラフィック パターンの確認
tail -1000 /var/log/nginx/access.log | awk '{print $1}' | sort | uniq -c | sort -nr | head -20

# 既知の攻撃パターンの確認
grep -i "union\|select\|drop\|script\|alert" /var/log/nginx/access.log | tail -10

# SSL セキュリティの確認
nmap --script ssl-enum-ciphers -p 443 yourdomain.com
```

### セキュリティ更新
```bash
#!/bin/bash
# /opt/scripts/security_updates.sh

# セキュリティ更新の確認
sudo apt list --upgradable | grep -i security

# セキュリティ更新の適用
sudo unattended-upgrades

# Docker イメージの更新
docker-compose -f /opt/ogp-service/docker-compose.prod.yml pull
docker-compose -f /opt/ogp-service/docker-compose.prod.yml up -d

# fail2ban ルールの更新
sudo fail2ban-client reload

# 実行中サービスの CVE 確認
# `docker scout cves` などのツールが利用可能な場合は使用
```

## 📋 ランブック チェックリスト

### 日次チェックリスト
- [ ] 監視によるサービス ヘルスの確認
- [ ] エラー ログの確認
- [ ] SSL 証明書状態の確認
- [ ] ディスク容量使用量の確認
- [ ] トラフィック パターンの確認

### 週次チェックリスト
- [ ] システム更新の適用
- [ ] Docker イメージのクリーンアップ
- [ ] パフォーマンス指標の確認
- [ ] バックアップ復旧のテスト
- [ ] 必要に応じてドキュメントの更新

### 月次チェックリスト
- [ ] 完全セキュリティ スキャン
- [ ] 容量計画のレビュー
- [ ] 災害復旧計画の更新
- [ ] 監視しきい値の見直しと更新
- [ ] パフォーマンス最適化のレビュー

## 📞 連絡先情報

### エスカレーション マトリックス
1. **オンコール エンジニア**: +81-XXX-XXXX-XXXX
2. **テクニカル リード**: tech-lead@yourdomain.com
3. **インフラストラクチャ チーム**: infra@yourdomain.com
4. **緊急連絡先**: emergency@yourdomain.com

### ベンダー連絡先
- **Sakura Cloud サポート**: +81-3-5332-7071
- **Cloudflare サポート**: Enterprise ダッシュボード
- **ドメイン レジストラ**: Web ポータル経由で連絡

## 📚 追加リソース

### ドキュメント リンク
- [セットアップ ガイド](SETUP.md)
- [デプロイメント ガイド](DEPLOYMENT.md)
- [Terraform セットアップ](TERRAFORM_SETUP.md)
- [API ドキュメント](../backend/api/README.md)

### 外部リソース
- [Docker ドキュメント](https://docs.docker.com/)
- [Nginx ドキュメント](https://nginx.org/en/docs/)
- [Let's Encrypt ドキュメント](https://letsencrypt.org/docs/)
- [Sakura Cloud ドキュメント](https://manual.sakura.ad.jp/cloud/)

---

**重要**: この運用マニュアルは、手順、連絡先情報、またはシステム アーキテクチャの変更に応じて最新の状態に保ってください。