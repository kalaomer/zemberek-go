package normalization

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/lm"
)

// TurkishSpellChecker provides spell checking and suggestion functionality
type TurkishSpellChecker struct {
	Decoder       *CharacterGraphDecoder
	CharMatcher   CharMatcher
	Morphology    interface{} // *morphology.TurkishMorphology
	stemWords     []string
	LanguageModel lm.LanguageModel
}

// NewTurkishSpellChecker creates a new spell checker
func NewTurkishSpellChecker(stemWords []string, endingsPath string, matcher CharMatcher) (*TurkishSpellChecker, error) {
	// Create stem-ending graph
	graph, err := NewStemEndingGraph(stemWords, endingsPath)
	if err != nil {
		return nil, err
	}

	decoder := NewCharacterGraphDecoder(graph.StemGraph)

	return &TurkishSpellChecker{
		Decoder:     decoder,
		CharMatcher: matcher,
		stemWords:   stemWords,
	}, nil
}

// SuggestForWord returns suggestions for a misspelled word
func (tsc *TurkishSpellChecker) SuggestForWord(word string) []string {
	unranked := tsc.getUnrankedSuggestions(word)
	return tsc.rankByEditDistance(word, unranked)
}

// SuggestForWordForNormalization returns suggestions for normalization (alias)
func (tsc *TurkishSpellChecker) SuggestForWordForNormalization(word string, leftContext string, rightContext string) []string {
	return tsc.SuggestForWordWithContext(word, leftContext, rightContext)
}

// SuggestForWordWithContext returns suggestions with context awareness
func (tsc *TurkishSpellChecker) SuggestForWordWithContext(word string, previous string, next string) []string {
	unranked := tsc.getUnrankedSuggestions(word)
	if len(unranked) == 0 {
		return unranked
	}

	if tsc.LanguageModel == nil {
		return tsc.rankByEditDistance(word, unranked)
	}

	order := tsc.LanguageModel.GetOrder()
	if order < 1 {
		return tsc.rankByEditDistance(word, unranked)
	}

	if order == 1 {
		return tsc.rankByUnigramProbability(word, unranked)
	}

	return tsc.rankByContextProbability(word, unranked, previous, next)
}

// getUnrankedSuggestions gets raw suggestions from decoder
func (tsc *TurkishSpellChecker) getUnrankedSuggestions(word string) []string {
	// Remove apostrophes
	re := regexp.MustCompile(`['']`)
	normalized := re.ReplaceAllString(word, "")

	// Normalize to Turkish alphabet
	normalized = turkish.Instance.Normalize(normalized)

	// Get suggestions from decoder
	suggestions := tsc.Decoder.GetSuggestions(normalized, tsc.CharMatcher)

	// Get case type
	caseType := guessCase(word)

	// Format results to match case
	results := make(map[string]bool)
	for _, suggestion := range suggestions {
		formatted := formatToCase(suggestion, caseType)
		results[formatted] = true
	}

	// Convert to slice
	resultSlice := make([]string, 0, len(results))
	for result := range results {
		resultSlice = append(resultSlice, result)
	}

	return resultSlice
}

// rankByEditDistance ranks suggestions by edit distance
func (tsc *TurkishSpellChecker) rankByEditDistance(original string, suggestions []string) []string {
	if len(suggestions) == 0 {
		return suggestions
	}

	type scoredSuggestion struct {
		word     string
		distance int
	}

	scored := make([]scoredSuggestion, len(suggestions))
	for i, suggestion := range suggestions {
		scored[i] = scoredSuggestion{
			word:     suggestion,
			distance: levenshteinDistance(turkish.Instance.ToLower(original), turkish.Instance.ToLower(suggestion)),
		}
	}

	// Sort by distance (lower is better)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].distance < scored[j].distance
	})

	// Extract words
	result := make([]string, len(scored))
	for i, s := range scored {
		result[i] = s.word
	}

	return result
}

func (tsc *TurkishSpellChecker) rankByUnigramProbability(original string, suggestions []string) []string {
	lmModel := tsc.LanguageModel
	if lmModel == nil || len(suggestions) == 0 {
		return suggestions
	}

	base := tsc.rankByEditDistance(original, append([]string(nil), suggestions...))
	baseRank := make(map[string]int, len(base))
	for i, s := range base {
		baseRank[s] = i
	}

	vocab := lmModel.GetVocabulary()
	type scored struct {
		candidate string
		score     float32
	}
	results := make([]scored, 0, len(suggestions))

	for _, cand := range suggestions {
		normalized := NormalizeForLM(cand)
		idx := vocab.IndexOf(normalized)
		score := lmModel.GetProbability([]int{idx})
		results = append(results, scored{candidate: cand, score: score})
	}

	sort.SliceStable(results, func(i, j int) bool {
		if math.Abs(float64(results[i].score-results[j].score)) < 1e-4 {
			return baseRank[results[i].candidate] < baseRank[results[j].candidate]
		}
		return results[i].score > results[j].score
	})

	ordered := make([]string, len(results))
	for i, r := range results {
		ordered[i] = r.candidate
	}
	return ordered
}

func (tsc *TurkishSpellChecker) rankByContextProbability(original string, suggestions []string, previous, next string) []string {
	lmModel := tsc.LanguageModel
	if lmModel == nil || len(suggestions) == 0 {
		return suggestions
	}

	base := tsc.rankByEditDistance(original, append([]string(nil), suggestions...))
	baseRank := make(map[string]int, len(base))
	for i, s := range base {
		baseRank[s] = i
	}

	vocab := lmModel.GetVocabulary()
	left := previous
	if strings.TrimSpace(left) == "" {
		left = vocab.SentenceStart
	} else {
		left = NormalizeForLM(left)
	}
	right := next
	if strings.TrimSpace(right) == "" {
		right = vocab.SentenceEnd
	} else {
		right = NormalizeForLM(right)
	}

	leftIdx := vocab.IndexOf(left)
	rightIdx := vocab.IndexOf(right)
	order := lmModel.GetOrder()

	type scored struct {
		candidate string
		score     float32
	}
	results := make([]scored, 0, len(suggestions))

	for _, cand := range suggestions {
		norm := NormalizeForLM(cand)
		idx := vocab.IndexOf(norm)
		var score float32
		if order >= 3 {
			score = lmModel.GetProbability([]int{leftIdx, idx, rightIdx})
		} else {
			score = lmModel.GetProbability([]int{leftIdx, idx}) + lmModel.GetProbability([]int{idx, rightIdx})
		}
		results = append(results, scored{candidate: cand, score: score})
	}

	sort.SliceStable(results, func(i, j int) bool {
		diff := results[i].score - results[j].score
		if math.Abs(float64(diff)) < 1e-4 {
			return baseRank[results[i].candidate] < baseRank[results[j].candidate]
		}
		return diff > 0
	})

	ordered := make([]string, len(results))
	for i, r := range results {
		ordered[i] = r.candidate
	}
	return ordered
}

// CaseType represents text case
type CaseType int

const (
	DefaultCase CaseType = iota
	LowerCase
	UpperCase
	TitleCase
	MixedCase
)

// guessCase guesses the case type of a word
func guessCase(word string) CaseType {
	if word == "" {
		return DefaultCase
	}

	runes := []rune(word)
	upperCount := 0
	lowerCount := 0

	for _, r := range runes {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				upperCount++
			} else {
				lowerCount++
			}
		}
	}

	if upperCount == 0 && lowerCount > 0 {
		return LowerCase
	}
	if lowerCount == 0 && upperCount > 0 {
		return UpperCase
	}
	if upperCount == 1 && lowerCount > 0 && unicode.IsUpper(runes[0]) {
		return TitleCase
	}
	if upperCount > 0 && lowerCount > 0 {
		return MixedCase
	}

	return DefaultCase
}

// formatToCase formats word to match case type
func formatToCase(word string, caseType CaseType) string {
	if word == "" {
		return word
	}

	switch caseType {
	case LowerCase:
		return turkish.Instance.ToLower(word)
	case UpperCase:
		return strings.ToUpper(word)
	case TitleCase:
		return turkish.Capitalize(word)
	case MixedCase:
		return word // Keep original for mixed case
	default:
		return word
	}
}

// levenshteinDistance calculates edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Create matrix
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len1][len2]
}

// min3 returns minimum of three integers
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// Check checks if word is spelled correctly
func (tsc *TurkishSpellChecker) Check(word string) bool {
	normalized := turkish.Instance.Normalize(word)
	for _, stem := range tsc.stemWords {
		if stem == normalized {
			return true
		}
	}
	return false
}

// NormalizeForLM normalizes word for language model
func NormalizeForLM(s string) string {
	// If has apostrophe, capitalize
	if strings.Contains(s, "'") || strings.Contains(s, "'") {
		return turkish.Capitalize(s)
	}
	// Otherwise lowercase
	return turkish.Instance.ToLower(s)
}

// GetApostrophe returns the apostrophe character used in word
func GetApostrophe(input string) string {
	// Right single quotation mark (U+2019)
	if strings.ContainsRune(input, '\u2019') {
		return "'"
	}
	// Standard apostrophe
	if strings.ContainsRune(input, '\'') {
		return "'"
	}
	return ""
}

// RankByFrequency ranks suggestions by frequency (requires frequency map)
func (tsc *TurkishSpellChecker) RankByFrequency(suggestions []string, frequencies map[string]int) []string {
	if len(suggestions) == 0 {
		return suggestions
	}

	type scoredSuggestion struct {
		word      string
		frequency int
	}

	scored := make([]scoredSuggestion, len(suggestions))
	for i, suggestion := range suggestions {
		freq := 0
		if f, exists := frequencies[suggestion]; exists {
			freq = f
		}
		scored[i] = scoredSuggestion{
			word:      suggestion,
			frequency: freq,
		}
	}

	// Sort by frequency (higher is better)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].frequency > scored[j].frequency
	})

	// Extract words
	result := make([]string, len(scored))
	for i, s := range scored {
		result[i] = s.word
	}

	return result
}
