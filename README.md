# Zemberek-Go

Go implementation of the original [zemberek-nlp](https://github.com/ahmetaa/zemberek-nlp) Java library for Turkish language processing.

## Features

Currently, the following modules have been ported:

### Core
- Turkish alphabet and phonetic attributes
- Multi-level perfect hash functions and compression primitives
- Text utilities for casing, diacritics and token helpers

### Tokenization
- Token/span types and sentence boundary detection

### Language Model (LM)
- Compressed vocabulary and n‑gram accessors
- SmoothLM reader with MPHFs

### Morphology
- Binary lexicon loader and dictionary items
- Morphotactics graph, analysis and generation helpers

### Normalization
- Full sentence normalizer with spell checker + LM ranking
- Deasciifier and ASCII tolerant utilities

## Installation

```bash
go get github.com/kalaomer/zemberek-go
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/kalaomer/zemberek-go/core/turkish"
    "github.com/kalaomer/zemberek-go/tokenization"
)

func main() {
    // Use Turkish alphabet
    alphabet := turkish.Instance
    fmt.Println("Is 'ı' a vowel?", alphabet.IsVowel('ı'))

    // Tokenize text
    extractor, _ := tokenization.NewTurkishSentenceExtractor(false, "")
    sentences := extractor.FromParagraph("Merhaba dünya! Bu bir test cümlesidir.")
    for _, sentence := range sentences {
        fmt.Println(sentence)
    }
}
```

### Sentence normalization

```go
package main

import (
    "fmt"
    "log"

    "github.com/kalaomer/zemberek-go/morphology"
    "github.com/kalaomer/zemberek-go/normalization"
)

func main() {
    morph := morphology.CreateWithDefaults()
    normalizer, err := normalization.NewTurkishSentenceNormalizerAdvanced(morph, "data")
    if err != nil {
        log.Fatalf("normalizer init: %v", err)
    }

    input := "Yrn okua gidicem"
    fmt.Println(normalizer.Normalize(input))
}
```

## Dependencies

- Go 1.18 or higher
- Standard library only (no external dependencies for core functionality)

### Resource data

Language resources (lexicon binaries, normalization tables, language models) are expected under `data/` by default. If you keep them elsewhere, export `ZEMBEREK_DATA_ROOT=/absolute/path/to/your/data` so both the examples and the advanced normalizer can locate them.

Example data bundles (LM and normalization folders) are available here: <https://drive.google.com/drive/folders/1tztjRiUs9BOTH-tb1v7FWyixl-iUpydW>. Download the archive, extract it to a directory of your choice, and point `ZEMBEREK_DATA_ROOT` to that directory before running the examples.

## Development Status

The port follows zemberek-nlp’s architecture module by module. Core components, tokenization, lexicon handling, language model loading and advanced normalization are functional; remaining work focuses on fine-tuning morphology generation/ambiguity resolution and extending test coverage as the Java baseline evolves.

## Notes

This port mirrors the Java implementation’s architecture while adapting to Go idioms:

- Java classes → Go structs/interfaces
- Java enums → Go iota constants
- Immutable data → Go value types and generated readers

## Credits

- Original Java implementation: [zemberek-nlp](https://github.com/ahmetaa/zemberek-nlp) by Ahmet A. Akın
- Go port: This repository and its contributors

## License

Apache License 2.0

## Contributing

Contributions are welcome! This is a large codebase and help with porting remaining modules would be appreciated.
