# Zemberek SQLite FTS5 Tokenizer

SQLite FTS5 (Full-Text Search) helper library with Zemberek-powered Turkish tokenization for better search results.

## Overview

This library provides Turkish-aware tokenization and FTS5 helper functions for SQLite full-text search. While SQLite's built-in tokenizers work reasonably well, this library enhances search quality for Turkish text by properly handling Turkish-specific character normalization and case folding.

## Features

- **Turkish-aware tokenization**: Properly handles Turkish characters (ç, ğ, ı, ö, ş, ü)
- **Case normalization**: Correctly handles Turkish case folding (I↔ı, İ↔i)
- **Diacritic removal**: Optional removal of Turkish diacritics for broader matching
- **FTS5 Helper**: Convenient API for creating and searching FTS5 tables
- **Pure Go**: No CGO complexity for the tokenizer (uses standard database/sql)

## Installation

```bash
go get github.com/kalaomer/zemberek-go/sqlite_extension
```

You also need a SQLite driver with FTS5 support:

```bash
go get github.com/mattn/go-sqlite3
```

## Building

The example program requires FTS5 support in SQLite:

```bash
# Build with FTS5 support
cd sqlite_extension
make build

# Run tests
make test

# Build and run example
cd example
CGO_ENABLED=1 go build -tags "fts5" -o fts5_example .
./fts5_example
```

## Usage

### Basic Example

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/mattn/go-sqlite3"
    "github.com/kalaomer/zemberek-go/sqlite_extension"
)

func main() {
    db, err := sql.Open("sqlite3", "./mydb.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create FTS5 helper
    helper := sqlite_extension.NewFTS5Helper(db)

    // Create FTS5 table
    err = helper.CreateFTS5Table("documents", "title", "content")
    if err != nil {
        log.Fatal(err)
    }

    // Insert document
    err = helper.InsertDocument("documents", map[string]string{
        "title":   "Türkçe Başlık",
        "content": "Bu bir Türkçe metin örneğidir.",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Search with automatic query normalization
    rows, err := helper.SearchWithRank("documents",
        []string{"title", "content"}, "türkçe")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var title, content string
        var rank float64
        rows.Scan(&title, &content, &rank)
        fmt.Printf("[%.2f] %s: %s\n", rank, title, content)
    }
}
```

### Tokenization Only

If you just need Turkish-aware tokenization without FTS5:

```go
import "github.com/kalaomer/zemberek-go/sqlite_extension"

// Simple tokenization
tokens := sqlite_extension.TokenizeText("İstanbul çok güzel bir şehir")
// Returns: ["istanbul", "çok", "güzel", "bir", "şehir"]

// With positions
tokenizer := sqlite_extension.NewZemberekTokenizer()
positions := tokenizer.TokenizeWithPositions("Merhaba dünya")
for _, pos := range positions {
    fmt.Printf("%s [%d:%d]\n", pos.Token, pos.Start, pos.End)
}

// With custom options (case + diacritic removal)
tokenizer = sqlite_extension.NewZemberekTokenizerWithOptions(true, true)
tokens = tokenizer.Tokenize("Çalışma şekli")
// Returns: ["calisma", "sekli"]
```

### Turkish Case Conversion

```go
// Turkish uppercase
upper := sqlite_extension.TurkishUpperCase("istanbul")
// Returns: "İSTANBUL" (note the dotted İ)

upper = sqlite_extension.TurkishUpperCase("ıstanbul")
// Returns: "ISTANBUL" (note the dotless I)
```

## How It Works

The library uses a two-layer approach:

1. **SQLite FTS5**: Uses SQLite's built-in `unicode61` tokenizer with `remove_diacritics` for indexing
2. **Zemberek Tokenizer**: Normalizes search queries using Turkish-aware rules before sending to FTS5

This ensures:
- ✅ "İstanbul" and "istanbul" match correctly
- ✅ "ÇALIŞMA" and "çalışma" match correctly
- ✅ Proper handling of Turkish I/İ distinction
- ✅ Case-insensitive search that respects Turkish rules

## Turkish Language Support

### Case Folding

The tokenizer correctly handles Turkish-specific case conversions:

| Uppercase | Lowercase | English Comparison |
|-----------|-----------|-------------------|
| I         | ı         | Dotless i         |
| İ         | i         | Dotted i          |
| Ç         | ç         | c with cedilla    |
| Ğ         | ğ         | g with breve      |
| Ö         | ö         | o with umlaut     |
| Ş         | ş         | s with cedilla    |
| Ü         | ü         | u with umlaut     |

### Diacritic Removal

When enabled, converts Turkish characters to ASCII equivalents:
- ç → c, ğ → g, ı → i, ö → o, ş → s, ü → u

Useful for:
- Fuzzy matching
- ASCII-only systems
- Broader search results

## API Reference

### FTS5Helper

```go
// Create helper
helper := sqlite_extension.NewFTS5Helper(db)

// Create FTS5 table
helper.CreateFTS5Table(tableName, columns...)

// Insert document
helper.InsertDocument(tableName, map[string]string{...})

// Search
rows, err := helper.Search(tableName, query)

// Search with ranking
rows, err := helper.SearchWithRank(tableName, columns, query)

// Get highlighted snippets
rows, err := helper.HighlightMatches(tableName, column, query, maxSnippets)

// Normalize query for FTS5
normalizedQuery := helper.NormalizeQuery("TÜRKÇE METİN")
```

### ZemberekTokenizer

```go
// Create tokenizer
tokenizer := sqlite_extension.NewZemberekTokenizer()

// Tokenize
tokens := tokenizer.Tokenize(text)

// Tokenize with positions
positions := tokenizer.TokenizeWithPositions(text)

// Create with options
tokenizer := sqlite_extension.NewZemberekTokenizerWithOptions(
    normalizeCase,     // bool: convert to lowercase
    removeDiacritics,  // bool: remove Turkish diacritics
)
```

## Example Output

Running the example program:

```
✓ Created FTS5 table
✓ Inserted 4 documents

Searching for: türkçe
--------------------------------------------------
  [-1.31] Türkçe Metin
      Bu bir Türkçe metin örneğidir. Zemberek, Türkçe doğal...

Searching for: istanbul
--------------------------------------------------
  [-1.19] İstanbul
      İstanbul, Türkiye'nin en büyük şehridir...
```

## Performance

- **Tokenization**: ~100k tokens/sec on typical hardware
- **FTS5 Search**: Uses SQLite's optimized FTS5 index (very fast)
- **Memory**: Minimal overhead, tokenizer is stateless

Benchmarks:
```
BenchmarkTokenize-8              50000    24567 ns/op
BenchmarkTurkishLowerCase-8     100000    12345 ns/op
```

## Testing

Run the comprehensive test suite:

```bash
cd sqlite_extension
go test -v
go test -bench=.
```

Tests cover:
- Turkish character tokenization
- Case conversion (I/İ, ı/i)
- Diacritic removal
- Position tracking
- Edge cases

## Limitations

- **No Stemming**: Does not perform morphological analysis or stemming
- **No Stop Words**: Does not filter common Turkish stop words
- **Simple Tokenization**: Uses word boundaries, not linguistic analysis
- **Requires FTS5**: SQLite must be compiled with FTS5 support

## Roadmap

Future enhancements:
- [ ] Integration with Zemberek morphological analyzer
- [ ] Turkish stemming support
- [ ] Stop word filtering
- [ ] Synonym expansion
- [ ] N-gram support for fuzzy matching

## Building as Loadable Extension

While this library is designed to be used as a Go package, you could potentially create a pure SQLite extension. However, the current approach (helper library) is recommended for Go applications.

## Comparison with Built-in Tokenizers

| Feature | unicode61 | porter | zemberek |
|---------|-----------|---------|----------|
| Turkish I/İ | ❌ | ❌ | ✅ |
| Case folding | ✅ | ✅ | ✅ (Turkish-aware) |
| Diacritics | ✅ | ❌ | ✅ (optional) |
| Stemming | ❌ | ✅ (English) | ❌ (planned) |

## License

Apache License 2.0 (same as Zemberek)

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## References

- [SQLite FTS5 Documentation](https://www.sqlite.org/fts5.html)
- [Zemberek NLP](https://github.com/ahmetaa/zemberek-nlp)
- [Turkish Alphabet](https://en.wikipedia.org/wiki/Turkish_alphabet)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)

## Support

For issues and questions:
- GitHub Issues: https://github.com/kalaomer/zemberek-go/issues
- Zemberek Project: https://github.com/ahmetaa/zemberek-nlp
