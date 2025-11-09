package turkish

import "testing"

// TestIsTurkishLetter tests the IsTurkishLetter utility function
func TestIsTurkishLetter(t *testing.T) {
	alphabet := Instance

	tests := []struct {
		name     string
		char     rune
		expected bool
	}{
		// ASCII lowercase
		{"lowercase a", 'a', true},
		{"lowercase z", 'z', true},
		{"lowercase m", 'm', true},

		// ASCII uppercase
		{"uppercase A", 'A', true},
		{"uppercase Z", 'Z', true},
		{"uppercase M", 'M', true},

		// Turkish-specific lowercase
		{"turkish ç", 'ç', true},
		{"turkish ğ", 'ğ', true},
		{"turkish ı", 'ı', true},
		{"turkish ö", 'ö', true},
		{"turkish ş", 'ş', true},
		{"turkish ü", 'ü', true},

		// Turkish-specific uppercase
		{"turkish Ç", 'Ç', true},
		{"turkish Ğ", 'Ğ', true},
		{"turkish İ", 'İ', true},
		{"turkish Ö", 'Ö', true},
		{"turkish Ş", 'Ş', true},
		{"turkish Ü", 'Ü', true},

		// Extended Turkish
		{"extended â", 'â', true},
		{"extended î", 'î', true},
		{"extended û", 'û', true},
		{"extended Â", 'Â', true},
		{"extended Î", 'Î', true},
		{"extended Û", 'Û', true},

		// Non-letters
		{"digit 0", '0', false},
		{"digit 5", '5', false},
		{"digit 9", '9', false},
		{"space", ' ', false},
		{"dot", '.', false},
		{"comma", ',', false},
		{"hyphen", '-', false},
		{"apostrophe", '\'', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alphabet.IsTurkishLetter(tt.char)
			if got != tt.expected {
				t.Errorf("IsTurkishLetter(%q) = %v, expected %v", tt.char, got, tt.expected)
			}
		})
	}
}

// TestIsDigit tests the IsDigit utility function
func TestIsDigit(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected bool
	}{
		// Digits
		{"zero", '0', true},
		{"five", '5', true},
		{"nine", '9', true},

		// Non-digits
		{"letter a", 'a', false},
		{"letter Z", 'Z', false},
		{"turkish ç", 'ç', false},
		{"space", ' ', false},
		{"dot", '.', false},
		{"slash", '/', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDigit(tt.char)
			if got != tt.expected {
				t.Errorf("IsDigit(%q) = %v, expected %v", tt.char, got, tt.expected)
			}
		})
	}
}

// TestIsTurkishLetterVsIsDigit tests that letters and digits are mutually exclusive
func TestIsTurkishLetterVsIsDigit(t *testing.T) {
	alphabet := Instance

	// Test all printable ASCII
	for r := rune(33); r <= rune(126); r++ {
		isLetter := alphabet.IsTurkishLetter(r)
		isDigit := IsDigit(r)

		// A character can be either letter or digit, but not both
		if isLetter && isDigit {
			t.Errorf("Character %q is both letter and digit", r)
		}

		// Specifically check a-z, A-Z are letters
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			if !isLetter {
				t.Errorf("Character %q should be a letter", r)
			}
		}

		// Specifically check 0-9 are digits
		if r >= '0' && r <= '9' {
			if !isDigit {
				t.Errorf("Character %q should be a digit", r)
			}
		}
	}

	// Test Turkish-specific characters
	turkishChars := []rune{'ç', 'ğ', 'ı', 'ö', 'ş', 'ü', 'Ç', 'Ğ', 'İ', 'Ö', 'Ş', 'Ü', 'â', 'î', 'û', 'Â', 'Î', 'Û'}
	for _, r := range turkishChars {
		if !alphabet.IsTurkishLetter(r) {
			t.Errorf("Turkish character %q should be recognized as letter", r)
		}
		if IsDigit(r) {
			t.Errorf("Turkish character %q should not be recognized as digit", r)
		}
	}
}
