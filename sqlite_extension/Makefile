.PHONY: all build test clean examples run-basic run-search help docker-build docker-test docker-run-basic docker-run-search docker-static docker-shell

# Default target
all: build

# Build tags
BUILD_TAGS = fts5
CGO_ENABLED = 1

# Go build flags
GO_BUILD_FLAGS = -tags "$(BUILD_TAGS)"
GO_TEST_FLAGS = -tags "$(BUILD_TAGS)" -v

help:
	@echo "Zemberek SQLite FTS5 Extension - Build Commands"
	@echo ""
	@echo "Local Build:"
	@echo "  make build        - Build all packages"
	@echo "  make test         - Run all tests"
	@echo "  make examples     - Build all examples"
	@echo "  make run-basic    - Run basic example"
	@echo "  make run-search   - Run search example"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Docker Build (Recommended):"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-test        - Run tests in Docker"
	@echo "  make docker-run-basic   - Run basic example in Docker"
	@echo "  make docker-run-search  - Run search example in Docker"
	@echo "  make docker-static      - Build static binaries (output: ./dist/)"
	@echo "  make docker-shell       - Open shell in dev container"
	@echo ""
	@echo "Requirements (Local):"
	@echo "  - CGO_ENABLED=1"
	@echo "  - GCC compiler"
	@echo "  - Go 1.18+"
	@echo ""
	@echo "Requirements (Docker):"
	@echo "  - Docker & Docker Compose"

# Build all packages
build:
	@echo "Building tokenizer package..."
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) ./tokenizer
	@echo "Building driver package..."
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) ./driver
	@echo "✓ Build successful"

# Run tests
test: test-tokenizer test-driver

test-tokenizer:
	@echo "Running tokenizer tests..."
	cd tokenizer && CGO_ENABLED=$(CGO_ENABLED) go test $(GO_TEST_FLAGS)

test-driver:
	@echo "Running driver tests..."
	cd driver && CGO_ENABLED=$(CGO_ENABLED) go test $(GO_TEST_FLAGS)

# Build examples
examples: build-basic build-search

build-basic:
	@echo "Building basic example..."
	cd examples/basic && CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o basic main.go
	@echo "✓ Basic example built: examples/basic/basic"

build-search:
	@echo "Building search example..."
	cd examples/search && CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o search main.go
	@echo "✓ Search example built: examples/search/search"

# Run examples
run-basic: build-basic
	@echo "Running basic example..."
	@echo "────────────────────────────────────────"
	cd examples/basic && ./basic

run-search: build-search
	@echo "Running search example..."
	@echo "────────────────────────────────────────"
	cd examples/search && ./search

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f examples/basic/basic
	rm -f examples/search/search
	go clean -cache -testcache
	@echo "✓ Clean complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies ready"

# Check environment
check:
	@echo "Checking build environment..."
	@echo "CGO_ENABLED: $${CGO_ENABLED:-not set}"
	@which gcc > /dev/null && echo "GCC: $$(gcc --version | head -1)" || echo "GCC: not found"
	@echo "Go version: $$(go version)"
	@echo ""
	@if [ "$${CGO_ENABLED}" != "1" ]; then \
		echo "⚠️  Warning: CGO_ENABLED is not set to 1"; \
		echo "   Run: export CGO_ENABLED=1"; \
	fi
	@which gcc > /dev/null || echo "⚠️  Warning: GCC not found. Install with: xcode-select --install (macOS)"

# Format code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "✓ Format complete"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Install dependencies and build
install: deps build
	@echo "✓ Installation complete"

# Full test with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	cd tokenizer && CGO_ENABLED=$(CGO_ENABLED) go test $(GO_TEST_FLAGS) -cover
	cd driver && CGO_ENABLED=$(CGO_ENABLED) go test $(GO_TEST_FLAGS) -cover

# ============================================================================
# Docker Commands
# ============================================================================

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker-compose build
	@echo "✓ Docker image built successfully"

# Run tests in Docker
docker-test:
	@echo "Running tests in Docker..."
	docker-compose run --rm test
	@echo "✓ Tests completed"

# Run basic example in Docker
docker-run-basic:
	@echo "Running basic example in Docker..."
	docker-compose run --rm basic

# Run search example in Docker
docker-run-search:
	@echo "Running search example in Docker..."
	docker-compose run --rm search

# Build static binaries in Docker and extract to ./dist/
docker-static:
	@echo "Building static binaries in Docker..."
	@mkdir -p dist
	docker-compose run --rm build-static
	@echo ""
	@echo "✓ Static binaries created in ./dist/"
	@echo "  - ./dist/basic"
	@echo "  - ./dist/search"
	@echo ""
	@echo "These binaries are statically linked and can run on any Linux system!"

# Open interactive shell in dev container
docker-shell:
	@echo "Opening shell in development container..."
	docker-compose run --rm dev

# Clean Docker resources
docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	@echo "✓ Docker resources cleaned"

# Full Docker workflow: build, test, create static binaries
docker-all: docker-build docker-test docker-static
	@echo "✓ Complete Docker workflow finished!"
