package sqlite_extension

import (
	"strings"
	"unicode"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/tokenization"
)

// AdvancedTokenizer provides more sophisticated tokenization using Zemberek's features
type AdvancedTokenizer struct {
	extractor      *tokenization.TurkishSentenceExtractor
	alphabet       *turkish.TurkishAlphabet
	normalizeCase  bool
	removeDiacritics bool
}

// NewAdvancedTokenizer creates a new advanced tokenizer
func NewAdvancedTokenizer(normalizeCase, removeDiacritics bool) (*AdvancedTokenizer, error) {
	extractor, err := tokenization.NewTurkishSentenceExtractor(false, "")
	if err != nil {
		return nil, err
	}

	return &AdvancedTokenizer{
		extractor:        extractor,
		alphabet:         turkish.Instance,
		normalizeCase:    normalizeCase,
		removeDiacritics: removeDiacritics,
	}, nil
}

// Tokenize tokenizes the input text and calls the callback for each token
func (t *AdvancedTokenizer) Tokenize(text string, callback func(flags int, token string, start, end int) int) int {
	tokens := tokenization.SimpleTokenize(text)

	offset := 0
	for _, token := range tokens {
		// Find token position in original text
		idx := strings.Index(text[offset:], token)
		if idx == -1 {
			continue
		}

		start := offset + idx
		end := start + len(token)

		// Skip whitespace-only tokens
		if strings.TrimSpace(token) == "" {
			offset = end
			continue
		}

		// Skip pure punctuation
		if isPunctuation(token) {
			offset = end
			continue
		}

		// Normalize token
		normalizedToken := t.normalizeToken(token)
		if normalizedToken == "" {
			offset = end
			continue
		}

		// Call callback with normalized token
		rc := callback(0, normalizedToken, start, end)
		if rc != 0 {
			return rc
		}

		offset = end
	}

	return 0 // SQLITE_OK
}

// normalizeToken normalizes a token according to tokenizer settings
func (t *AdvancedTokenizer) normalizeToken(token string) string {
	result := token

	// Remove diacritics if configured
	if t.removeDiacritics {
		result = t.removeTurkishDiacritics(result)
	}

	// Normalize case if configured
	if t.normalizeCase {
		result = t.turkishLowerCase(result)
	}

	return result
}

// turkishLowerCase converts string to lowercase using Turkish rules
func (t *AdvancedTokenizer) turkishLowerCase(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch r {
		case 'I':
			result.WriteRune('ı')
		case 'İ':
			result.WriteRune('i')
		default:
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

// removeTurkishDiacritics removes Turkish diacritical marks
func (t *AdvancedTokenizer) removeTurkishDiacritics(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch r {
		case 'ç':
			result.WriteRune('c')
		case 'Ç':
			result.WriteRune('C')
		case 'ğ':
			result.WriteRune('g')
		case 'Ğ':
			result.WriteRune('G')
		case 'ı':
			result.WriteRune('i')
		case 'I':
			result.WriteRune('I')
		case 'İ':
			result.WriteRune('I')
		case 'ö':
			result.WriteRune('o')
		case 'Ö':
			result.WriteRune('O')
		case 'ş':
			result.WriteRune('s')
		case 'Ş':
			result.WriteRune('S')
		case 'ü':
			result.WriteRune('u')
		case 'Ü':
			result.WriteRune('U')
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

// isPunctuation checks if a string contains only punctuation
func isPunctuation(s string) bool {
	for _, r := range s {
		if !unicode.IsPunct(r) && !unicode.IsSymbol(r) {
			return false
		}
	}
	return true
}
