package normalization

import (
	"bufio"
	"os"
	"strings"

	"github.com/kalaomer/zemberek-go/lm"
	"github.com/kalaomer/zemberek-go/morphology"
	"github.com/kalaomer/zemberek-go/morphology/analysis"
	"github.com/kalaomer/zemberek-go/normalization/deasciifier"
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
func NewTurkishSentenceNormalizerAdvanced(morph *morphology.TurkishMorphology, resourcesPath string) (*TurkishSentenceNormalizerAdvanced, error) {
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

	// Load language model (simplified for now)
	lm, _ := lm.LoadFromFile(resourcesPath + "/lm.2gram.slm")
	tsn.LanguageModel = lm

	// Load resources
	if resourcesPath == "" {
		resourcesPath = "resources/normalization"
	}

	// Load all resource files
	tsn.loadReplacements(resourcesPath + "/multi-word-replacements.txt")
	tsn.loadNoSplit(resourcesPath + "/no-split.txt")
	tsn.loadCommonSplits(resourcesPath + "/split.txt")
	tsn.loadConnectedSuffixes(resourcesPath + "/question-suffixes.txt")
	tsn.loadMultimap(resourcesPath+"/candidates-manual.txt", tsn.LookupManual)
	tsn.loadMultimap(resourcesPath+"/lookup-from-graph.txt", tsn.LookupFromGraph)
	tsn.loadMultimap(resourcesPath+"/ascii-map.txt", tsn.LookupFromASCII)

	// Create spell checker with morphology
	graph := NewCharacterGraph()
	decoder := NewCharacterGraphDecoder(graph)
	tsn.SpellChecker = &TurkishSpellChecker{
		Morphology:  morph,
		Decoder:     decoder,
		CharMatcher: DiacriticsIgnoringMatcherInstance,
	}

	return tsn, nil
}

// Normalize normalizes a sentence using full morphological analysis and beam search
func (tsn *TurkishSentenceNormalizerAdvanced) Normalize(sentence string) string {
	processed := tsn.preProcess(sentence)
	tokens := tokenize(processed)

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
	return strings.Join(result, " ")
}

// getCandidatesAdvanced gets normalization candidates using morphological analysis
func (tsn *TurkishSentenceNormalizerAdvanced) getCandidatesAdvanced(word, previous, next string) []string {
	candidates := make([]string, 0)
	seen := make(map[string]bool)

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

	// Analyze with informal morphology
	analyses := tsn.InformalMorphology.Analyze(word)
	for _, sa := range analyses.AnalysisResults {
		if sa.ContainsInformalMorpheme() {
			// Convert informal to formal
			result := tsn.AnalysisConverter.Convert(word, sa)
			if result != nil && !seen[result.Surface] {
				candidates = append(candidates, result.Surface)
				seen[result.Surface] = true
			}
		} else {
			// Generate from morphemes
			results := tsn.Morphology.WordGenerator.Generate(sa.Item, sa.GetMorphemes())
			for _, r := range results {
				if !seen[r.Surface] {
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

	// Add original if it has correct analysis or if no candidates
	if len(candidates) == 0 || tsn.Morphology.Analyze(word).IsCorrect() {
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

				// Calculate score using LM
				indexes := make([]int, lmOrder)
				vocab := tsn.LanguageModel.GetVocabulary()
				for j := 0; j < lmOrder-1; j++ {
					indexes[j] = vocab.IndexOf(hist[j].Content)
				}
				indexes[lmOrder-1] = vocab.IndexOf(c.Content)

				score := tsn.LanguageModel.GetProbability(indexes)
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

	// Extract sequence
	seq := make([]string, 0)
	h := best.Previous
	for h != nil && h.Current != GetStartCandidate() {
		seq = append([]string{h.Current.Content}, seq...)
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
func (tsn *TurkishSentenceNormalizerAdvanced) preProcess(sentence string) string {
	sentence = strings.ToLower(sentence)

	tokens := tokenize(sentence)
	sentence = tsn.replaceCommon(tokens)

	tokens = tokenize(sentence)
	sentence = tsn.combineNecessaryWords(tokens)

	tokens = tokenize(sentence)
	sentence = tsn.splitNecessaryWords(tokens, false)

	if tsn.AlwaysApplyDeasciifier || probablyRequiresDeasciifier(sentence) {
		d := deasciifier.NewDeasciifier(sentence)
		sentence = d.ConvertToTurkish()
	}

	tokens = tokenize(sentence)
	sentence = tsn.combineNecessaryWords(tokens)

	tokens = tokenize(sentence)
	sentence = tsn.splitNecessaryWords(tokens, true)

	return sentence
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

	// Check with morphology
	if tsn.Morphology.HasAnalysis(combined) {
		return combined
	}

	if strings.HasPrefix(w2, "'") || strings.HasPrefix(w2, "bil") {
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

func (tsn *TurkishSentenceNormalizerAdvanced) loadConnectedSuffixes(path string) error {
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
