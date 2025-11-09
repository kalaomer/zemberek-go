package normalization

import (
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

// TurkishSentenceNormalizerWithLexicon uses full lexicon for normalization
type TurkishSentenceNormalizerWithLexicon struct {
	LookupManual   map[string][]string
	Lexicon        *lexicon.RootLexicon
	WordDictionary map[string]bool
	SpellChecker   *CharacterGraphDecoder
	Graph          *CharacterGraph
}

// NewTurkishSentenceNormalizerWithLexicon creates normalizer with full lexicon
func NewTurkishSentenceNormalizerWithLexicon() (*TurkishSentenceNormalizerWithLexicon, error) {
	// Load manual lookup map
	lookupManual := GetDefaultLookupMap()

	// Try to load from file if exists
	fileLookup, err := LoadLookupMap("resources/normalization/candidates-manual.txt")
	if err == nil {
		for k, v := range fileLookup {
			lookupManual[k] = v
		}
	}

	// Load full lexicon (94K+ words)
	lex, err := lexicon.LoadDefaultLexicon()
	if err != nil {
		return nil, err
	}

	// Build word dictionary from lexicon
	allItems := lex.GetAllItems()
	wordDict := make(map[string]bool, len(allItems)*2)

	for _, item := range allItems {
		if item.Lemma != "" {
			wordDict[item.Lemma] = true
			wordDict[turkish.Instance.ToLower(item.Lemma)] = true
		}
	}

	// Add lookup values to dictionary
	for _, candidates := range lookupManual {
		for _, word := range candidates {
			wordDict[word] = true
			wordDict[turkish.Instance.ToLower(word)] = true
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

	return &TurkishSentenceNormalizerWithLexicon{
		LookupManual:   lookupManual,
		Lexicon:        lex,
		WordDictionary: wordDict,
		SpellChecker:   decoder,
		Graph:          graph,
	}, nil
}

// Normalize normalizes a Turkish sentence
func (tsnl *TurkishSentenceNormalizerWithLexicon) Normalize(sentence string) string {
	// Tokenize sentence
	words := tokenizeSentence(sentence)
	normalized := make([]string, 0, len(words))

	for _, word := range words {
		if word == "" {
			continue
		}

		// Preserve punctuation
		if isPunctuation(word) {
			normalized = append(normalized, word)
			continue
		}

		// Preserve email addresses and URLs
		if strings.Contains(word, "@") || isLikelyDomain(word) {
			normalized = append(normalized, word)
			continue
		}

		// Normalize the word
		normalizedWord := tsnl.normalizeWord(word)
		normalized = append(normalized, normalizedWord)
	}

	// Join and clean
	result := strings.Join(normalized, " ")
	result = cleanPunctuation(result)

	return result
}

// normalizeWord normalizes a single word
func (tsnl *TurkishSentenceNormalizerWithLexicon) normalizeWord(word string) string {
	isCapitalized := len(word) > 0 && isUpper(rune(word[0]))
	lowerWord := turkish.Instance.ToLower(word)

	// 1. Manual lookup (highest priority)
	if candidates, ok := tsnl.LookupManual[lowerWord]; ok && len(candidates) > 0 {
		normalized := candidates[0]
		if isCapitalized {
			return capitalize(normalized)
		}
		return normalized
	}

	// 2. Check if word exists in lexicon
	if tsnl.WordDictionary[lowerWord] {
		return word // Already correct
	}

	// 3. Spell checker
	suggestions := tsnl.SpellChecker.GetSuggestions(lowerWord, DiacriticsIgnoringMatcherInstance)
	if len(suggestions) > 0 {
		normalized := suggestions[0]
		if isCapitalized {
			return capitalize(normalized)
		}
		return normalized
	}

	// 4. Return original
	return word
}
