package tokenization

import (
	"regexp"
)

// TurkishTokenizer tokenizes Turkish text into typed tokens
type TurkishTokenizer struct {
	acceptedTypes map[TokenType]bool
}

// Pattern matchers with priority order
var tokenPatterns = []struct {
	pattern   *regexp.Regexp
	tokenType TokenType
}{
	// Whitespace (highest priority - single char)
	{regexp.MustCompile(`^\s`), SpaceTab},

	// URL (before email - can contain @)
	{urlPattern, URL},

	// Email
	{emailPattern, Email},

	// Date (before Time - more specific)
	{datePattern, Date},

	// Time
	{timePattern, Time},

	// Mention
	{mentionPattern, Mention},

	// HashTag
	{hashTagPattern, HashTag},

	// MetaTag
	{metaTagPattern, MetaTag},

	// Emoticon (before Number - catches "8-)")
	{emoticonPattern, Emoticon},

	// Percent
	{percentPattern, PercentNumeral},

	// Number (all formats)
	{numberExpPattern, Number},
	{numberFractionPattern, Number},
	{numberThousandDotPattern, Number},
	{numberThousandCommaPattern, Number},
	{numberDecimalPattern, Number},
	{numberOrdinalPattern, Number},
	{numberIntegerPattern, Number},

	// Abbreviation with dots
	{abbreviationWithDotsPattern, AbbreviationWithDots},

	// Roman numeral
	{romanNumeralPattern, RomanNumeral},

	// Word with symbol
	{wordWithSymbolPattern, WordWithSymbol},

	// Word alphanumerical
	{wordAlphanumericalPattern, WordAlphanumerical},

	// Pure word
	{wordPattern, Word},

	// Punctuation
	{punctuationPattern, Punctuation},

	// Unknown word
	{unknownWordPattern, UnknownWord},
}

// Tokenize tokenizes text and returns array of tokens with types and positions
func (t *TurkishTokenizer) Tokenize(text string) []*Token {
	tokens := make([]*Token, 0)
	runes := []rune(text)
	pos := 0

	for pos < len(runes) {
		remaining := string(runes[pos:])
		matched := false

		// Try all patterns and find the longest match
		var longestMatch struct {
			text      string
			tokenType TokenType
			length    int
		}

		for _, pt := range tokenPatterns {
			if loc := pt.pattern.FindStringIndex(remaining); loc != nil && loc[0] == 0 {
				// Pattern matched at start of remaining text
				matchedText := remaining[:loc[1]]
				tokenType := pt.tokenType
				matchLen := loc[1]

				// Special handling for whitespace types
				if tokenType == SpaceTab {
					if matchedText == "\n" || matchedText == "\r" {
						tokenType = NewLine
					}
				}

				// Special case: Word + "." might be abbreviation
				if (tokenType == Word || tokenType == WordAlphanumerical) && pos+matchLen < len(runes) && runes[pos+matchLen] == '.' {
					withDot := matchedText + "."
					if IsAbbreviation(withDot) {
						matchedText = withDot
						tokenType = Abbreviation
						matchLen++ // Include the dot in length
					}
				}

				// Keep track of longest match
				if matchLen > longestMatch.length {
					longestMatch.text = matchedText
					longestMatch.tokenType = tokenType
					longestMatch.length = matchLen
				}
			}
		}

		// Use the longest match
		if longestMatch.length > 0 {
			matched = true
			if t.acceptedTypes[longestMatch.tokenType] {
				tokens = append(tokens, &Token{
					Content:    longestMatch.text,
					Type:       longestMatch.tokenType,
					Start:      pos,
					End:        pos + len([]rune(longestMatch.text)) - 1,
					Normalized: NormalizeApostrophe(longestMatch.text),
				})
			}
			pos += len([]rune(longestMatch.text))
		}

		if !matched {
			// No pattern matched - treat as unknown single character
			if t.acceptedTypes[Unknown] {
				tokens = append(tokens, &Token{
					Content:    string(runes[pos]),
					Type:       Unknown,
					Start:      pos,
					End:        pos,
					Normalized: string(runes[pos]),
				})
			}
			pos++
		}
	}

	return tokens
}

// TokenizeToStrings is a convenience method that returns only token contents
func (t *TurkishTokenizer) TokenizeToStrings(text string) []string {
	tokens := t.Tokenize(text)
	result := make([]string, len(tokens))
	for i, token := range tokens {
		result[i] = token.Content
	}
	return result
}

// Predefined tokenizer instances (matches Java)
var (
	// ALL tokenizer accepts all token types
	ALL = NewBuilder().AcceptAll().Build()

	// DEFAULT tokenizer ignores whitespace (NewLine, SpaceTab)
	DEFAULT = NewBuilder().
		AcceptAll().
		IgnoreTypes(NewLine, SpaceTab).
		Build()
)
