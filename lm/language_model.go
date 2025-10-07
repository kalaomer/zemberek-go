package lm

import (
	"fmt"
	"math"
)

// Vocabulary is an alias for LmVocabulary
type Vocabulary = LmVocabulary

// NewVocabulary creates an empty vocabulary
func NewVocabulary() *Vocabulary {
	return &LmVocabulary{
		VocabularyIndexMap: make(map[string]int),
		Vocabulary:         make([]string, 0),
		UnknownWord:        DefaultUnknownWord,
		SentenceStart:      DefaultSentenceBeginMarker,
		SentenceEnd:        DefaultSentenceEndMarker,
	}
}

// LanguageModel is an interface for language models
type LanguageModel interface {
	// GetProbability returns log probability for given indexes
	GetProbability(indexes []int) float32
	// GetOrder returns n-gram order
	GetOrder() int
	// GetVocabulary returns vocabulary
	GetVocabulary() *Vocabulary
	// NgramExists checks if n-gram exists
	NgramExists(indexes []int) bool
}

// SimpleLM is a simplified language model implementation
type SimpleLM struct {
	Order      int
	Vocabulary *Vocabulary
	// Unigram probabilities (simplified)
	UnigramProbs map[string]float32
}

// NewSimpleLM creates a simple LM
func NewSimpleLM(order int) *SimpleLM {
	return &SimpleLM{
		Order:        order,
		Vocabulary:   NewVocabulary(),
		UnigramProbs: make(map[string]float32),
	}
}

// GetProbability returns log probability for given indexes
func (slm *SimpleLM) GetProbability(indexes []int) float32 {
	if len(indexes) == 0 {
		return float32(math.Log(1.0))
	}

	// For now, return uniform probability
	// Full implementation would look up n-grams
	return float32(math.Log(0.01))
}

// GetOrder returns n-gram order
func (slm *SimpleLM) GetOrder() int {
	return slm.Order
}

// GetVocabulary returns vocabulary
func (slm *SimpleLM) GetVocabulary() *Vocabulary {
	return slm.Vocabulary
}

// NgramExists checks if n-gram exists
func (slm *SimpleLM) NgramExists(indexes []int) bool {
	// Simplified: check if all indexes are valid
	for _, idx := range indexes {
		if idx < 0 || idx >= len(slm.Vocabulary.Vocabulary) {
			return false
		}
	}
	return true
}

// LoadFromFile loads LM from file (stub for now)
func LoadFromFile(path string) (LanguageModel, error) {
	lm, err := LoadSmoothLM(path)
	if err == nil {
		return lm, nil
	}
	smoothErr := fmt.Errorf("smoothlm: %w", err)

	if uni, uniErr := TryLoadUnigram(path); uniErr == nil {
		return uni, nil
	} else {
		return nil, fmt.Errorf("%v; unigram: %w", smoothErr, uniErr)
	}
}
