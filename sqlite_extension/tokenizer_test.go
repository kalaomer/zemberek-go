package sqlite_extension

import (
	"testing"
)

func TestTokenizeWithZemberek(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple Turkish text",
			input:    "Merhaba dünya",
			expected: []string{"merhaba", "dünya"},
		},
		{
			name:     "Text with punctuation",
			input:    "Merhaba, dünya!",
			expected: []string{"merhaba", "dünya"},
		},
		{
			name:     "Turkish characters",
			input:    "çalışma şekli ğöüşıi",
			expected: []string{"çalışma", "şekli", "ğöüşıi"},
		},
		{
			name:     "Mixed case",
			input:    "İstanbul ANKARA",
			expected: []string{"istanbul", "ankara"},
		},
		{
			name:     "Numbers and text",
			input:    "2024 yılı",
			expected: []string{"2024", "yılı"},
		},
		{
			name:     "Multiple spaces",
			input:    "bir  iki   üç",
			expected: []string{"bir", "iki", "üç"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tokens []string
			callback := func(flags int, token string, start, end int) int {
				tokens = append(tokens, token)
				return 0 // SQLITE_OK
			}

			rc := tokenizeWithZemberek(tt.input, callback)
			if rc != 0 {
				t.Errorf("tokenizeWithZemberek() returned %d, want 0", rc)
			}

			if len(tokens) != len(tt.expected) {
				t.Errorf("got %d tokens, want %d: %v", len(tokens), len(tt.expected), tokens)
				return
			}

			for i, token := range tokens {
				if token != tt.expected[i] {
					t.Errorf("token[%d] = %q, want %q", i, token, tt.expected[i])
				}
			}
		})
	}
}

func TestAdvancedTokenizer_Tokenize(t *testing.T) {
	tokenizer, err := NewAdvancedTokenizer(true, false)
	if err != nil {
		t.Fatalf("NewAdvancedTokenizer() error = %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Basic tokenization",
			input:    "Türkçe metin",
			expected: []string{"türkçe", "metin"},
		},
		{
			name:     "With punctuation",
			input:    "Merhaba, nasılsın?",
			expected: []string{"merhaba", "nasılsın"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tokens []string
			callback := func(flags int, token string, start, end int) int {
				tokens = append(tokens, token)
				return 0
			}

			rc := tokenizer.Tokenize(tt.input, callback)
			if rc != 0 {
				t.Errorf("Tokenize() returned %d, want 0", rc)
			}

			if len(tokens) != len(tt.expected) {
				t.Errorf("got %d tokens, want %d: %v", len(tokens), len(tt.expected), tokens)
				return
			}

			for i, token := range tokens {
				if token != tt.expected[i] {
					t.Errorf("token[%d] = %q, want %q", i, token, tt.expected[i])
				}
			}
		})
	}
}

func TestAdvancedTokenizer_TurkishLowerCase(t *testing.T) {
	tokenizer, err := NewAdvancedTokenizer(true, false)
	if err != nil {
		t.Fatalf("NewAdvancedTokenizer() error = %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"ISTANBUL", "ıstanbul"},
		{"İSTANBUL", "istanbul"},
		{"ÇALIŞMA", "çalışma"},
		{"ÖĞRENCI", "öğrenci"},
		{"ABC", "abc"},
		{"IiİıĞğÜüŞşÇçÖö", "iıiığğüüşşççöö"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tokenizer.turkishLowerCase(tt.input)
			if result != tt.expected {
				t.Errorf("turkishLowerCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAdvancedTokenizer_RemoveTurkishDiacritics(t *testing.T) {
	tokenizer, err := NewAdvancedTokenizer(false, true)
	if err != nil {
		t.Fatalf("NewAdvancedTokenizer() error = %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"çalışma", "calisma"},
		{"Şişli", "Sisli"},
		{"ğöüı", "goui"},
		{"ÇĞIÖŞÜ", "CGIOSU"},
		{"İstanbul", "Istanbul"},
		{"normal", "normal"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tokenizer.removeTurkishDiacritics(tt.input)
			if result != tt.expected {
				t.Errorf("removeTurkishDiacritics(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAdvancedTokenizer_WithDiacriticRemoval(t *testing.T) {
	tokenizer, err := NewAdvancedTokenizer(true, true)
	if err != nil {
		t.Fatalf("NewAdvancedTokenizer() error = %v", err)
	}

	input := "Çalışma öğrenci şehir"
	expected := []string{"calisma", "ogrenci", "sehir"}

	var tokens []string
	callback := func(flags int, token string, start, end int) int {
		tokens = append(tokens, token)
		return 0
	}

	rc := tokenizer.Tokenize(input, callback)
	if rc != 0 {
		t.Errorf("Tokenize() returned %d, want 0", rc)
	}

	if len(tokens) != len(expected) {
		t.Errorf("got %d tokens, want %d: %v", len(tokens), len(expected), tokens)
		return
	}

	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("token[%d] = %q, want %q", i, token, expected[i])
		}
	}
}

func TestIsPunctuation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{".", true},
		{",", true},
		{"!", true},
		{"?", true},
		{"...", true},
		{"!?", true},
		{"word", false},
		{"word.", false},
		{".word", false},
		{"123", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isPunctuation(tt.input)
			if result != tt.expected {
				t.Errorf("isPunctuation(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkTokenizeWithZemberek(b *testing.B) {
	text := "Bu bir Türkçe metin örneğidir. Zemberek, Türkçe doğal dil işleme kütüphanesidir. İstanbul, Türkiye'nin en büyük şehridir."
	callback := func(flags int, token string, start, end int) int {
		return 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizeWithZemberek(text, callback)
	}
}

func BenchmarkAdvancedTokenizer(b *testing.B) {
	tokenizer, err := NewAdvancedTokenizer(true, true)
	if err != nil {
		b.Fatalf("NewAdvancedTokenizer() error = %v", err)
	}

	text := "Bu bir Türkçe metin örneğidir. Zemberek, Türkçe doğal dil işleme kütüphanesidir. İstanbul, Türkiye'nin en büyük şehridir."
	callback := func(flags int, token string, start, end int) int {
		return 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize(text, callback)
	}
}
