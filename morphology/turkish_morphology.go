package morphology

import (
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/analysis"
	"github.com/kalaomer/zemberek-go/morphology/generator"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

// TurkishMorphology is the main morphological analyzer for Turkish
type TurkishMorphology struct {
	Lexicon            *lexicon.RootLexicon
	Morphotactics      *morphotactics.TurkishMorphotactics
	Analyzer           *analysis.RuleBasedAnalyzer
	WordGenerator      *generator.WordGenerator
	InformalAnalysis   bool
	IgnoreDiacritics   bool
}

// Builder for TurkishMorphology
type Builder struct {
	lexicon                     *lexicon.RootLexicon
	informalAnalysis            bool
	ignoreDiacriticsInAnalysis  bool
}

// NewBuilder creates a new builder with lexicon
func NewBuilder(lex *lexicon.RootLexicon) *Builder {
	return &Builder{
		lexicon:                    lex,
		informalAnalysis:           false,
		ignoreDiacriticsInAnalysis: false,
	}
}

// UseInformalAnalysis enables informal analysis
func (b *Builder) UseInformalAnalysis() *Builder {
	b.informalAnalysis = true
	return b
}

// IgnoreDiacriticsInAnalysis enables diacritics ignoring
func (b *Builder) IgnoreDiacriticsInAnalysis() *Builder {
	b.ignoreDiacriticsInAnalysis = true
	return b
}

// Build creates TurkishMorphology instance
func (b *Builder) Build() *TurkishMorphology {
	morph := morphotactics.NewTurkishMorphotactics(b.lexicon)

	var analyzer *analysis.RuleBasedAnalyzer
	if b.ignoreDiacriticsInAnalysis {
		analyzer = analysis.NewIgnoreDiacriticsAnalyzer(morph)
	} else {
		analyzer = analysis.NewRuleBasedAnalyzer(morph)
	}

	return &TurkishMorphology{
		Lexicon:          b.lexicon,
		Morphotactics:    morph,
		Analyzer:         analyzer,
		WordGenerator:    generator.NewWordGenerator(morph),
		InformalAnalysis: b.informalAnalysis,
		IgnoreDiacritics: b.ignoreDiacriticsInAnalysis,
	}
}

// CreateWithDefaults creates a morphology with default settings
// Loads from binary lexicon (lexicon.bin) like Java does
func CreateWithDefaults() *TurkishMorphology {
	items, err := lexicon.LoadBinaryLexicon()
	if err != nil {
		panic(err)
	}

	lex := lexicon.NewRootLexicon(items)
	return NewBuilder(lex).Build()
}

// Analyze analyzes a word
func (tm *TurkishMorphology) Analyze(word string) *analysis.WordAnalysis {
	if word == "" {
		return analysis.EmptyInputResult
	}

	normalized := tm.NormalizeForAnalysis(word)
	if normalized == "" {
		return analysis.EmptyInputResult
	}

	// Handle apostrophe
	if turkish.Instance.ContainsApostrophe(normalized) {
		normalized = turkish.Instance.NormalizeApostrophe(normalized)
		results := tm.analyzeWordsWithApostrophe(normalized)
		return analysis.NewWordAnalysis(word, results, normalized)
	}

	// Normal analysis
	results := tm.Analyzer.Analyze(normalized)

	// Filter unknown single results
	if len(results) == 1 && results[0].IsUnknown() {
		results = make([]*analysis.SingleAnalysis, 0)
	}

	return analysis.NewWordAnalysis(word, results, normalized)
}

// NormalizeForAnalysis normalizes word for analysis
func (tm *TurkishMorphology) NormalizeForAnalysis(word string) string {
	// Convert to lowercase using Turkish rules
	s := strings.ToLower(word)

	// Remove dots
	noDot := strings.ReplaceAll(s, ".", "")
	if noDot == "" {
		noDot = s
	}

	return noDot
}

// analyzeWordsWithApostrophe analyzes words containing apostrophe
func (tm *TurkishMorphology) analyzeWordsWithApostrophe(word string) []*analysis.SingleAnalysis {
	index := strings.IndexRune(word, '\'')
	if index <= 0 || index == len(word)-1 {
		return make([]*analysis.SingleAnalysis, 0)
	}

	// Remove apostrophe and analyze
	withoutQuote := strings.ReplaceAll(word, "'", "")
	noQuotesParses := tm.Analyzer.Analyze(withoutQuote)

	if len(noQuotesParses) == 0 {
		return make([]*analysis.SingleAnalysis, 0)
	}

	// Filter for noun analyses
	results := make([]*analysis.SingleAnalysis, 0)
	for _, parse := range noQuotesParses {
		// Would check: parse.item.primaryPos == PrimaryPos.Noun
		// For now, accept all
		results = append(results, parse)
	}

	return results
}

// AnalyzeSentence analyzes all words in a sentence
func (tm *TurkishMorphology) AnalyzeSentence(sentence string) []*analysis.WordAnalysis {
	// Simple tokenization by spaces
	tokens := strings.Fields(sentence)
	results := make([]*analysis.WordAnalysis, len(tokens))

	for i, token := range tokens {
		results[i] = tm.Analyze(token)
	}

	return results
}

// HasRegularAnalysis checks if word has regular analysis (not unknown, not runtime)
func (tm *TurkishMorphology) HasRegularAnalysis(word string) bool {
	wa := tm.Analyze(word)
	for _, sa := range wa.AnalysisResults {
		if !sa.IsUnknown() && !sa.IsRuntime() {
			// Would also check: sa.item.secondaryPos != SecondaryPos.ProperNoun
			return true
		}
	}
	return false
}

// HasAnalysis checks if word has any analysis (not runtime, not unknown)
func (tm *TurkishMorphology) HasAnalysis(word string) bool {
	wa := tm.Analyze(word)
	for _, sa := range wa.AnalysisResults {
		if !sa.IsRuntime() && !sa.IsUnknown() {
			return true
		}
	}
	return false
}
