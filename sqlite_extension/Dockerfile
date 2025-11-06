# Multi-stage Dockerfile for Zemberek SQLite FTS5 Extension
# Produces a static binary that can run anywhere

# ============================================================================
# Stage 1: Builder - Compile the Go application with CGO
# ============================================================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    make \
    git

# Set working directory
WORKDIR /build/sqlite_extension

# Copy parent go.mod first (Zemberek root module)
COPY go.mod go.sum /build/

# Copy sqlite_extension module files
COPY sqlite_extension/go.mod sqlite_extension/go.sum ./

# Download dependencies
RUN cd /build && go mod download
RUN go mod download

# Copy Zemberek source (needed for morphology imports)
COPY morphology /build/morphology
COPY core /build/core

# Copy sqlite_extension source
COPY sqlite_extension .

# Copy Zemberek data files (lexicon, etc.)
COPY morphology/lexicon/data /build/data

# Build binaries for examples
RUN mkdir -p /build/bin && \
    CGO_ENABLED=1 go build \
    -tags "fts5" \
    -ldflags '-w -s' \
    -o /build/bin/basic \
    ./examples/basic/main.go

RUN CGO_ENABLED=1 go build \
    -tags "fts5" \
    -ldflags '-w -s' \
    -o /build/bin/search \
    ./examples/search/main.go

# ============================================================================
# Stage 2: Runtime - Minimal image with just the binary
# ============================================================================
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    sqlite-libs

# Create non-root user
RUN addgroup -g 1000 zemberek && \
    adduser -D -u 1000 -G zemberek zemberek

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/bin/* /app/
COPY --from=builder /build/data /app/data

# Set ownership
RUN chown -R zemberek:zemberek /app

# Switch to non-root user
USER zemberek

# Default command
ENTRYPOINT ["/app/basic"]

# Metadata
LABEL maintainer="Zemberek Project"
LABEL description="Zemberek Turkish Stemmer for SQLite FTS5"
LABEL version="1.0"
