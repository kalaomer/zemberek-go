package sqlite_extension

import (
	"strings"
	"unicode"

	"github.com/kalaomer/zemberek-go/tokenization"
)

// ZemberekTokenizer provides Turkish-aware tokenization for SQLite FTS5
type ZemberekTokenizer struct {
	normalizeCase    bool
	removeDiacritics bool
}

// NewZemberekTokenizer creates a new tokenizer with default settings
func NewZemberekTokenizer() *ZemberekTokenizer {
	return &ZemberekTokenizer{
		normalizeCase:    true,
		removeDiacritics: false,
	}
}

// NewZemberekTokenizerWithOptions creates a tokenizer with custom options
func NewZemberekTokenizerWithOptions(normalizeCase, removeDiacritics bool) *ZemberekTokenizer {
	return &ZemberekTokenizer{
		normalizeCase:    normalizeCase,
		removeDiacritics: removeDiacritics,
	}
}

// Tokenize tokenizes the input text and returns tokens
func (z *ZemberekTokenizer) Tokenize(text string) []string {
	tokens := tokenization.SimpleTokenize(text)
	result := make([]string, 0, len(tokens))

	for _, token := range tokens {
		// Skip whitespace-only tokens
		if strings.TrimSpace(token) == "" {
			continue
		}

		// Skip pure punctuation
		if isPunctuation(token) {
			continue
		}

		// Normalize token
		normalized := z.normalizeToken(token)
		if normalized != "" {
			result = append(result, normalized)
		}
	}

	return result
}

// TokenizeWithPositions returns tokens with their positions in the original text
func (z *ZemberekTokenizer) TokenizeWithPositions(text string) []TokenPosition {
	tokens := tokenization.SimpleTokenize(text)
	result := make([]TokenPosition, 0, len(tokens))

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
		normalized := z.normalizeToken(token)
		if normalized != "" {
			result = append(result, TokenPosition{
				Token: normalized,
				Start: start,
				End:   end,
			})
		}

		offset = end
	}

	return result
}

// TokenPosition represents a token with its position in the text
type TokenPosition struct {
	Token string
	Start int
	End   int
}

// normalizeToken normalizes a token according to tokenizer settings
func (z *ZemberekTokenizer) normalizeToken(token string) string {
	result := token

	// Remove diacritics if configured
	if z.removeDiacritics {
		result = removeTurkishDiacritics(result)
	}

	// Normalize case if configured
	if z.normalizeCase {
		result = turkishLowerCase(result)
	}

	return result
}

// turkishLowerCase converts string to lowercase using Turkish rules
func turkishLowerCase(s string) string {
	var result strings.Builder
	result.Grow(len(s))
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

// TurkishUpperCase converts string to uppercase using Turkish rules
func TurkishUpperCase(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for _, r := range s {
		switch r {
		case 'i':
			result.WriteRune('İ')
		case 'ı':
			result.WriteRune('I')
		default:
			result.WriteRune(unicode.ToUpper(r))
		}
	}
	return result.String()
}

// removeTurkishDiacritics removes Turkish diacritical marks
func removeTurkishDiacritics(s string) string {
	var result strings.Builder
	result.Grow(len(s))
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
	return len(s) > 0
}

// TokenizeText is a convenience function for quick tokenization
func TokenizeText(text string) []string {
	tokenizer := NewZemberekTokenizer()
	return tokenizer.Tokenize(text)
}

// NormalizeForSearch normalizes text for FTS5 search queries
// This ensures the search terms match how the text was indexed
func NormalizeForSearch(query string) string {
	tokenizer := NewZemberekTokenizer()
	tokens := tokenizer.Tokenize(query)
	return strings.Join(tokens, " ")
}
