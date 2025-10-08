# Variables
APP_NAME=mir-url-shortener
DOCKER_COMPOSE=docker compose
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Run all the unit test cases with verbose output
test:
	@echo "Running test cases..."
	go test ./... -v

# Run gofmt and show diffs for format checking
fmt-check:
	@echo "Checking code formatting..."
	@gofmt -d $(GO_FILES) | tee /dev/stderr | (! read)

# Format all the Go files
fmt:
	@echo "Formatting all the Go code files..."
	gofmt -w $(GO_FILES)

# Run go vet for static analysis
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint for linter checking
lint:
	@echo "Running lint..."
	golangci-lint run

# Build the Go binary locally
build:
	@echo "Building binary..."
	go build -o $(APP_NAME) ./

# Start the url shortener app using Docker Compose
up:
	@echo "Starting Docker Compose..."
	$(DOCKER_COMPOSE) up -d

# Stop Docker Compose services
down:
	@echo "Stopping Docker Compose..."
	$(DOCKER_COMPOSE) down

# Full check: format, vet, lint, test
check: fmt-check vet lint test

.PHONY: test fmt fmt-check vet lint build up down check
