package tokenization

import (
	_ "embed"
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

//go:embed data/abbreviations.txt
var abbreviationsData string

// Global abbreviations set loaded from abbreviations.txt
var abbreviations map[string]bool

func init() {
	abbreviations = loadAbbreviations()
}

// loadAbbreviations loads abbreviations from embedded file
// Matches Java implementation:
// - Only loads abbreviations that end with dot (Prof., Dr., etc.)
// - Adds lowercase versions (English and Turkish)
func loadAbbreviations() map[string]bool {
	abbr := make(map[string]bool)

	lines := strings.Split(abbreviationsData, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove any spaces from abbreviation
		line = strings.ReplaceAll(line, " ", "")

		// Only process abbreviations that end with dot (matches Java behavior)
		if strings.HasSuffix(line, ".") {
			// Add as-is
			abbr[line] = true

			// Add Turkish lowercase
			abbr[turkish.Instance.ToLower(line)] = true
		}
	}

	return abbr
}

// IsAbbreviation checks if a word is a known abbreviation
// Case-insensitive check (Prof., prof., PROF. all match)
// Example: IsAbbreviation("Prof.") -> true
func IsAbbreviation(word string) bool {
	// Check as-is first (fast path)
	if abbreviations[word] {
		return true
	}
	// Check Turkish lowercase version
	return abbreviations[turkish.Instance.ToLower(word)]
}
