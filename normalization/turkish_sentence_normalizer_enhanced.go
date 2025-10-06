package normalization

import (
	"strings"
)

// TurkishSentenceNormalizerEnhanced is an improved normalizer with lookup tables and dictionary support
type TurkishSentenceNormalizerEnhanced struct {
	LookupManual   map[string][]string
	WordDictionary map[string]bool
	SpellChecker   *CharacterGraphDecoder
	Graph          *CharacterGraph
}

// NewTurkishSentenceNormalizerEnhanced creates a new enhanced normalizer
func NewTurkishSentenceNormalizerEnhanced() (*TurkishSentenceNormalizerEnhanced, error) {
	// Load manual lookup map (embedded default)
	lookupManual := GetDefaultLookupMap()

	// Try to load from file if exists
	fileLookup, err := LoadLookupMap("resources/normalization/candidates-manual.txt")
	if err == nil {
		// Merge file lookups with defaults
		for k, v := range fileLookup {
			lookupManual[k] = v
		}
	}

	// Load word dictionary
	wordDict := make(map[string]bool)
	words, err := LoadWordList("resources/dictionary/basic-turkish.txt")
	if err == nil {
		for _, word := range words {
			wordDict[word] = true
		}
	}

	// Add lookup values to dictionary
	for _, candidates := range lookupManual {
		for _, word := range candidates {
			wordDict[word] = true
		}
	}

	// Build character graph for spell checking
	graph := NewCharacterGraph()
	for word := range wordDict {
		if word != "" {
			graph.AddWord(word, TypeWord)
		}
	}

	decoder := NewCharacterGraphDecoder(graph)

	return &TurkishSentenceNormalizerEnhanced{
		LookupManual:   lookupManual,
		WordDictionary: wordDict,
		SpellChecker:   decoder,
		Graph:          graph,
	}, nil
}

// Normalize normalizes a Turkish sentence
func (tsne *TurkishSentenceNormalizerEnhanced) Normalize(sentence string) string {
	// Tokenize sentence (simple whitespace split, preserving punctuation)
	words := tokenizeSentence(sentence)
	normalized := make([]string, 0, len(words))

	for _, word := range words {
		// Skip empty words
		if word == "" {
			continue
		}

		// Preserve punctuation and special characters
		if isPunctuation(word) {
			normalized = append(normalized, word)
			continue
		}

		// Preserve email addresses and URLs (contain @ or multiple dots)
		if strings.Contains(word, "@") || isLikelyDomain(word) {
			normalized = append(normalized, word)
			continue
		}

		// Normalize the word
		normalizedWord := tsne.normalizeWord(word)
		normalized = append(normalized, normalizedWord)
	}

	// Join with spaces
	result := strings.Join(normalized, " ")

	// Clean up spacing around punctuation
	result = cleanPunctuation(result)

	return result
}

// normalizeWord normalizes a single word
func (tsne *TurkishSentenceNormalizerEnhanced) normalizeWord(word string) string {
	// Keep original casing info
	isCapitalized := len(word) > 0 && isUpper(rune(word[0]))
	lowerWord := strings.ToLower(word)

	// 1. Check manual lookup first (highest priority)
	if candidates, ok := tsne.LookupManual[lowerWord]; ok && len(candidates) > 0 {
		normalized := candidates[0] // Take first candidate
		if isCapitalized {
			return capitalize(normalized)
		}
		return normalized
	}

	// 2. Check if word exists in dictionary
	if tsne.WordDictionary[lowerWord] {
		return word // Already correct
	}

	// 3. Try spell checker
	suggestions := tsne.SpellChecker.GetSuggestions(lowerWord, DiacriticsIgnoringMatcherInstance)
	if len(suggestions) > 0 {
		normalized := suggestions[0]
		if isCapitalized {
			return capitalize(normalized)
		}
		return normalized
	}

	// 4. Return original if no normalization found
	return word
}

// Helper functions

func tokenizeSentence(sentence string) []string {
	var tokens []string
	var current strings.Builder

	runes := []rune(sentence)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == ' ' || r == '\t' || r == '\n' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else if r == '@' || r == '_' {
			// Keep email addresses and underscores as part of word
			current.WriteRune(r)
		} else if r == '.' && current.Len() > 0 && (strings.Contains(current.String(), "@") || i+3 < len(runes)) {
			// Keep dots in email addresses or potential domains (e.g., .com)
			current.WriteRune(r)
		} else if isPunctuationRune(r) {
			// Check for emoticons like :) or ;)
			if (r == ':' || r == ';') && i+1 < len(runes) && runes[i+1] == ')' {
				// Emit current word
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
				// Emit emoticon
				tokens = append(tokens, string(r)+string(runes[i+1]))
				i++ // Skip the ')'
				continue
			}

			// Emit current word
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

			// Emit punctuation as separate token
			tokens = append(tokens, string(r))
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func isPunctuation(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Check for emoticons
	if s == ":)" || s == ":(" || s == ":D" || s == ";)" {
		return true
	}

	if len(s) == 1 {
		return isPunctuationRune(rune(s[0]))
	}
	return false
}

func isPunctuationRune(r rune) bool {
	return r == '.' || r == ',' || r == '!' || r == '?' ||
		   r == ':' || r == ';' || r == ')' || r == '(' ||
		   r == '"' || r == '\'' || r == '-' || r == '…'
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z' || r == 'İ' || r == 'Ş' || r == 'Ğ' ||
		   r == 'Ü' || r == 'Ö' || r == 'Ç'
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = toUpper(runes[0])
	return string(runes)
}

func toUpper(r rune) rune {
	switch r {
	case 'a':
		return 'A'
	case 'b':
		return 'B'
	case 'c':
		return 'C'
	case 'd':
		return 'D'
	case 'e':
		return 'E'
	case 'f':
		return 'F'
	case 'g':
		return 'G'
	case 'h':
		return 'H'
	case 'i':
		return 'İ' // Turkish i
	case 'j':
		return 'J'
	case 'k':
		return 'K'
	case 'l':
		return 'L'
	case 'm':
		return 'M'
	case 'n':
		return 'N'
	case 'o':
		return 'O'
	case 'p':
		return 'P'
	case 'r':
		return 'R'
	case 's':
		return 'S'
	case 't':
		return 'T'
	case 'u':
		return 'U'
	case 'v':
		return 'V'
	case 'y':
		return 'Y'
	case 'z':
		return 'Z'
	case 'ş':
		return 'Ş'
	case 'ğ':
		return 'Ğ'
	case 'ü':
		return 'Ü'
	case 'ö':
		return 'Ö'
	case 'ç':
		return 'Ç'
	default:
		return r
	}
}

func isLikelyDomain(s string) bool {
	// Check if string looks like a domain (contains .com, .org, etc.)
	lowerS := strings.ToLower(s)
	return strings.HasSuffix(lowerS, ".com") ||
		   strings.HasSuffix(lowerS, ".org") ||
		   strings.HasSuffix(lowerS, ".net") ||
		   strings.HasSuffix(lowerS, ".tr") ||
		   strings.HasSuffix(lowerS, ".io") ||
		   strings.HasSuffix(lowerS, ".co")
}

func cleanPunctuation(s string) string {
	// Remove spaces before punctuation
	s = strings.ReplaceAll(s, " .", ".")
	s = strings.ReplaceAll(s, " ,", ",")
	s = strings.ReplaceAll(s, " !", "!")
	s = strings.ReplaceAll(s, " ?", "?")
	s = strings.ReplaceAll(s, " :", ":")
	s = strings.ReplaceAll(s, " ;", ";")
	s = strings.ReplaceAll(s, " )", ")")
	s = strings.ReplaceAll(s, "( ", "(")

	// Ensure space after punctuation (except for special cases)
	s = strings.ReplaceAll(s, ".,", ", ")
	s = strings.ReplaceAll(s, ".!", "!")

	// Handle emoticons - don't add space
	s = strings.ReplaceAll(s, ": )", ":)")
	s = strings.ReplaceAll(s, "; )", ";)")

	return s
}
