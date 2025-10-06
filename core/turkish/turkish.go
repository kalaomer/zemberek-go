package turkish

import (
	"unicode"
)

// Capitalize capitalizes a Turkish word correctly
func Capitalize(word string) string {
	if len(word) == 0 {
		return word
	}

	runes := []rune(word)
	if runes[0] == 'i' {
		runes[0] = 'İ'
	} else if runes[0] == 'I' {
		runes[0] = 'I'
	} else {
		runes[0] = unicode.ToUpper(runes[0])
	}

	// Lowercase the rest
	for i := 1; i < len(runes); i++ {
		if runes[i] == 'I' {
			runes[i] = 'ı'
		} else if runes[i] == 'İ' {
			runes[i] = 'i'
		} else {
			runes[i] = unicode.ToLower(runes[i])
		}
	}

	return string(runes)
}
