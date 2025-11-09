package tokenization

import (
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

// DetermineTokenType analyzes a word and returns its token type.
// This function uses pattern matching with priority order.
//
// Priority order (highest to lowest):
// 1. Whitespace (SpaceTab, NewLine)
// 2. URL (must be before Email - can contain @)
// 3. Email
// 4. Time
// 5. Date
// 6. Mention (@user)
// 7. HashTag (#tag)
// 8. MetaTag (<tag>)
// 9. Percent (%100)
// 10. Number (various formats)
// 11. Emoticon
// 12. Roman Numeral
// 13. Abbreviation with dots (I.B.M.)
// 14. Word with symbol (F-16)
// 15. Word alphanumerical (F16)
// 16. Word
// 17. Punctuation
// 18. Unknown word
// 19. Unknown
func DetermineTokenType(word string) TokenType {
	if word == "" {
		return Unknown
	}

	// 1. Whitespace
	if spaceTabPattern.MatchString(word) {
		return SpaceTab
	}
	if newLinePattern.MatchString(word) {
		return NewLine
	}

	// 2. URL (before Email - can contain @)
	if urlPattern.MatchString(word) {
		return URL
	}

	// 3. Email
	if emailPattern.MatchString(word) {
		return Email
	}

	// 4. Date (before Time - more specific pattern)
	if datePattern.MatchString(word) {
		return Date
	}

	// 5. Time
	if timePattern.MatchString(word) {
		return Time
	}

	// 6. Mention
	if mentionPattern.MatchString(word) {
		return Mention
	}

	// 7. HashTag
	if hashTagPattern.MatchString(word) {
		return HashTag
	}

	// 8. MetaTag
	if metaTagPattern.MatchString(word) {
		return MetaTag
	}

	// 9. Emoticon (before Number - to catch patterns like "8-)")
	if emoticonPattern.MatchString(word) {
		return Emoticon
	}

	// 10. Percent
	if percentPattern.MatchString(word) {
		return PercentNumeral
	}

	// 11. Number (check specific patterns in order)
	// Order matters: most specific patterns first
	if numberExpPattern.MatchString(word) {
		return Number
	}
	if numberFractionPattern.MatchString(word) {
		return Number
	}
	if numberThousandDotPattern.MatchString(word) {
		return Number
	}
	if numberThousandCommaPattern.MatchString(word) {
		return Number
	}
	if numberDecimalPattern.MatchString(word) {
		return Number
	}
	if numberOrdinalPattern.MatchString(word) {
		return Number
	}
	if numberIntegerPattern.MatchString(word) {
		return Number
	}

	// 12. Abbreviation with dots (before RomanNumeral - more specific)
	if abbreviationWithDotsPattern.MatchString(word) {
		return AbbreviationWithDots
	}

	// 13. Roman Numeral
	if romanNumeralPattern.MatchString(word) {
		return RomanNumeral
	}

	// 14. Word with symbol (F-16'yı)
	if wordWithSymbolPattern.MatchString(word) {
		return WordWithSymbol
	}

	// 15. Word alphanumerical (F16, H1N1)
	// Must check if contains both letters and digits
	hasLetter := false
	hasDigit := false
	alphabet := turkish.Instance
	for _, r := range word {
		if alphabet.IsTurkishLetter(r) {
			hasLetter = true
		}
		if turkish.IsDigit(r) {
			hasDigit = true
		}
	}
	if hasLetter && hasDigit && wordAlphanumericalPattern.MatchString(word) {
		return WordAlphanumerical
	}

	// 16. Pure word (Turkish letters only)
	if wordPattern.MatchString(word) {
		return Word
	}

	// 17. Punctuation
	if punctuationPattern.MatchString(word) {
		return Punctuation
	}

	// 18. Unknown word (some characters but not matching other patterns)
	if unknownWordPattern.MatchString(word) {
		return UnknownWord
	}

	// 19. Unknown (default)
	return Unknown
}

// NormalizeApostrophe normalizes different apostrophe types to standard apostrophe.
// Converts U+2019 (') to U+0027 (')
func NormalizeApostrophe(text string) string {
	return strings.ReplaceAll(text, "'", "'")
}

// stripTurkishSuffix removes Turkish suffix from token if present.
// Returns the base form without suffix.
// Example: "100'e" -> "100", "www.foo.com'da" -> "www.foo.com"
func stripTurkishSuffix(word string) string {
	// Find apostrophe position (U+0027 or U+2019)
	apostrophePos := -1
	for i, r := range word {
		if r == '\'' || r == 0x2019 { // U+0027 or U+2019
			apostrophePos = i
			break
		}
	}

	if apostrophePos > 0 {
		return word[:apostrophePos]
	}

	return word
}
