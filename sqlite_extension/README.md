# Zemberek SQLite FTS5 Extension

Turkish language stemmer tokenizer for SQLite Full-Text Search (FTS5), powered by Zemberek morphology.

## Overview

This extension integrates Zemberek's Turkish morphological analysis into SQLite's FTS5 full-text search engine as a custom tokenizer. It enables proper stemming-based search for Turkish text, so searching for "kitap" will match "kitaplar", "kitapları", "kitabım", etc.

## Features

- **Morphological Stemming**: Uses Zemberek's advanced Turkish morphology for accurate stem extraction
- **FTS5 Native Integration**: Works as a built-in FTS5 tokenizer, no preprocessing needed
- **UTF-8 Aware**: Properly handles Turkish characters (ı, ğ, ş, ö, ü, ç) and byte offsets
- **Thread-Safe**: Safe for concurrent use across multiple connections
- **Easy to Use**: Simple drop-in replacement for standard SQLite driver

## Requirements

- Go 1.18 or higher
- CGO enabled (`CGO_ENABLED=1`)
- GCC compiler (macOS: Xcode Command Line Tools, Linux: build-essential)
- Zemberek data files (`lexicon.bin`)

## Installation

```bash
go get github.com/kalaomer/zemberek-go/sqlite_extension/driver
```

## Quick Start with Docker (Recommended)

The easiest way to use this extension is with Docker, which handles all dependencies:

```bash
# Build Docker image
make docker-build

# Run basic example
make docker-run-basic

# Run search example
make docker-run-search

# Build static binaries (works on any Linux)
make docker-static
# Output: ./dist/basic and ./dist/search
```

## Build from Source

The extension requires CGO and the FTS5 build tag:

### Local Build
```bash
CGO_ENABLED=1 go build -tags "fts5" your_program.go
```

### Using Docker
```bash
# Development shell
make docker-shell

# Run tests
make docker-test

# Complete workflow (build + test + static binaries)
make docker-all
```

## Usage

### Basic Example

```go
package main

import (
    "database/sql"
    "log"

    // Import the custom driver
    _ "github.com/kalaomer/zemberek-go/sqlite_extension/driver"
)

func main() {
    // Use "sqlite3_turkish" driver instead of "sqlite3"
    db, err := sql.Open("sqlite3_turkish", "mydb.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create FTS5 table with turkish_stem tokenizer
    _, err = db.Exec(`
        CREATE VIRTUAL TABLE documents USING fts5(
            title,
            content,
            tokenize='turkish_stem'
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Insert Turkish text
    _, err = db.Exec(
        "INSERT INTO documents (title, content) VALUES (?, ?)",
        "Kitaplar",
        "Kitapları okuyorum ve çok seviyorum",
    )

    // Search - "kitap" matches "Kitaplar" and "kitapları"
    rows, err := db.Query(
        "SELECT title, content FROM documents WHERE documents MATCH 'kitap'",
    )
    // Process results...
}
```

### Search Examples

The stemmer allows flexible searching:

```go
// Search for "kitap" matches:
// - kitap, kitaplar, kitapları, kitabım, kitabı, etc.
rows, _ := db.Query("SELECT * FROM docs WHERE docs MATCH 'kitap'")

// Search for "yaz" matches:
// - yaz, yazıyor, yazıyorum, yazdım, yazılım, yazı, etc.
rows, _ := db.Query("SELECT * FROM docs WHERE docs MATCH 'yaz'")

// Complex queries work too
rows, _ := db.Query("SELECT * FROM docs WHERE docs MATCH 'kitap OR yazılım'")
```

## How It Works

1. **Custom Driver**: Wraps `mattn/go-sqlite3` with a `ConnectHook` that registers the tokenizer
2. **FTS5 Registration**: On each connection, registers `turkish_stem` tokenizer via FTS5 API
3. **C Bridge**: C tokenizer callbacks delegate to Go's `StemTextWithPositions()`
4. **Morphological Analysis**: Zemberek analyzes each word and extracts the stem
5. **Token Emission**: Stems are returned to FTS5 with proper UTF-8 byte offsets

## Architecture

```
┌─────────────────────────────────────────────┐
│  Your Application (Go)                      │
├─────────────────────────────────────────────┤
│  database/sql with "sqlite3_turkish" driver │
├─────────────────────────────────────────────┤
│  driver.go (ConnectHook)                    │
│    └─> RegisterTurkishTokenizer()           │
├─────────────────────────────────────────────┤
│  registration.go (CGO)                      │
│    └─> C: registerZemberekTokenizer()       │
├─────────────────────────────────────────────┤
│  SQLite FTS5 Engine                         │
│    └─> turkish_stem tokenizer               │
├─────────────────────────────────────────────┤
│  zemberek_tokenizer.c                       │
│    └─> xTokenize() callback                 │
├─────────────────────────────────────────────┤
│  bridge.go (CGO)                            │
│    └─> goTokenizeText() [exported to C]     │
├─────────────────────────────────────────────┤
│  Zemberek Morphology (Go)                   │
│    └─> StemTextWithPositions()              │
└─────────────────────────────────────────────┘
```

## Performance

- **Initialization**: Morphology instance created once per process (lazy singleton)
- **Thread Safety**: Read-locked access to shared morphology instance
- **Memory**: ~50-100MB for lexicon data (loaded once, shared across connections)
- **Speed**: Comparable to other morphological FTS5 tokenizers

## Examples

See the [examples](examples/) directory:

- [basic](examples/basic/main.go): Simple usage demonstration
- [search](examples/search/main.go): Advanced search queries with Turkish text

### Run Examples Locally

```bash
cd examples/basic
CGO_ENABLED=1 go run -tags "fts5" main.go

cd examples/search
CGO_ENABLED=1 go run -tags "fts5" main.go
```

### Run Examples with Docker

```bash
# Run basic example
make docker-run-basic

# Run search example
make docker-run-search
```

### Build Static Binaries

Create standalone binaries that work on any Linux system:

```bash
make docker-static
```

This creates `./dist/basic` and `./dist/search` - statically linked binaries with no dependencies!

## Testing

### Local Testing

```bash
cd tokenizer
CGO_ENABLED=1 go test -tags "fts5" -v

cd ../driver
CGO_ENABLED=1 go test -tags "fts5" -v
```

### Docker Testing

```bash
# Run all tests in Docker
make docker-test

# Interactive testing in container
make docker-shell
go test -tags "fts5" -v ./...
```

## Troubleshooting

### "undefined: fts5_api"

**Solution**: Build with the FTS5 tag:
```bash
go build -tags "fts5"
```

### "C compiler not found"

**Solution**: Install GCC
- macOS: `xcode-select --install`
- Ubuntu/Debian: `apt-get install build-essential`
- Windows: Install MinGW or TDM-GCC

### "failed to register turkish_stem tokenizer"

**Causes**:
- FTS5 not enabled (missing `-tags "fts5"`)
- Missing lexicon.bin file
- SQLite version too old (need 3.20.0+)

**Solution**: Verify build tags and check lexicon.bin exists at expected path

### "CGO_ENABLED=0" errors

**Solution**: Enable CGO
```bash
export CGO_ENABLED=1
go build -tags "fts5"
```

## Deployment Options

### Option 1: Docker (Recommended)

Build a Docker image with your application:

```dockerfile
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=1 go build -tags "fts5" -o myapp

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/myapp /
ENTRYPOINT ["/myapp"]
```

### Option 2: Static Binary

Use our Docker setup to create a static binary:

```bash
make docker-static
# Deploy ./dist/basic or ./dist/search to any Linux server
```

### Option 3: Native Build

Build directly on your target platform:

```bash
CGO_ENABLED=1 go build -tags "fts5" -o myapp
```

## Limitations

- Requires CGO (cannot easily cross-compile without Docker)
- ~50-100MB memory overhead for lexicon
- Slower than simple ASCII tokenizers (due to morphological analysis)
- Docker recommended for consistent builds across platforms

## Contributing

Contributions welcome! Areas for improvement:

- Benchmark optimizations
- Additional language support
- Better error handling
- Documentation improvements

## License

Apache License 2.0 (same as Zemberek)

## Credits

- [Zemberek-NLP](https://github.com/ahmetaa/zemberek-nlp) - Original Java library by Ahmet A. Akın
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for Go by Yasuhiro Matsumoto
- [SQLite FTS5](https://www.sqlite.org/fts5.html) - Full-Text Search extension by SQLite team

## See Also

- [Zemberek-Go](../) - Parent project with morphology, tokenization, normalization
- [SQLite FTS5 Documentation](https://www.sqlite.org/fts5.html)
- [Custom Tokenizers Guide](https://www.sqlite.org/fts5.html#custom_tokenizers)
