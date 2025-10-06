package tokenization

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var (
	webWords         = []string{"http:", ".html", "www", ".tr", ".edu", ".com", ".net", ".gov", ".org", "@"}
	lowercaseVowels  = map[rune]bool{'a': true, 'e': true, 'ı': true, 'i': true, 'o': true, 'ö': true, 'u': true, 'ü': true, 'â': true, 'î': true, 'û': true}
	uppercaseVowels  = map[rune]bool{'A': true, 'E': true, 'I': true, 'İ': true, 'O': true, 'Ö': true, 'U': true, 'Ü': true, 'Â': true, 'Î': true, 'Û': true}
	whitespaceRegex  = regexp.MustCompile(`\s+`)
	trailingDotRegex = regexp.MustCompile(`\.$`)
)

// PerceptronSegmenter loads Binary Averaged Perceptron model for sentence boundary detection
type PerceptronSegmenter struct {
	TurkishAbbreviationSet map[string]bool
}

// NewPerceptronSegmenter creates a new PerceptronSegmenter
func NewPerceptronSegmenter() *PerceptronSegmenter {
	return &PerceptronSegmenter{
		TurkishAbbreviationSet: LoadAbbreviations(""),
	}
}

// LoadWeightsFromCSV loads model weights from CSV file
func LoadWeightsFromCSV(path string) (map[string]float64, error) {
	if path == "" {
		// Default path would be set here
		path = "resources/sentence_boundary_model_weights.csv"
	}

	weights := make(map[string]float64)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) >= 2 {
			var value float64
			_, err := fmt.Sscanf(record[1], "%f", &value)
			if err == nil {
				weights[record[0]] = value
			}
		}
	}

	return weights, nil
}

// LoadAbbreviations loads Turkish abbreviations from a text file
func LoadAbbreviations(path string) map[string]bool {
	if path == "" {
		path = "resources/abbreviations.txt"
	}

	abbrSet := make(map[string]bool)
	file, err := os.Open(path)
	if err != nil {
		return abbrSet
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			abbr := whitespaceRegex.ReplaceAllString(line, "")
			abbr = trailingDotRegex.ReplaceAllString(abbr, "")
			abbrSet[abbr] = true
			abbrSet[turkishLower(abbr)] = true
		}
	}

	return abbrSet
}

func turkishLower(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r == 'I' {
			result.WriteRune('ı')
		} else if r == 'İ' {
			result.WriteRune('i')
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

// PotentialWebsite checks if a string potentially represents a website
func PotentialWebsite(s string) bool {
	for _, word := range webWords {
		if strings.Contains(s, word) {
			return true
		}
	}
	return false
}

// GetMetaChar returns the meta character for a letter
func GetMetaChar(letter string) string {
	if len(letter) == 0 {
		return ""
	}

	r := []rune(letter)[0]

	if unicode.IsUpper(r) {
		if uppercaseVowels[r] {
			return "V"
		}
		return "C"
	} else if unicode.IsLower(r) {
		if lowercaseVowels[r] {
			return "v"
		}
		return "c"
	} else if unicode.IsDigit(r) {
		return "d"
	} else if unicode.IsSpace(r) {
		return " "
	} else if r == '.' || r == '!' || r == '?' {
		return string(r)
	}

	return "-"
}
