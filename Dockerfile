# Stage 1-> Build the Go binary
FROM golang:1.25.1 AS builder

WORKDIR /app

# Copy the Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy entire source code
COPY . .

# Build the Go binary from cmd/main.go
RUN go build -o mir-url-shortener ./cmd/main.go

# Stage 2-> Minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from builder stage
COPY --from=builder /app/mir-url-shortener .

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./mir-url-shortener"]
