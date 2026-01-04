# Techblog Backend

A high-performance Go backend service for a tech blog platform that allows developers to create, publish, and share technology posts. Built with clean architecture principles, featuring graceful shutdown, comprehensive logging, and production-ready configurations.

## 🚀 Features

- **RESTful API** - Clean and well-structured API endpoints
- **JWT Authentication** - Secure token-based authentication with refresh tokens
- **Rate Limiting** - Configurable rate limiting (100 requests/minute per IP)
- **CORS Support** - Cross-origin resource sharing configuration
- **Graceful Shutdown** - Context-based graceful shutdown handling
- **Structured Logging** - Logrus-based logging with environment-specific formatting
- **Database Support** - PostgreSQL with connection pooling
- **Redis Support** - Caching and session management
- **Configuration Management** - YAML-based configuration with environment support
- **Docker Support** - Containerized deployment ready

## 🛠 Tech Stack

- **Language**: Go 1.24.4
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Logging**: [Logrus](https://github.com/sirupsen/logrus)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Container**: Docker

## 📁 Project Structure

```
be-techblog/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                   # Configuration management
│   ├── logger/                   # Logging utilities
│   ├── server/                   # HTTP server setup
│   ├── handlers/                 # HTTP request handlers
│   ├── services/                 # Business logic layer
│   ├── repositories/             # Data access layer
│   ├── daos/                     # Data access objects
│   ├── models/                   # Domain models
│   └── middlewares/              # HTTP middlewares
│       ├── auth_middleware.go
│       ├── logger_middleware.go
│       ├── ratelimit_middleware.go
│       ├── security_middleware.go
│       └── recovery_middleware.go
├── configs/                      # Configuration files
│   ├── env.local.yaml
│   └── sample.config
├── migrations/                   # Database migrations
├── docs/                         # Documentation
│   └── architecture/
│       └── SystemDesign.md
├── scripts/                      # Utility scripts
├── Dockerfile                    # Docker configuration
├── Makefile                      # Build automation
├── go.mod                        # Go dependencies
└── README.md                     # This file
```

## 📋 Prerequisites

- Go 1.24.4 or higher
- PostgreSQL 12+ (for database)
- Redis 6+ (for caching)
- Make (optional, for using Makefile commands)

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/datpham2001/be-techblog.git
   cd be-techblog
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up configuration**
   ```bash
   cp configs/sample.config configs/env.local.yaml
   # Edit configs/env.local.yaml with your settings
   ```

4. **Set up database**
   ```bash
   # Create database
   createdb techblog
   
   # Run migrations (if available)
   # make migrate-up
   ```

## ⚙️ Configuration

Configuration is managed through YAML files in the `configs/` directory. Copy `sample.config` to `env.local.yaml` and customize:

### Key Configuration Sections

- **Server**: Host, port, environment, TLS settings
- **Database**: PostgreSQL connection settings
- **Redis**: Redis connection settings
- **JWT Auth**: Secret key, token expiration times
- **Rate Limit**: Request limits and periods
- **CORS**: Allowed origins, methods, and headers

Example configuration:
```yaml
server:
  env: development
  host: "0.0.0.0"
  port: "8080"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  db_name: "techblog"
```

## 🚀 Running the Application

### Development Mode

```bash
# Run directly
go run cmd/app/main.go

# Or using Make (if Makefile has run target)
make run
```

### Production Mode

```bash
# Build the application
go build -o bin/techblog cmd/app/main.go

# Run the binary
./bin/techblog
```

### Using Docker

```bash
# Build Docker image
docker build -t techblog:latest .

# Run container
docker run -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  techblog:latest
```

## 📡 API Endpoints

API documentation will be available at `/swagger` (when Swagger is configured).

### Example Endpoints (to be implemented)

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/posts` - List posts
- `POST /api/v1/posts` - Create post
- `GET /api/v1/posts/:id` - Get post by ID
- `PUT /api/v1/posts/:id` - Update post
- `DELETE /api/v1/posts/:id` - Delete post

## 🧪 Development

### Code Structure

The project follows clean architecture principles:

- **Handlers**: Handle HTTP requests/responses
- **Services**: Business logic implementation
- **Repositories**: Data persistence abstraction
- **Models**: Domain entities
- **Middlewares**: Cross-cutting concerns (auth, logging, rate limiting)

### Adding New Features

1. Define models in `internal/models/`
2. Create repository interface in `internal/repositories/`
3. Implement repository in `internal/daos/`
4. Create service in `internal/services/`
5. Create handler in `internal/handlers/`
6. Register routes in `internal/server/server.go`

## 🧹 Graceful Shutdown

The application implements graceful shutdown using context:

- Listens for `SIGINT` and `SIGTERM` signals
- Stops accepting new connections
- Waits for existing requests to complete (10-second timeout)
- Closes database and Redis connections
- Exits cleanly

## 📝 Logging

Logging is configured based on environment:

- **Development**: Colored text output with timestamps
- **Production**: JSON format for log aggregation

Log levels can be configured in the logger initialization.

## 🔒 Security Features

- JWT-based authentication
- Rate limiting per IP
- CORS configuration
- Security headers middleware
- Input validation
- SQL injection protection (via parameterized queries)

## 📊 Performance Requirements

- Homepage load (100 CCUs): P95 < 500ms
- API single posts: P99 < 200ms
- Image delivery: < 100ms via CDN
- Concurrent readers: 500 users sustained

## 🐳 Docker

The project includes a `Dockerfile` for containerized deployment. Build and run:

```bash
docker build -t techblog .
docker run -p 8080:8080 techblog
```

## 📚 Documentation

- [System Design](./docs/architecture/SystemDesign.md) - Detailed architecture documentation

## 📄 License

This project is licensed under the MIT License.