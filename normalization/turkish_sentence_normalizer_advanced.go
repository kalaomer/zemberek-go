package normalization

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/lm"
	"github.com/kalaomer/zemberek-go/morphology"
	"github.com/kalaomer/zemberek-go/morphology/analysis"
	"github.com/kalaomer/zemberek-go/normalization/deasciifier"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TurkishSentenceNormalizerAdvanced normalizes informal Turkish sentences with full morphology support
type TurkishSentenceNormalizerAdvanced struct {
	SpellChecker            *TurkishSpellChecker
	Morphology              *morphology.TurkishMorphology
	InformalMorphology      *morphology.TurkishMorphology
	AnalysisConverter       *analysis.InformalAnalysisConverter
	LanguageModel           lm.LanguageModel
	Replacements            map[string]string
	NoSplitWords            map[string]bool
	CommonSplits            map[string]string
	CommonConnectedSuffixes map[string]bool
	LookupManual            map[string][]string
	LookupFromGraph         map[string][]string
	LookupFromASCII         map[string][]string
	AlwaysApplyDeasciifier  bool
}

// NewTurkishSentenceNormalizerAdvanced creates a new advanced sentence normalizer with morphology
func NewTurkishSentenceNormalizerAdvanced(morph *morphology.TurkishMorphology, dataRoot string) (*TurkishSentenceNormalizerAdvanced, error) {
	tsn := &TurkishSentenceNormalizerAdvanced{
		Morphology:              morph,
		Replacements:            make(map[string]string),
		NoSplitWords:            make(map[string]bool),
		CommonSplits:            make(map[string]string),
		CommonConnectedSuffixes: make(map[string]bool),
		LookupManual:            make(map[string][]string),
		LookupFromGraph:         make(map[string][]string),
		LookupFromASCII:         make(map[string][]string),
		AlwaysApplyDeasciifier:  false,
	}

	// Create informal morphology
	informalMorph := morphology.NewBuilder(morph.Lexicon).
		UseInformalAnalysis().
		IgnoreDiacriticsInAnalysis().
		Build()
	tsn.InformalMorphology = informalMorph

	// Create analysis converter
	tsn.AnalysisConverter = analysis.NewInformalAnalysisConverter(morph.WordGenerator)

	// Resolve paths
	if dataRoot == "" {
		// default to repo resources root
		dataRoot = "resources"
	}
	normRoot := filepath.Join(dataRoot, "normalization")
	// Language model path: prefer dataRoot/lm.2gram.slm then dataRoot/lm/lm.2gram.slm
	lmPath := filepath.Join(dataRoot, "lm.2gram.slm")
	if _, err := os.Stat(lmPath); err != nil {
		alt := filepath.Join(dataRoot, "lm", "lm.2gram.slm")
		if _, err2 := os.Stat(alt); err2 == nil {
			lmPath = alt
		}
	}

	// Load language model (required for advanced decoding)
	langModel, err := lm.LoadFromFile(lmPath)
	if err != nil {
		return nil, fmt.Errorf("load language model from %s: %w", lmPath, err)
	}
	tsn.LanguageModel = langModel

	// Load all resource files (support extensionless names as in Java data)
	tsn.loadReplacements(firstExisting(normRoot, "multi-word-replacements.txt", "multi-word-replacements"))
	tsn.loadNoSplit(firstExisting(normRoot, "no-split.txt", "no-split"))
	tsn.loadCommonSplits(firstExisting(normRoot, "split.txt", "split"))
	tsn.loadConnectedSuffixes(firstExisting(normRoot, "question-suffixes.txt", "question-suffixes"))
	tsn.loadMultimap(firstExisting(normRoot, "candidates-manual.txt", "candidates-manual"), tsn.LookupManual)
	tsn.loadMultimap(firstExisting(normRoot, "lookup-from-graph.txt", "lookup-from-graph"), tsn.LookupFromGraph)
	tsn.loadMultimap(firstExisting(normRoot, "ascii-map.txt", "ascii-map"), tsn.LookupFromASCII)

	skipManualDefaults := map[string]struct{}{
		"annemde": {},
	}

	for key, values := range GetDefaultLookupMap() {
		if len(values) == 0 {
			continue
		}
		if tsn.Morphology.HasRegularAnalysis(key) {
			if _, blocked := skipManualDefaults[key]; blocked {
				continue
			}
		}
		existing := append([]string{}, tsn.LookupManual[key]...)
	valueLoop:
		for _, v := range values {
			for _, ex := range existing {
				if ex == v {
					continue valueLoop
				}
			}
			existing = append(existing, v)
		}
		tsn.LookupManual[key] = existing
	}

	// Built-in minimal manual expansions (Java has candidates-manual in resources; data/ may not include it)
	if _, ok := tsn.LookupManual["tmm"]; !ok {
		tsn.LookupManual["tmm"] = []string{"tamam"}
	} else {
		tsn.LookupManual["tmm"] = append(tsn.LookupManual["tmm"], "tamam")
	}

	// Create spell checker with morphology
	graph := NewCharacterGraph()
	decoder := NewCharacterGraphDecoder(graph)
	tsn.SpellChecker = &TurkishSpellChecker{
		Morphology:    morph,
		Decoder:       decoder,
		CharMatcher:   DiacriticsIgnoringMatcherInstance,
		LanguageModel: langModel,
	}

	return tsn, nil
}

// Normalize normalizes a sentence using full morphological analysis and beam search
func (tsn *TurkishSentenceNormalizerAdvanced) Normalize(sentence string) string {
	processed := tsn.preProcess(sentence)
	tokens := tokenizeAdvanced(processed)

	// Get candidates for each token
	candidatesList := make([]*Candidates, 0, len(tokens))

	for i, token := range tokens {
		var previous, next string
		if i > 0 {
			previous = tokens[i-1]
		}
		if i < len(tokens)-1 {
			next = tokens[i+1]
		}

		candidateStrs := tsn.getCandidatesAdvanced(token, previous, next)
		candidates := make([]*Candidate, len(candidateStrs))
		for j, c := range candidateStrs {
			candidates[j] = NewCandidate(c)
		}
		candidatesList = append(candidatesList, NewCandidates(token, candidates))
	}

	// Decode using beam search with language model
	result := tsn.decode(candidatesList)
	return joinTokensWithPunct(result)
}

// getCandidatesAdvanced gets normalization candidates using morphological analysis
func (tsn *TurkishSentenceNormalizerAdvanced) getCandidatesAdvanced(word, previous, next string) []string {
	// Keep punctuation as-is
	if !isWordAdvanced(word) {
		return []string{word}
	}
	candidates := make([]string, 0)
	seen := make(map[string]bool)

	// Heuristic formalization for informal future (only if base is a verb)
	for _, cand := range tsn.expandInformalFutureVerbOnly(word) {
		if cand != word && !seen[cand] {
			candidates = append(candidates, cand)
			seen[cand] = true
		}
	}
	for _, cand := range tsn.expandInformalProgressive(word) {
		if cand != word && !seen[cand] {
			candidates = append(candidates, cand)
			seen[cand] = true
		}
	}
	for _, cand := range tsn.expandQuestionParticle(word, previous) {
		if cand != word && !seen[cand] {
			candidates = append(candidates, cand)
			seen[cand] = true
		}
	}

	// Add from lookup maps (highest priority)
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

	// Always include spell-checker suggestions (not only when no analysis), limited
	if tsn.SpellChecker != nil && len(word) > 3 {
		suggestions := tsn.SpellChecker.SuggestForWordWithContext(word, previous, next)
		if len(suggestions) > 3 {
			suggestions = suggestions[:3]
		}
		for _, s := range suggestions {
			if s != word && !seen[s] {
				candidates = append(candidates, s)
				seen[s] = true
			}
		}
	}

	// Heuristic: single 'l' insertion before final 'a' or 'e' if that yields a valid word (e.g., okua → okula)
	// Apply only when original word has no regular analysis (to avoid changing valid forms like "havuza").
	if len(word) > 3 && !tsn.Morphology.HasRegularAnalysis(word) {
		r := []rune(word)
		last := r[len(r)-1]
		if last == 'a' || last == 'e' {
			cand := string(r[:len(r)-1]) + "l" + string(last)
			if cand != word && !seen[cand] && tsn.Morphology.HasRegularAnalysis(cand) {
				candidates = append(candidates, cand)
				seen[cand] = true
			}
		}
	}

	// Analyze with informal morphology
	analyses := tsn.InformalMorphology.Analyze(word)
	for _, sa := range analyses.AnalysisResults {
		if sa.ContainsInformalMorpheme() {
			// Convert informal to formal
			result := tsn.AnalysisConverter.Convert(word, sa)
			if result != nil && result.Surface != word && !seen[result.Surface] {
				candidates = append(candidates, result.Surface)
				seen[result.Surface] = true
			}
		} else {
			// Generate from morphemes
			results := tsn.Morphology.WordGenerator.Generate(sa.Item, sa.GetMorphemes())
			for _, r := range results {
				if r.Surface != word && !seen[r.Surface] {
					candidates = append(candidates, r.Surface)
					seen[r.Surface] = true
				}
			}
		}
	}

	// If no analysis and word > 3 chars, use spell checker
	if len(analyses.AnalysisResults) == 0 && len(word) > 3 && tsn.SpellChecker != nil {
		suggestions := tsn.SpellChecker.SuggestForWordWithContext(word, previous, next)
		if len(suggestions) > 3 {
			suggestions = suggestions[:3]
		}
		for _, sugg := range suggestions {
			if !seen[sugg] {
				candidates = append(candidates, sugg)
				seen[sugg] = true
			}
		}
	}

	// Add original only if no candidate was generated
	if len(candidates) == 0 {
		if !seen[word] {
			candidates = append(candidates, word)
			seen[word] = true
		}
	}

	return candidates
}

// decode performs beam search decoding with language model
func (tsn *TurkishSentenceNormalizerAdvanced) decode(candidatesList []*Candidates) []string {
	if tsn.LanguageModel == nil {
		// Fallback to simple decoding
		return tsn.decodeSimple(candidatesList)
	}

	current := make([]*Hypothesis, 0)
	next := make([]*Hypothesis, 0)

	// Add END candidates
	candidatesList = append(candidatesList, GetEndCandidates())

	// Initial hypothesis
	lmOrder := tsn.LanguageModel.GetOrder()
	initial := NewHypothesis()
	initial.History = make([]*Candidate, lmOrder-1)
	for i := 0; i < lmOrder-1; i++ {
		initial.History[i] = GetStartCandidate()
	}
	initial.Current = GetStartCandidate()
	initial.Score = 0.0
	current = append(current, initial)

	// Process each candidate set
	for _, candidates := range candidatesList {
		for _, h := range current {
			for _, c := range candidates.Candidates {
				newHyp := NewHypothesis()

				// Update history
				hist := make([]*Candidate, lmOrder-1)
				if lmOrder > 2 {
					copy(hist, h.History[1:])
				}
				hist[len(hist)-1] = h.Current
				newHyp.History = hist
				newHyp.Current = c
				newHyp.Previous = h

				// Calculate score using LM + morphological prior
				indexes := make([]int, lmOrder)
				vocab := tsn.LanguageModel.GetVocabulary()
				for j := 0; j < lmOrder-1; j++ {
					indexes[j] = vocab.IndexOf(hist[j].Content)
				}
				indexes[lmOrder-1] = vocab.IndexOf(c.Content)

				lmScore := tsn.LanguageModel.GetProbability(indexes)

				// Morphological prior: prefer tokens with regular analysis, penalize unknowns
				var morphPrior float32 = 0.0
				if tsn.Morphology.HasRegularAnalysis(c.Content) {
					morphPrior += 1.5 // strong preference for valid words
				} else {
					morphPrior -= 0.8 // penalize invalid forms
				}
				// Diacritic prior: prefer Turkish-specific letters; penalize circumflex
				if turkish.Instance.ContainsASCIIRelated(c.Content) && tsn.Morphology.HasRegularAnalysis(c.Content) {
					morphPrior += 0.5
				}
				if strings.ContainsAny(c.Content, "âîûÂÎÛ") {
					morphPrior -= 0.6
				}

				// Future-form preference vs informal '-cam'
				if isFormalFuture(c.Content) {
					morphPrior += 2.5
				}
				if isInformalCam(c.Content) {
					morphPrior -= 1.5
				}
				// Simple contextual prior: immediate previous token is a time/adverb (e.g., 'yarın') or 'kadar'
				prevToken := ""
				if hist[len(hist)-1] != nil {
					prevToken = hist[len(hist)-1].Content
				}
				if prevToken == "kadar" || prevToken == "yarın" || prevToken == "akşama" {
					if isFormalFuture(c.Content) {
						morphPrior += 0.7
					}
					if isInformalCam(c.Content) {
						morphPrior -= 0.2
					}
				}

				score := lmScore + morphPrior
				newHyp.Score = h.Score + score

				// Add or update in next list
				found := false
				for i, existing := range next {
					if newHyp.Equals(existing) {
						if newHyp.Score > existing.Score {
							next[i] = newHyp
						}
						found = true
						break
					}
				}
				if !found {
					next = append(next, newHyp)
				}
			}
		}

		current = next
		next = make([]*Hypothesis, 0)
	}

	// Get best hypothesis
	best := GetBestHypothesis(current)
	if best == nil {
		return make([]string, 0)
	}

	// Extract sequence (skip START sentinel by content)
	seq := make([]string, 0)
	h := best.Previous
	for h != nil {
		if h.Current != nil && h.Current.Content != "<s>" {
			seq = append([]string{h.Current.Content}, seq...)
		}
		h = h.Previous
	}

	return seq
}

// decodeSimple fallback simple decoding
func (tsn *TurkishSentenceNormalizerAdvanced) decodeSimple(candidatesList []*Candidates) []string {
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

// preProcess performs preprocessing on sentence
var trLowerAdvanced = cases.Lower(language.Turkish)

func (tsn *TurkishSentenceNormalizerAdvanced) preProcess(sentence string) string {
	sentence = trLowerAdvanced.String(sentence)

	tokens := tokenizeAdvanced(sentence)
	sentence = tsn.replaceCommon(tokens)

	tokens = tokenizeAdvanced(sentence)
	sentence = tsn.combineNecessaryWords(tokens)

	tokens = tokenizeAdvanced(sentence)
	sentence = tsn.splitNecessaryWords(tokens, false)

	if tsn.AlwaysApplyDeasciifier || probablyRequiresDeasciifier(sentence) {
		origTokens := tokenizeAdvanced(sentence)
		d := deasciifier.NewDeasciifier(sentence)
		converted := d.ConvertToTurkish()
		replaced := false
		restored := tokenizeAdvanced(converted)
		if len(restored) == len(origTokens) {
			for i := range restored {
				if hasNonASCII(origTokens[i]) {
					restored[i] = origTokens[i]
					replaced = true
				}
			}
		}
		if replaced {
			sentence = joinTokensWithPunct(restored)
		} else {
			sentence = converted
		}
	}

	tokens = tokenizeAdvanced(sentence)
	sentence = tsn.combineNecessaryWords(tokens)

	tokens = tokenizeAdvanced(sentence)
	sentence = tsn.splitNecessaryWords(tokens, true)

	return sentence
}

// expandInformalFutureVerbOnly returns formal future form
// Only returns a candidate if the resulting word has a regular analysis (to avoid false positives like "hocam"→"hocağım").
func (tsn *TurkishSentenceNormalizerAdvanced) expandInformalFutureVerbOnly(word string) []string {
	w := trLowerAdvanced.String(word)
	if len(w) < 5 {
		return nil
	}
	var base string
	var target string
	if strings.HasSuffix(w, "icem") || strings.HasSuffix(w, "ücem") || strings.HasSuffix(w, "ecem") || strings.HasSuffix(w, "ucem") {
		base = w[:len(w)-4]
		// Choose eceğim /acağım by last vowel harmony of base, ignoring a trailing buffer vowel if any
		hv := lastVowelIgnoreTrailingHigh(base)
		if hv == 0 {
			return nil
		}
		if isFrontVowel(hv) {
			target = base + "eceğim"
		} else {
			target = base + "acağım"
		}
		return []string{target}
	} else if strings.HasSuffix(w, "acam") || strings.HasSuffix(w, "ucam") || strings.HasSuffix(w, "ıcam") || strings.HasSuffix(w, "ocam") || strings.HasSuffix(w, "icam") || strings.HasSuffix(w, "ücam") {
		base = w[:len(w)-4]
		// If base ends with a high vowel (i,ı,u,ü) drop it (yatı → yat)
		if endsWithHighVowel(base) {
			r := []rune(base)
			base = string(r[:len(r)-1])
		}
		// Also try removing assimilation chunk like "ic/uc/oc/ac" (…ic → …)
		if stem, ok := stripAssimilationChunk(base); ok {
			base = stem
		}
		hv := lastVowelIgnoreTrailingHigh(base)
		if hv == 0 {
			return nil
		}
		if isFrontVowel(hv) {
			target = base + "eceğim"
		} else {
			target = base + "acağım"
		}
	} else {
		return nil
	}
	// Return target as candidate; beam + morphological prior will down-rank invalid forms
	return []string{target}
}

func (tsn *TurkishSentenceNormalizerAdvanced) expandInformalProgressive(word string) []string {
	w := trLowerAdvanced.String(word)
	if len(w) < 3 {
		return nil
	}
	if strings.Contains(w, "yor") || !strings.Contains(w, "yo") {
		return nil
	}
	idx := strings.LastIndex(w, "yo")
	if idx < 0 || idx+2 > len(w) {
		return nil
	}
	candidate := w[:idx+2] + "r" + w[idx+2:]
	if candidate == w {
		return nil
	}
	analyses := tsn.Morphology.Analyze(candidate)
	for _, sa := range analyses.AnalysisResults {
		if sa.Item != nil && sa.Item.PrimaryPos == turkish.Verb {
			return []string{candidate}
		}
	}
	return nil
}

func (tsn *TurkishSentenceNormalizerAdvanced) expandQuestionParticle(word, previous string) []string {
	w := trLowerAdvanced.String(word)
	if w != "mi" && w != "mı" && w != "mu" && w != "mü" {
		return nil
	}
	prev := trLowerAdvanced.String(strings.TrimSpace(previous))
	if prev == "" {
		return nil
	}
	hv := lastVowelIgnoreTrailingHigh(prev)
	if hv == 0 {
		return nil
	}
	var target string
	if isFrontVowel(hv) {
		if isRoundedVowel(hv) {
			target = "mü"
		} else {
			target = "mi"
		}
	} else {
		if isRoundedVowel(hv) {
			target = "mu"
		} else {
			target = "mı"
		}
	}
	if target == "" || target == w {
		return nil
	}
	return []string{target}
}

func lastVowel(s string) rune {
	if s == "" {
		return 0
	}
	alpha := turkish.Instance
	// Using provided API: get last vowel; if not available, scan runes
	r := []rune(s)
	for i := len(r) - 1; i >= 0; i-- {
		if alpha.Vowels[r[i]] {
			return r[i]
		}
	}
	return 0
}

// lastVowelIgnoreTrailingHigh finds last vowel ignoring a final high vowel (i, ı, u, ü)
func lastVowelIgnoreTrailingHigh(s string) rune {
	if s == "" {
		return 0
	}
	alpha := turkish.Instance
	r := []rune(s)
	end := len(r) - 1
	if end >= 0 {
		switch r[end] {
		case 'i', 'ı', 'u', 'ü', 'İ', 'I', 'U', 'Ü':
			end--
		}
	}
	for i := end; i >= 0; i-- {
		if alpha.Vowels[r[i]] {
			return r[i]
		}
	}
	return 0
}

func isFrontVowel(r rune) bool {
	switch r {
	case 'e', 'i', 'ö', 'ü', 'E', 'İ', 'Ö', 'Ü':
		return true
	default:
		return false
	}
}

func isRoundedVowel(r rune) bool {
	switch r {
	case 'o', 'ö', 'u', 'ü', 'O', 'Ö', 'U', 'Ü':
		return true
	default:
		return false
	}
}

func isFormalFuture(s string) bool {
	return strings.HasSuffix(s, "eceğim") || strings.HasSuffix(s, "acağım") || strings.HasSuffix(s, "ecegim") || strings.HasSuffix(s, "acagim")
}

func isInformalCam(s string) bool {
	if len(s) < 4 {
		return false
	}
	r := []rune(s)
	// ends with 'cam'
	if strings.HasSuffix(s, "cam") {
		// check previous vowel is any of a/e/ı/i/u/ü/o/ö
		for i := len(r) - 4; i >= 0 && i >= len(r)-6; i-- {
			switch r[i] {
			case 'a', 'e', 'ı', 'i', 'u', 'ü', 'o', 'ö', 'A', 'E', 'I', 'İ', 'U', 'Ü', 'O', 'Ö':
				return true
			}
		}
	}
	return false
}

// stripAssimilationChunk tries to remove final consonant+high-vowel assimilation before -cam (e.g., ...+"ic" → stem)
func stripAssimilationChunk(base string) (string, bool) {
	r := []rune(base)
	if len(r) < 2 {
		return base, false
	}
	// common chunks: ic, uc, oc, ac (ASCII forms)
	last := r[len(r)-1]
	prev := r[len(r)-2]
	if (prev == 'i' || prev == 'u' || prev == 'o' || prev == 'a' || prev == 'ı' || prev == 'ü' || prev == 'ö') && (last == 'c') {
		return string(r[:len(r)-2]), true
	}
	return base, false
}

func endsWithHighVowel(s string) bool {
	if s == "" {
		return false
	}
	r := []rune(s)
	last := r[len(r)-1]
	switch last {
	case 'i', 'ı', 'u', 'ü', 'İ', 'I', 'U', 'Ü':
		return true
	default:
		return false
	}
}

// Helper methods (reuse from basic normalizer)

func (tsn *TurkishSentenceNormalizerAdvanced) replaceCommon(tokens []string) string {
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

func (tsn *TurkishSentenceNormalizerAdvanced) combineNecessaryWords(tokens []string) string {
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

func (tsn *TurkishSentenceNormalizerAdvanced) combineCommon(w1, w2 string) string {
	combined := w1 + w2

	// Java logic:
	// - If second starts with apostrophe or "bil", and combined has analysis → combine.
	if strings.HasPrefix(w2, "'") || strings.HasPrefix(w2, "bil") {
		if tsn.Morphology.HasAnalysis(combined) {
			return combined
		}
	}
	// - Else if second does NOT have regular analysis, but combined has analysis → combine.
	if !tsn.Morphology.HasRegularAnalysis(w2) {
		if tsn.Morphology.HasAnalysis(combined) {
			return combined
		}
	}
	return ""
}

func (tsn *TurkishSentenceNormalizerAdvanced) splitNecessaryWords(tokens []string, useLookup bool) string {
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

func (tsn *TurkishSentenceNormalizerAdvanced) separateCommon(word string, useLookup bool) string {
	if tsn.NoSplitWords[word] {
		return word
	}

	if useLookup {
		if split, exists := tsn.CommonSplits[word]; exists {
			return split
		}
	}

	if !tsn.Morphology.HasRegularAnalysis(word) {
		for i := 1; i < len(word); i++ {
			tail := word[i:]
			if tsn.CommonConnectedSuffixes[tail] {
				head := word[:i]

				if len(tail) < 3 && tsn.LanguageModel != nil {
					vocab := tsn.LanguageModel.GetVocabulary()
					indexes := []int{vocab.IndexOf(head), vocab.IndexOf(tail)}
					if !tsn.LanguageModel.NgramExists(indexes) {
						return word
					}
				}

				if tsn.Morphology.HasRegularAnalysis(head) {
					return head + " " + tail
				} else {
					return word
				}
			}
		}
	}

	return word
}

func (tsn *TurkishSentenceNormalizerAdvanced) loadReplacements(path string) error {
	if path == "" {
		return nil
	}
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

func (tsn *TurkishSentenceNormalizerAdvanced) loadNoSplit(path string) error {
	if path == "" {
		return nil
	}
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

func (tsn *TurkishSentenceNormalizerAdvanced) loadCommonSplits(path string) error {
	if path == "" {
		return nil
	}
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
		if len(parts) != 2 {
			parts = strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
		}
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(parts[1])
		if left == "" || right == "" {
			continue
		}
		tsn.CommonSplits[left] = right
	}
	return scanner.Err()
}

func (tsn *TurkishSentenceNormalizerAdvanced) loadConnectedSuffixes(path string) error {
	if path == "" {
		return nil
	}
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

func (tsn *TurkishSentenceNormalizerAdvanced) loadMultimap(path string, target map[string][]string) error {
	if path == "" {
		return nil
	}
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

// --- Tokenization helpers (punctuation-aware) ---

func tokenizeAdvanced(s string) []string {
	runes := []rune(s)
	tokens := make([]string, 0)
	var buf []rune
	flush := func() {
		if len(buf) > 0 {
			tokens = append(tokens, string(buf))
			buf = buf[:0]
		}
	}
	containsAt := func(rs []rune) bool {
		for _, r := range rs {
			if r == '@' {
				return true
			}
		}
		return false
	}
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case isWordRune(r) || r == '\'' || r == '-':
			buf = append(buf, r)
		case r == '.':
			rest := runes[i+1:]
			if shouldTreatDotAsWord(buf, rest, containsAt) {
				buf = append(buf, r)
			} else {
				flush()
				tokens = append(tokens, ".")
			}
		case r == ' ' || r == '\t' || r == '\n':
			flush()
		default:
			flush()
			tokens = append(tokens, string(r))
		}
	}
	flush()
	return tokens
}

func shouldTreatDotAsWord(buf, rest []rune, containsAt func([]rune) bool) bool {
	if len(buf) == 0 {
		return false
	}
	if unicode.IsDigit(buf[len(buf)-1]) {
		if len(rest) > 0 && unicode.IsDigit(rest[0]) {
			return true
		}
	}
	if containsAt(buf) {
		return true
	}
	for _, r := range rest {
		if r == '@' {
			return true
		}
		if r == '.' {
			continue
		}
		if !isEmailContinuationRune(r) {
			break
		}
	}
	return false
}

func isEmailContinuationRune(r rune) bool {
	if isWordRune(r) {
		return true
	}
	switch r {
	case '-', '+', '%', '_':
		return true
	}
	return false
}

func isWordRune(r rune) bool {
	switch {
	case unicode.IsLetter(r):
		return true
	case unicode.IsDigit(r):
		return true
	case unicode.IsMark(r):
		return true
	}
	// Email/handle characters
	switch r {
	case '_', '@', '%', '+':
		return true
	}
	return false
}

func isWordAdvanced(token string) bool {
	if token == "" {
		return false
	}
	hasAt := strings.ContainsRune(token, '@')
	runes := []rune(token)
	for i, r := range runes {
		if isWordRune(r) || r == '\'' || r == '-' {
			continue
		}
		if r == '.' {
			if hasAt {
				continue
			}
			prevDigit := i > 0 && unicode.IsDigit(runes[i-1])
			nextDigit := i < len(runes)-1 && unicode.IsDigit(runes[i+1])
			if prevDigit && nextDigit {
				continue
			}
		}
		return false
	}
	return true
}

// joinTokensWithPunct joins tokens preserving punctuation adjacency (no extra space before ,.;:!?))
func joinTokensWithPunct(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}
	expanded := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if strings.ContainsAny(t, " \t\n\r") {
			parts := strings.Fields(t)
			if len(parts) == 0 {
				continue
			}
			expanded = append(expanded, parts...)
		} else {
			expanded = append(expanded, t)
		}
	}
	tokens = expanded
	if len(tokens) == 0 {
		return ""
	}
	var b strings.Builder
	prevWord := false
	for i, t := range tokens {
		isPunct := !isWordAdvanced(t)
		if i > 0 {
			if isPunct {
				// special-case emoticons like ":)": add a space before ':'
				if strings.HasPrefix(t, ":") {
					b.WriteRune(' ')
				}
				// otherwise no space before punctuation
			} else if !prevWord {
				// previous was punct -> no leading space if it was opening
				if t == ")" || t == "]" || t == "}" {
					// rare; add space anyway
					b.WriteRune(' ')
				} else {
					b.WriteRune(' ')
				}
			} else {
				b.WriteRune(' ')
			}
		}
		b.WriteString(t)
		prevWord = !isPunct
	}
	return b.String()
}

func hasNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
		switch r {
		case 'ç', 'Ç', 'ğ', 'Ğ', 'ı', 'İ', 'ö', 'Ö', 'ş', 'Ş', 'ü', 'Ü':
			return true
		}
	}
	return false
}

// firstExisting tries dataRoot/name1 or name2 and returns first existing path.
// Note: firstExisting helper is defined in turkish_sentence_normalizer.go for this package.
