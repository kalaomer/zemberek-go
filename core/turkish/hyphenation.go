package turkish

// TurkishSyllableExtractor handles syllable extraction for Turkish
type TurkishSyllableExtractor struct {
	Alphabet *TurkishAlphabet
	Strict   bool
}

// Strict syllable extractor instance
var Strict *TurkishSyllableExtractor

func init() {
	Strict = NewTurkishSyllableExtractor(true)
}

// NewTurkishSyllableExtractor creates a new syllable extractor
func NewTurkishSyllableExtractor(strict bool) *TurkishSyllableExtractor {
	return &TurkishSyllableExtractor{
		Alphabet: Instance,
		Strict:   strict,
	}
}
