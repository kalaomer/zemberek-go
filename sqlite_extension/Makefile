.PHONY: all build test clean example

# Build variables
GO := go
GOFLAGS := -v
CGO_ENABLED := 1

# Targets
all: build

build:
	@echo "Building SQLite FTS5 extension..."
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) .

test:
	@echo "Running tests..."
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test $(GOFLAGS) -race .

example:
	@echo "Building example..."
	cd example && CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -o fts5_example .

run-example: example
	@echo "Running example..."
	cd example && ./fts5_example

clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f example/fts5_example
	rm -f example/*.db

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

vet:
	@echo "Vetting code..."
	$(GO) vet ./...

lint: fmt vet
	@echo "Code quality checks complete"

help:
	@echo "Available targets:"
	@echo "  all         - Build the extension (default)"
	@echo "  build       - Build the extension"
	@echo "  test        - Run tests"
	@echo "  example     - Build example program"
	@echo "  run-example - Build and run example"
	@echo "  clean       - Clean build artifacts"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  lint        - Run formatting and vetting"
	@echo "  help        - Show this help message"
