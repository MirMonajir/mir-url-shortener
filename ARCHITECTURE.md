# Architecture & Design

## Overview

The URL Shortener is built with a layered, domain-driven architecture following SOLID principles for maintainability, testability, and scalability.

```
┌─────────────────────────────────────────────────────┐
│              HTTP Layer (REST API)                  │
│  • Handler: HTTP request/response processing       │
│  • Middleware: Security, CORS, logging              │
│  • Error: Standardized error responses              │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│           Application Layer                         │
│  • ShortenerService: Business logic                │
│  • Validation: Input validation                     │
│  • Middleware: Cross-cutting concerns              │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│            Domain Layer (Models)                    │
│  • URL: Domain model                               │
│  • Shortener: Domain interface                      │
│  • Storage: Storage abstraction                     │
│  • Errors: Domain error types                       │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│         Infrastructure Layer                       │
│  • InMemoryStore: Storage implementation           │
│  • Config: Configuration management                │
│  • Database: Future: SQL/NoSQL implementations     │
└─────────────────────────────────────────────────────┘
```

## Design Patterns

### 1. Domain-Driven Design (DDD)
- **Domain Layer**: Contains core business logic and models
  - `domain/url.go`: URL entity
  - `domain/shortener.go`: Shortener interface
  - `domain/storage.go`: Storage abstraction
  - `domain/errors.go`: Domain-specific errors
  - `domain/validator.go`: Business validation rules

- **Application Layer**: Orchestrates domain logic
  - `application/shortener_service_uc.go`: Use case implementation
  - `application/handler.go`: HTTP endpoint handlers
  - `application/middleware.go`: Cross-cutting concerns

- **Infrastructure Layer**: Technical implementations
  - `infrastructure/inmemoryStore.go`: Storage implementation
  - `infrastructure/config.go`: Configuration handling

### 2. Dependency Injection
- Services depend on interfaces, not concrete implementations
- `Storage` interface allows swapping implementations
- Makes testing easier with mock implementations

### 3. Middleware Pattern
```
Request
  ↓
RecoveryMiddleware (panic handling)
  ↓
RequestLoggerMiddleware (request logging)
  ↓
ErrorLoggerMiddleware (error logging)
  ↓
SecurityHeadersMiddleware (security headers)
  ↓
CORSMiddleware (CORS handling)
  ↓
Router Handler
  ↓
Response
```

### 4. Error Handling
- Structured error types with error codes
- HTTP status codes match semantic meaning
- Errors don't expose internal implementation details
- Validation errors include field-level details

### 5. Configuration Management
- All configuration from environment variables
- Validation on startup (fail-fast)
- Type-safe configuration struct
- Sensible defaults for all values

## Key Components

### Handler (`cmd/handler.go`)
- Routes HTTP requests to service methods
- Validates input
- Formats responses
- Handles errors with proper HTTP codes

### ShortenerService (`application/shortener_service_uc.go`)
- Orchestrates URL shortening workflow
- Validates URLs
- Persists to storage
- Returns formatted URLs

### InMemoryStore (`infrastructure/inmemoryStore.go`)
- Maps original URLs to shortened codes
- Maintains bidirectional lookup
- Tracks domain metrics
- Thread-safe with RWMutex
- Optimized for fast lookups

### URLValidator (`domain/validator.go`)
- Validates URL format
- Checks for SSRF attempts
- Enforces length constraints
- Ensures valid schemes (HTTP/HTTPS)

## API Endpoints

### POST `/shortenurl`
Shortens a URL.

**Request:**
```json
{
  "url": "https://www.example.com/very/long/path"
}
```

**Response (Success - 200):**
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

**Response (Error - 400):**
```json
{
  "error_type": "invalid_url",
  "message": "URL must be a valid HTTP/HTTPS URL",
  "details": "Provided scheme: ftp",
  "code": 400
}
```

### GET `/:shortenedurl`
Redirects to original URL.

**Response (Success - 302):**
- Redirect Location header set to original URL

**Response (Not Found - 404):**
```json
{
  "error_type": "url_not_found",
  "message": "Shortened URL not found",
  "code": 404
}
```

### GET `/appmetrics`
Returns top domains by usage.

**Response:**
```json
{
  "top_domains": {
    "example.com": 42,
    "google.com": 28,
    "github.com": 15
  }
}
```

## Security Architecture

### Input Validation
1. **Schema Validation**: JSON binding checks required fields
2. **Format Validation**: URL must be valid format
3. **Scheme Validation**: Only HTTP/HTTPS allowed
4. **SSRF Prevention**: Blocks private IPs and localhost
5. **Length Validation**: Enforces max URL length

### Output Security
1. **Error Sanitization**: No internal details in errors
2. **Security Headers**: CSP, X-Frame-Options, etc.
3. **CORS Control**: Whitelisted origins only
4. **HTTPS Enforcement**: HSTS header for secure transport

### Runtime Security
1. **Panic Recovery**: Middleware catches panics
2. **Concurrency Safety**: Mutex protection for shared state
3. **Race Detection**: All tests run with `-race` flag
4. **Type Safety**: Leverages Go's type system

## Concurrency Model

### Thread Safety
- `InMemoryStore` uses `sync.RWMutex` for concurrent access
- Read operations use `RLock()` for parallelism
- Write operations use `Lock()` for exclusivity
- No goroutine leaks or deadlocks

### Performance Characteristics
- **Read**: O(1) HashMap lookup
- **Write**: O(1) with mutex
- **Metrics**: O(n log n) sort on every query (can be optimized)

## Testing Strategy

### Unit Tests
- Validator tests: 11 test cases covering edge cases
- Handler tests: HTTP request/response validation
- Service tests: Business logic verification
- Config tests: Configuration validation
- Storage tests: Concurrency and correctness

### Integration Tests
- End-to-end shorten/resolve workflow
- Error handling paths
- Middleware stack

### Quality Assurance
- Race condition detection: `go test -race`
- Code coverage: >80% target
- Linting: golangci-lint
- Format checking: gofmt -d

## Scalability Considerations

### Current Limitations
- **In-Memory Only**: Data lost on restart
- **Single Instance**: No distributed state
- **No Persistence**: No durability guarantees
- **No Caching**: Every request hits memory

### Future Improvements
1. **Persistent Storage**: PostgreSQL/MongoDB backend
2. **Distributed Cache**: Redis for hot URLs
3. **Load Balancing**: Multiple instances with shared state
4. **API Keys**: Authentication and rate limiting
5. **Analytics**: Advanced metrics and reporting
6. **Custom Domains**: Branded short URLs

## Dependencies

### Core
- **gin-gonic/gin**: REST API framework
- **golang.org/x/net**: URL/domain utilities
- **golang.org/x/crypto**: Future cryptography

### Testing
- **stretchr/testify**: Assertions and mocking

### Utilities
- **joho/godotenv**: Environment file loading

See [DEPENDENCY_POLICY.md](DEPENDENCY_POLICY.md) for detailed information.

## Best Practices Implemented

✅ **Code Organization**
- Clear separation of concerns (domain/application/infrastructure)
- Logical package structure matching responsibilities

✅ **Error Handling**
- Structured error types
- Proper HTTP status codes
- Non-leaking error messages

✅ **Testing**
- Unit test coverage >80%
- Race condition detection enabled
- Mock implementations for dependencies

✅ **Security**
- SSRF prevention
- Input validation
- Security headers
- Non-root container execution

✅ **Configuration**
- Environment-based configuration
- Validation on startup
- Sensible defaults

✅ **Documentation**
- Code comments for non-obvious logic
- API documentation
- Architecture diagrams

---

## References

- [Domain-Driven Design](https://www.domainlanguage.com/ddd/)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Effective Go](https://golang.org/doc/effective_go)
- [Gin Web Framework](https://gin-gonic.com/)

---

**Version**: 1.0.0  
**Last Updated**: September 2026  
**Status**: Stable
