# Zemberek SQLite FTS5 Extension

SQLite FTS5 (Full-Text Search) tokenizer extension using Zemberek-Go for Turkish language support.

## Overview

This extension provides a custom tokenizer for SQLite's FTS5 full-text search that uses Zemberek's Turkish language processing capabilities. It enables better search results for Turkish text by properly handling Turkish-specific characters, case folding, and tokenization rules.

## Features

- **Turkish-aware tokenization**: Properly handles Turkish characters (ç, ğ, ı, ö, ş, ü)
- **Case normalization**: Correctly handles Turkish case folding (I↔ı, İ↔i)
- **Diacritic removal**: Optional removal of Turkish diacritics for broader matching
- **FTS5 integration**: Seamlessly integrates with SQLite's FTS5 full-text search

## Building

### Prerequisites

- Go 1.18 or higher
- CGO enabled (required for SQLite extension)
- SQLite 3.20.0+ with FTS5 support
- GCC or compatible C compiler

### Build Commands

```bash
# Build the extension
make build

# Run tests
make test

# Build and run example
make run-example
```

## Usage

### Basic Example

```go
package main

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./mydb.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create FTS5 table with Zemberek tokenizer
    _, err = db.Exec(`
        CREATE VIRTUAL TABLE documents USING fts5(
            title,
            content,
            tokenize='zemberek'
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Insert Turkish text
    _, err = db.Exec(`
        INSERT INTO documents(title, content)
        VALUES ('Türkçe Başlık', 'Bu bir Türkçe metin örneğidir.')
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Search
    rows, err := db.Query(`
        SELECT title, content
        FROM documents
        WHERE documents MATCH 'türkçe'
    `)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var title, content string
        rows.Scan(&title, &content)
        fmt.Printf("Found: %s - %s\n", title, content)
    }
}
```

### Creating an FTS5 Table

```sql
-- Basic usage
CREATE VIRTUAL TABLE documents USING fts5(
    title,
    content,
    tokenize='zemberek'
);

-- With multiple columns
CREATE VIRTUAL TABLE articles USING fts5(
    headline,
    body,
    author,
    tokenize='zemberek'
);
```

### Searching

```sql
-- Simple search
SELECT * FROM documents WHERE documents MATCH 'türkçe';

-- Phrase search
SELECT * FROM documents WHERE documents MATCH '"doğal dil işleme"';

-- Boolean search
SELECT * FROM documents WHERE documents MATCH 'türkçe OR yazılım';

-- Column-specific search
SELECT * FROM documents WHERE documents MATCH 'title:istanbul';

-- With ranking
SELECT title, content, rank
FROM documents
WHERE documents MATCH 'yazılım'
ORDER BY rank;
```

## How It Works

The tokenizer performs the following steps:

1. **Tokenization**: Splits text into words using Zemberek's tokenization logic
2. **Normalization**:
   - Converts to lowercase using Turkish case-folding rules
   - Optionally removes diacritics (ç→c, ğ→g, ı→i, ö→o, ş→s, ü→u)
3. **Filtering**: Removes punctuation and whitespace-only tokens
4. **Indexing**: Passes normalized tokens to SQLite's FTS5 indexer

## Turkish Case Folding

The tokenizer correctly handles Turkish-specific case conversions:

| Uppercase | Lowercase |
|-----------|-----------|
| I         | ı         |
| İ         | i         |
| Ç         | ç         |
| Ğ         | ğ         |
| Ö         | ö         |
| Ş         | ş         |
| Ü         | ü         |

## Advanced Usage

### Custom Tokenizer Configuration

```go
import "github.com/kalaomer/zemberek-go/sqlite_extension"

// Create advanced tokenizer with custom settings
tokenizer, err := sqlite_extension.NewAdvancedTokenizer(
    true,  // normalizeCase
    true,  // removeDiacritics
)
if err != nil {
    log.Fatal(err)
}
```

## Example Program

See the `example/` directory for a complete working example:

```bash
cd example
go build
./fts5_example
```

This will:
1. Create a sample database with Turkish documents
2. Perform various searches
3. Display ranked results

## Performance Considerations

- The tokenizer is designed for correctness rather than maximum speed
- For large datasets, consider using SQLite's built-in tokenizers if Turkish-specific handling is not critical
- The tokenizer uses Go's Unicode libraries which are well-optimized

## Limitations

- Currently uses simple word-based tokenization
- Does not perform stemming or lemmatization
- Does not handle compound words specially

## Future Enhancements

- [ ] Morphological analysis integration
- [ ] Stemming support
- [ ] Compound word handling
- [ ] Stop word filtering
- [ ] Configurable tokenizer options via SQL

## Building as a Loadable Extension

To build as a loadable SQLite extension (`.so` or `.dylib`):

```bash
# Linux
go build -buildmode=c-shared -o zemberek_fts5.so

# macOS
go build -buildmode=c-shared -o zemberek_fts5.dylib
```

Then load in SQLite:

```sql
.load ./zemberek_fts5
```

## Testing

Run the test suite:

```bash
make test
```

## Dependencies

- [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - CGO SQLite3 driver
- github.com/kalaomer/zemberek-go - Zemberek Turkish NLP library

## License

Apache License 2.0 (same as Zemberek)

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## References

- [SQLite FTS5 Documentation](https://www.sqlite.org/fts5.html)
- [FTS5 Tokenizers](https://www.sqlite.org/fts5.html#tokenizers)
- [Zemberek NLP](https://github.com/ahmetaa/zemberek-nlp)
- [Turkish Alphabet](https://en.wikipedia.org/wiki/Turkish_alphabet)

## Support

For issues and questions:
- GitHub Issues: https://github.com/kalaomer/zemberek-go/issues
- Original Zemberek: https://github.com/ahmetaa/zemberek-nlp
