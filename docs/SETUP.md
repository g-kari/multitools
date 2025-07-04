# OGP Verification Service - Setup Guide

This guide will help you set up and run the OGP Verification Service locally for development or production use.

## 📋 Prerequisites

### Required Software
- **Docker** (v20.0+) and **Docker Compose** (v2.0+)
- **Git** for cloning the repository

### Optional (for native development)
- **Go** (v1.21+) for backend development
- **Node.js** (v18+) and **npm** for frontend development
- **Terraform** (v1.0+) for infrastructure deployment

### System Requirements
- **Memory**: 2GB RAM minimum, 4GB recommended
- **Storage**: 1GB free space
- **Network**: Internet connection for fetching OGP data

## 🚀 Quick Start (Docker)

### 1. Clone the Repository

```bash
git clone <repository-url>
cd ogp-verification-service
```

### 2. Start Services with Docker Compose

```bash
# Start backend and frontend services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

### 3. Verify Installation

```bash
# Test backend health
curl http://localhost:8080/health

# Test OGP verification
curl -X POST http://localhost:8080/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'

# Access frontend
open http://localhost:3000
```

### 4. Stop Services

```bash
docker-compose down
```

## 🛠️ Development Setup

### Backend Development

1. **Prerequisites**
   ```bash
   # Install Go 1.21+
   go version
   ```

2. **Setup Backend**
   ```bash
   cd backend
   
   # Install dependencies
   go mod download
   
   # Run tests
   go test ./...
   
   # Run development server
   go run cmd/main.go
   ```

3. **Backend Configuration**
   ```bash
   # Environment variables (optional)
   export PORT=8080
   export CORS_ORIGINS="*"
   export RATE_LIMIT=10
   ```

### Frontend Development

1. **Prerequisites**
   ```bash
   # Install Node.js 18+ and npm
   node --version
   npm --version
   ```

2. **Setup Frontend**
   ```bash
   cd frontend
   
   # Install dependencies
   npm install
   
   # Start development server
   npm run dev
   
   # Build for production
   npm run build
   ```

3. **Frontend Configuration**
   ```bash
   # Create .env file
   cat > .env << EOF
   VITE_API_URL=http://localhost:8080
   VITE_ENV=development
   EOF
   ```

## 🔧 Configuration

### Environment Variables

#### Backend Variables
| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PORT` | Server port | `8080` | `8080` |
| `CORS_ORIGINS` | Allowed CORS origins | `*` | `https://example.com` |
| `RATE_LIMIT` | Requests per minute per IP | `10` | `20` |

#### Frontend Variables
| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` | `https://api.example.com` |
| `VITE_ENV` | Environment name | `development` | `production` |

### Docker Compose Configuration

Edit `docker-compose.yml` to customize:

```yaml
services:
  backend:
    environment:
      - PORT=8080
      - CORS_ORIGINS=*
      - RATE_LIMIT=10
    ports:
      - "8080:8080"  # Change external port if needed
  
  frontend:
    environment:
      - VITE_API_URL=http://localhost:8080
    ports:
      - "3000:3000"  # Change external port if needed
```

## 🧪 Testing

### Run All Tests

```bash
# Backend tests
cd backend
go test ./... -v -cover

# Frontend tests (if implemented)
cd frontend
npm test

# Integration tests
cd backend
docker build -f Dockerfile.test -t ogp-backend-test .
docker run --rm ogp-backend-test
```

### E2E Tests

```bash
cd backend/tests/e2e
./run-e2e.sh
```

### Manual Testing

```bash
# Test various websites
curl -X POST http://localhost:8080/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'

curl -X POST http://localhost:8080/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.wikipedia.org"}'

# Test rate limiting (run 11 times quickly)
for i in {1..11}; do
  curl -X POST http://localhost:8080/api/v1/ogp/verify \
    -H "Content-Type: application/json" \
    -d '{"url":"https://example.com"}' \
    -w "Status: %{http_code}\n"
done
```

## 📚 API Documentation

### Swagger UI (Recommended)

```bash
# Start Swagger UI server
cd backend
go run cmd/swagger/main.go

# Open in browser
open http://localhost:8081
```

### API Endpoints

- **GET** `/health` - Health check
- **POST** `/api/v1/ogp/verify` - Verify OGP metadata
- **OPTIONS** `/api/v1/ogp/verify` - CORS preflight

### Example Request

```bash
curl -X POST http://localhost:8080/api/v1/ogp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com"
  }'
```

### Example Response

```json
{
  "url": "https://github.com",
  "ogp_data": {
    "title": "GitHub · Build and ship software on a single, collaborative platform",
    "description": "Join the world's most widely adopted...",
    "image": "https://github.githubassets.com/assets/home24-5939032587c9.jpg",
    "url": "https://github.com/",
    "type": "object",
    "site_name": "GitHub"
  },
  "validation": {
    "is_valid": true,
    "warnings": [],
    "errors": [],
    "checks": {
      "has_title": true,
      "has_description": true,
      "has_image": true,
      "image_valid": true,
      "url_valid": true
    }
  },
  "previews": {
    "twitter": {
      "platform": "twitter",
      "title": "GitHub · Build and ship software...",
      "description": "Join the world's most widely adopted...",
      "image": "https://github.githubassets.com/assets/home24-5939032587c9.jpg",
      "is_valid": true,
      "warnings": [],
      "title_length": 69,
      "desc_length": 186,
      "max_title_len": 70,
      "max_desc_len": 200
    },
    "facebook": { ... },
    "discord": { ... }
  },
  "timestamp": "2025-07-03T18:00:00Z"
}
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
docker-compose up -d --scale backend=0
docker run -p 8081:8080 multitools-backend
```

#### 2. Docker Permission Issues
```bash
# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Or use sudo
sudo docker-compose up -d
```

#### 3. Go Module Issues
```bash
cd backend
go clean -modcache
go mod download
go mod tidy
```

#### 4. Frontend Build Issues
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

#### 5. CORS Issues
```bash
# Check CORS configuration
curl -X OPTIONS http://localhost:8080/api/v1/ogp/verify -v

# Update docker-compose.yml
environment:
  - CORS_ORIGINS=http://localhost:3000
```

### Debug Commands

```bash
# View detailed logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Check container status
docker-compose ps

# Enter container for debugging
docker-compose exec backend sh
docker-compose exec frontend sh

# Check resource usage
docker stats

# Reset everything
docker-compose down --volumes --remove-orphans
docker system prune -f
```

### Log Analysis

```bash
# Backend logs
docker-compose logs backend | grep ERROR
docker-compose logs backend | grep "Rate limit"

# Test connectivity
docker-compose exec backend ping google.com
docker-compose exec backend curl -I https://github.com
```

## 🔒 Security Considerations

### Development Environment
- Default configuration allows all CORS origins (`*`)
- Rate limiting is set to 10 requests/minute per IP
- Private IP addresses are blocked by default

### Production Recommendations
1. **Configure specific CORS origins**
   ```yaml
   environment:
     - CORS_ORIGINS=https://yourdomain.com
   ```

2. **Use HTTPS in production**
3. **Set up proper firewall rules**
4. **Monitor rate limiting logs**
5. **Regular security updates**

## 📊 Monitoring

### Health Checks

```bash
# Manual health check
curl http://localhost:8080/health

# Automated monitoring (add to cron)
#!/bin/bash
if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
  echo "Service is down" | mail -s "OGP Service Alert" admin@example.com
fi
```

### Metrics Collection

The service exposes basic metrics through logs. For production, consider:
- **Prometheus** for metrics collection
- **Grafana** for visualization
- **ELK Stack** for log analysis

## 🆘 Getting Help

### Resources
- **API Documentation**: http://localhost:8081 (Swagger UI)
- **Project Repository**: [GitHub Repository]
- **Issue Tracker**: [GitHub Issues]

### Support Channels
1. Check this setup guide
2. Review troubleshooting section
3. Search existing GitHub issues
4. Create new GitHub issue with:
   - OS and version
   - Docker version
   - Error messages
   - Steps to reproduce

### Debug Information to Include

```bash
# System information
docker --version
docker-compose --version
uname -a

# Service status
docker-compose ps
docker-compose logs --tail=50

# Network connectivity
curl -I https://github.com
ping google.com
```

---

**Next Steps**: Once setup is complete, see [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment instructions.