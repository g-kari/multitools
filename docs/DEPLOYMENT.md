# OGP Verification Service - Deployment Guide

This guide covers deploying the OGP Verification Service to production environments using various deployment strategies.

## 📋 Overview

### Deployment Options
1. **Docker Compose** (Recommended for single server)
2. **Kubernetes** (For container orchestration)
3. **Sakura VPS + Cloudflare** (Cost-effective production)
4. **Cloud Platforms** (AWS, GCP, Azure)

### Architecture
- **Backend**: Go API server (Docker container)
- **Frontend**: React SPA (Static files served by Nginx)
- **Proxy**: Nginx (SSL termination, reverse proxy)
- **DNS**: Cloudflare (CDN, DDoS protection)

## 🚀 Production Deployment (Docker Compose)

### Prerequisites
- Ubuntu 22.04 LTS server
- Docker and Docker Compose installed
- Domain name with DNS access
- SSL certificate (Let's Encrypt recommended)

### 1. Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Nginx
sudo apt install -y nginx certbot python3-certbot-nginx

# Configure firewall
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 2. Application Deployment

```bash
# Clone repository
git clone <repository-url> /opt/ogp-service
cd /opt/ogp-service

# Create production environment file
cat > .env.production << EOF
# Backend Configuration
PORT=8080
CORS_ORIGINS=https://yourdomain.com
RATE_LIMIT=20

# Frontend Configuration
VITE_API_URL=https://api.yourdomain.com
VITE_ENV=production

# Database Configuration (if needed)
# DATABASE_URL=postgresql://user:pass@localhost/ogp_service
EOF

# Create production Docker Compose file
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

### 3. Frontend Production Build

Create `frontend/Dockerfile.prod`:

```dockerfile
# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Copy source and build
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built files
COPY --from=build /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 4. Nginx Configuration

Create `nginx/prod.conf`:

```nginx
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;

    # Basic settings
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    server_tokens off;

    # Gzip compression
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

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=30r/m;
    limit_req_zone $binary_remote_addr zone=static:10m rate=60r/m;

    # API server
    upstream backend {
        server backend:8080;
        keepalive 32;
    }

    # Frontend server
    upstream frontend {
        server frontend:80;
        keepalive 32;
    }

    # Redirect HTTP to HTTPS
    server {
        listen 80;
        server_name yourdomain.com www.yourdomain.com api.yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    # Main website (Frontend)
    server {
        listen 443 ssl http2;
        server_name yourdomain.com www.yourdomain.com;

        # SSL configuration
        ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
        ssl_session_timeout 1d;
        ssl_session_cache shared:SSL:50m;
        ssl_session_tickets off;

        # Modern configuration
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # HSTS
        add_header Strict-Transport-Security "max-age=63072000" always;

        # Security headers
        add_header X-Content-Type-Options nosniff;
        add_header X-Frame-Options DENY;
        add_header X-XSS-Protection "1; mode=block";
        add_header Referrer-Policy "strict-origin-when-cross-origin";

        # Rate limiting
        limit_req zone=static burst=20 nodelay;

        location / {
            proxy_pass http://frontend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Health check
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }

    # API server
    server {
        listen 443 ssl http2;
        server_name api.yourdomain.com;

        # SSL configuration (same as above)
        ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
        ssl_session_timeout 1d;
        ssl_session_cache shared:SSL:50m;
        ssl_session_tickets off;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # Security headers
        add_header Strict-Transport-Security "max-age=63072000" always;
        add_header X-Content-Type-Options nosniff;
        add_header X-Frame-Options DENY;
        add_header X-XSS-Protection "1; mode=block";

        # CORS headers
        add_header Access-Control-Allow-Origin "https://yourdomain.com" always;
        add_header Access-Control-Allow-Methods "GET, POST, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Content-Type, Authorization" always;

        # Rate limiting for API
        limit_req zone=api burst=10 nodelay;

        location / {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # Timeout settings
            proxy_connect_timeout 30s;
            proxy_send_timeout 30s;
            proxy_read_timeout 30s;
        }

        # Health check
        location /health {
            proxy_pass http://backend/health;
            access_log off;
        }
    }
}
```

### 5. SSL Certificate Setup

```bash
# Install certificate for your domain
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com -d api.yourdomain.com

# Test automatic renewal
sudo certbot renew --dry-run

# Add renewal to crontab
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

### 6. Deploy Application

```bash
# Build and start services
docker-compose -f docker-compose.prod.yml up --build -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 7. Verification

```bash
# Test API
curl https://api.yourdomain.com/health
curl -X POST https://api.yourdomain.com/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'

# Test frontend
curl -I https://yourdomain.com

# Check SSL
curl -I https://yourdomain.com 2>&1 | grep -i ssl
```

## ☸️ Kubernetes Deployment

### Prerequisites
- Kubernetes cluster (v1.24+)
- kubectl configured
- Ingress controller (nginx-ingress)
- cert-manager for SSL

### 1. Namespace and ConfigMap

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

### 2. Backend Deployment

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

### 3. Frontend Deployment

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

### 4. Ingress Configuration

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

### 5. Deploy to Kubernetes

```bash
# Apply all configurations
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n ogp-service
kubectl get services -n ogp-service
kubectl get ingress -n ogp-service

# View logs
kubectl logs -f deployment/ogp-backend -n ogp-service
kubectl logs -f deployment/ogp-frontend -n ogp-service
```

## 🔄 CI/CD Pipeline

### GitHub Actions Workflow

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

## 📊 Monitoring & Logging

### Application Monitoring

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

### Log Aggregation

```bash
# Using journald for system logs
sudo journalctl -u docker -f

# Using docker logs
docker-compose logs -f --tail=100

# Centralized logging with ELK stack
docker run -d --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  elasticsearch:7.14.0
```

## 🔒 Security Hardening

### Server Security

```bash
# Disable root login
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# Setup fail2ban
sudo apt install -y fail2ban
sudo systemctl enable fail2ban

# Configure automatic updates
sudo apt install -y unattended-upgrades
echo 'Unattended-Upgrade::Automatic-Reboot "true";' | sudo tee -a /etc/apt/apt.conf.d/50unattended-upgrades
```

### Container Security

```dockerfile
# Use non-root user in containers
FROM golang:1.21-alpine AS builder
RUN adduser -D -s /bin/sh appuser

FROM alpine:latest
RUN adduser -D -s /bin/sh appuser
USER appuser
COPY --from=builder --chown=appuser:appuser /app/ogp-service .
```

### Network Security

```bash
# Configure firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Setup VPN access for admin (optional)
sudo apt install -y wireguard
```

## 🚨 Disaster Recovery

### Backup Strategy

```bash
#!/bin/bash
# backup.sh - Daily backup script

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/ogp-service"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup application configuration
tar -czf $BACKUP_DIR/config-$DATE.tar.gz /opt/ogp-service/

# Backup SSL certificates
tar -czf $BACKUP_DIR/ssl-$DATE.tar.gz /etc/letsencrypt/

# Upload to cloud storage (optional)
# aws s3 cp $BACKUP_DIR/ s3://your-backup-bucket/ --recursive

# Cleanup old backups (keep last 7 days)
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### Recovery Procedures

```bash
# Restore from backup
DATE=20250703_120000  # Replace with actual backup date

# Stop services
docker-compose -f docker-compose.prod.yml down

# Restore configuration
tar -xzf /backup/ogp-service/config-$DATE.tar.gz -C /

# Restore SSL certificates
tar -xzf /backup/ogp-service/ssl-$DATE.tar.gz -C /

# Restart services
docker-compose -f docker-compose.prod.yml up -d
```

## 🔍 Troubleshooting

### Common Deployment Issues

#### 1. Container Won't Start
```bash
# Check logs
docker-compose logs backend
docker-compose logs frontend

# Check resource usage
docker stats

# Rebuild images
docker-compose build --no-cache
```

#### 2. SSL Certificate Issues
```bash
# Check certificate status
sudo certbot certificates

# Renew certificate
sudo certbot renew

# Test nginx configuration
sudo nginx -t
```

#### 3. High Memory Usage
```bash
# Check memory usage
free -h
docker stats

# Restart services
docker-compose restart

# Add swap if needed
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

#### 4. Performance Issues
```bash
# Check system load
htop
iostat -x 1

# Optimize nginx
# Increase worker_processes in nginx.conf
# Enable gzip compression
# Add caching headers
```

---

**Next Steps**: After deployment, refer to [OPERATIONS.md](OPERATIONS.md) for ongoing maintenance and monitoring procedures.