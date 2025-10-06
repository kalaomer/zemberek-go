package normalization

import (
	"bufio"
	"os"
	"strings"

	"github.com/kalaomer/zemberek-go/normalization/deasciifier"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// turkishLowerCaser is a reusable Turkish lowercase converter
var turkishLowerCaser = cases.Lower(language.Turkish)

// turkishToLower converts a string to lowercase using Turkish locale rules
// This is important because Turkish has special handling for I/İ and ı/i
func turkishToLower(s string) string {
	return turkishLowerCaser.String(s)
}

// TurkishSentenceNormalizer normalizes informal Turkish sentences
type TurkishSentenceNormalizer struct {
	SpellChecker          *TurkishSpellChecker
	Replacements          map[string]string
	NoSplitWords          map[string]bool
	CommonSplits          map[string]string
	CommonConnectedSuffixes map[string]bool
	LookupManual          map[string][]string
	LookupFromGraph       map[string][]string
	LookupFromASCII       map[string][]string
	AlwaysApplyDeasciifier bool
	stemWords             []string
}

// NewTurkishSentenceNormalizer creates a new sentence normalizer
func NewTurkishSentenceNormalizer(stemWords []string, resourcesPath string) (*TurkishSentenceNormalizer, error) {
	tsn := &TurkishSentenceNormalizer{
		Replacements:          make(map[string]string),
		NoSplitWords:          make(map[string]bool),
		CommonSplits:          make(map[string]string),
		CommonConnectedSuffixes: make(map[string]bool),
		LookupManual:          make(map[string][]string),
		LookupFromGraph:       make(map[string][]string),
		LookupFromASCII:       make(map[string][]string),
		AlwaysApplyDeasciifier: false,
		stemWords:             stemWords,
	}

	// Load resources
	if resourcesPath == "" {
		resourcesPath = "resources/normalization"
	}

	// Load replacements
	if err := tsn.loadReplacements(resourcesPath + "/multi-word-replacements.txt"); err == nil {
		// Ignore error if file doesn't exist
	}

	// Load no-split words
	if err := tsn.loadNoSplit(resourcesPath + "/no-split.txt"); err == nil {
		// Ignore error
	}

	// Load common splits
	if err := tsn.loadCommonSplits(resourcesPath + "/split.txt"); err == nil {
		// Ignore error
	}

	// Load connected suffixes
	if err := tsn.loadConnectedSuffixes(resourcesPath + "/question-suffixes.txt"); err == nil {
		// Ignore error
	}

	// Load lookup maps
	if err := tsn.loadMultimap(resourcesPath+"/candidates-manual.txt", tsn.LookupManual); err == nil {
		// Ignore error
	}
	if err := tsn.loadMultimap(resourcesPath+"/lookup-from-graph.txt", tsn.LookupFromGraph); err == nil {
		// Ignore error
	}
	if err := tsn.loadMultimap(resourcesPath+"/ascii-map.txt", tsn.LookupFromASCII); err == nil {
		// Ignore error
	}

	// Create spell checker
	spellChecker, err := NewTurkishSpellChecker(stemWords, resourcesPath+"/endings.txt", DiacriticsIgnoringMatcherInstance)
	if err != nil {
		// Use nil spell checker if creation fails
		spellChecker = nil
	}
	tsn.SpellChecker = spellChecker

	return tsn, nil
}

// Normalize normalizes a sentence
func (tsn *TurkishSentenceNormalizer) Normalize(sentence string) string {
	processed := tsn.preProcess(sentence)
	tokens := tokenize(processed)

	normalized := make([]string, 0, len(tokens))
	for _, token := range tokens {
		candidates := tsn.getCandidates(token)
		if len(candidates) > 0 {
			normalized = append(normalized, candidates[0]) // Take best candidate
		} else {
			normalized = append(normalized, token)
		}
	}

	return strings.Join(normalized, " ")
}

// preProcess performs preprocessing on sentence
func (tsn *TurkishSentenceNormalizer) preProcess(sentence string) string {
	// Convert to lowercase using Turkish locale (important for I/İ and ı/i)
	sentence = turkishToLower(sentence)

	// Replace common phrases
	tokens := tokenize(sentence)
	sentence = tsn.replaceCommon(tokens)

	// Combine necessary words
	tokens = tokenize(sentence)
	sentence = tsn.combineNecessaryWords(tokens)

	// Split necessary words
	tokens = tokenize(sentence)
	sentence = tsn.splitNecessaryWords(tokens, false)

	// Apply deasciifier if needed
	if tsn.AlwaysApplyDeasciifier || probablyRequiresDeasciifier(sentence) {
		d := deasciifier.NewDeasciifier(sentence)
		sentence = d.ConvertToTurkish()
	}

	// Combine again after deasciification
	tokens = tokenize(sentence)
	sentence = tsn.combineNecessaryWords(tokens)

	// Split with lookup
	tokens = tokenize(sentence)
	sentence = tsn.splitNecessaryWords(tokens, true)

	return sentence
}

// getCandidates gets normalization candidates for a word
func (tsn *TurkishSentenceNormalizer) getCandidates(word string) []string {
	candidates := make([]string, 0)
	seen := make(map[string]bool)

	// Add from lookup maps
	for _, candidate := range tsn.LookupManual[word] {
		if !seen[candidate] {
			candidates = append(candidates, candidate)
			seen[candidate] = true
		}
	}
	for _, candidate := range tsn.LookupFromGraph[word] {
		if !seen[candidate] {
			candidates = append(candidates, candidate)
			seen[candidate] = true
		}
	}
	for _, candidate := range tsn.LookupFromASCII[word] {
		if !seen[candidate] {
			candidates = append(candidates, candidate)
			seen[candidate] = true
		}
	}

	// Add from spell checker
	if tsn.SpellChecker != nil && len(word) > 3 {
		suggestions := tsn.SpellChecker.SuggestForWord(word)
		if len(suggestions) > 3 {
			suggestions = suggestions[:3]
		}
		for _, suggestion := range suggestions {
			if !seen[suggestion] {
				candidates = append(candidates, suggestion)
				seen[suggestion] = true
			}
		}
	}

	// Add original if nothing found
	if len(candidates) == 0 {
		candidates = append(candidates, word)
	}

	return candidates
}

// replaceCommon replaces common multi-word phrases
func (tsn *TurkishSentenceNormalizer) replaceCommon(tokens []string) string {
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if replacement, exists := tsn.Replacements[token]; exists {
			result = append(result, replacement)
		} else {
			result = append(result, token)
		}
	}
	return strings.Join(result, " ")
}

// combineNecessaryWords combines words that should be together
func (tsn *TurkishSentenceNormalizer) combineNecessaryWords(tokens []string) string {
	if len(tokens) < 2 {
		return strings.Join(tokens, " ")
	}

	result := make([]string, 0, len(tokens))
	combined := false

	for i := 0; i < len(tokens)-1; i++ {
		if combined {
			combined = false
			continue
		}

		first := tokens[i]
		second := tokens[i+1]

		// Try to combine
		if isWord(first) && isWord(second) {
			comb := tsn.combineCommon(first, second)
			if comb != "" {
				result = append(result, comb)
				combined = true
				continue
			}
		}

		result = append(result, first)
	}

	if !combined {
		result = append(result, tokens[len(tokens)-1])
	}

	return strings.Join(result, " ")
}

// combineCommon tries to combine two words
func (tsn *TurkishSentenceNormalizer) combineCommon(w1, w2 string) string {
	combined := w1 + w2

	// Check if combined form is valid
	for _, stem := range tsn.stemWords {
		if stem == combined {
			return combined
		}
	}

	// Check if starts with apostrophe or "bil"
	if strings.HasPrefix(w2, "'") || strings.HasPrefix(w2, "bil") {
		return combined
	}

	return ""
}

// splitNecessaryWords splits words that should be separated
func (tsn *TurkishSentenceNormalizer) splitNecessaryWords(tokens []string, useLookup bool) string {
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if isWord(token) {
			result = append(result, tsn.separateCommon(token, useLookup))
		} else {
			result = append(result, token)
		}
	}
	return strings.Join(result, " ")
}

// separateCommon separates common suffixes
func (tsn *TurkishSentenceNormalizer) separateCommon(word string, useLookup bool) string {
	// Check no-split list
	if tsn.NoSplitWords[word] {
		return word
	}

	// Check common splits
	if useLookup {
		if split, exists := tsn.CommonSplits[word]; exists {
			return split
		}
	}

	// Try to separate common suffixes
	for i := 1; i < len(word); i++ {
		tail := word[i:]
		if tsn.CommonConnectedSuffixes[tail] {
			head := word[:i]
			// Check if head is valid
			for _, stem := range tsn.stemWords {
				if stem == head {
					return head + " " + tail
				}
			}
		}
	}

	return word
}

// Helper functions

// loadReplacements loads replacement map
func (tsn *TurkishSentenceNormalizer) loadReplacements(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			tsn.Replacements[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return scanner.Err()
}

// loadNoSplit loads no-split words
func (tsn *TurkishSentenceNormalizer) loadNoSplit(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tsn.NoSplitWords[line] = true
		}
	}
	return scanner.Err()
}

// loadCommonSplits loads common splits
func (tsn *TurkishSentenceNormalizer) loadCommonSplits(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "-")
		if len(parts) == 2 {
			tsn.CommonSplits[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return scanner.Err()
}

// loadConnectedSuffixes loads connected suffixes
func (tsn *TurkishSentenceNormalizer) loadConnectedSuffixes(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tsn.CommonConnectedSuffixes[line] = true
		}
	}
	return scanner.Err()
}

// loadMultimap loads multimap
func (tsn *TurkishSentenceNormalizer) loadMultimap(path string, target map[string][]string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		if strings.Contains(value, ",") {
			values := strings.Split(value, ",")
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}
			target[key] = values
		} else {
			if existing, exists := target[key]; exists {
				target[key] = append(existing, value)
			} else {
				target[key] = []string{value}
			}
		}
	}
	return scanner.Err()
}

// tokenize splits sentence into tokens
func tokenize(sentence string) []string {
	// Simple tokenization by spaces
	return strings.Fields(sentence)
}

// isWord checks if token is a word
func isWord(token string) bool {
	if token == "" {
		return false
	}
	for _, r := range token {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			r == 'ç' || r == 'Ç' || r == 'ğ' || r == 'Ğ' || r == 'ı' || r == 'İ' ||
			r == 'ö' || r == 'Ö' || r == 'ş' || r == 'Ş' || r == 'ü' || r == 'Ü' ||
			r == '\'' || r == '-') {
			return false
		}
	}
	return true
}

// probablyRequiresDeasciifier checks if sentence needs deasciification
func probablyRequiresDeasciifier(sentence string) bool {
	turkishSpecCount := 0
	total := 0

	for _, c := range sentence {
		if c == ' ' {
			continue
		}
		total++
		if c == 'ç' || c == 'Ç' || c == 'ğ' || c == 'Ğ' || c == 'ö' || c == 'Ö' ||
			c == 'ş' || c == 'Ş' || c == 'ü' || c == 'Ü' {
			turkishSpecCount++
		}
	}

	if total == 0 {
		return false
	}

	ratio := float64(turkishSpecCount) / float64(total)
	return ratio < 0.1
}

// Candidate represents a normalization candidate
type Candidate struct {
	Content string
	Score   float32
}

// NewCandidate creates a new candidate
func NewCandidate(content string) *Candidate {
	return &Candidate{
		Content: content,
		Score:   1.0,
	}
}

// Candidates represents multiple candidates for a word
type Candidates struct {
	Word       string
	Candidates []*Candidate
}

// NewCandidates creates a new candidates structure
func NewCandidates(word string, candidates []*Candidate) *Candidates {
	return &Candidates{
		Word:       word,
		Candidates: candidates,
	}
}

// Hypothesis represents a normalization hypothesis in beam search
type Hypothesis struct {
	History  []*Candidate
	Current  *Candidate
	Previous *Hypothesis
	Score    float32
}

// NewHypothesis creates a new hypothesis
func NewHypothesis() *Hypothesis {
	return &Hypothesis{
		History:  nil,
		Current:  nil,
		Previous: nil,
		Score:    0.0,
	}
}

// Equals checks if two hypotheses are equal
func (h *Hypothesis) Equals(other *Hypothesis) bool {
	if h == other {
		return true
	}
	if other == nil {
		return false
	}

	// Check history
	if len(h.History) != len(other.History) {
		return false
	}
	for i, c := range h.History {
		if c != other.History[i] {
			return false
		}
	}

	// Check current
	if h.Current != other.Current {
		return false
	}

	return true
}

// Hash returns hash for hypothesis
func (h *Hypothesis) Hash() int {
	result := 0
	for _, c := range h.History {
		if c != nil {
			for _, r := range c.Content {
				result = 31*result + int(r)
			}
		}
	}
	if h.Current != nil {
		for _, r := range h.Current.Content {
			result = 31*result + int(r)
		}
	}
	return result
}

// GetStartCandidate returns the START sentinel candidate
func GetStartCandidate() *Candidate {
	return &Candidate{Content: "<s>", Score: 1.0}
}

// GetEndCandidate returns the END sentinel candidate
func GetEndCandidate() *Candidate {
	return &Candidate{Content: "</s>", Score: 1.0}
}

// GetEndCandidates returns the END candidates structure
func GetEndCandidates() *Candidates {
	return &Candidates{
		Word:       "</s>",
		Candidates: []*Candidate{GetEndCandidate()},
	}
}

// NormalizeWithBeamSearch normalizes sentence using beam search (simplified without LM)
func (tsn *TurkishSentenceNormalizer) NormalizeWithBeamSearch(sentence string) string {
	processed := tsn.preProcess(sentence)
	tokens := tokenize(processed)

	// Get candidates for each token
	candidatesList := make([]*Candidates, 0, len(tokens))
	for _, token := range tokens {
		candidateStrs := tsn.getCandidates(token)
		candidates := make([]*Candidate, len(candidateStrs))
		for i, c := range candidateStrs {
			candidates[i] = NewCandidate(c)
		}
		candidatesList = append(candidatesList, NewCandidates(token, candidates))
	}

	// Decode using beam search (simplified - without language model)
	result := tsn.decodeSimple(candidatesList)
	return strings.Join(result, " ")
}

// decodeSimple performs simple decoding (takes first/best candidate)
func (tsn *TurkishSentenceNormalizer) decodeSimple(candidatesList []*Candidates) []string {
	result := make([]string, len(candidatesList))
	for i, candidates := range candidatesList {
		if len(candidates.Candidates) > 0 {
			result[i] = candidates.Candidates[0].Content
		} else {
			result[i] = candidates.Word
		}
	}
	return result
}

// GetBestHypothesis returns best hypothesis from list
func GetBestHypothesis(hypotheses []*Hypothesis) *Hypothesis {
	if len(hypotheses) == 0 {
		return nil
	}

	best := hypotheses[0]
	for _, h := range hypotheses[1:] {
		if h.Score > best.Score {
			best = h
		}
	}
	return best
}
