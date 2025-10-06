package tokenization

import (
	"regexp"
	"strings"
	"unicode"

	
)

var (
	boundaryChars  = map[rune]bool{'.': true, '!': true, '?': true, '…': true}
	doubleQuotes   = map[rune]bool{34: true, 8220: true, 8221: true, '»': true, '«': true}
	punctuationReg = regexp.MustCompile(`[.!?…]`)
)

// TurkishSentenceExtractor separates sentences using perceptron model and rule-based approaches
type TurkishSentenceExtractor struct {
	*PerceptronSegmenter
	Weights                      map[string]float64
	DoNotSplitInDoubleQuotes     bool
	AbbrSet                      map[string]bool
}

// NewTurkishSentenceExtractor creates a new sentence extractor
func NewTurkishSentenceExtractor(doNotSplitInDoubleQuotes bool, weightsPath string) (*TurkishSentenceExtractor, error) {
	weights, err := LoadWeightsFromCSV(weightsPath)
	if err != nil {
		// Create empty weights if file not found
		weights = make(map[string]float64)
	}

	return &TurkishSentenceExtractor{
		PerceptronSegmenter:      NewPerceptronSegmenter(),
		Weights:                  weights,
		DoNotSplitInDoubleQuotes: doNotSplitInDoubleQuotes,
		AbbrSet:                  LoadAbbreviations(""),
	}, nil
}

// ExtractToSpans divides paragraph into spans
func (t *TurkishSentenceExtractor) ExtractToSpans(paragraph string) []*Span {
	spans := make([]*Span, 0)
	begin := 0

	runes := []rune(paragraph)

	for j, ch := range runes {
		if boundaryChars[ch] {
			// Include the boundary character in the sentence
			end := j + 1

			// Check if it's end of sentence
			if j < len(runes)-1 {
				nextChar := runes[j+1]
				if unicode.IsSpace(nextChar) || unicode.IsUpper(nextChar) {
					span, _ := NewSpan(begin, end)
					if span.GetLength() > 0 {
						spans = append(spans, span)
					}
					// Skip spaces after boundary
					for end < len(runes) && unicode.IsSpace(runes[end]) {
						end++
					}
					begin = end
				}
			} else {
				// Last character is a boundary
				span, _ := NewSpan(begin, end)
				if span.GetLength() > 0 {
					spans = append(spans, span)
				}
				begin = end
			}
		}
	}

	if begin < len(runes) {
		span, _ := NewSpan(begin, len(runes))
		if span.GetLength() > 0 {
			spans = append(spans, span)
		}
	}

	return spans
}

// FromParagraph extracts sentences from a paragraph
func (t *TurkishSentenceExtractor) FromParagraph(paragraph string) []string {
	spans := t.ExtractToSpans(paragraph)
	sentences := make([]string, 0)
	runes := []rune(paragraph)

	for _, span := range spans {
		// Use rune-based substring since spans are in rune indices
		if span.Start >= 0 && span.End <= len(runes) {
			sentence := strings.TrimSpace(string(runes[span.Start:span.End]))
			if len(sentence) > 0 {
				sentences = append(sentences, sentence)
			}
		}
	}

	return sentences
}

// GetWeight returns the weight for a feature
func (t *TurkishSentenceExtractor) GetWeight(key string) float64 {
	if weight, ok := t.Weights[key]; ok {
		return weight
	}
	return 0.0
}
