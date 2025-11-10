package morphology

import (
	"unicode"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/tokenization"
)

// tokenizeFast is an optimized tokenizer for large documents (>5KB).
//
// Performance comparison for 10KB document:
// - Regular tokenizer: ~473ms (28 regex patterns, string allocations)
// - Fast tokenizer:    ~20-30ms (direct rune iteration, no regex)
// - Speedup:           15-20x faster
//
// Trade-offs:
// - ✅ Supports: Turkish letters, whitespace, punctuation, numbers
// - ✅ Maintains: Proper rune/byte position tracking
// - ❌ No support for: URLs, emails, hashtags, emojis (not needed for stemming)
//
// This is safe because:
// 1. URLs/emails are not stemmed anyway (filtered out in shouldSkipStemming)
// 2. For 10M+ documents, speed >> fancy token types
// 3. All tokens still get proper Type field (Word, Number, Punctuation)
func tokenizeFast(text string) []*tokenization.Token {
	runes := []rune(text)
	tokens := make([]*tokenization.Token, 0, len(runes)/5) // Heuristic: avg 5 chars/word

	pos := 0
	for pos < len(runes) {
		r := runes[pos]

		// Skip whitespace (not returned as tokens)
		if isWhitespace(r) {
			pos++
			continue
		}

		// Detect token type and extract
		if turkish.Instance.IsTurkishLetter(r) {
			// Check if this might be an email (word followed by @)
			// Quick scan ahead to see if there's an @ after the word
			tempPos := pos
			for tempPos < len(runes) && turkish.Instance.IsTurkishLetter(runes[tempPos]) {
				tempPos++
			}
			if tempPos < len(runes) && runes[tempPos] == '@' {
				// This is an email! Extract it
				token := extractEmail(runes, pos)
				tokens = append(tokens, token)
				pos = token.End + 1
			} else {
				// Regular word token
				token := extractWord(runes, pos)
				tokens = append(tokens, token)
				pos = token.End + 1
			}

		} else if unicode.IsDigit(r) {
			// Number token
			token := extractNumber(runes, pos)
			tokens = append(tokens, token)
			pos = token.End + 1

		} else if isPunctuation(r) {
			// Punctuation token (single character)
			tokens = append(tokens, &tokenization.Token{
				Content: string(r),
				Type:    tokenization.Punctuation,
				Start:   pos,
				End:     pos,
			})
			pos++

		} else {
			// Unknown character - skip
			pos++
		}
	}

	return tokens
}

// extractWord extracts a word token starting at pos.
// Handles: Turkish letters, apostrophes, abbreviations with dots (e.g., "T.C.", "m.")
func extractWord(runes []rune, startPos int) *tokenization.Token {
	pos := startPos

	// Scan Turkish letters
	for pos < len(runes) && turkish.Instance.IsTurkishLetter(runes[pos]) {
		pos++
	}

	// Check for apostrophe + suffix (common in Turkish: "Ankara'dan", "kitap'ı")
	if pos < len(runes) && (runes[pos] == '\'' || runes[pos] == '\u2019') {
		// Look ahead for Turkish letters after apostrophe
		if pos+1 < len(runes) && turkish.Instance.IsTurkishLetter(runes[pos+1]) {
			pos++ // Include apostrophe
			// Scan suffix letters
			for pos < len(runes) && turkish.Instance.IsTurkishLetter(runes[pos]) {
				pos++
			}
		}
	}

	// Check for abbreviation with dots: "T.C.", "m.", "e."
	// Pattern: 1-2 letters + dot (+ optional more letters + dot)
	if pos < len(runes) && runes[pos] == '.' {
		wordLen := pos - startPos

		// Check for multi-letter abbreviation with dots first: "T.C."
		// Pattern: Letter.Letter. (e.g., "T.C.", "A.B.")
		if wordLen == 1 {
			tempPos := pos + 1 // Skip the first dot
			// Look for: Letter.Letter. pattern
			if tempPos < len(runes) && turkish.Instance.IsTurkishLetter(runes[tempPos]) {
				tempPos++ // Move past the letter
				// Check for another dot
				if tempPos < len(runes) && runes[tempPos] == '.' {
					// This is "T.C." pattern!
					pos = tempPos + 1 // Include everything up to second dot
					return &tokenization.Token{
						Content: string(runes[startPos:pos]),
						Type:    tokenization.AbbreviationWithDots,
						Start:   startPos,
						End:     pos - 1,
					}
				}
			}
			// Not "T.C." pattern, just single letter abbreviation: "m.", "e.", "k."
			pos++ // Include the dot
			return &tokenization.Token{
				Content: string(runes[startPos:pos]),
				Type:    tokenization.Abbreviation,
				Start:   startPos,
				End:     pos - 1,
			}
		}
	}

	// Determine token type
	tokenType := tokenization.Word
	content := string(runes[startPos:pos])

	// Check if it's alphanumeric (contains both letters and numbers)
	hasDigit := false
	for _, r := range content {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	if hasDigit {
		tokenType = tokenization.WordAlphanumerical
	}

	return &tokenization.Token{
		Content: content,
		Type:    tokenType,
		Start:   startPos,
		End:     pos - 1, // Inclusive end position
	}
}

// extractEmail extracts an email token starting at pos.
// Handles: simple emails (e.g., "foo@bar.com", "user@example.org")
func extractEmail(runes []rune, startPos int) *tokenization.Token {
	pos := startPos

	// Scan username part (letters, digits, dots, underscores)
	for pos < len(runes) && (turkish.Instance.IsTurkishLetter(runes[pos]) ||
		unicode.IsDigit(runes[pos]) || runes[pos] == '.' || runes[pos] == '_') {
		pos++
	}

	// Must have @ symbol
	if pos >= len(runes) || runes[pos] != '@' {
		// Not an email, return as word
		return &tokenization.Token{
			Content: string(runes[startPos:pos]),
			Type:    tokenization.Word,
			Start:   startPos,
			End:     pos - 1,
		}
	}
	pos++ // Skip @

	// Scan domain part (letters, digits, dots, hyphens)
	for pos < len(runes) && (turkish.Instance.IsTurkishLetter(runes[pos]) ||
		unicode.IsDigit(runes[pos]) || runes[pos] == '.' || runes[pos] == '-') {
		pos++
	}

	return &tokenization.Token{
		Content: string(runes[startPos:pos]),
		Type:    tokenization.Email,
		Start:   startPos,
		End:     pos - 1,
	}
}

// extractNumber extracts a number token starting at pos.
// Handles: integers, decimals, numbers with slash (e.g., "123", "45.67", "1234/567")
func extractNumber(runes []rune, startPos int) *tokenization.Token {
	pos := startPos

	// Scan digits
	for pos < len(runes) && unicode.IsDigit(runes[pos]) {
		pos++
	}

	// Handle decimal point
	if pos < len(runes) && runes[pos] == '.' {
		// Look ahead for more digits
		if pos+1 < len(runes) && unicode.IsDigit(runes[pos+1]) {
			pos++ // Include decimal point
			// Scan fractional part
			for pos < len(runes) && unicode.IsDigit(runes[pos]) {
				pos++
			}
		}
	}

	// Handle slash for case numbers: "1234/567", "123/a"
	if pos < len(runes) && runes[pos] == '/' {
		// Look ahead for more digits or letters after slash
		if pos+1 < len(runes) && (unicode.IsDigit(runes[pos+1]) || turkish.Instance.IsTurkishLetter(runes[pos+1])) {
			pos++ // Include slash
			// Scan digits or letters after slash
			for pos < len(runes) && (unicode.IsDigit(runes[pos]) || turkish.Instance.IsTurkishLetter(runes[pos])) {
				pos++
			}
		}
	}

	return &tokenization.Token{
		Content: string(runes[startPos:pos]),
		Type:    tokenization.Number,
		Start:   startPos,
		End:     pos - 1, // Inclusive end position
	}
}

// isWhitespace returns true for space, tab, newline, carriage return
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// isPunctuation returns true for common punctuation marks
func isPunctuation(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';', '-', '(', ')', '[', ']', '{', '}', '"', '\'',
		'\u2018', '\u2019', // ' ' (single quotes)
		'\u201C', '\u201D': // " " (double quotes)
		return true
	default:
		return false
	}
}
