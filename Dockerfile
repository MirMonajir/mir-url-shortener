# Stage 1: Build static Go binary
FROM golang:1.25.1-alpine3.21 AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

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

