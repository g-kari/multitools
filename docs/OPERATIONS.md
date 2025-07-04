# OGP Verification Service - Operations Manual

This manual provides comprehensive guidance for operating and maintaining the OGP Verification Service in production.

## 📋 Overview

### Service Architecture
- **Backend**: Go API server (Port 8080)
- **Frontend**: React SPA served by Nginx
- **Proxy**: Nginx (SSL termination, load balancing)
- **Infrastructure**: Sakura VPS + Cloudflare CDN

### Key Components
- API endpoint: `/api/v1/ogp/verify`
- Health check: `/health`
- Rate limiting: 10 requests/minute per IP
- SSL: Let's Encrypt certificates

## 🔧 Daily Operations

### Service Health Monitoring

#### Automated Health Checks
```bash
#!/bin/bash
# /opt/scripts/health_check.sh

API_URL="https://api.yourdomain.com"
FRONTEND_URL="https://yourdomain.com"
LOG_FILE="/var/log/ogp-service/health.log"

# Create log directory
mkdir -p /var/log/ogp-service

# Function to log with timestamp
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> $LOG_FILE
}

# Check API health
if curl -f -s "$API_URL/health" > /dev/null; then
    log "API health check: OK"
else
    log "API health check: FAILED"
    # Send alert
    echo "API health check failed" | mail -s "OGP Service Alert" admin@yourdomain.com
fi

# Check frontend
if curl -f -s "$FRONTEND_URL" > /dev/null; then
    log "Frontend health check: OK"
else
    log "Frontend health check: FAILED"
    echo "Frontend health check failed" | mail -s "OGP Service Alert" admin@yourdomain.com
fi

# Check API functionality
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

#### Manual Health Verification
```bash
# Quick health check
curl https://api.yourdomain.com/health

# Detailed API test
curl -X POST https://api.yourdomain.com/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}' \
  | jq '.validation.is_valid'

# Check response time
time curl -s https://api.yourdomain.com/health

# Check SSL certificate
openssl s_client -connect api.yourdomain.com:443 -servername api.yourdomain.com < /dev/null 2>/dev/null | openssl x509 -noout -dates
```

### Log Management

#### Viewing Logs
```bash
# Application logs
sudo docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs -f

# Nginx access logs
sudo tail -f /var/log/nginx/access.log

# Nginx error logs
sudo tail -f /var/log/nginx/error.log

# System logs
sudo journalctl -u docker -f

# SSL certificate renewal logs
sudo tail -f /var/log/letsencrypt/letsencrypt.log
```

#### Log Rotation Configuration
```bash
# Configure logrotate for application logs
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

# Test logrotate
sudo logrotate -d /etc/logrotate.d/ogp-service
```

## 📊 Performance Monitoring

### Resource Monitoring

#### System Resources
```bash
#!/bin/bash
# /opt/scripts/system_monitor.sh

# CPU usage
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
echo "CPU Usage: ${CPU_USAGE}%"

# Memory usage
MEMORY_USAGE=$(free | grep Mem | awk '{printf("%.1f", $3/$2 * 100.0)}')
echo "Memory Usage: ${MEMORY_USAGE}%"

# Disk usage
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | cut -d'%' -f1)
echo "Disk Usage: ${DISK_USAGE}%"

# Load average
LOAD_AVG=$(uptime | awk -F'load average:' '{print $2}')
echo "Load Average:${LOAD_AVG}"

# Alert if any metric exceeds threshold
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

#### Docker Container Monitoring
```bash
# Container resource usage
docker stats --no-stream

# Container health status
docker-compose -f /opt/ogp-service/docker-compose.prod.yml ps

# Detailed container inspection
docker inspect ogp-service_backend_1 | jq '.State.Health'
```

#### Application Performance Metrics
```bash
#!/bin/bash
# /opt/scripts/performance_monitor.sh

API_URL="https://api.yourdomain.com"

# Measure response time
RESPONSE_TIME=$(curl -o /dev/null -s -w '%{time_total}' "$API_URL/health")
echo "API Response Time: ${RESPONSE_TIME}s"

# Test rate limiting
echo "Testing rate limiting..."
for i in {1..15}; do
    STATUS=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API_URL/api/v1/ogp/verify" \
      -H "Content-Type: application/json" \
      -d '{"url":"https://example.com"}')
    echo "Request $i: HTTP $STATUS"
    sleep 1
done

# Alert if response time is too high
if (( $(echo "$RESPONSE_TIME > 5.0" | bc -l) )); then
    echo "Slow API response: ${RESPONSE_TIME}s" | mail -s "Performance Alert" admin@yourdomain.com
fi
```

## 🔄 Maintenance Procedures

### Regular Maintenance Tasks

#### Weekly Tasks
```bash
#!/bin/bash
# /opt/scripts/weekly_maintenance.sh

echo "Starting weekly maintenance..."

# Update system packages
sudo apt update && sudo apt upgrade -y

# Clean up Docker
docker system prune -f
docker image prune -f

# Rotate logs manually if needed
sudo logrotate -f /etc/logrotate.d/ogp-service

# Check SSL certificate expiry
openssl x509 -in /etc/letsencrypt/live/yourdomain.com/cert.pem -noout -dates

# Backup configuration
tar -czf /backup/ogp-config-$(date +%Y%m%d).tar.gz /opt/ogp-service/

# Check disk space
df -h

echo "Weekly maintenance completed."
```

#### Monthly Tasks
```bash
#!/bin/bash
# /opt/scripts/monthly_maintenance.sh

echo "Starting monthly maintenance..."

# Full system backup
rsync -av /opt/ogp-service/ /backup/monthly/ogp-service-$(date +%Y%m)/

# Security updates
sudo unattended-upgrades

# Review logs for errors
grep -i error /var/log/nginx/error.log | tail -20
sudo docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs --since="30d" | grep -i error

# Performance analysis
# Generate monthly report
echo "Monthly Performance Report - $(date)" > /tmp/monthly_report.txt
echo "==================================" >> /tmp/monthly_report.txt
echo "" >> /tmp/monthly_report.txt

# Average response time over the month
grep "health" /var/log/nginx/access.log | awk '{print $NF}' | awk '{sum+=$1; count++} END {print "Average response time: " sum/count "s"}' >> /tmp/monthly_report.txt

# Error rate
TOTAL_REQUESTS=$(grep -c "api/v1/ogp/verify" /var/log/nginx/access.log)
ERROR_REQUESTS=$(grep "api/v1/ogp/verify" /var/log/nginx/access.log | grep -c " 5[0-9][0-9] ")
ERROR_RATE=$(echo "scale=2; $ERROR_REQUESTS * 100 / $TOTAL_REQUESTS" | bc)
echo "Error rate: ${ERROR_RATE}%" >> /tmp/monthly_report.txt

# Send report
mail -s "Monthly Performance Report" admin@yourdomain.com < /tmp/monthly_report.txt

echo "Monthly maintenance completed."
```

### Application Updates

#### Backend Update Process
```bash
#!/bin/bash
# /opt/scripts/update_backend.sh

cd /opt/ogp-service

echo "Starting backend update..."

# Pull latest code
git pull origin main

# Backup current state
docker-compose -f docker-compose.prod.yml stop backend
docker tag ogp-service_backend:latest ogp-service_backend:backup-$(date +%Y%m%d)

# Build new image
docker-compose -f docker-compose.prod.yml build backend

# Start with new image
docker-compose -f docker-compose.prod.yml up -d backend

# Wait for health check
sleep 30

# Verify health
if curl -f https://api.yourdomain.com/health; then
    echo "Backend update successful"
    # Clean up old backup
    docker rmi ogp-service_backend:backup-$(date --date="7 days ago" +%Y%m%d) 2>/dev/null || true
else
    echo "Backend update failed, rolling back..."
    docker-compose -f docker-compose.prod.yml stop backend
    docker tag ogp-service_backend:backup-$(date +%Y%m%d) ogp-service_backend:latest
    docker-compose -f docker-compose.prod.yml up -d backend
    exit 1
fi
```

#### Frontend Update Process
```bash
#!/bin/bash
# /opt/scripts/update_frontend.sh

cd /opt/ogp-service

echo "Starting frontend update..."

# Pull latest code
git pull origin main

# Build new frontend
docker-compose -f docker-compose.prod.yml build frontend

# Rolling update
docker-compose -f docker-compose.prod.yml up -d frontend

echo "Frontend update completed."
```

## 🚨 Incident Response

### Alert Categories

#### Critical Alerts (Immediate Response)
- API completely down (health check fails)
- SSL certificate expired
- Server unresponsive
- Disk space >95%

#### Warning Alerts (Response within 1 hour)
- High response times (>5 seconds)
- Error rate >10%
- Memory usage >85%
- Rate limiting not working

#### Info Alerts (Response within 24 hours)
- SSL certificate expiring in 7 days
- High CPU usage (>80%)
- Log errors

### Incident Response Procedures

#### API Down Incident
```bash
# 1. Check service status
docker-compose -f /opt/ogp-service/docker-compose.prod.yml ps

# 2. Check logs for errors
docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs --tail=50

# 3. Restart backend service
docker-compose -f /opt/ogp-service/docker-compose.prod.yml restart backend

# 4. If restart fails, check system resources
free -h
df -h
docker system df

# 5. Emergency restart
sudo systemctl restart docker
docker-compose -f /opt/ogp-service/docker-compose.prod.yml up -d

# 6. Verify recovery
curl https://api.yourdomain.com/health
```

#### High Error Rate Incident
```bash
# 1. Check recent errors in logs
tail -100 /var/log/nginx/error.log

# 2. Check application errors
docker-compose -f /opt/ogp-service/docker-compose.prod.yml logs backend --tail=100 | grep -i error

# 3. Analyze traffic patterns
tail -100 /var/log/nginx/access.log | awk '{print $1}' | sort | uniq -c | sort -nr | head -10

# 4. Check for DDoS or abuse
fail2ban-client status nginx-req-limit

# 5. Temporarily increase rate limiting if needed
# Edit nginx configuration and reload
sudo nginx -s reload
```

#### SSL Certificate Expiry
```bash
# 1. Check certificate status
sudo certbot certificates

# 2. Attempt renewal
sudo certbot renew --force-renewal

# 3. If renewal fails, check DNS
dig yourdomain.com
dig api.yourdomain.com

# 4. Manual certificate installation (emergency)
sudo certbot certonly --manual -d yourdomain.com -d api.yourdomain.com

# 5. Restart nginx
sudo systemctl restart nginx
```

### Recovery Procedures

#### Database Recovery (if applicable)
```bash
# 1. Stop application
docker-compose -f /opt/ogp-service/docker-compose.prod.yml stop

# 2. Restore from backup
sudo tar -xzf /backup/database-backup-YYYYMMDD.tar.gz -C /

# 3. Start application
docker-compose -f /opt/ogp-service/docker-compose.prod.yml up -d

# 4. Verify data integrity
# Run application-specific checks
```

#### Full System Recovery
```bash
# 1. Boot from rescue media (if needed)
# 2. Mount filesystems
# 3. Restore from backup
sudo tar -xzf /backup/full-system-backup.tar.gz -C /

# 4. Reinstall Docker if needed
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 5. Start services
cd /opt/ogp-service
docker-compose -f docker-compose.prod.yml up -d

# 6. Verify all services
./scripts/health_check.sh
```

## 📈 Capacity Planning

### Traffic Analysis
```bash
#!/bin/bash
# /opt/scripts/traffic_analysis.sh

# Daily request count
TODAY=$(date +%Y-%m-%d)
REQUESTS_TODAY=$(grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | wc -l)
echo "Requests today: $REQUESTS_TODAY"

# Peak hour analysis
grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | awk '{print $4}' | cut -d: -f2 | sort | uniq -c | sort -nr | head -5

# Response time analysis
grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | awk '{print $NF}' | sort -n | awk '{all[NR] = $0} END{print "Min: " all[1] "s, Max: " all[NR] "s, Median: " all[int(NR/2)] "s"}'

# Error rate
TOTAL=$(grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | wc -l)
ERRORS=$(grep "$TODAY" /var/log/nginx/access.log | grep "api/v1/ogp/verify" | grep -c " 5[0-9][0-9] ")
ERROR_RATE=$(echo "scale=2; $ERRORS * 100 / $TOTAL" | bc)
echo "Error rate: ${ERROR_RATE}%"
```

### Scaling Recommendations

#### Horizontal Scaling Indicators
- CPU usage consistently >70%
- Memory usage consistently >80%
- Response time >3 seconds
- Request volume >1000/hour

#### Vertical Scaling Steps
```bash
# 1. Monitor current usage for 7 days
# 2. If upgrade needed, schedule maintenance window
# 3. Create server snapshot (if supported)
# 4. Upgrade server plan via Sakura Cloud console
# 5. Restart services after upgrade
# 6. Monitor performance for 24 hours
```

## 🔒 Security Operations

### Security Monitoring
```bash
#!/bin/bash
# /opt/scripts/security_monitor.sh

# Check for failed login attempts
sudo grep "Failed password" /var/log/auth.log | tail -10

# Check fail2ban status
sudo fail2ban-client status
sudo fail2ban-client status sshd

# Check for unusual traffic patterns
tail -1000 /var/log/nginx/access.log | awk '{print $1}' | sort | uniq -c | sort -nr | head -20

# Check for known attack patterns
grep -i "union\|select\|drop\|script\|alert" /var/log/nginx/access.log | tail -10

# Check SSL security
nmap --script ssl-enum-ciphers -p 443 yourdomain.com
```

### Security Updates
```bash
#!/bin/bash
# /opt/scripts/security_updates.sh

# Check for security updates
sudo apt list --upgradable | grep -i security

# Apply security updates
sudo unattended-upgrades

# Update Docker images
docker-compose -f /opt/ogp-service/docker-compose.prod.yml pull
docker-compose -f /opt/ogp-service/docker-compose.prod.yml up -d

# Update fail2ban rules
sudo fail2ban-client reload

# Check for CVEs in running services
# Use tools like `docker scout cves` if available
```

## 📋 Runbook Checklists

### Daily Checklist
- [ ] Check service health via monitoring
- [ ] Review error logs
- [ ] Verify SSL certificate status
- [ ] Check disk space usage
- [ ] Review traffic patterns

### Weekly Checklist
- [ ] Apply system updates
- [ ] Clean up Docker images
- [ ] Review performance metrics
- [ ] Test backup restoration
- [ ] Update documentation if needed

### Monthly Checklist
- [ ] Full security scan
- [ ] Capacity planning review
- [ ] Update disaster recovery plan
- [ ] Review and update monitoring thresholds
- [ ] Performance optimization review

## 📞 Contact Information

### Escalation Matrix
1. **On-call Engineer**: +81-XXX-XXXX-XXXX
2. **Technical Lead**: tech-lead@yourdomain.com
3. **Infrastructure Team**: infra@yourdomain.com
4. **Emergency Contact**: emergency@yourdomain.com

### Vendor Contacts
- **Sakura Cloud Support**: +81-3-5332-7071
- **Cloudflare Support**: Enterprise dashboard
- **Domain Registrar**: Contact via web portal

## 📚 Additional Resources

### Documentation Links
- [Setup Guide](SETUP.md)
- [Deployment Guide](DEPLOYMENT.md)
- [Terraform Setup](TERRAFORM_SETUP.md)
- [API Documentation](../backend/api/README.md)

### External Resources
- [Docker Documentation](https://docs.docker.com/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [Sakura Cloud Documentation](https://manual.sakura.ad.jp/cloud/)

---

**Important**: Keep this operations manual updated with any changes to procedures, contact information, or system architecture.