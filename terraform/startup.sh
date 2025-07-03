#!/bin/bash

# OGP Verification Service Server Setup Script
set -e

# Update system
apt-get update
apt-get upgrade -y

# Install required packages
apt-get install -y \
    curl \
    wget \
    unzip \
    git \
    nginx \
    certbot \
    python3-certbot-nginx \
    ufw \
    fail2ban \
    htop \
    tree

# Install Go
GO_VERSION="1.21.5"
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
usermod -aG docker ubuntu

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Create application directory
mkdir -p /opt/ogp-service
chown ubuntu:ubuntu /opt/ogp-service

# Configure Nginx
cat > /etc/nginx/sites-available/ogp-service <<EOF
server {
    listen 80;
    server_name _;
    
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    location /health {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    location / {
        return 200 '{"message": "OGP Verification Service", "status": "ok"}';
        add_header Content-Type application/json;
    }
}
EOF

# Enable the site
ln -sf /etc/nginx/sites-available/ogp-service /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test and reload nginx
nginx -t
systemctl reload nginx

# Configure UFW firewall
ufw --force enable
ufw allow ssh
ufw allow 'Nginx Full'
ufw allow 8080

# Configure fail2ban
systemctl enable fail2ban
systemctl start fail2ban

# Create systemd service for OGP service
cat > /etc/systemd/system/ogp-service.service <<EOF
[Unit]
Description=OGP Verification Service
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/opt/ogp-service
ExecStart=/opt/ogp-service/ogp-service
Restart=always
RestartSec=3
Environment=PORT=8080
Environment=CORS_ORIGINS=*

[Install]
WantedBy=multi-user.target
EOF

# Enable the service
systemctl daemon-reload
systemctl enable ogp-service

# Create deployment script
cat > /opt/ogp-service/deploy.sh <<EOF
#!/bin/bash
set -e

echo "Starting OGP Service deployment..."

# Stop the service if running
sudo systemctl stop ogp-service || true

# Backup current binary if exists
if [ -f /opt/ogp-service/ogp-service ]; then
    sudo cp /opt/ogp-service/ogp-service /opt/ogp-service/ogp-service.backup
fi

# Build the service
cd /opt/ogp-service
/usr/local/go/bin/go build -o ogp-service ./cmd/main.go

# Set permissions
sudo chown ubuntu:ubuntu /opt/ogp-service/ogp-service
chmod +x /opt/ogp-service/ogp-service

# Start the service
sudo systemctl start ogp-service
sudo systemctl status ogp-service

echo "OGP Service deployment completed successfully!"
EOF

chmod +x /opt/ogp-service/deploy.sh
chown ubuntu:ubuntu /opt/ogp-service/deploy.sh

# Create log directory
mkdir -p /var/log/ogp-service
chown ubuntu:ubuntu /var/log/ogp-service

# Set up log rotation
cat > /etc/logrotate.d/ogp-service <<EOF
/var/log/ogp-service/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 ubuntu ubuntu
    postrotate
        systemctl reload ogp-service
    endscript
}
EOF

echo "OGP Verification Service server setup completed!"
echo "Server is ready for deployment."
echo "Use /opt/ogp-service/deploy.sh to deploy the application."