package morphology

import (
	"testing"
)

func TestStemTextWithPositions_Basic(t *testing.T) {
	morph := CreateWithDefaults()
	text := "Kitapları okuyorum"

	tokens := StemTextWithPositions(text, morph)

	if len(tokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d", len(tokens))
	}

	// First token: "Kitapları"
	if tokens[0].Original != "Kitapları" {
		t.Errorf("Expected original 'Kitapları', got '%s'", tokens[0].Original)
	}
	if tokens[0].StartByte != 0 {
		t.Errorf("Expected start byte 0, got %d", tokens[0].StartByte)
	}
	// "Kitapları" = 10 bytes in UTF-8 (ı is 2 bytes: K=1,i=1,t=1,a=1,p=1,l=1,a=1,r=1,ı=2)
	if tokens[0].EndByte != 10 {
		t.Errorf("Expected end byte 10, got %d", tokens[0].EndByte)
	}

	// Second token: "okuyorum"
	if tokens[1].Original != "okuyorum" {
		t.Errorf("Expected original 'okuyorum', got '%s'", tokens[1].Original)
	}
	// Space is at byte 10, "okuyorum" starts at byte 11
	if tokens[1].StartByte != 11 {
		t.Errorf("Expected start byte 11, got %d", tokens[1].StartByte)
	}
	// "okuyorum" = 8 bytes, so ends at 11+8=19
	if tokens[1].EndByte != 19 {
		t.Errorf("Expected end byte 19, got %d", tokens[1].EndByte)
	}
}

func TestStemTextWithPositions_UTF8Handling(t *testing.T) {
	morph := CreateWithDefaults()
	// Turkish chars: ş=2bytes, ı=2bytes, ğ=2bytes, ö=2bytes, ü=2bytes, ç=2bytes
	text := "şöyle güzel"

	tokens := StemTextWithPositions(text, morph)

	if len(tokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d", len(tokens))
	}

	// "şöyle" = 7 bytes (ş=2, ö=2, y=1, l=1, e=1)
	if tokens[0].Original != "şöyle" {
		t.Errorf("Expected original 'şöyle', got '%s'", tokens[0].Original)
	}
	if tokens[0].StartByte != 0 {
		t.Errorf("Expected start byte 0, got %d", tokens[0].StartByte)
	}
	if tokens[0].EndByte != 7 {
		t.Errorf("Expected end byte 7, got %d", tokens[0].EndByte)
	}

	// "güzel" starts at byte 8 (after şöyle + space)
	// "güzel" = 6 bytes (g=1, ü=2, z=1, e=1, l=1)
	if tokens[1].Original != "güzel" {
		t.Errorf("Expected original 'güzel', got '%s'", tokens[1].Original)
	}
	if tokens[1].StartByte != 8 {
		t.Errorf("Expected start byte 8, got %d", tokens[1].StartByte)
	}
	if tokens[1].EndByte != 14 {
		t.Errorf("Expected end byte 14, got %d", tokens[1].EndByte)
	}
}

func TestStemTextWithPositions_WithPunctuation(t *testing.T) {
	morph := CreateWithDefaults()
	text := "Merhaba, dünya!"

	tokens := StemTextWithPositions(text, morph)

	// Should only return word tokens, not punctuation
	if len(tokens) != 2 {
		t.Fatalf("Expected 2 word tokens, got %d", len(tokens))
	}

	if tokens[0].Original != "Merhaba" {
		t.Errorf("Expected 'Merhaba', got '%s'", tokens[0].Original)
	}

	if tokens[1].Original != "dünya" {
		t.Errorf("Expected 'dünya', got '%s'", tokens[1].Original)
	}
}

func TestStemTextWithPositions_NoMorphology(t *testing.T) {
	// Test without morphology (nil)
	text := "test kelimeleri"

	tokens := StemTextWithPositions(text, nil)

	if len(tokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d", len(tokens))
	}

	// Without morphology, stem should equal original
	if tokens[0].Stem != tokens[0].Original {
		t.Errorf("Expected stem to equal original without morphology")
	}
	if tokens[1].Stem != tokens[1].Original {
		t.Errorf("Expected stem to equal original without morphology")
	}
}

func TestStemTextWithPositions_EmptyText(t *testing.T) {
	morph := CreateWithDefaults()
	text := ""

	tokens := StemTextWithPositions(text, morph)

	if len(tokens) != 0 {
		t.Errorf("Expected 0 tokens for empty text, got %d", len(tokens))
	}
}

func TestStemTextWithPositions_OnlyPunctuation(t *testing.T) {
	morph := CreateWithDefaults()
	text := "... !!! ???"

	tokens := StemTextWithPositions(text, morph)

	// Should return 0 tokens (no words)
	if len(tokens) != 0 {
		t.Errorf("Expected 0 word tokens, got %d", len(tokens))
	}
}

func TestStemTextWithPositions_Apostrophe(t *testing.T) {
	morph := CreateWithDefaults()
	text := "Ali'nin kitabı"

	tokens := StemTextWithPositions(text, morph)

	if len(tokens) < 1 {
		t.Fatalf("Expected at least 1 token, got %d", len(tokens))
	}

	// "Ali'nin" should be treated as one word
	if tokens[0].Original != "Ali'nin" {
		t.Errorf("Expected 'Ali'nin' as one token, got '%s'", tokens[0].Original)
	}
}

func TestIsWordChar(t *testing.T) {
	tests := []struct {
		char     rune
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'ş', true},
		{'İ', true},
		{'5', true},
		{'\'', true},
		{' ', false},
		{'.', false},
		{',', false},
		{'!', false},
	}

	for _, tt := range tests {
		result := isWordChar(tt.char)
		if result != tt.expected {
			t.Errorf("isWordChar('%c') = %v, expected %v", tt.char, result, tt.expected)
		}
	}
}

func TestIsWordToken(t *testing.T) {
	tests := []struct {
		token    string
		expected bool
	}{
		{"kitap", true},
		{"123abc", true},
		{"Ali'nin", true},
		{"123", false},    // only digits
		{"...", false},    // only punctuation
		{"", false},       // empty
		{"test123", true}, // mixed
	}

	for _, tt := range tests {
		result := isWordToken(tt.token)
		if result != tt.expected {
			t.Errorf("isWordToken('%s') = %v, expected %v", tt.token, result, tt.expected)
		}
	}
}

func BenchmarkStemTextWithPositions(b *testing.B) {
	morph := CreateWithDefaults()
	text := "Kitapları okuyorum ve yazıyorum. Bugün hava çok güzel."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StemTextWithPositions(text, morph)
	}
}

func TestStemTextWithPositions_ByteOffsetValidation(t *testing.T) {
	morph := CreateWithDefaults()
	text := "Türkçe metinleri işliyorum"

	tokens := StemTextWithPositions(text, morph)

	// Validate that we can extract original words using byte offsets
	for _, token := range tokens {
		extracted := text[token.StartByte:token.EndByte]
		if extracted != token.Original {
			t.Errorf("Byte offset mismatch: extracted '%s', expected '%s'", extracted, token.Original)
		}
	}
}
