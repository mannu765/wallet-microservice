.PHONY: test test-unit test-integration test-docker clean

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	go test -v ./... -tags=unit

# Run integration tests only (requires database)
test-integration:
	go test -v ./... -tags=integration

# Run tests in Docker containers
test-docker:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Run tests in Docker containers and clean up
test-docker-clean: test-docker
	docker-compose -f docker-compose.test.yml down -v

# Clean up Docker containers and volumes
clean:
	docker-compose -f docker-compose.test.yml down -v
	docker system prune -f

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	go test -race -v ./...

# Install test dependencies
test-deps:
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/mock
	go get github.com/stretchr/testify/suite

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Build the application
build:
	go build -o bin/wallet-microservice ./cmd/main.go

# Run the application
run:
	go run ./cmd/main.go
