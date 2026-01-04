# 🚀 Understanding the Problem
<strong>Techblog</strong> is an online web application that allow developers to make posts to share about the technology.

## Functional Requirements
### Core requirements
1. Users should be able to create and publish posts.
2. Users should be able to search posts.
3. Users should be able to view and react posts.

## Non-functional Requirements

| Requirement | Scenario | Requirement | Verification
|------------|-------------|--------|---------
| **Performance** | Homepage load (100 CCUs) | P95 < 500ms | k6 load test
| **Performance** | API single posts | P99 < 200ms | Grafana dashboard
| **Performance** | Image delivery | < 100ms via CDN | Cloudflare analytics
| **Availability & Reliability** | Production uptime | 99.5% monthly (3.6h downtime) | UptimeRobot
| **Availability & Reliability** | Graceful degradation | Read-only mode on DB failure | Health-check endpoint
| **Scalability** | Concurrent readers | 500 users sustained | k6 load test
| **Scalability** | Post storage | 10k posts, 1GB images | Disk monitoring
| **Scalability** | CPU/Memory | <80% under peak load> | Prometheus
| **Security** | Authentication | JWT expiry 24h + ratelimit 100reqs/min/IP | Postman + Burp Suite
| **Security** | Content input | XSS/SQL injection protection | OWASP ZAP scan
| **Security** | HTTPs | Enforced everywhere | SSL Labs test
| **Maintainability** | Deployment time | <10min via CI/CD | Github Action logs
| **Maintainability** | Error logging | 100% API errors captures | Sentry dashboard
| **Cost & Operations** | Monthly infrastructure | <$20 USD | AWS Billing
| **Cost & Operations** | Monitoring alerts | Latency > 400ms, error rate > 2% | Grafana alerts

# 🔑 High-Level Architecture

## Architecture Overview

The Techblog application follows a **monolithic architecture** with a Go backend service, leveraging external providers for CDN and object storage. The system is designed for cost-effectiveness, scalability, and high availability.

![Overall Architecture](/techblog.png)

## System Components

### 1. **Client Layer**
- **Web Application**: React/Vue.js SPA served via CDN
- **Mobile Application**: Native or React Native (future)
- **API Clients**: Third-party integrations

### 2. **CDN Layer (Cloudflare)**
- **Purpose**: Static asset delivery, image optimization, DDoS protection
- **Features**:
  - Global edge caching for static files
  - Image resizing and optimization
  - SSL/TLS termination
  - Rate limiting at edge
  - WAF (Web Application Firewall)

### 3. **Load Balancer**
- **Function**: Distribute traffic across multiple Go API instances
- **Features**:
  - Health checks
  - SSL termination
  - Session affinity (if needed)

### 4. **Go Monolith Service**
The core backend service handling all business logic:

#### **API Layer**
- RESTful API endpoints
- GraphQL (optional, for future)
- WebSocket support (for real-time reactions)

#### **Application Layer**
- **User Service**: Authentication, authorization, user management
- **Post Service**: CRUD operations, publishing, drafts
- **Search Service**: Full-text search, filtering, pagination
- **Reaction Service**: Like, comment, bookmark functionality
- **Media Service**: Image upload, processing, CDN integration

#### **Domain Layer**
- Business logic and domain models
- Validation rules
- Business rules enforcement

#### **Infrastructure Layer**
- **Database Repository**: PostgreSQL data access
- **Cache Repository**: Redis operations
- **Storage Repository**: S3 integration
- **External Services**: Email, notifications, etc.

### 5. **Data Layer**

#### **PostgreSQL Database**
- **Primary Database**: All structured data
- **Tables**:
  - `users`: User accounts, profiles
  - `posts`: Blog posts, metadata
  - `reactions`: Likes, comments, bookmarks
  - `tags`: Post categorization
  - `sessions`: JWT refresh tokens
- **Features**:
  - Full-text search (PostgreSQL tsvector)
  - Connection pooling
  - Read replicas (for scaling reads)

#### **Redis Cache**
- **Purpose**: Hot data caching, session storage
- **Use Cases**:
  - Post metadata cache (popular posts)
  - User session cache
  - Rate limiting counters
  - Search result cache
  - Real-time reaction counts

#### **S3 Object Storage**
- **Purpose**: Image and media file storage
- **Structure**:
  - `posts/{postId}/images/` - Post images
  - `users/{userId}/avatars/` - User avatars
  - `uploads/temp/` - Temporary uploads
- **Integration**: Direct upload from client to S3 (presigned URLs)

## Data Flow

### **Read Flow (Post Viewing)**
```
1. Client → CDN (check cache) → Go API
2. Go API → Redis (check cache)
3. If miss: Go API → PostgreSQL
4. Go API → Redis (cache result)
5. Go API → Client
```

### **Write Flow (Post Creation)**
```
1. Client → Go API (with auth)
2. Go API → Validate & Process
3. Go API → S3 (upload images via presigned URL)
4. Go API → PostgreSQL (save post)
5. Go API → Redis (invalidate cache)
6. Go API → Client (return post)
```

### **Image Upload Flow**
```
1. Client → Go API (request presigned URL)
2. Go API → Generate S3 presigned URL
3. Client → S3 (direct upload)
4. S3 → CDN (automatic sync)
5. Client → Go API (confirm upload, save post)
```

## Technology Stack

### **Backend**
- **Language**: Go 1.24+
- **Web Framework**: Gin/Echo/Chi (TBD)
- **ORM/Database**: GORM or sqlx
- **Cache**: go-redis
- **Storage**: AWS SDK for S3
- **Authentication**: golang-jwt/jwt

### **Database**
- **Primary**: PostgreSQL 15+
- **Cache**: Redis 7+

### **Infrastructure**
- **CDN**: Cloudflare
- **Object Storage**: AWS S3 (or compatible)
- **Container**: Docker
- **Orchestration**: Docker Compose (dev) / ECS/EKS (prod)

### **Monitoring & Observability**
- **Metrics**: Prometheus
- **Logging**: Structured JSON logs → CloudWatch/ELK
- **Error Tracking**: Sentry
- **APM**: Optional (Datadog, New Relic)
- **Uptime Monitoring**: UptimeRobot

### **CI/CD**
- **Version Control**: GitHub
- **CI/CD**: GitHub Actions
- **Container Registry**: Docker Hub / ECR

## Security Architecture

### **Security Layers**
1. **Edge Security (CDN)**
   - DDoS protection
   - WAF rules
   - Rate limiting

2. **Application Security**
   - JWT authentication (24h expiry)
   - Rate limiting (100 req/min per IP)
   - Input validation & sanitization
   - SQL injection prevention (parameterized queries)
   - XSS protection (content sanitization)

3. **Network Security**
   - HTTPS/TLS everywhere
   - VPC isolation (production)
   - Security groups

4. **Data Security**
   - Encrypted connections (TLS)
   - Encrypted at rest (S3, RDS)
   - Secrets management (environment variables / AWS Secrets Manager)

## Caching Strategy

### **Cache Layers**
1. **CDN Cache** (Cloudflare)
   - Static assets: 1 year
   - API responses: 5 minutes (public endpoints only)

2. **Application Cache** (Redis)
   - Post metadata: 15 minutes
   - Popular posts: 1 hour
   - User sessions: 24 hours
   - Search results: 5 minutes

### **Cache Invalidation**
- Write-through: Update cache on write
- Write-behind: Invalidate on write
- TTL-based: Automatic expiration

## Scalability Considerations

### **Horizontal Scaling**
- **Stateless API**: Multiple Go instances behind load balancer
- **Database**: Read replicas for read-heavy workloads
- **Cache**: Redis cluster (if needed)

### **Vertical Scaling**
- Start with single instance, scale up as needed
- Database: Upgrade instance size
- Cache: Increase Redis memory

### **Cost Optimization**
- Use CDN for static assets (reduces origin load)
- Redis caching (reduces database queries)
- S3 for images (cheaper than database storage)
- Auto-scaling based on metrics

## Deployment Architecture

### **Development**
```
Docker Compose:
- Go API (1 instance)
- PostgreSQL
- Redis
- Local S3 (MinIO)
```

### **Production**
```
AWS / Cloud Provider:
- ECS/EKS: Go API (2-3 instances, auto-scaling)
- RDS: PostgreSQL (Multi-AZ for HA)
- ElastiCache: Redis
- S3: Object storage
- CloudFront/Cloudflare: CDN
- ALB: Application Load Balancer
```

## Monitoring & Health Checks

### **Health Check Endpoint**
- `/health`: Basic health check
- `/health/ready`: Readiness probe (DB connection)
- `/health/live`: Liveness probe

### **Metrics Endpoints**
- `/metrics`: Prometheus metrics
- Track: Request latency, error rates, cache hit rates, DB connection pool

### **Alerting**
- Latency > 400ms (P95)
- Error rate > 2%
- Database connection failures
- Cache hit rate < 80%

## API Design

### **RESTful Endpoints**
```
Authentication:
  POST   /api/v1/auth/register
  POST   /api/v1/auth/login
  POST   /api/v1/auth/refresh
  POST   /api/v1/auth/logout

Posts:
  GET    /api/v1/posts              # List posts (paginated, searchable)
  GET    /api/v1/posts/:id          # Get single post
  POST   /api/v1/posts              # Create post (auth required)
  PUT    /api/v1/posts/:id          # Update post (auth required)
  DELETE /api/v1/posts/:id          # Delete post (auth required)

Reactions:
  POST   /api/v1/posts/:id/like     # Like post
  POST   /api/v1/posts/:id/comment  # Add comment
  GET    /api/v1/posts/:id/comments # Get comments

Media:
  POST   /api/v1/media/upload-url   # Get presigned S3 URL
  GET    /api/v1/media/:id          # Get media metadata

Health:
  GET    /health                    # Health check
  GET    /metrics                   # Prometheus metrics
```

## Future Enhancements

1. **Microservices Migration** (if needed)
   - Split by domain (User, Post, Search services)
   - Message queue (RabbitMQ/Kafka) for async processing

2. **Advanced Features**
   - Real-time notifications (WebSocket)
   - Full-text search (Elasticsearch)
   - Analytics service
   - Recommendation engine

3. **Performance**
   - Database read replicas
   - Redis cluster
   - GraphQL API
   - Edge computing (Cloudflare Workers)
