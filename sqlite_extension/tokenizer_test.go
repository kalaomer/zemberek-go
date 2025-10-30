package sqlite_extension

import (
	"testing"
)

func TestZemberekTokenizer_Tokenize(t *testing.T) {
	tokenizer := NewZemberekTokenizer()

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
			tokens := tokenizer.Tokenize(tt.input)

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

func TestZemberekTokenizer_TokenizeWithPositions(t *testing.T) {
	tokenizer := NewZemberekTokenizer()

	input := "Merhaba dünya"
	positions := tokenizer.TokenizeWithPositions(input)

	if len(positions) != 2 {
		t.Errorf("got %d tokens, want 2", len(positions))
		return
	}

	// Check first token
	if positions[0].Token != "merhaba" {
		t.Errorf("positions[0].Token = %q, want %q", positions[0].Token, "merhaba")
	}
	if positions[0].Start != 0 || positions[0].End != 7 {
		t.Errorf("positions[0] = {%d, %d}, want {0, 7}", positions[0].Start, positions[0].End)
	}

	// Check second token
	if positions[1].Token != "dünya" {
		t.Errorf("positions[1].Token = %q, want %q", positions[1].Token, "dünya")
	}
}

func TestTurkishLowerCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ISTANBUL", "ıstanbul"},     // I (U+0049) -> ı
		{"İSTANBUL", "istanbul"},     // İ (U+0130) -> i
		{"ÇALIŞMA", "çalışma"},       // Turkish chars
		{"ÖĞRENCI", "öğrencı"},       // Last char is I (U+0049) -> ı
		{"ÖĞRENCİ", "öğrenci"},       // Last char is İ (U+0130) -> i
		{"ABC", "abc"},               // Normal ASCII
		{"IiİıĞğÜüŞşÇçÖö", "ıiiığğüüşşççöö"}, // All combinations
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := turkishLowerCase(tt.input)
			if result != tt.expected {
				t.Errorf("turkishLowerCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTurkishUpperCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"istanbul", "İSTANBUL"},
		{"ıstanbul", "ISTANBUL"},
		{"çalışma", "ÇALIŞMA"},
		{"öğrenci", "ÖĞRENCİ"},
		{"abc", "ABC"},
		{"iıİIğĞüÜşŞçÇöÖ", "İIİIĞĞÜÜŞŞÇÇÖÖ"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TurkishUpperCase(tt.input)
			if result != tt.expected {
				t.Errorf("TurkishUpperCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveTurkishDiacritics(t *testing.T) {
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
			result := removeTurkishDiacritics(tt.input)
			if result != tt.expected {
				t.Errorf("removeTurkishDiacritics(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestZemberekTokenizer_WithDiacriticRemoval(t *testing.T) {
	tokenizer := NewZemberekTokenizerWithOptions(true, true)

	input := "Çalışma öğrenci şehir"
	expected := []string{"calisma", "ogrenci", "sehir"}

	tokens := tokenizer.Tokenize(input)

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

func TestTokenizeText(t *testing.T) {
	input := "Merhaba dünya"
	expected := []string{"merhaba", "dünya"}

	tokens := TokenizeText(input)

	if len(tokens) != len(expected) {
		t.Errorf("got %d tokens, want %d", len(tokens), len(expected))
		return
	}

	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("token[%d] = %q, want %q", i, token, expected[i])
		}
	}
}

func TestNormalizeForSearch(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Merhaba Dünya", "merhaba dünya"},
		{"TÜRKÇE", "türkçe"},
		{"İstanbul, Ankara!", "istanbul ankara"},
		{"  çok   boşluk  ", "çok boşluk"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeForSearch(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeForSearch(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkTokenize(b *testing.B) {
	tokenizer := NewZemberekTokenizer()
	text := "Bu bir Türkçe metin örneğidir. Zemberek, Türkçe doğal dil işleme kütüphanesidir. İstanbul, Türkiye'nin en büyük şehridir."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize(text)
	}
}

func BenchmarkTokenizeWithPositions(b *testing.B) {
	tokenizer := NewZemberekTokenizer()
	text := "Bu bir Türkçe metin örneğidir. Zemberek, Türkçe doğal dil işleme kütüphanesidir. İstanbul, Türkiye'nin en büyük şehridir."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.TokenizeWithPositions(text)
	}
}

func BenchmarkTurkishLowerCase(b *testing.B) {
	text := "ISTANBUL İSTANBUL ÇALIŞMA ÖĞRENCI"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		turkishLowerCase(text)
	}
}

func BenchmarkRemoveDiacritics(b *testing.B) {
	text := "çalışma öğrenci şehir ığüöç"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		removeTurkishDiacritics(text)
	}
}
