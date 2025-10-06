package turkish

import "strings"

// StemAndEnding represents a word split as stem and ending.
// If the word is a stem, ending is empty string
type StemAndEnding struct {
	Stem   string
	Ending string
}

// NewStemAndEnding creates a new StemAndEnding instance
func NewStemAndEnding(stem, ending string) *StemAndEnding {
	if !hasText(ending) {
		ending = ""
	}
	return &StemAndEnding{
		Stem:   stem,
		Ending: ending,
	}
}

func hasText(s string) bool {
	return len(s) > 0 && len(strings.TrimSpace(s)) > 0
}
