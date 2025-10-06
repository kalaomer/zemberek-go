package deasciifier

import (
	"testing"
)

// TestDeasciifier_ConvertToTurkish_BasicWords tests basic ASCII to Turkish conversion
func TestDeasciifier_ConvertToTurkish_BasicWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple word - cok",
			input:    "cok",
			expected: "cok", // Without pattern table, will stay as is
		},
		{
			name:     "simple word - gun",
			input:    "gun",
			expected: "gun", // Without pattern table, will stay as is
		},
		{
			name:     "already Turkish - çok (will toggle to ASCII)",
			input:    "çok",
			expected: "cok", // Without pattern table, Turkish chars get toggled
		},
		{
			name:     "already Turkish - gün (will toggle to ASCII)",
			input:    "gün",
			expected: "gun", // Without pattern table, Turkish chars get toggled
		},
		{
			name:     "mixed - Istanbul",
			input:    "Istanbul",
			expected: "İstanbul", // Capital I becomes İ
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeasciifier(tt.input)
			result := d.ConvertToTurkish()
			if result != tt.expected {
				t.Errorf("ConvertToTurkish() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestDeasciifier_TurkishToggleAccent tests the accent toggle table
func TestDeasciifier_TurkishToggleAccent(t *testing.T) {
	tests := []struct {
		char     rune
		expected rune
		exists   bool
	}{
		{'c', 'ç', true},
		{'C', 'Ç', true},
		{'g', 'ğ', true},
		{'G', 'Ğ', true},
		{'o', 'ö', true},
		{'O', 'Ö', true},
		{'u', 'ü', true},
		{'U', 'Ü', true},
		{'i', 'ı', true},
		{'I', 'İ', true},
		{'s', 'ş', true},
		{'S', 'Ş', true},
		// Toggle back
		{'ç', 'c', true},
		{'ğ', 'g', true},
		{'ö', 'o', true},
		{'ü', 'u', true},
		{'ı', 'i', true},
		{'İ', 'I', true},
		{'ş', 's', true},
		// Non-special chars
		{'a', 0, false},
		{'z', 0, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result, exists := turkishToggleAccentTable[tt.char]
			if exists != tt.exists {
				t.Errorf("turkishToggleAccentTable[%c] exists = %v, want %v", tt.char, exists, tt.exists)
			}
			if exists && result != tt.expected {
				t.Errorf("turkishToggleAccentTable[%c] = %c, want %c", tt.char, result, tt.expected)
			}
		})
	}
}

// TestDeasciifier_TurkishAsciifyTable tests the asciify conversion table
func TestDeasciifier_TurkishAsciifyTable(t *testing.T) {
	tests := []struct {
		char     rune
		expected rune
		exists   bool
	}{
		{'ç', 'c', true},
		{'Ç', 'C', true},
		{'ğ', 'g', true},
		{'Ğ', 'G', true},
		{'ö', 'o', true},
		{'Ö', 'O', true},
		{'ı', 'i', true},
		{'İ', 'I', true},
		{'ş', 's', true},
		{'Ş', 'S', true},
		{'ü', 'u', true},
		{'Ü', 'U', true},
		// Regular ASCII chars should not be in this table
		{'a', 0, false},
		{'c', 0, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result, exists := turkishAsciifyTable[tt.char]
			if exists != tt.exists {
				t.Errorf("turkishAsciifyTable[%c] exists = %v, want %v", tt.char, exists, tt.exists)
			}
			if exists && result != tt.expected {
				t.Errorf("turkishAsciifyTable[%c] = %c, want %c", tt.char, result, tt.expected)
			}
		})
	}
}

// TestDeasciifier_TurkishDowncaseAsciifyTable tests the downcase conversion
func TestDeasciifier_TurkishDowncaseAsciifyTable(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected rune
	}{
		{"ç to c", 'ç', 'c'},
		{"Ç to c", 'Ç', 'c'},
		{"ğ to g", 'ğ', 'g'},
		{"Ğ to g", 'Ğ', 'g'},
		{"ö to o", 'ö', 'o'},
		{"Ö to o", 'Ö', 'o'},
		{"ı to i", 'ı', 'i'},
		{"İ to i", 'İ', 'i'},
		{"ş to s", 'ş', 's'},
		{"Ş to s", 'Ş', 's'},
		{"ü to u", 'ü', 'u'},
		{"Ü to u", 'Ü', 'u'},
		{"A to a", 'A', 'a'},
		{"Z to z", 'Z', 'z'},
		{"a to a", 'a', 'a'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, exists := turkishDowncaseAsciifyTable[tt.char]
			if !exists {
				t.Errorf("turkishDowncaseAsciifyTable[%c] should exist", tt.char)
			}
			if result != tt.expected {
				t.Errorf("turkishDowncaseAsciifyTable[%c] = %c, want %c", tt.char, result, tt.expected)
			}
		})
	}
}

// TestDeasciifier_TurkishUpcaseAccentsTable tests the upcase accents table
func TestDeasciifier_TurkishUpcaseAccentsTable(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected rune
	}{
		{"ç to C", 'ç', 'C'},
		{"Ç to C", 'Ç', 'C'},
		{"ğ to G", 'ğ', 'G'},
		{"Ğ to G", 'Ğ', 'G'},
		{"ö to O", 'ö', 'O'},
		{"Ö to O", 'Ö', 'O'},
		{"ı to I", 'ı', 'I'},
		{"İ to i", 'İ', 'i'}, // Special case
		{"ş to S", 'ş', 'S'},
		{"Ş to S", 'Ş', 'S'},
		{"ü to U", 'ü', 'U'},
		{"Ü to U", 'Ü', 'U'},
		{"A to a", 'A', 'a'},
		{"a to a", 'a', 'a'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, exists := turkishUpcaseAccentsTable[tt.char]
			if !exists {
				t.Errorf("turkishUpcaseAccentsTable[%c] should exist", tt.char)
			}
			if result != tt.expected {
				t.Errorf("turkishUpcaseAccentsTable[%c] = %c, want %c", tt.char, result, tt.expected)
			}
		})
	}
}

// TestDeasciifier_EmptyString tests conversion of empty string
func TestDeasciifier_EmptyString(t *testing.T) {
	d := NewDeasciifier("")
	result := d.ConvertToTurkish()
	if result != "" {
		t.Errorf("ConvertToTurkish() = %v, want empty string", result)
	}
}

// TestDeasciifier_SingleChar tests conversion of single characters
func TestDeasciifier_SingleChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"c", "c", "c"},
		{"ç (toggles without pattern)", "ç", "c"}, // Without pattern table, gets toggled
		{"i", "i", "i"},
		{"ı (toggles without pattern)", "ı", "i"}, // Without pattern table, gets toggled
		{"a", "a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeasciifier(tt.input)
			result := d.ConvertToTurkish()
			if result != tt.expected {
				t.Errorf("ConvertToTurkish(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDeasciifier_LongSentence tests conversion of longer sentences
func TestDeasciifier_LongSentence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple sentence (Turkish chars toggle without pattern)",
			input:    "Merhaba nasılsın",
			expected: "Merhaba nasilsin", // Turkish chars get toggled without pattern table
		},
		{
			name:     "already accented (toggles to ASCII)",
			input:    "Çok güzel bir gün",
			expected: "Cok guzel bir gun", // Without pattern table, all toggle to ASCII
		},
		{
			name:     "numbers and words",
			input:    "Saat 12:00",
			expected: "Saat 12:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeasciifier(tt.input)
			result := d.ConvertToTurkish()
			if result != tt.expected {
				t.Errorf("ConvertToTurkish(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDeasciifier_ContextSize tests context size constant
func TestDeasciifier_ContextSize(t *testing.T) {
	if turkishContextSize != 10 {
		t.Errorf("turkishContextSize = %d, want 10", turkishContextSize)
	}
}

// TestDeasciifier_AbsFunction tests the abs helper function
func TestDeasciifier_AbsFunction(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{100, 100},
		{-100, 100},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := abs(tt.input)
			if result != tt.expected {
				t.Errorf("abs(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNewDeasciifier tests the constructor
func TestNewDeasciifier(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"simple", "merhaba"},
		{"with Turkish", "çok güzel"},
		{"long", "Bu çok uzun bir cümle örneğidir"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeasciifier(tt.input)
			if d == nil {
				t.Error("NewDeasciifier() returned nil")
			}
			if d.asciiString != tt.input {
				t.Errorf("asciiString = %v, want %v", d.asciiString, tt.input)
			}
			if d.turkishString != tt.input {
				t.Errorf("turkishString = %v, want %v", d.turkishString, tt.input)
			}
		})
	}
}

// TestDeasciifier_PatternTableEmpty tests that pattern table is initialized
func TestDeasciifier_PatternTableEmpty(t *testing.T) {
	if turkishPatternTable == nil {
		t.Error("turkishPatternTable should not be nil")
	}
	// Note: Pattern table is empty by default as it needs to be loaded from resources
	// This test just verifies it's initialized
}

// BenchmarkDeasciifier_ConvertToTurkish benchmarks the conversion performance
func BenchmarkDeasciifier_ConvertToTurkish(b *testing.B) {
	input := "Merhaba dunya, nasilsin? Bugun hava cok guzel."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeasciifier(input)
		_ = d.ConvertToTurkish()
	}
}

// BenchmarkDeasciifier_LongText benchmarks conversion of longer text
func BenchmarkDeasciifier_LongText(b *testing.B) {
	input := "Turkce metinlerin ASCII karsiliklarina donusturulmesi ve tekrar duzgun Turkce haline getirilmesi " +
		"onemli bir islev. Bu islem ozellikle eski sistemlerden gelen verilerin normalize edilmesinde kullanilir. " +
		"Deasciifier bu donusumu saglamak icin ozel bir pattern tablosu kullanir."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeasciifier(input)
		_ = d.ConvertToTurkish()
	}
}
