# OGP Verification Service - Architecture Documentation

This document describes the system architecture, components, and data flow of the OGP Verification Service.

## 🏗️ System Overview

The OGP Verification Service is a distributed web application that analyzes websites for Open Graph Protocol (OGP) metadata and provides validation results with platform-specific previews.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLOUDFLARE CDN                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   DNS & Proxy   │  │  DDoS Protection │  │   SSL/TLS       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     SAKURA VPS SERVER                          │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                     NGINX PROXY                             ││
│  │  • SSL Termination  • Load Balancing  • Rate Limiting      ││
│  │  • Static Serving   • CORS Headers    • Security Headers   ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│           ┌────────────────────┼────────────────────┐           │
│           ▼                    ▼                    ▼           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   FRONTEND      │  │    BACKEND      │  │  MONITORING     │ │
│  │   (React SPA)   │  │   (Go API)      │  │   & LOGGING     │ │
│  │                 │  │                 │  │                 │ │
│  │ • React 18      │  │ • Go 1.21       │  │ • Docker Logs   │ │
│  │ • TypeScript    │  │ • HTTP Server   │  │ • System Logs   │ │
│  │ • Tailwind CSS  │  │ • OGP Parser    │  │ • Health Checks │ │
│  │ • Vite          │  │ • Rate Limiter  │  │ • Metrics       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                │                                │
│                                ▼                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  EXTERNAL APIS                              ││
│  │  • Target Websites (OGP Data Sources)                      ││
│  │  • DNS Resolution Services                                 ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## 🔧 Component Architecture

### Frontend Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     REACT FRONTEND                             │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  PRESENTATION LAYER                         ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │  URL Input  │  │ OGP Display │  │  Platform   │        ││
│  │  │  Component  │  │  Component  │  │  Preview    │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   SERVICE LAYER                             ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │ API Client  │  │ Validation  │  │ State Mgmt  │        ││
│  │  │  Service    │  │  Service    │  │   (React)   │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                     BUILD TOOLS                             ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │    Vite     │  │ TypeScript  │  │ Tailwind    │        ││
│  │  │   Bundler   │  │   Compiler  │  │     CSS     │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

### Backend Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      GO BACKEND                                 │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   HTTP LAYER                                ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │   Router    │  │ Middleware  │  │  Handlers   │        ││
│  │  │ (Gorilla)   │  │  (CORS,     │  │ (OGP, Health)│        ││
│  │  │             │  │   Logging)  │  │             │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  SERVICE LAYER                              ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │ OGP Parser  │  │ Validator   │  │ Rate Limiter│        ││
│  │  │   Service   │  │  Service    │  │   Service   │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   DATA LAYER                                ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │   Models    │  │ HTTP Client │  │  Utilities  │        ││
│  │  │ (Structs)   │  │   (net/http)│  │  (Parsing)  │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## 🌐 Network Architecture

### DNS and CDN Configuration

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     USER        │    │   CLOUDFLARE    │    │  SAKURA VPS     │
│   (Browser)     │    │      CDN        │    │    SERVER       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ 1. DNS Query          │                       │
         │ yourdomain.com        │                       │
         ├──────────────────────►│                       │
         │                       │                       │
         │ 2. DNS Response       │                       │
         │ (Cloudflare IP)       │                       │
         │◄──────────────────────┤                       │
         │                       │                       │
         │ 3. HTTP/HTTPS Request │                       │
         ├──────────────────────►│                       │
         │                       │ 4. Origin Request     │
         │                       ├──────────────────────►│
         │                       │                       │
         │                       │ 5. Origin Response    │
         │                       │◄──────────────────────┤
         │ 6. Cached Response    │                       │
         │◄──────────────────────┤                       │
```

### Load Balancing and Routing

```
┌─────────────────────────────────────────────────────────────────┐
│                      NGINX CONFIGURATION                       │
│                                                                 │
│  Frontend Domain: yourdomain.com                               │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │  Location: /                                                ││
│  │  ┌─────────────────┐                                       ││
│  │  │  React SPA      │ ◄── Static Files & SPA Routing       ││
│  │  │  (Port 3000)    │                                       ││
│  │  └─────────────────┘                                       ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                 │
│  API Domain: api.yourdomain.com                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │  Location: /                                                ││
│  │  ┌─────────────────┐                                       ││
│  │  │  Go Backend     │ ◄── API Requests & Health Checks     ││
│  │  │  (Port 8080)    │                                       ││
│  │  └─────────────────┘                                       ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## 📊 Data Flow Architecture

### Request Flow Diagram

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ Browser │    │ Cloudflare│   │  Nginx  │    │  Backend│    │ Target  │
│         │    │   CDN   │    │  Proxy  │    │   API   │    │ Website │
└─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘
     │              │              │              │              │
     │ 1. POST      │              │              │              │
     │ /ogp/verify  │              │              │              │
     ├─────────────►│              │              │              │
     │              │ 2. Forward   │              │              │
     │              ├─────────────►│              │              │
     │              │              │ 3. Route     │              │
     │              │              │ to Backend   │              │
     │              │              ├─────────────►│              │
     │              │              │              │ 4. Fetch     │
     │              │              │              │ OGP Data     │
     │              │              │              ├─────────────►│
     │              │              │              │              │
     │              │              │              │ 5. HTML      │
     │              │              │              │ Response     │
     │              │              │              │◄─────────────┤
     │              │              │              │              │
     │              │              │ 6. JSON      │              │
     │              │              │ Response     │              │
     │              │              │◄─────────────┤              │
     │              │ 7. Response  │              │              │
     │              │◄─────────────┤              │              │
     │ 8. JSON      │              │              │              │
     │ with OGP     │              │              │              │
     │◄─────────────┤              │              │              │
```

### OGP Processing Flow

```
┌─────────────────┐
│   URL Input     │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│ Input Validation│
│ • URL Format    │
│ • Private IP    │
│ • Rate Limit    │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  HTTP Request   │
│ • Fetch HTML    │
│ • 10s Timeout   │
│ • User-Agent    │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  HTML Parsing   │
│ • Extract OGP   │
│ • Meta Tags     │
│ • Validation    │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│ Platform Preview│
│ • Twitter/X     │
│ • Facebook      │
│ • Discord       │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│ JSON Response   │
│ • OGP Data      │
│ • Validation    │
│ • Previews      │
└─────────────────┘
```

## 🔐 Security Architecture

### Security Layers

```
┌─────────────────────────────────────────────────────────────────┐
│                       SECURITY LAYERS                          │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  NETWORK SECURITY                           ││
│  │  • Cloudflare DDoS Protection                              ││
│  │  • WAF (Web Application Firewall)                          ││
│  │  • Rate Limiting (30 req/min per IP)                       ││
│  │  • Geo-blocking (if enabled)                               ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                 TRANSPORT SECURITY                          ││
│  │  • TLS 1.2/1.3 Only                                        ││
│  │  • HSTS Headers                                             ││
│  │  • Strong Cipher Suites                                    ││
│  │  • Certificate Pinning                                     ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                APPLICATION SECURITY                         ││
│  │  • Input Validation                                         ││
│  │  • Private IP Blocking                                      ││
│  │  • CORS Configuration                                       ││
│  │  • Security Headers                                         ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                 SYSTEM SECURITY                             ││
│  │  • Docker Container Isolation                              ││
│  │  • Non-root User Processes                                 ││
│  │  • UFW Firewall                                            ││
│  │  • fail2ban Intrusion Prevention                           ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

### Rate Limiting Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLOUDFLARE    │    │     NGINX       │    │   APPLICATION   │
│  Rate Limiting  │    │  Rate Limiting  │    │  Rate Limiting  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ 100 req/min           │ 30 req/min            │ 10 req/min
         │ per IP                │ per IP                │ per IP
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Global DDoS    │    │  Burst Control  │    │ Business Logic  │
│  Protection     │    │  & Throttling   │    │  Rate Limiting  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📦 Deployment Architecture

### Container Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     DOCKER COMPOSE                             │
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   NGINX         │  │   BACKEND       │  │   FRONTEND      │ │
│  │   Container     │  │   Container     │  │   Container     │ │
│  │                 │  │                 │  │                 │ │
│  │ • nginx:alpine  │  │ • golang:1.21   │  │ • node:18       │ │
│  │ • Port 80/443   │  │ • Port 8080     │  │ • Port 3000     │ │
│  │ • SSL Certs     │  │ • Health Check  │  │ • Static Build  │ │
│  │ • Config Files  │  │ • Log Output    │  │ • Nginx Serve   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│           │                     │                     │        │
│           └─────────────────────┼─────────────────────┘        │
│                                 │                              │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    DOCKER NETWORK                           ││
│  │  • ogp-production network                                  ││
│  │  • Internal container communication                        ││
│  │  • DNS resolution between containers                       ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

### Infrastructure as Code

```
┌─────────────────────────────────────────────────────────────────┐
│                      TERRAFORM                                 │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                 SAKURA CLOUD                                ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │   Server    │  │   Disk      │  │  Network    │        ││
│  │  │ (1core-1gb) │  │  (20GB SSD) │  │  (Public IP)│        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  CLOUDFLARE                                 ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │ DNS Records │  │ SSL Proxy   │  │  Firewall   │        ││
│  │  │ (A, CNAME)  │  │ (TLS 1.3)   │  │   Rules     │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## 📊 Monitoring Architecture

### Observability Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                     MONITORING                                 │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    LOGGING                                  ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │Docker Logs  │  │System Logs  │  │Access Logs  │        ││
│  │  │ (stdout)    │  │ (journald)  │  │  (nginx)    │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   METRICS                                   ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │System Stats │  │ App Metrics │  │HTTP Metrics │        ││
│  │  │ (CPU, RAM)  │  │(Response    │  │(Status      │        ││
│  │  │             │  │ Time)       │  │ Codes)      │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                 HEALTH CHECKS                               ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        ││
│  │  │Cloudflare   │  │ Docker      │  │ Custom      │        ││
│  │  │Health Check │  │ Health      │  │ Scripts     │        ││
│  │  │             │  │ Check       │  │             │        ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘        ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## 🚀 Performance Architecture

### Caching Strategy

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   BROWSER       │    │   CLOUDFLARE    │    │   APPLICATION   │
│    CACHE        │    │     CACHE       │    │     CACHE       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ Static Assets         │ Static + API          │ In-Memory
         │ (24 hours)            │ (5 minutes)           │ Rate Limits
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ • HTML          │    │ • CSS/JS        │    │ • Client IPs    │
│ • CSS/JS        │    │ • Images        │    │ • Request Count │
│ • Images        │    │ • API Responses │    │ • Timestamps    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Scaling Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                     SCALING APPROACH                           │
│                                                                 │
│  Vertical Scaling (Current)                                    │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │  • Single Server Instance                                  ││
│  │  • Scale up CPU/Memory as needed                           ││
│  │  • Cost effective for moderate load                        ││
│  └─────────────────────────────────────────────────────────────┘│
│                                │                                │
│  Horizontal Scaling (Future)                                   │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │  • Multiple Server Instances                               ││
│  │  • Load Balancer Distribution                              ││
│  │  • Database for shared state                               ││
│  │  • Container orchestration (Kubernetes)                    ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## 📋 Technology Stack Summary

### Backend Technologies
- **Language**: Go 1.21
- **HTTP Framework**: Standard library + Gorilla Mux
- **HTML Parser**: golang.org/x/net/html
- **Container**: Docker with Alpine Linux
- **Build Tool**: Go modules

### Frontend Technologies
- **Framework**: React 18
- **Language**: TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Container**: Node.js 18 with Nginx

### Infrastructure Technologies
- **Cloud Provider**: Sakura Cloud VPS
- **CDN**: Cloudflare
- **Web Server**: Nginx
- **Containerization**: Docker + Docker Compose
- **IaC**: Terraform
- **CI/CD**: GitHub Actions

### Security Technologies
- **SSL/TLS**: Let's Encrypt + Cloudflare
- **Firewall**: UFW + Cloudflare WAF
- **Intrusion Detection**: fail2ban
- **DDoS Protection**: Cloudflare

### Monitoring Technologies
- **Logs**: Docker logs + journald
- **Health Checks**: Docker healthcheck + Cloudflare
- **Metrics**: Custom scripts + system monitoring
- **Alerting**: Email notifications

## 🔄 Future Architecture Considerations

### Potential Enhancements
1. **Database Integration**: PostgreSQL for persistent data
2. **Message Queue**: Redis for async processing
3. **Microservices**: Split into smaller services
4. **API Gateway**: Centralized API management
5. **Kubernetes**: Container orchestration
6. **Observability**: Prometheus + Grafana

### Scalability Roadmap
1. **Phase 1**: Current single-server setup
2. **Phase 2**: Add database and caching layer
3. **Phase 3**: Horizontal scaling with load balancer
4. **Phase 4**: Microservices architecture
5. **Phase 5**: Multi-region deployment

---

This architecture documentation should be updated as the system evolves and new components are added.