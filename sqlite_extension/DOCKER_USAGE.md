# Docker Usage Guide

Complete guide for using Zemberek SQLite FTS5 extension with Docker.

## Quick Reference

```bash
# Build everything
make docker-build

# Run examples
make docker-run-basic
make docker-run-search

# Run tests
make docker-test

# Create static binaries
make docker-static

# Development shell
make docker-shell

# Complete workflow
make docker-all
```

---

## Detailed Workflows

### 1. First Time Setup

```bash
# Clone repository
git clone https://github.com/kalaomer/zemberek-go.git
cd zemberek-go/sqlite_extension

# Build Docker image
make docker-build
```

This creates a Docker image with:
- Go 1.24
- GCC compiler
- SQLite with FTS5 support
- All dependencies pre-installed

### 2. Running Examples

#### Basic Example
```bash
make docker-run-basic
```

Output:
```
✓ FTS5 table created successfully with turkish_stem tokenizer
✓ Test data inserted successfully

Searching for 'kitap':
  Found: Kitaplar - Kitapları okuyorum ve çok seviyorum

✓ Basic example completed successfully!
```

#### Search Example
```bash
make docker-run-search
```

Shows advanced search capabilities with Turkish stemming.

### 3. Testing

#### Run All Tests
```bash
make docker-test
```

This runs:
- Tokenizer unit tests
- Driver integration tests
- Coverage reports

#### Interactive Testing
```bash
# Open shell in container
make docker-shell

# Inside container, run specific tests
cd tokenizer
go test -tags "fts5" -v -run TestGoTokenizeText

# Or test everything
go test -tags "fts5" -v ./...
```

### 4. Building Static Binaries

Create standalone binaries that work on any Linux:

```bash
make docker-static
```

Output:
```
✓ Static binaries created in ./dist/
  - ./dist/basic
  - ./dist/search

These binaries are statically linked and can run on any Linux system!
```

#### Deploy to Server

```bash
# Copy to server
scp dist/basic user@server:/usr/local/bin/zemberek-basic
scp dist/search user@server:/usr/local/bin/zemberek-search

# Run on server (no dependencies needed!)
ssh user@server
zemberek-basic
```

### 5. Development Workflow

#### Edit and Test Cycle

```bash
# 1. Make changes to source code (locally)
vim tokenizer/bridge.go

# 2. Rebuild image
make docker-build

# 3. Run tests
make docker-test

# 4. Test manually
make docker-shell
```

#### Live Development

For faster iteration, mount source as volume:

```bash
docker-compose run --rm -v $(pwd):/build dev
```

Inside container:
```bash
# Your changes are live!
go test -tags "fts5" -v ./tokenizer
go build -tags "fts5" ./examples/basic
```

---

## Docker Compose Services

### Available Services

1. **dev** - Development shell
   ```bash
   docker-compose run --rm dev
   ```

2. **test** - Run tests
   ```bash
   docker-compose run --rm test
   ```

3. **basic** - Run basic example
   ```bash
   docker-compose run --rm basic
   ```

4. **search** - Run search example
   ```bash
   docker-compose run --rm search
   ```

5. **build-static** - Build static binaries
   ```bash
   docker-compose run --rm build-static
   ```

### Custom Commands

Run any command in the dev container:

```bash
# Build a specific example
docker-compose run --rm dev go build -tags "fts5" ./examples/basic

# Run benchmarks
docker-compose run --rm dev go test -tags "fts5" -bench=. ./tokenizer

# Check code coverage
docker-compose run --rm dev go test -tags "fts5" -cover ./...
```

---

## Advanced Usage

### Multi-Platform Builds

Build for multiple architectures:

```bash
# Build for ARM64 (e.g., Raspberry Pi)
docker buildx build --platform linux/arm64 -t zemberek-fts5:arm64 .

# Build for AMD64 (standard servers)
docker buildx build --platform linux/amd64 -t zemberek-fts5:amd64 .

# Build for both
docker buildx build --platform linux/amd64,linux/arm64 -t zemberek-fts5:latest .
```

### Custom SQLite Version

Edit `Dockerfile` to use specific SQLite version:

```dockerfile
# In builder stage
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev=3.45.0-r0  # Specific version
```

### Optimize Binary Size

Already optimized with:
```bash
-ldflags '-w -s -extldflags "-static"'
```

- `-w`: Remove DWARF debugging info
- `-s`: Remove symbol table
- `-extldflags "-static"`: Static linking

Result: ~5-6MB binary instead of 50MB+

### Production Dockerfile

For your own application:

```dockerfile
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY . .

# Build your app with Zemberek FTS5
RUN CGO_ENABLED=1 go build \
    -tags "fts5" \
    -ldflags '-w -s -extldflags "-static"' \
    -o /app/myapp \
    ./cmd/myapp

# Minimal runtime
FROM scratch
COPY --from=builder /app/myapp /
COPY --from=builder /app/data /data  # Zemberek data files
ENTRYPOINT ["/myapp"]
```

---

## Troubleshooting

### Build Fails: "gcc: not found"

**Solution:** Rebuild Docker image
```bash
docker-compose build --no-cache
```

### Tests Fail: "lexicon.bin not found"

**Solution:** Ensure data files are copied in Dockerfile
```dockerfile
COPY ../morphology/lexicon/data /build/data
```

### Static Binary Doesn't Run: "no such file or directory"

**Cause:** Missing shared libraries (not truly static)

**Solution:** Verify static build:
```bash
ldd dist/basic
# Should show: "not a dynamic executable"
```

If it shows libraries, rebuild with proper flags.

### Permission Denied in Container

**Solution:** Run as current user
```bash
docker-compose run --rm --user $(id -u):$(id -g) dev
```

### Out of Disk Space

**Solution:** Clean up Docker resources
```bash
make docker-clean
docker system prune -a
```

---

## Performance Tips

### Build Cache

Docker caches layers. Order matters:

```dockerfile
# ✓ Good: Dependencies first (cached)
COPY go.mod go.sum ./
RUN go mod download

# Then source code (changes frequently)
COPY . .
RUN go build ...
```

### Multi-Stage Efficiency

Our Dockerfile uses multi-stage builds:

- **Stage 1 (builder)**: 800MB+ with tools
- **Stage 2 (runtime)**: ~10MB with just binary

This keeps final image small.

### Volume Mounts

For development, mount source:

```yaml
volumes:
  - .:/app/sqlite_extension  # Live code changes
  - go-cache:/go/pkg         # Persistent cache
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Build and Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build Docker image
        run: make docker-build

      - name: Run tests
        run: make docker-test

      - name: Build static binaries
        run: make docker-static

      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: binaries
          path: dist/
```

### GitLab CI

```yaml
test:
  image: docker:latest
  services:
    - docker:dind
  script:
    - make docker-build
    - make docker-test

build:
  image: docker:latest
  services:
    - docker:dind
  script:
    - make docker-static
  artifacts:
    paths:
      - dist/
```

---

## Clean Up

### Remove Everything

```bash
# Stop and remove containers
docker-compose down

# Remove images
docker rmi zemberek-fts5

# Remove volumes
docker volume prune

# Or use make command
make docker-clean
```

### Keep Images, Remove Containers

```bash
docker-compose down
```

---

## Summary

| Task | Command |
|------|---------|
| Build | `make docker-build` |
| Test | `make docker-test` |
| Run basic | `make docker-run-basic` |
| Run search | `make docker-run-search` |
| Static binary | `make docker-static` |
| Dev shell | `make docker-shell` |
| Everything | `make docker-all` |
| Clean up | `make docker-clean` |

**Recommended workflow:**
1. `make docker-build` (once)
2. Make code changes
3. `make docker-test` (verify)
4. `make docker-static` (deploy)

Happy coding! 🚀
