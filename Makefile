.PHONY: build run test clean

# Default build target
build:
	go build -o bin/securevault ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Lint the code
lint:
	go vet ./...

# Run all checks
check: test lint