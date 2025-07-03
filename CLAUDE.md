# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an **OGP (Open Graph Protocol) Verification Service** project that analyzes websites for OGP metadata and provides validation results with platform-specific previews for Twitter/X, Facebook, and Discord.

## Technology Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Standard Go HTTP server with JSON API
- **Deployment**: Sakura VPS (512MB) with Ubuntu 22.04 LTS

### Frontend
- **Framework**: React with TypeScript
- **Runtime**: Bun (package manager and build tool)
- **Styling**: Tailwind CSS (recommended)
- **Deployment**: Cloudflare Pages

### Infrastructure
- **IaC**: Terraform
- **CI/CD**: GitHub Actions
- **Monitoring**: Cloudflare Analytics + custom health checks

## Development Commands

### Backend (Go)
```bash
# Initialize Go module
go mod init ogp-verification-service

# Run development server
go run main.go

# Build for production
go build -o ogp-service

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Frontend (React + Bun)
```bash
# Initialize Bun project
bun create react-app frontend --template typescript

# Install dependencies
bun install

# Development server
bun dev

# Build for production
bun run build

# Run tests
bun test

# Type checking
bun run type-check
```

### Infrastructure (Terraform)
```bash
# Initialize Terraform
terraform init

# Plan infrastructure changes
terraform plan

# Apply infrastructure changes
terraform apply

# Destroy infrastructure
terraform destroy
```

### Docker Development
```bash
# Start development environment
docker-compose up -d

# Build and start services
docker-compose up --build

# Stop services
docker-compose down
```

## Project Structure

```
├── backend/              # Go backend application
│   ├── cmd/             # Application entrypoints
│   ├── internal/        # Private application code
│   │   ├── handlers/    # HTTP handlers
│   │   ├── models/      # Data models
│   │   ├── services/    # Business logic
│   │   └── validators/  # Input validation
│   ├── pkg/            # Public library code
│   └── go.mod          # Go module definition
├── frontend/            # React + Bun frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── hooks/       # Custom React hooks
│   │   ├── services/    # API services
│   │   └── types/       # TypeScript types
│   ├── package.json
│   └── bun.lockb
├── terraform/           # Infrastructure as Code
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── docker-compose.yml   # Local development environment
└── .github/workflows/   # GitHub Actions CI/CD
```

## API Specification

### Main Endpoint
- **POST** `/api/v1/ogp/verify`
- **Request**: `{"url": "https://example.com"}`
- **Response**: JSON with OGP data, validation results, and platform previews

### Platform Support
- **Twitter/X**: Title (70 chars), Description (200 chars), Image (1200x630px)
- **Facebook**: Title (100 chars), Description (300 chars), Image (1200x630px)
- **Discord**: Title (256 chars), Description (2048 chars), Image (flexible)

## Security Requirements

- CORS configuration for frontend domain
- Rate limiting: 10 requests/minute per IP
- Private IP address blocking
- Input validation and sanitization
- No sensitive data in logs or responses

## Performance Requirements

- Response time: < 3 seconds
- Concurrent requests: 100 req/sec
- Request timeout: 10 seconds
- Test coverage: 80%+

## Development Workflow

1. **Always commit completed work**: When finishing any task or making significant progress, commit changes with descriptive messages
2. **Use feature branches**: Create branches for new features or major changes
3. **Write tests**: Implement unit tests for new functionality
4. **Document changes**: Update relevant documentation when making changes
5. **Test before deployment**: Run all tests and type checks before pushing

## Commit Guidelines

- Use conventional commit messages
- Always commit when completing tasks
- Include 🤖 emoji for AI-generated commits
- Example: `feat: implement OGP validation service 🤖`

## Environment Variables

### Backend
- `PORT`: Server port (default: 8080)
- `CORS_ORIGINS`: Allowed CORS origins
- `RATE_LIMIT`: Requests per minute per IP

### Frontend
- `REACT_APP_API_URL`: Backend API URL
- `REACT_APP_ENV`: Environment (development/production)

## Testing

### Backend Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/handlers -v
```

### Frontend Testing
```bash
# Run all tests
bun test

# Run in watch mode
bun test --watch

# Run with coverage
bun test --coverage
```

## Deployment

### Production Deployment
1. Backend: Build Go binary and deploy to Sakura VPS
2. Frontend: Push to GitHub (auto-deploys to Cloudflare Pages)
3. Infrastructure: Apply Terraform changes

### Development Environment
Use Docker Compose for local development with hot reload enabled.

## Monitoring

- **Health Check**: `/health` endpoint
- **Metrics**: Response times, error rates, request counts
- **Logs**: Structured JSON logging
- **Alerts**: Configure for high error rates or downtime

## Architecture Notes

This is a distributed system with:
- **Stateless backend** for horizontal scaling
- **Static frontend** for CDN delivery
- **Terraform IaC** for reproducible infrastructure
- **GitHub Actions** for automated CI/CD

The service focuses on OGP validation and preview generation for social media platforms, with emphasis on performance, security, and reliability.