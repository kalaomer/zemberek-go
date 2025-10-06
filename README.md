# Zemberek-Go

Go implementation of Zemberek NLP library for Turkish language processing.

This is a port of [zemberek-python](https://github.com/Loodos/zemberek-python) which itself is a Python port of the original [zemberek-nlp](https://github.com/ahmetaa/zemberek-nlp) Java library.

## Features

Currently, the following modules have been ported:

### Core
- **Turkish Language Support**: Turkish alphabet, letters, phonetic attributes
- **Hash Functions**: Multi-level perfect hash functions for compression
- **Text Processing**: Text normalization utilities
- **Compression**: Lossy integer lookup, quantization
- **Data Structures**: Weight lookups, compressed weights

### Tokenization
- **Token**: Token types and structures
- **Span**: Text span handling
- **Sentence Extraction**: Turkish sentence boundary detection using perceptron models
- **Perceptron Segmenter**: Rule-based and ML-based sentence segmentation

### Language Model (LM)
- **Vocabulary**: Language model vocabulary handling
- **N-gram Data**: Compressed n-gram storage
- **Gram Data Array**: Efficient n-gram data access

### Morphology
- **Lexicon**: Dictionary items and root lexicon
- **Morphemes**: Morpheme definitions and structures
- **Morphotactics**: Turkish morphological rules (in progress)
- **Analysis**: Word analysis (in progress)
- **Generation**: Word generation (in progress)

### Normalization
- **Spell Checking**: Turkish spell checking (in progress)
- **Text Normalization**: Noisy text normalization (in progress)
- **Deasciifier**: Turkish diacritics restoration (in progress)

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

## Project Structure

```
zemberek-go/
├── core/
│   ├── turkish/      # Turkish language core
│   ├── text/         # Text utilities
│   ├── hash/         # Hash functions
│   ├── compression/  # Compression algorithms
│   ├── quantization/ # Quantization
│   ├── data/         # Data structures
│   └── utils/        # Utilities
├── tokenization/     # Tokenization
├── lm/              # Language models
│   └── compression/ # LM compression
├── morphology/      # Morphological analysis
│   ├── lexicon/     # Dictionary
│   ├── morphotactics/ # Morphological rules
│   ├── analysis/    # Word analysis
│   ├── generator/   # Word generation
│   └── ambiguity/   # Disambiguation
├── normalization/   # Text normalization
│   └── deasciifier/ # Diacritics restoration
└── resources/       # Data files

```

## Dependencies

- Go 1.18 or higher
- Standard library only (no external dependencies for core functionality)

## Development Status

This is an ongoing port of the Python version. The core functionality has been implemented, but some modules are still in progress:

- ✅ Core modules (Turkish, Hash, Compression, Text)
- ✅ Tokenization (Token, Span, Sentence Extraction)
- ✅ LM Vocabulary and basic structures
- ✅ Morphology Lexicon
- 🚧 Morphology Analysis and Generation
- 🚧 Normalization modules
- 🚧 Complete LM implementation

## Notes

This port maintains the architecture and approach of the original Python implementation while adapting to Go's idioms and best practices:

- Python classes → Go structs with methods
- Python enums → Go iota constants
- Python dictionaries → Go maps
- Python sets → Go maps with bool values
- Python inheritance → Go composition and interfaces

## Credits

- Original Java implementation: [zemberek-nlp](https://github.com/ahmetaa/zemberek-nlp) by Ahmet A. Akın
- Python port: [zemberek-python](https://github.com/Loodos/zemberek-python) by Loodos
- Go port: This repository

## License

Apache License 2.0

## Contributing

Contributions are welcome! This is a large codebase and help with porting remaining modules would be appreciated.
