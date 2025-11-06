package tokenizer

import (
	"testing"

	"github.com/kalaomer/zemberek-go/morphology"
)

func TestGetMorphology(t *testing.T) {
	morph1 := getMorphology()
	if morph1 == nil {
		t.Fatal("getMorphology() returned nil")
	}

	morph2 := getMorphology()
	if morph2 == nil {
		t.Fatal("getMorphology() returned nil on second call")
	}

	// Should return same instance (singleton)
	if morph1 != morph2 {
		t.Error("getMorphology() did not return the same instance")
	}
}

func TestGoTokenizeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string // Expected stems
	}{
		{
			name:     "Simple word",
			input:    "kitap",
			expected: []string{"kitap"},
		},
		{
			name:     "Inflected word",
			input:    "kitaplar",
			expected: []string{"kitap"},
		},
		{
			name:     "Multiple words",
			input:    "kitapları okuyorum",
			expected: []string{"kitap", "oku"},
		},
		{
			name:     "Turkish characters",
			input:    "şehirler",
			expected: []string{"şehir"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			morph := getMorphology()
			tokens := morphology.StemTextWithPositions(tt.input, morph)

			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				return
			}

			for i, token := range tokens {
				if token.Stem != tt.expected[i] {
					t.Errorf("Token %d: expected stem '%s', got '%s'",
						i, tt.expected[i], token.Stem)
				}
			}
		})
	}
}

func TestGetTokenizerStruct(t *testing.T) {
	ptr := GetTokenizerStruct()
	if ptr == nil {
		t.Error("GetTokenizerStruct() returned nil")
	}
}

func BenchmarkTokenization(b *testing.B) {
	text := "Kitapları okuyorum ve çok seviyorum. Yazılım geliştiriyorum."
	morph := getMorphology()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = morphology.StemTextWithPositions(text, morph)
	}
}
