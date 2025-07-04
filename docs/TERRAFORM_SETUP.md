# Terraform Infrastructure Setup Guide

This guide explains how to deploy the OGP Verification Service to production using Terraform on Sakura VPS and Cloudflare.

## 📋 Prerequisites

### Required Accounts & Services
- **Sakura VPS** account with API access
- **Cloudflare** account with domain management
- **GitHub** repository for CI/CD
- **Domain name** managed by Cloudflare

### Required Software
- **Terraform** (v1.0+)
- **Git** for version control
- **SSH client** for server access

### Required Credentials
- Sakura Cloud API key and secret
- Cloudflare API token
- SSH key pair for server access

## 🚀 Quick Start

### 1. Configure Credentials

```bash
# Sakura Cloud credentials
export SAKURACLOUD_ACCESS_TOKEN="your-access-token"
export SAKURACLOUD_ACCESS_TOKEN_SECRET="your-access-secret"
export SAKURACLOUD_ZONE="is1b"  # or your preferred zone

# Cloudflare credentials
export CLOUDFLARE_API_TOKEN="your-cloudflare-token"

# SSH key for server access
ssh-keygen -t rsa -b 4096 -f ~/.ssh/ogp-service-key
```

### 2. Initialize Terraform

```bash
cd terraform

# Initialize Terraform
terraform init

# Copy example variables
cp terraform.tfvars.example terraform.tfvars

# Edit configuration
nano terraform.tfvars
```

### 3. Configure Variables

Edit `terraform.tfvars`:

```hcl
# Domain and DNS
domain_name = "yourdomain.com"
subdomain   = "ogp-api"  # Will create ogp-api.yourdomain.com

# Server configuration
server_name = "ogp-service-prod"
server_plan = "1core-1gb"  # or "2core-2gb" for higher load

# SSH access
ssh_key_name = "ogp-service-key"
ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2E..."  # Your public key

# Application settings
app_environment = "production"
cors_origins = "https://yourdomain.com"
rate_limit = "20"  # requests per minute

# Cloudflare settings
cloudflare_zone_id = "your-zone-id"

# Optional: Database settings (if using external DB)
# database_host = "your-db-host"
# database_name = "ogp_service"
```

### 4. Plan and Apply

```bash
# Review planned changes
terraform plan

# Apply infrastructure
terraform apply

# Note the outputs
terraform output
```

## 🛠️ Detailed Configuration

### Sakura VPS Configuration

The Terraform configuration creates:

#### Server Specifications
- **OS**: Ubuntu 22.04 LTS
- **Memory**: 512MB (1core-1gb) or 2GB (2core-2gb)
- **Storage**: 20GB SSD
- **Network**: Public IP with firewall rules

#### Security Configuration
```hcl
# Firewall rules
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

#### Automated Setup Script
The server runs a cloud-init script that:
1. Updates system packages
2. Installs Docker and Docker Compose
3. Configures firewall (UFW)
4. Sets up SSL certificates (Let's Encrypt)
5. Deploys the OGP service
6. Configures monitoring

### Cloudflare Configuration

#### DNS Records
- **A Record**: `ogp-api.yourdomain.com` → Server IP
- **CNAME Record**: `www.ogp-api.yourdomain.com` → `ogp-api.yourdomain.com`

#### Security Settings
- **SSL/TLS**: Full (strict)
- **Security Level**: Medium
- **Bot Fight Mode**: Enabled
- **Rate Limiting**: 100 requests/minute per IP

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

## 🔧 Advanced Configuration

### Variables Reference

#### Required Variables
| Variable | Description | Example |
|----------|-------------|---------|
| `domain_name` | Your domain name | `"example.com"` |
| `subdomain` | API subdomain | `"ogp-api"` |
| `ssh_public_key` | SSH public key | `"ssh-rsa AAAAB3..."` |
| `cloudflare_zone_id` | Cloudflare zone ID | `"abcd1234..."` |

#### Optional Variables
| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `server_plan` | Sakura VPS plan | `"1core-1gb"` | `"2core-2gb"` |
| `server_name` | Server hostname | `"ogp-service"` | `"ogp-prod"` |
| `app_environment` | Environment name | `"production"` | `"staging"` |
| `cors_origins` | Allowed CORS origins | `"*"` | `"https://example.com"` |
| `rate_limit` | Requests per minute | `"10"` | `"20"` |

### Server Plans

#### 1core-1gb (Recommended for development)
- **CPU**: 1 vCPU
- **Memory**: 1GB
- **Storage**: 20GB SSD
- **Cost**: ~¥680/month

#### 2core-2gb (Recommended for production)
- **CPU**: 2 vCPU
- **Memory**: 2GB
- **Storage**: 20GB SSD
- **Cost**: ~¥1,580/month

#### 4core-4gb (High load scenarios)
- **CPU**: 4 vCPU
- **Memory**: 4GB
- **Storage**: 20GB SSD
- **Cost**: ~¥3,200/month

### Custom Server Configuration

Create `terraform/modules/server/user-data.sh` for custom setup:

```bash
#!/bin/bash
set -e

# Custom application setup
echo "Setting up OGP Verification Service..."

# Install additional packages
apt-get update
apt-get install -y htop nginx-utils

# Configure custom monitoring
cat > /opt/ogp-service/monitoring.sh << 'EOF'
#!/bin/bash
# Custom monitoring script
curl -f http://localhost:8080/health || systemctl restart ogp-service
EOF

chmod +x /opt/ogp-service/monitoring.sh

# Add to crontab for health monitoring
echo "*/5 * * * * /opt/ogp-service/monitoring.sh" | crontab -

echo "Custom setup completed!"
```

## 🚀 Deployment Process

### 1. Pre-deployment Checks

```bash
# Validate Terraform configuration
terraform validate

# Check formatting
terraform fmt -check

# Security scan (optional)
terraform plan -out=tfplan
terraform show -json tfplan | jq > tfplan.json
# Use tools like tfsec or checkov for security analysis
```

### 2. Staging Deployment

```bash
# Create staging workspace
terraform workspace new staging
terraform workspace select staging

# Deploy to staging
terraform apply -var="app_environment=staging" -var="subdomain=ogp-api-staging"
```

### 3. Production Deployment

```bash
# Switch to production workspace
terraform workspace new production
terraform workspace select production

# Deploy to production
terraform apply
```

### 4. Post-deployment Verification

```bash
# Get server information
terraform output server_ip
terraform output server_fqdn

# Test deployment
FQDN=$(terraform output -raw server_fqdn)
curl https://${FQDN}/health
curl -X POST https://${FQDN}/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'
```

## 🔒 Security Best Practices

### SSH Access
```bash
# Connect to server
SERVER_IP=$(terraform output -raw server_ip)
ssh -i ~/.ssh/ogp-service-key ubuntu@${SERVER_IP}

# Disable password authentication (done automatically)
# Configure fail2ban (included in cloud-init)
# Set up automatic security updates (included)
```

### Firewall Rules
```hcl
# Only allow necessary ports
resource "sakuracloud_simple_monitor" "ogp_service" {
  # HTTP (80) - redirects to HTTPS
  # HTTPS (443) - main service
  # SSH (22) - admin access only
  # All other ports blocked by default
}
```

### SSL/TLS Configuration
- **Certificate**: Let's Encrypt (auto-renewal)
- **TLS Version**: 1.2+ only
- **Ciphers**: Strong ciphers only
- **HSTS**: Enabled with 1-year max-age

### Monitoring & Alerting
```bash
# Built-in monitoring checks:
# - Service health (/health endpoint)
# - SSL certificate expiry
# - Disk space usage
# - Memory usage
# - Load average
```

## 📊 Monitoring & Maintenance

### Health Checks

Cloudflare monitors the service with:
- **Endpoint**: `https://ogp-api.yourdomain.com/health`
- **Frequency**: Every 60 seconds
- **Timeout**: 10 seconds
- **Expected Status**: 200 OK

### Log Management

```bash
# SSH into server
ssh -i ~/.ssh/ogp-service-key ubuntu@${SERVER_IP}

# View application logs
sudo docker-compose -f /opt/ogp-service/docker-compose.yml logs -f

# View system logs
sudo journalctl -fu ogp-service

# View access logs
sudo tail -f /var/log/nginx/access.log
```

### Backup Strategy

```bash
# Automated daily backups (configured in cloud-init)
# - Application configuration: /opt/ogp-service/
# - Docker images and data
# - System configuration
# - SSL certificates

# Manual backup
sudo tar -czf /tmp/ogp-backup-$(date +%Y%m%d).tar.gz \
  /opt/ogp-service/ \
  /etc/letsencrypt/ \
  /etc/nginx/sites-available/
```

### Updates & Maintenance

```bash
# SSH into server
ssh -i ~/.ssh/ogp-service-key ubuntu@${SERVER_IP}

# Update application (zero-downtime)
cd /opt/ogp-service
sudo docker-compose pull
sudo docker-compose up -d

# Update system packages
sudo apt update && sudo apt upgrade -y
sudo reboot  # if kernel updates
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Terraform Apply Fails

```bash
# Check credentials
echo $SAKURACLOUD_ACCESS_TOKEN
echo $CLOUDFLARE_API_TOKEN

# Verify zone ID
curl -X GET "https://api.cloudflare.com/client/v4/zones" \
  -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
  -H "Content-Type: application/json"

# Check terraform state
terraform show
terraform refresh
```

#### 2. Server Not Accessible

```bash
# Check server status in Sakura Cloud console
# Verify security group rules
# Check cloud-init logs on server:
ssh ubuntu@${SERVER_IP} sudo cat /var/log/cloud-init-output.log
```

#### 3. SSL Certificate Issues

```bash
# SSH into server and check certificate
ssh ubuntu@${SERVER_IP}
sudo certbot certificates
sudo nginx -t
sudo systemctl status nginx
```

#### 4. Service Not Starting

```bash
# Check Docker status
sudo docker ps -a
sudo docker-compose logs

# Check system resources
free -h
df -h
sudo systemctl status ogp-service
```

### Debugging Commands

```bash
# Terraform debugging
export TF_LOG=DEBUG
terraform apply

# Server diagnostics
curl -I https://ogp-api.yourdomain.com
dig ogp-api.yourdomain.com
nslookup ogp-api.yourdomain.com

# Port connectivity
telnet ogp-api.yourdomain.com 443
nc -zv ogp-api.yourdomain.com 80 443
```

## 🔄 CI/CD Integration

### GitHub Actions

Add to `.github/workflows/terraform.yml`:

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

### Required Secrets

Add these to GitHub repository secrets:
- `SAKURACLOUD_ACCESS_TOKEN`
- `SAKURACLOUD_ACCESS_TOKEN_SECRET`
- `CLOUDFLARE_API_TOKEN`

## 🧹 Cleanup

### Remove Infrastructure

```bash
# Destroy all resources
terraform destroy

# Confirm deletion
terraform state list  # Should be empty

# Clean up local files
rm -rf .terraform/
rm terraform.tfstate*
```

### Cost Optimization

```bash
# Stop server (retains data, stops billing for compute)
terraform apply -var="server_power_state=down"

# Reduce server size
terraform apply -var="server_plan=1core-1gb"

# Remove monitoring (optional cost saving)
terraform apply -var="enable_monitoring=false"
```

## 📚 Additional Resources

### Terraform Documentation
- [Sakura Cloud Provider](https://registry.terraform.io/providers/sacloud/sakuracloud/latest/docs)
- [Cloudflare Provider](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs)

### Sakura Cloud Documentation
- [API Reference](https://manual.sakura.ad.jp/cloud/api/)
- [Server Plans](https://cloud.sakura.ad.jp/specification/server-disk/)

### Security Resources
- [Let's Encrypt](https://letsencrypt.org/)
- [UFW Documentation](https://help.ubuntu.com/community/UFW)
- [fail2ban](https://www.fail2ban.org/)

---

**Next Steps**: After infrastructure deployment, see [DEPLOYMENT.md](DEPLOYMENT.md) for application deployment and [OPERATIONS.md](OPERATIONS.md) for ongoing maintenance.