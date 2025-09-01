# URL Shortener

A robust, scalable URL shortening service built with Go, following clean architecture principles and onion architecture design patterns.

## 🏗️ Architecture

This project follows a **multi-layered onion architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                        │
├─────────────────────────────────────────────────────────────┤
│                   Handler Layer                            │
├─────────────────────────────────────────────────────────────┤
│                   Service Layer                            │
├─────────────────────────────────────────────────────────────┤
│                  Repository Layer                          │
├─────────────────────────────────────────────────────────────┤
│                    Data Layer (Redis)                      │
└─────────────────────────────────────────────────────────────┘
```

### Layers Explained

- **HTTP Layer**: Gin web framework for routing and HTTP handling
- **Handler Layer**: HTTP request/response handling, input validation
- **Service Layer**: Business logic, orchestration
- **Repository Layer**: Data access abstraction with interfaces
- **Data Layer**: Redis storage implementation

## 🚀 Features

- ✅ **Clean Architecture**: Follows SOLID principles and dependency injection
- ✅ **URL Shortening**: Generate unique short URLs from long URLs
- ✅ **User Isolation**: URLs are scoped to specific users
- ✅ **Redis Storage**: Fast, in-memory storage with persistence
- ✅ **RESTful API**: Clean HTTP endpoints
- ✅ **Error Handling**: Comprehensive error handling and validation
- ✅ **Testable**: Easy to unit test with mock interfaces

## 📁 Project Structure

```
url-shortner/
├── main.go                 # Application entry point
├── handler/
│   └── handlers.go        # HTTP request handlers
├── service/
│   ├── url_service.go     # Business logic layer
│   └── storage_service.go # Storage service implementation
│   ├── interfaces.go      # Repository interfaces
│   └── storage_service.go # Redis implementation
├── shortener/
│   └── shortener.go       # URL shortening algorithm
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

## 🛠️ Prerequisites

- **Go 1.24+** - [Download Go](https://golang.org/dl/)
- **Redis** - [Download Redis](https://redis.io/download)
- **Git** - For cloning the repository

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd url-shortner
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Start Redis

```bash
# macOS (using Homebrew)
brew services start redis

# Linux
sudo systemctl start redis

# Windows
redis-server
```

### 4. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:9808`

## 📡 API Endpoints

### Create Short URL

```http
POST /create-short-url
Content-Type: application/json

{
  "long_url": "https://example.com/very-long-url-that-needs-shortening",
  "user_id": "user-123"
}
```

**Response:**

```json
{
  "message": "short url created successfully",
  "short_url": "http://localhost:9808/Jsz4k57oAX"
}
```

### Redirect to Original URL

```http
GET /:shortUrl
```

**Example:** `GET /Jsz4k57oAX` → Redirects to the original long URL

### Health Check

```http
GET /
```

**Response:**

```json
{
  "message": "Welcome to the URL Shortener API"
}
```

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./store
go test ./service
```

## 🔧 Configuration

### Redis Configuration

Default Redis settings in `store/storage_service.go`:

```go
redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",  // Redis server address
    Password: "",                 // Redis password (if any)
    DB:       0,                 // Redis database number
})
```

### Cache Duration

URL mappings expire after 6 hours by default:

```go
const CacheDuration = 6 * time.Hour
```

## 🏗️ Design Principles

### 1. **Dependency Inversion**

- High-level modules don't depend on low-level modules
- Both depend on abstractions (interfaces)

### 2. **Single Responsibility**

- Each layer has one reason to change
- Handlers handle HTTP, services handle business logic

### 3. **Interface Segregation**

- `URLRepository` interface defines only what's needed
- Easy to implement different storage backends

### 4. **Open/Closed Principle**

- Open for extension (new storage backends)
- Closed for modification (existing code)

## 🔄 Adding New Features

### Adding a New Storage Backend

1. Implement the `URLRepository` interface
2. Update dependency injection in `main.go`
3. No changes needed in service or handler layers

### Adding New Business Logic

1. Add methods to the `URLService`
2. Update handlers to use new service methods
3. Maintain separation of concerns

## 🚀 Deployment

### Docker (Recommended)

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 9808
CMD ["./main"]
```

### Environment Variables

```bash
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=
export REDIS_DB=0
export SERVER_PORT=9808
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues:

1. Check the Redis connection
2. Verify all dependencies are installed
3. Check the logs for error messages
4. Open an issue with detailed error information

## 🔮 Future Enhancements

- [ ] **Rate Limiting**: Prevent abuse
- [ ] **Analytics**: Track URL clicks and usage
- [ ] **Custom Domains**: Support for custom short domains
- [ ] **Authentication**: User management and security
- [ ] **Database Backend**: PostgreSQL/MySQL support
- [ ] **Monitoring**: Prometheus metrics and health checks
- [ ] **Load Balancing**: Multiple Redis instances
- [ ] **Caching**: Additional caching layers

---

Built with ❤️ using Go and clean architecture principles.
