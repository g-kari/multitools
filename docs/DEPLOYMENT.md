# OGP 検証サービス - デプロイメントガイド

このガイドでは、さまざまなデプロイメント戦略を使用して OGP 検証サービスを本番環境にデプロイする方法について説明します。

## 📋 概要

### デプロイメント オプション
1. **Docker Compose**（単一サーバー推奨）
2. **Kubernetes**（コンテナーオーケストレーション用）
3. **Sakura VPS + Cloudflare**（コスト効率的な本番環境）
4. **クラウドプラットフォーム**（AWS、GCP、Azure）

### アーキテクチャ
- **バックエンド**: Go API サーバー（Docker コンテナー）
- **フロントエンド**: React SPA（Nginx で配信される静的ファイル）
- **プロキシ**: Nginx（SSL 終端、リバースプロキシ）
- **DNS**: Cloudflare（CDN、DDoS 保護）

## 🚀 本番デプロイメント（Docker Compose）

### 前提条件
- Ubuntu 22.04 LTS サーバー
- Docker と Docker Compose がインストール済み
- DNS アクセス可能なドメイン名
- SSL 証明書（Let's Encrypt 推奨）

### 1. サーバーの準備

```bash
# システムの更新
sudo apt update && sudo apt upgrade -y

# Docker のインストール
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Docker Compose のインストール
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Nginx のインストール
sudo apt install -y nginx certbot python3-certbot-nginx

# ファイアウォールの設定
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 2. アプリケーションのデプロイメント

```bash
# リポジトリのクローン
git clone <repository-url> /opt/ogp-service
cd /opt/ogp-service

# 本番環境ファイルの作成
cat > .env.production << EOF
# バックエンド設定
PORT=8080
CORS_ORIGINS=https://yourdomain.com
RATE_LIMIT=20

# フロントエンド設定
VITE_API_URL=https://api.yourdomain.com
VITE_ENV=production

# データベース設定（必要に応じて）
# DATABASE_URL=postgresql://user:pass@localhost/ogp_service
EOF

# 本番用 Docker Compose ファイルの作成
cat > docker-compose.prod.yml << EOF
services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    restart: unless-stopped
    environment:
      - PORT=8080
      - CORS_ORIGINS=https://yourdomain.com
      - RATE_LIMIT=20
    expose:
      - "8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.prod
      args:
        - VITE_API_URL=https://api.yourdomain.com
    restart: unless-stopped
    expose:
      - "80"
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  nginx:
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/prod.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
      - /etc/letsencrypt:/etc/letsencrypt:ro
    depends_on:
      backend:
        condition: service_healthy
      frontend:
        condition: service_started
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  default:
    name: ogp-production
EOF
```

### 3. フロントエンドの本番ビルド

`frontend/Dockerfile.prod` を作成：

```dockerfile
# ビルドステージ
FROM node:18-alpine AS build

WORKDIR /app

# パッケージファイルのコピー
COPY package*.json ./
RUN npm ci --only=production

# ソースのコピーとビルド
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# 本番ステージ
FROM nginx:alpine

# ビルドされたファイルのコピー
COPY --from=build /app/dist /usr/share/nginx/html

# nginx 設定のコピー
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 4. Nginx 設定

`nginx/prod.conf` を作成：

```nginx
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # ログ設定
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;

    # 基本設定
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    server_tokens off;

    # Gzip 圧縮
    gzip on;
    gzip_vary on;
    gzip_min_length 10240;
    gzip_proxied expired no-cache no-store private must-revalidate;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/x-javascript
        application/xml+rss
        application/javascript
        application/json;

    # レート制限
    limit_req_zone $binary_remote_addr zone=api:10m rate=30r/m;
    limit_req_zone $binary_remote_addr zone=static:10m rate=60r/m;

    # API サーバー
    upstream backend {
        server backend:8080;
        keepalive 32;
    }

    # フロントエンド サーバー
    upstream frontend {
        server frontend:80;
        keepalive 32;
    }

    # HTTP から HTTPS へのリダイレクト
    server {
        listen 80;
        server_name yourdomain.com www.yourdomain.com api.yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    # メインウェブサイト（フロントエンド）
    server {
        listen 443 ssl http2;
        server_name yourdomain.com www.yourdomain.com;

        # SSL 設定
        ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
        ssl_session_timeout 1d;
        ssl_session_cache shared:SSL:50m;
        ssl_session_tickets off;

        # モダンな設定
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # HSTS
        add_header Strict-Transport-Security "max-age=63072000" always;

        # セキュリティ ヘッダー
        add_header X-Content-Type-Options nosniff;
        add_header X-Frame-Options DENY;
        add_header X-XSS-Protection "1; mode=block";
        add_header Referrer-Policy "strict-origin-when-cross-origin";

        # レート制限
        limit_req zone=static burst=20 nodelay;

        location / {
            proxy_pass http://frontend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # ヘルスチェック
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }

    # API サーバー
    server {
        listen 443 ssl http2;
        server_name api.yourdomain.com;

        # SSL 設定（上記と同じ）
        ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
        ssl_session_timeout 1d;
        ssl_session_cache shared:SSL:50m;
        ssl_session_tickets off;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # セキュリティ ヘッダー
        add_header Strict-Transport-Security "max-age=63072000" always;
        add_header X-Content-Type-Options nosniff;
        add_header X-Frame-Options DENY;
        add_header X-XSS-Protection "1; mode=block";

        # CORS ヘッダー
        add_header Access-Control-Allow-Origin "https://yourdomain.com" always;
        add_header Access-Control-Allow-Methods "GET, POST, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Content-Type, Authorization" always;

        # API 用レート制限
        limit_req zone=api burst=10 nodelay;

        location / {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # タイムアウト設定
            proxy_connect_timeout 30s;
            proxy_send_timeout 30s;
            proxy_read_timeout 30s;
        }

        # ヘルスチェック
        location /health {
            proxy_pass http://backend/health;
            access_log off;
        }
    }
}
```

### 5. SSL 証明書の設定

```bash
# ドメイン用証明書のインストール
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com -d api.yourdomain.com

# 自動更新のテスト
sudo certbot renew --dry-run

# 更新を crontab に追加
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

### 6. アプリケーションのデプロイ

```bash
# サービスのビルドと開始
docker-compose -f docker-compose.prod.yml up --build -d

# ステータスの確認
docker-compose -f docker-compose.prod.yml ps

# ログの表示
docker-compose -f docker-compose.prod.yml logs -f
```

### 7. 検証

```bash
# API のテスト
curl https://api.yourdomain.com/health
curl -X POST https://api.yourdomain.com/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'

# フロントエンドのテスト
curl -I https://yourdomain.com

# SSL の確認
curl -I https://yourdomain.com 2>&1 | grep -i ssl
```

## ☸️ Kubernetes デプロイメント

### 前提条件
- Kubernetes クラスター（v1.24+）
- kubectl が設定済み
- Ingress コントローラー（nginx-ingress）
- SSL 用の cert-manager

### 1. 名前空間と ConfigMap

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: ogp-service

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ogp-config
  namespace: ogp-service
data:
  PORT: "8080"
  CORS_ORIGINS: "https://yourdomain.com"
  RATE_LIMIT: "20"
  VITE_API_URL: "https://api.yourdomain.com"
```

### 2. バックエンドのデプロイメント

```yaml
# k8s/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ogp-backend
  namespace: ogp-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ogp-backend
  template:
    metadata:
      labels:
        app: ogp-backend
    spec:
      containers:
      - name: backend
        image: ogp-backend:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: ogp-config
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"

---
apiVersion: v1
kind: Service
metadata:
  name: ogp-backend-service
  namespace: ogp-service
spec:
  selector:
    app: ogp-backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

### 3. フロントエンドのデプロイメント

```yaml
# k8s/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ogp-frontend
  namespace: ogp-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ogp-frontend
  template:
    metadata:
      labels:
        app: ogp-frontend
    spec:
      containers:
      - name: frontend
        image: ogp-frontend:latest
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"

---
apiVersion: v1
kind: Service
metadata:
  name: ogp-frontend-service
  namespace: ogp-service
spec:
  selector:
    app: ogp-frontend
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
```

### 4. Ingress 設定

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ogp-ingress
  namespace: ogp-service
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "30"
    nginx.ingress.kubernetes.io/rate-limit-burst: "10"
spec:
  tls:
  - hosts:
    - yourdomain.com
    - api.yourdomain.com
    secretName: ogp-tls
  rules:
  - host: yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: ogp-frontend-service
            port:
              number: 80
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: ogp-backend-service
            port:
              number: 8080
```

### 5. Kubernetes へのデプロイ

```bash
# すべての設定を適用
kubectl apply -f k8s/

# デプロイメント状態の確認
kubectl get pods -n ogp-service
kubectl get services -n ogp-service
kubectl get ingress -n ogp-service

# ログの表示
kubectl logs -f deployment/ogp-backend -n ogp-service
kubectl logs -f deployment/ogp-frontend -n ogp-service
```

## 🔄 CI/CD パイプライン

### GitHub Actions ワークフロー

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Run Backend Tests
      run: |
        cd backend
        go test ./... -v -cover
    
    - name: Run E2E Tests
      run: |
        cd backend/tests/e2e
        ./run-e2e.sh

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Build Backend Image
      run: |
        docker build -t ${{ secrets.REGISTRY_URL }}/ogp-backend:${{ github.sha }} backend/
        docker push ${{ secrets.REGISTRY_URL }}/ogp-backend:${{ github.sha }}
    
    - name: Build Frontend Image
      run: |
        docker build -t ${{ secrets.REGISTRY_URL }}/ogp-frontend:${{ github.sha }} \
          --build-arg VITE_API_URL=${{ secrets.API_URL }} frontend/
        docker push ${{ secrets.REGISTRY_URL }}/ogp-frontend:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - name: Deploy to Server
      run: |
        ssh ${{ secrets.DEPLOY_USER }}@${{ secrets.DEPLOY_HOST }} "
          cd /opt/ogp-service
          export BACKEND_IMAGE=${{ secrets.REGISTRY_URL }}/ogp-backend:${{ github.sha }}
          export FRONTEND_IMAGE=${{ secrets.REGISTRY_URL }}/ogp-frontend:${{ github.sha }}
          docker-compose -f docker-compose.prod.yml pull
          docker-compose -f docker-compose.prod.yml up -d
          docker system prune -f
        "
```

## 📊 監視とログ記録

### アプリケーション監視

```yaml
# docker-compose.monitoring.yml
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana

volumes:
  grafana-storage:
```

### ログ集約

```bash
# システムログに journald を使用
sudo journalctl -u docker -f

# docker logs を使用
docker-compose logs -f --tail=100

# ELK スタックによる集中ログ記録
docker run -d --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  elasticsearch:7.14.0
```

## 🔒 セキュリティ強化

### サーバーセキュリティ

```bash
# root ログインの無効化
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# fail2ban の設定
sudo apt install -y fail2ban
sudo systemctl enable fail2ban

# 自動更新の設定
sudo apt install -y unattended-upgrades
echo 'Unattended-Upgrade::Automatic-Reboot "true";' | sudo tee -a /etc/apt/apt.conf.d/50unattended-upgrades
```

### コンテナーセキュリティ

```dockerfile
# コンテナー内で非root ユーザーを使用
FROM golang:1.21-alpine AS builder
RUN adduser -D -s /bin/sh appuser

FROM alpine:latest
RUN adduser -D -s /bin/sh appuser
USER appuser
COPY --from=builder --chown=appuser:appuser /app/ogp-service .
```

### ネットワークセキュリティ

```bash
# ファイアウォールの設定
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# 管理者用 VPN アクセスの設定（オプション）
sudo apt install -y wireguard
```

## 🚨 災害復旧

### バックアップ戦略

```bash
#!/bin/bash
# backup.sh - 日次バックアップ スクリプト

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/ogp-service"

# バックアップ ディレクトリの作成
mkdir -p $BACKUP_DIR

# アプリケーション設定のバックアップ
tar -czf $BACKUP_DIR/config-$DATE.tar.gz /opt/ogp-service/

# SSL 証明書のバックアップ
tar -czf $BACKUP_DIR/ssl-$DATE.tar.gz /etc/letsencrypt/

# クラウドストレージへのアップロード（オプション）
# aws s3 cp $BACKUP_DIR/ s3://your-backup-bucket/ --recursive

# 古いバックアップのクリーンアップ（最新7日分を保持）
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### 復旧手順

```bash
# バックアップから復旧
DATE=20250703_120000  # 実際のバックアップ日付に置き換え

# サービスの停止
docker-compose -f docker-compose.prod.yml down

# 設定の復旧
tar -xzf /backup/ogp-service/config-$DATE.tar.gz -C /

# SSL 証明書の復旧
tar -xzf /backup/ogp-service/ssl-$DATE.tar.gz -C /

# サービスの再起動
docker-compose -f docker-compose.prod.yml up -d
```

## 🔍 トラブルシューティング

### 一般的なデプロイメントの問題

#### 1. コンテナーが起動しない
```bash
# ログの確認
docker-compose logs backend
docker-compose logs frontend

# リソース使用量の確認
docker stats

# イメージの再ビルド
docker-compose build --no-cache
```

#### 2. SSL 証明書の問題
```bash
# 証明書の状態確認
sudo certbot certificates

# 証明書の更新
sudo certbot renew

# nginx 設定のテスト
sudo nginx -t
```

#### 3. メモリ使用量が高い
```bash
# メモリ使用量の確認
free -h
docker stats

# サービスの再起動
docker-compose restart

# 必要に応じてスワップを追加
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

#### 4. パフォーマンスの問題
```bash
# システム負荷の確認
htop
iostat -x 1

# nginx の最適化
# nginx.conf で worker_processes を増やす
# gzip 圧縮を有効化
# キャッシュヘッダーを追加
```

---

**次のステップ**: デプロイメント後は、継続的なメンテナンスと監視手順について [OPERATIONS.md](OPERATIONS.md) を参照してください。