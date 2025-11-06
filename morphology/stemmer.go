package morphology

import (
	"unicode"
	"unicode/utf8"
)

// StemToken represents a stemmed token with its byte position in the original text
type StemToken struct {
	Stem      string // Stemmed form: "kitap"
	Original  string // Original word: "kitapları"
	StartByte int    // UTF-8 byte offset start
	EndByte   int    // UTF-8 byte offset end
}

// StemTextWithPositions extracts stems from text with byte positions.
//
// This is the MAIN function for FTS5 integration. It:
// 1. Tokenizes text into words (tracking byte offsets)
// 2. Performs morphological analysis on each word
// 3. Extracts stems using GetStem()
// 4. Returns list of stems with their byte positions
//
// Example:
//
//	morph := CreateWithDefaults()
//	text := "Kitapları okuyorum"
//	tokens := StemTextWithPositions(text, morph)
//	// tokens[0] = {Stem: "kitap", Original: "Kitapları", StartByte: 0, EndByte: 11}
//	// tokens[1] = {Stem: "oku", Original: "okuyorum", StartByte: 12, EndByte: 20}
func StemTextWithPositions(text string, morphology *TurkishMorphology) []StemToken {
	if morphology == nil {
		// Fallback: no stemming, just tokenize
		return tokenizeWithoutStemming(text)
	}

	result := make([]StemToken, 0)
	wordInfos := tokenizeWithByteOffsets(text)

	for _, info := range wordInfos {
		// Skip non-word tokens (punctuation, whitespace, etc.)
		if !isWordToken(info.text) {
			continue
		}

		// Perform morphological analysis
		analysis := morphology.Analyze(info.text)

		stem := info.text // Default: use original if no stem found

		// Extract stem from first analysis result
		if len(analysis.AnalysisResults) > 0 {
			extractedStem := analysis.AnalysisResults[0].GetStem()
			if extractedStem != "" {
				stem = extractedStem
			}
		}

		result = append(result, StemToken{
			Stem:      stem,
			Original:  info.text,
			StartByte: info.startByte,
			EndByte:   info.endByte,
		})
	}

	return result
}

// wordInfo holds tokenization info with byte positions
type wordInfo struct {
	text      string
	startByte int
	endByte   int
}

// tokenizeWithByteOffsets tokenizes text tracking UTF-8 byte offsets
func tokenizeWithByteOffsets(text string) []wordInfo {
	words := make([]wordInfo, 0)

	currentWord := ""
	wordStartByte := 0
	byteOffset := 0
	inWord := false

	for _, r := range text {
		runeByteLen := utf8.RuneLen(r)

		if isWordChar(r) {
			// Start or continue word
			if !inWord {
				wordStartByte = byteOffset
				inWord = true
			}
			currentWord += string(r)
		} else {
			// End of word
			if inWord {
				words = append(words, wordInfo{
					text:      currentWord,
					startByte: wordStartByte,
					endByte:   byteOffset,
				})
				currentWord = ""
				inWord = false
			}
			// Note: We could also emit punctuation here if needed
		}

		byteOffset += runeByteLen
	}

	// Final word if text ends with word char
	if inWord && currentWord != "" {
		words = append(words, wordInfo{
			text:      currentWord,
			startByte: wordStartByte,
			endByte:   byteOffset,
		})
	}

	return words
}

// tokenizeWithoutStemming returns tokens without morphological analysis
func tokenizeWithoutStemming(text string) []StemToken {
	wordInfos := tokenizeWithByteOffsets(text)
	result := make([]StemToken, 0, len(wordInfos))

	for _, info := range wordInfos {
		if !isWordToken(info.text) {
			continue
		}
		result = append(result, StemToken{
			Stem:      info.text,
			Original:  info.text,
			StartByte: info.startByte,
			EndByte:   info.endByte,
		})
	}

	return result
}

// isWordChar checks if rune is part of a word
func isWordChar(r rune) bool {
	// Letters (including Turkish)
	if unicode.IsLetter(r) {
		return true
	}
	// Digits
	if unicode.IsDigit(r) {
		return true
	}
	// Apostrophe for possessives: "Ali'nin"
	if r == '\'' {
		return true
	}
	return false
}

// isWordToken checks if token is a word (not just punctuation/numbers)
func isWordToken(token string) bool {
	if token == "" {
		return false
	}

	// Must contain at least one letter
	for _, r := range token {
		if unicode.IsLetter(r) {
			return true
		}
	}

	return false
}
