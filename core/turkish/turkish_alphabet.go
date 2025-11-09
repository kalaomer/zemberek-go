package turkish

import (
	"strings"
	"unicode"

	"github.com/kalaomer/zemberek-go/core/text"
)

// TurkishAlphabet represents Turkish alphabet with all letters and special characters
type TurkishAlphabet struct {
	Lowercase             string
	Uppercase             string
	AllLetters            string
	Digits                string // "0123456789"
	AllLettersAndDigits   string // Digits + AllLetters for regex patterns
	AllLettersDigitsUnderscore string // Digits + AllLetters + "_" for regex patterns
	VowelsLowercase       string
	VowelsUppercase       string
	Vowels                map[rune]bool
	Circumflex            string
	CircumflexUpper       string
	Circumflexes          map[rune]bool
	Apostrophe            map[rune]bool
	StopConsonants        string
	VoicelessConsonants   string
	TurkishSpecific       string
	TurkishSpecificLookup map[rune]bool
	TurkishASCII          string
	ASCIIEqTr             string
	ASCIIEqTrSet          map[rune]bool
	ASCIIEq               string
	ForeignDiacritics     string
	DiacriticsToTurkish   string

	VoicingMap           map[rune]rune
	DevoicingMap         map[rune]rune
	CircumflexMap        map[rune]rune
	LetterMap            map[rune]*TurkicLetter
	ASCIIEqualMap        map[rune]rune
	TurkishToASCIIMap    map[rune]rune
	ForeignDiacriticsMap map[rune]rune
}

var Instance *TurkishAlphabet

func init() {
	Instance = NewTurkishAlphabet()
}

// NewTurkishAlphabet creates a new TurkishAlphabet instance
func NewTurkishAlphabet() *TurkishAlphabet {
	ta := &TurkishAlphabet{
		Lowercase:            "abcçdefgğhıijklmnoöprsştuüvyzxwqâîû",
		VowelsLowercase:      "aeıioöuüâîû",
		Circumflex:           "âîû",
		CircumflexUpper:      "ÂÎÛ",
		StopConsonants:       "çkptÇKPT",
		VoicelessConsonants:  "çfhkpsştÇFHKPSŞT",
		TurkishSpecific:      "çÇğĞıİöÖşŞüÜâîûÂÎÛ",
		TurkishASCII:         "cCgGiIoOsSuUaiuAIU",
		ASCIIEqTr:            "cCgGiIoOsSuUçÇğĞıİöÖşŞüÜ",
		ASCIIEq:              "çÇğĞıİöÖşŞüÜcCgGiIoOsSuU",
		ForeignDiacritics:    "ÀÁÂÃÄÅÈÉÊËÌÍÎÏÑÒÓÔÕÙÚÛàáâãäåèéêëìíîïñòóôõùúû",
		DiacriticsToTurkish:  "AAAAAAEEEEIIIINOOOOUUUaaaaaaeeeeiiiinoooouuu",
		VoicingMap:           make(map[rune]rune),
		DevoicingMap:         make(map[rune]rune),
		CircumflexMap:        make(map[rune]rune),
		LetterMap:            make(map[rune]*TurkicLetter),
		ASCIIEqualMap:        make(map[rune]rune),
		TurkishToASCIIMap:    make(map[rune]rune),
		ForeignDiacriticsMap: make(map[rune]rune),
	}

	// Generate uppercase and combined character sets
	ta.Uppercase = turkishUpper(ta.Lowercase)
	ta.AllLetters = ta.Lowercase + ta.Uppercase
	ta.Digits = "0123456789"
	ta.AllLettersAndDigits = ta.Digits + ta.AllLetters
	ta.AllLettersDigitsUnderscore = ta.Digits + ta.AllLetters + "_"
	ta.VowelsUppercase = turkishUpper(ta.VowelsLowercase)

	// Create vowels map
	ta.Vowels = make(map[rune]bool)
	for _, c := range ta.VowelsLowercase + ta.VowelsUppercase {
		ta.Vowels[c] = true
	}

	// Create circumflex map
	ta.Circumflexes = make(map[rune]bool)
	for _, c := range ta.Circumflex + ta.CircumflexUpper {
		ta.Circumflexes[c] = true
	}

	// Create apostrophe map
	ta.Apostrophe = make(map[rune]bool)
	for _, c := range "′´`'''" {
		ta.Apostrophe[c] = true
	}

	// Turkish specific lookup
	ta.TurkishSpecificLookup = make(map[rune]bool)
	for _, c := range ta.TurkishSpecific {
		ta.TurkishSpecificLookup[c] = true
	}

	// ASCII equal tr set
	ta.ASCIIEqTrSet = make(map[rune]bool)
	for _, c := range ta.ASCIIEqTr {
		ta.ASCIIEqTrSet[c] = true
	}

	// Generate letters
	letters := generateLetters()
	for _, letter := range letters {
		ta.LetterMap[letter.CharValue] = letter
	}

	// ASCII equal map
	for i, c1 := range []rune(ta.ASCIIEqTr) {
		c2 := []rune(ta.ASCIIEq)[i]
		ta.ASCIIEqualMap[c1] = c2
	}

	// Generate voicing/devoicing lookups
	ta.generateVoicingDevoicingLookups()

	// Turkish to ASCII map
	populateDict(ta.TurkishToASCIIMap, ta.TurkishSpecific, ta.TurkishASCII)

	// Foreign diacritics map
	populateDict(ta.ForeignDiacriticsMap, ta.ForeignDiacritics, ta.DiacriticsToTurkish)

	return ta
}

func turkishUpper(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r == 'i' {
			result.WriteRune('İ')
		} else {
			result.WriteRune(unicode.ToUpper(r))
		}
	}
	return result.String()
}

// ToLower converts a string to lowercase using Turkish-specific rules
func (ta *TurkishAlphabet) ToLower(s string) string {
	return turkishLower(s)
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

// IsTurkishSpecific returns true if the character is Turkish specific
func (ta *TurkishAlphabet) IsTurkishSpecific(c rune) bool {
	return ta.TurkishSpecificLookup[c]
}

// IsTurkishLetter returns true if the rune is a Turkish letter (including ASCII and Turkish-specific)
// Covers: a-z, A-Z, ç, ğ, ı, İ, ö, ş, ü, Ç, Ğ, Ö, Ş, Ü, â, î, û, Â, Î, Û
func (ta *TurkishAlphabet) IsTurkishLetter(r rune) bool {
	// Check ASCII letters first (fast path)
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	// Use existing TurkishSpecificLookup map for Turkish-specific letters
	return ta.TurkishSpecificLookup[r]
}

// IsDigit returns true if the rune is a digit (0-9)
func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// ContainsASCIIRelated returns true if the string contains ASCII related characters
func (ta *TurkishAlphabet) ContainsASCIIRelated(s string) bool {
	for _, c := range s {
		if ta.ASCIIEqTrSet[c] {
			return true
		}
	}
	return false
}

// ToASCII converts Turkish characters to ASCII equivalents
func (ta *TurkishAlphabet) ToASCII(inp string) string {
	var sb strings.Builder
	for _, c := range inp {
		if res, ok := ta.TurkishToASCIIMap[c]; ok {
			sb.WriteRune(res)
		} else {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

// IsASCIIEqual returns true if two characters are ASCII equal
func (ta *TurkishAlphabet) IsASCIIEqual(c1, c2 rune) bool {
	if c1 == c2 {
		return true
	}
	if a1, ok := ta.ASCIIEqualMap[c1]; ok {
		return a1 == c2
	}
	return false
}

// EqualsIgnoreDiacritics returns true if two strings are equal ignoring diacritics
func (ta *TurkishAlphabet) EqualsIgnoreDiacritics(s1, s2 string) bool {
	r1 := []rune(s1)
	r2 := []rune(s2)
	if len(r1) != len(r2) {
		return false
	}
	for i := range r1 {
		if !ta.IsASCIIEqual(r1[i], r2[i]) {
			return false
		}
	}
	return true
}

// StartsWithIgnoreDiacritics returns true if s1 starts with s2 ignoring diacritics
func (ta *TurkishAlphabet) StartsWithIgnoreDiacritics(s1, s2 string) bool {
	r1 := []rune(s1)
	r2 := []rune(s2)
	if len(r1) < len(r2) {
		return false
	}
	for i := range r2 {
		if !ta.IsASCIIEqual(r1[i], r2[i]) {
			return false
		}
	}
	return true
}

// ContainsDigit returns true if the string contains a digit
func ContainsDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

// ContainsApostrophe returns true if the string contains an apostrophe
func (ta *TurkishAlphabet) ContainsApostrophe(s string) bool {
	for _, c := range s {
		if ta.Apostrophe[c] {
			return true
		}
	}
	return false
}

// NormalizeApostrophe normalizes apostrophes in a string
func (ta *TurkishAlphabet) NormalizeApostrophe(s string) string {
	if !ta.ContainsApostrophe(s) {
		return s
	}
	var sb strings.Builder
	for _, c := range s {
		if ta.Apostrophe[c] {
			sb.WriteRune('\'')
		} else {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

// ContainsForeignDiacritics returns true if the string contains foreign diacritics
func (ta *TurkishAlphabet) ContainsForeignDiacritics(s string) bool {
	for _, c := range s {
		for _, d := range ta.ForeignDiacritics {
			if c == d {
				return true
			}
		}
	}
	return false
}

// ForeignDiacriticsToTurkish converts foreign diacritics to Turkish equivalents
func (ta *TurkishAlphabet) ForeignDiacriticsToTurkish(inp string) string {
	var sb strings.Builder
	for _, c := range inp {
		if res, ok := ta.ForeignDiacriticsMap[c]; ok {
			sb.WriteRune(res)
		} else {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

// ContainsCircumflex returns true if the string contains a circumflex
func (ta *TurkishAlphabet) ContainsCircumflex(s string) bool {
	for _, c := range s {
		if ta.Circumflexes[c] {
			return true
		}
	}
	return false
}

// NormalizeCircumflex normalizes circumflex characters
func (ta *TurkishAlphabet) NormalizeCircumflex(s string) string {
	runes := []rune(s)
	if len(runes) == 1 {
		if res, ok := ta.CircumflexMap[runes[0]]; ok {
			return string(res)
		}
		return s
	}
	if !ta.ContainsCircumflex(s) {
		return s
	}
	var sb strings.Builder
	for _, c := range s {
		if ta.Circumflexes[c] {
			sb.WriteRune(ta.CircumflexMap[c])
		} else {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

// Normalize normalizes the input string
func (ta *TurkishAlphabet) Normalize(inp string) string {
	inp = text.NormalizeApostrophes(turkishLower(inp))
	var sb strings.Builder
	for _, c := range inp {
		// Keep Turkish letters, spaces, and common punctuation
		if _, ok := ta.LetterMap[c]; ok {
			sb.WriteRune(c)
		} else if c == ' ' || c == '.' || c == '!' || c == '?' || c == ',' || c == '-' || c == ':' || c == ';' {
			sb.WriteRune(c)
		}
		// Skip other characters (don't replace with '?')
	}
	return sb.String()
}

// IsVowel returns true if the character is a vowel
func (ta *TurkishAlphabet) IsVowel(c rune) bool {
	return ta.Vowels[c]
}

// ContainsVowel returns true if the string contains a vowel
func (ta *TurkishAlphabet) ContainsVowel(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if ta.IsVowel(c) {
			return true
		}
	}
	return false
}

func (ta *TurkishAlphabet) generateVoicingDevoicingLookups() {
	voicingIn := "çgkpt"
	voicingOut := "cğğbd"
	devoicingIn := "bcdgğ"
	devoicingOut := "pçtkk"

	populateDict(ta.VoicingMap, voicingIn+turkishUpper(voicingIn), voicingOut+turkishUpper(voicingOut))
	populateDict(ta.DevoicingMap, devoicingIn+turkishUpper(devoicingIn), devoicingOut+turkishUpper(devoicingOut))

	circumflexNormalized := "aiu"
	populateDict(ta.CircumflexMap, ta.Circumflex+turkishUpper(ta.Circumflex),
		circumflexNormalized+turkishUpper(circumflexNormalized))
}

func populateDict(dict map[rune]rune, inStr, outStr string) {
	inRunes := []rune(inStr)
	outRunes := []rune(outStr)
	for i := range inRunes {
		dict[inRunes[i]] = outRunes[i]
	}
}

func generateLetters() []*TurkicLetter {
	letters := []*TurkicLetter{
		NewTurkicLetter('a', true, false, false, false, false),
		NewTurkicLetter('e', true, true, false, false, false),
		NewTurkicLetter('ı', true, false, false, false, false),
		NewTurkicLetter('i', true, true, false, false, false),
		NewTurkicLetter('o', true, false, true, false, false),
		NewTurkicLetter('ö', true, true, true, false, false),
		NewTurkicLetter('u', true, false, true, false, false),
		NewTurkicLetter('ü', true, true, true, false, false),
		NewTurkicLetter('â', true, false, false, false, false),
		NewTurkicLetter('î', true, true, false, false, false),
		NewTurkicLetter('û', true, true, true, false, false),
		NewTurkicLetter('b', false, false, false, false, false),
		NewTurkicLetter('c', false, false, false, false, false),
		NewTurkicLetter('ç', false, false, false, true, false),
		NewTurkicLetter('d', false, false, false, false, false),
		NewTurkicLetter('f', false, false, false, true, true),
		NewTurkicLetter('g', false, false, false, false, false),
		NewTurkicLetter('ğ', false, false, false, false, true),
		NewTurkicLetter('h', false, false, false, true, true),
		NewTurkicLetter('j', false, false, false, false, true),
		NewTurkicLetter('k', false, false, false, true, false),
		NewTurkicLetter('l', false, false, false, false, true),
		NewTurkicLetter('m', false, false, false, false, true),
		NewTurkicLetter('n', false, false, false, false, true),
		NewTurkicLetter('p', false, false, false, true, false),
		NewTurkicLetter('r', false, false, false, false, true),
		NewTurkicLetter('s', false, false, false, true, true),
		NewTurkicLetter('ş', false, false, false, true, true),
		NewTurkicLetter('t', false, false, false, true, false),
		NewTurkicLetter('v', false, false, false, false, true),
		NewTurkicLetter('y', false, false, false, false, true),
		NewTurkicLetter('z', false, false, false, false, true),
		NewTurkicLetter('q', false, false, false, false, false),
		NewTurkicLetter('w', false, false, false, false, false),
		NewTurkicLetter('x', false, false, false, false, false),
	}

	// Add capitals
	capitals := make([]*TurkicLetter, 0, len(letters))
	for _, letter := range letters {
		var upper rune
		if letter.CharValue == 'i' {
			upper = 'İ'
		} else {
			upper = []rune(turkishUpper(string(letter.CharValue)))[0]
		}
		capitals = append(capitals, letter.CopyFor(upper))
	}

	letters = append(letters, capitals...)
	return letters
}

// GetLastLetter returns the last letter of the input as TurkicLetter
func (ta *TurkishAlphabet) GetLastLetter(s string) *TurkicLetter {
	if len(s) == 0 {
		return Undefined
	}
	runes := []rune(s)
	return ta.GetLetter(runes[len(runes)-1])
}

// GetLetter returns the letter for the given character
func (ta *TurkishAlphabet) GetLetter(c rune) *TurkicLetter {
	if letter, ok := ta.LetterMap[c]; ok {
		return letter
	}
	return Undefined
}

// GetLastVowel returns the last vowel in the string
func (ta *TurkishAlphabet) GetLastVowel(s string) *TurkicLetter {
	if len(s) == 0 {
		return Undefined
	}
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if ta.IsVowel(runes[i]) {
			return ta.GetLetter(runes[i])
		}
	}
	return Undefined
}

// GetFirstLetter returns the first letter of the input as TurkicLetter
func (ta *TurkishAlphabet) GetFirstLetter(s string) *TurkicLetter {
	if len(s) == 0 {
		return Undefined
	}
	runes := []rune(s)
	return ta.GetLetter(runes[0])
}

// LastChar returns the last character of the string
func LastChar(s string) rune {
	runes := []rune(s)
	return runes[len(runes)-1]
}

// Voice applies voicing to a character
func (ta *TurkishAlphabet) Voice(c rune) rune {
	if res, ok := ta.VoicingMap[c]; ok {
		return res
	}
	return c
}

// Devoice applies devoicing to a character
func (ta *TurkishAlphabet) Devoice(c rune) rune {
	if res, ok := ta.DevoicingMap[c]; ok {
		return res
	}
	return c
}
