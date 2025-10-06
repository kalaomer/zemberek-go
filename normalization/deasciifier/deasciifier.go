package deasciifier

import (
	"unicode"
)

const turkishContextSize = 10

// Deasciifier converts ASCII Turkish text to properly accented Turkish
type Deasciifier struct {
	asciiString   string
	turkishString string
}

// Lookup tables
var (
	turkishAsciifyTable = map[rune]rune{
		'ç': 'c', 'Ç': 'C', 'ğ': 'g', 'Ğ': 'G', 'ö': 'o',
		'Ö': 'O', 'ı': 'i', 'İ': 'I', 'ş': 's', 'Ş': 'S',
		'ü': 'u', 'Ü': 'U',
	}

	turkishToggleAccentTable = map[rune]rune{
		'c': 'ç', 'C': 'Ç', 'g': 'ğ', 'G': 'Ğ', 'o': 'ö', 'O': 'Ö',
		'u': 'ü', 'U': 'Ü', 'i': 'ı', 'I': 'İ', 's': 'ş', 'S': 'Ş',
		'ç': 'c', 'Ç': 'C', 'ğ': 'g', 'Ğ': 'G', 'ö': 'o', 'Ö': 'O',
		'ü': 'u', 'Ü': 'U', 'ı': 'i', 'İ': 'I', 'ş': 's', 'Ş': 'S',
	}

	turkishDowncaseAsciifyTable map[rune]rune
	turkishUpcaseAccentsTable   map[rune]rune
	turkishPatternTable         map[string]map[string]int
)

func init() {
	// Initialize downcase table
	turkishDowncaseAsciifyTable = make(map[rune]rune)
	for k, v := range map[rune]rune{
		'ç': 'c', 'Ç': 'c', 'ğ': 'g', 'Ğ': 'g', 'ö': 'o', 'Ö': 'o',
		'ı': 'i', 'İ': 'i', 'ş': 's', 'Ş': 's', 'ü': 'u', 'Ü': 'u',
	} {
		turkishDowncaseAsciifyTable[k] = v
	}

	for r := 'A'; r <= 'Z'; r++ {
		turkishDowncaseAsciifyTable[r] = unicode.ToLower(r)
		turkishDowncaseAsciifyTable[unicode.ToLower(r)] = unicode.ToLower(r)
	}

	// Initialize upcase accents table
	turkishUpcaseAccentsTable = make(map[rune]rune)
	for r := 'A'; r <= 'Z'; r++ {
		turkishUpcaseAccentsTable[r] = unicode.ToLower(r)
		turkishUpcaseAccentsTable[unicode.ToLower(r)] = unicode.ToLower(r)
	}

	turkishUpcaseAccentsTable['ç'] = 'C'
	turkishUpcaseAccentsTable['Ç'] = 'C'
	turkishUpcaseAccentsTable['ğ'] = 'G'
	turkishUpcaseAccentsTable['Ğ'] = 'G'
	turkishUpcaseAccentsTable['ö'] = 'O'
	turkishUpcaseAccentsTable['Ö'] = 'O'
	turkishUpcaseAccentsTable['ı'] = 'I'
	turkishUpcaseAccentsTable['İ'] = 'i'
	turkishUpcaseAccentsTable['ş'] = 'S'
	turkishUpcaseAccentsTable['Ş'] = 'S'
	turkishUpcaseAccentsTable['ü'] = 'U'
	turkishUpcaseAccentsTable['Ü'] = 'U'

	// Pattern table would be loaded from resources/normalization/turkish_pattern_table.pickle
	// For now, initialize as empty - this needs to be loaded from the serialized pattern file
	turkishPatternTable = make(map[string]map[string]int)
}

// NewDeasciifier creates a new Deasciifier
func NewDeasciifier(asciiString string) *Deasciifier {
	return &Deasciifier{
		asciiString:   asciiString,
		turkishString: asciiString,
	}
}

// ConvertToTurkish converts ASCII Turkish to properly accented Turkish
func (d *Deasciifier) ConvertToTurkish() string {
	runes := []rune(d.turkishString)

	for i, c := range runes {
		if d.turkishNeedCorrection(c, i) {
			if toggle, ok := turkishToggleAccentTable[c]; ok {
				runes[i] = toggle
			}
		}
	}

	d.turkishString = string(runes)
	return d.turkishString
}

func (d *Deasciifier) turkishNeedCorrection(c rune, point int) bool {
	tr := c
	if asciified, ok := turkishAsciifyTable[c]; ok {
		tr = asciified
	}

	trLower := unicode.ToLower(tr)
	patterns, ok := turkishPatternTable[string(trLower)]

	m := false
	if ok {
		m = d.turkishMatchPattern(patterns, point)
	}

	if tr == 'I' {
		if c == tr {
			return !m
		}
		return m
	}

	if c == tr {
		return m
	}
	return !m
}

func (d *Deasciifier) turkishMatchPattern(dlist map[string]int, point int) bool {
	rank := len(dlist) * 2
	str := d.turkishGetContext(turkishContextSize, point)

	start := 0
	length := len([]rune(str))

	for start <= turkishContextSize {
		end := turkishContextSize + 1
		for end <= length {
			runes := []rune(str)
			s := string(runes[start:end])
			if r, ok := dlist[s]; ok {
				if abs(r) < abs(rank) {
					rank = r
				}
			}
			end++
		}
		start++
	}

	return rank > 0
}

func (d *Deasciifier) turkishGetContext(size int, point int) string {
	runes := make([]rune, 1+(2*size))
	for i := range runes {
		runes[i] = ' '
	}
	runes[size] = 'X'

	i := size + 1
	space := false
	index := point + 1
	turkishRunes := []rune(d.turkishString)

	for i < len(runes) && !space && index < len(turkishRunes) {
		currentChar := turkishRunes[index]
		if x, ok := turkishDowncaseAsciifyTable[currentChar]; ok {
			runes[i] = x
			i++
			space = false
		} else if !space {
			i++
			space = true
		}
		index++
	}

	result := string(runes[:i])
	index = point - 1
	i = size - 1
	space = false

	resultRunes := []rune(result)
	for i >= 0 && index >= 0 {
		currentChar := turkishRunes[index]
		if x, ok := turkishUpcaseAccentsTable[currentChar]; ok {
			resultRunes[i] = x
			i--
			space = false
		} else if !space {
			i--
			space = true
		}
		index--
	}

	return string(resultRunes)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// LoadPatternTable loads the pattern table from file
// This should be called during initialization with the pickle file
func LoadPatternTable(filename string) error {
	// TODO: Implement loading from pickle/serialized format
	// For now, this is a placeholder
	return nil
}
