# Stage 1: Build static Go binary
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mir-url-shortener ./cmd

# Stage 2: Minimal Alpine image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/mir-url-shortener .

EXPOSE 8080

CMD ["./mir-url-shortener"]
