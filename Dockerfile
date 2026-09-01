# Stage 1: Build static Go binary
# Use Go 1.23 to match go.mod and avoid builder/runtime mismatches
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
# Ensure Go can fetch modules reliably using the default proxy and checksum DB
# Set the env vars explicitly and retry module download to handle transient CI network/sumdb issues
RUN set -eux; \
    go env -w GOPROXY="https://proxy.golang.org,direct"; \
    go env -w GOSUMDB="sum.golang.org"; \
    # Try a few times to download modules (handles transient network issues)
    n=0; until [ "$n" -ge 3 ] || go mod download; do n=$((n+1)); echo "go mod download attempt $n failed, retrying..."; sleep $((n*2)); done; \
    if ! go list -m all >/dev/null 2>&1; then \
        echo "Module download/verification failed after retries; retrying with GOSUMDB=off (fallback)"; \
        go env -w GOSUMDB="off"; \
        go mod download; \
    fi

COPY . .

# Build with static linking and security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o mir-url-shortener ./cmd

# Stage 2: Minimal Alpine image with security hardening
FROM alpine:3.21

# Install CA certificates and wget for healthcheck
RUN apk add --no-cache ca-certificates tzdata wget

# Create non-root user for container execution
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

WORKDIR /home/app

# Copy binary from builder
COPY --from=builder /app/mir-url-shortener .

# Set ownership and permissions
RUN chown -R app:app /home/app && \
    chmod 755 mir-url-shortener

# Add security labels
LABEL maintainer="MirMonajir" \
      description="URL Shortener Service - High-performance Go REST API" \
      version="1.0.0"

# Switch to non-root user
USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD test -z "$URL" || wget --quiet --tries=1 --spider http://localhost:8080/appmetrics || exit 1

CMD ["./mir-url-shortener"]

