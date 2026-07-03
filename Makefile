.PHONY: build test test-race wire lint run clean check-loc proto

BINARY_NAME=ml-elec
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) ./cmd/ml-elec

# Generate protobuf Go code
PROTOC=$(shell which protoc 2>/dev/null || echo ~/go/bin/protoc)
proto:
	@echo "Generating protobuf code..."
	$(PROTOC) --go_out=pkg/sdk/v1 --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/sdk/v1 --go-grpc_opt=paths=source_relative \
	       -I pkg/sdk/v1/proto \
	       pkg/sdk/v1/proto/lifecycle.proto \
	       pkg/sdk/v1/proto/sensor.proto

# Run all tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -race ./...

# Generate Wire code
wire:
	@echo "Generating Wire code..."
	wire ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/ coverage.out coverage.html

# Check LOC constraint (must be < 5000, excluding tests and generated code)
check-loc:
	@echo "Checking LOC constraint (< 5000)..."
	@LOC=$$(find . -name '*.go' ! -name '*_test.go' ! -path './generated/*' | xargs wc -l 2>/dev/null | tail -1 | awk '{print $$1}'); \
	if [ $$LOC -lt 5000 ]; then \
		echo "✓ LOC: $$LOC (< 5000)"; \
	else \
		echo "✗ LOC: $$LOC (exceeds 5000 limit)"; \
		exit 1; \
	fi
