package morphotactics

import (
	"testing"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

// Test GetPrefixMatches with simple words
func TestStemTransitions_GetPrefixMatches_Simple(t *testing.T) {
	// Create test lexicon
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("kalem", "kalem", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("kitap", "kitap", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("ev", "ev", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	tests := []struct {
		name          string
		input         string
		expectedCount int
		expectedStems []string
	}{
		{
			name:          "exact match - kalem",
			input:         "kalem",
			expectedCount: 1,
			expectedStems: []string{"kalem"},
		},
		{
			name:          "exact match - kitap",
			input:         "kitap",
			expectedCount: 1,
			expectedStems: []string{"kitap"},
		},
		{
			name:          "exact match - ev",
			input:         "ev",
			expectedCount: 1,
			expectedStems: []string{"ev"},
		},
		{
			name:          "inflected form - kalemin",
			input:         "kalemin",
			expectedCount: 1,
			expectedStems: []string{"kalem"},
		},
		{
			name:          "inflected form - kitaplar",
			input:         "kitaplar",
			expectedCount: 1,
			expectedStems: []string{"kitap"},
		},
		{
			name:          "inflected form - evler",
			input:         "evler",
			expectedCount: 1,
			expectedStems: []string{"ev"},
		},
		{
			name:          "no match - xyz",
			input:         "xyz",
			expectedCount: 0,
			expectedStems: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := stemTrans.GetPrefixMatches(tt.input, false)

			if len(matches) != tt.expectedCount {
				t.Errorf("GetPrefixMatches(%q) returned %d matches, expected %d",
					tt.input, len(matches), tt.expectedCount)
			}

			for i, expectedStem := range tt.expectedStems {
				if i >= len(matches) {
					t.Errorf("Expected stem %q at index %d, but not enough matches", expectedStem, i)
					continue
				}
				if matches[i].Surface != expectedStem {
					t.Errorf("Expected stem %q, got %q", expectedStem, matches[i].Surface)
				}
			}
		})
	}
}

// Test GetPrefixMatches with multiple possible stems
func TestStemTransitions_GetPrefixMatches_Ambiguous(t *testing.T) {
	// Create test lexicon with ambiguous words
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("ara", "ara", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("ara", "ara", turkish.Verb, turkish.NonePos, nil, "", 1),
		lexicon.NewDictionaryItem("kar", "kar", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	// Test ambiguous word "ara" (both noun and verb)
	matches := stemTrans.GetPrefixMatches("ara", false)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'ara', got %d", len(matches))
	}

	// Both should have surface "ara"
	for _, match := range matches {
		if match.Surface != "ara" {
			t.Errorf("Expected surface 'ara', got %q", match.Surface)
		}
	}
}

// Test GetPrefixMatches with compound words
func TestStemTransitions_GetPrefixMatches_Compounds(t *testing.T) {
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("masa", "masa", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("masaüstü", "masaüstü", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	// Test "masaüstü" - should match both "masa" and "masaüstü"
	matches := stemTrans.GetPrefixMatches("masaüstü", false)

	// Should find both stems (longest first in our implementation)
	if len(matches) < 1 {
		t.Errorf("Expected at least 1 match for 'masaüstü', got %d", len(matches))
	}

	found := make(map[string]bool)
	for _, match := range matches {
		found[match.Surface] = true
	}

	if !found["masaüstü"] && !found["masa"] {
		t.Error("Expected to find either 'masaüstü' or 'masa' as stem")
	}
}

// Test GetTransitions for exact surface form
func TestStemTransitions_GetTransitions(t *testing.T) {
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("kalem", "kalem", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("kitap", "kitap", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	tests := []struct {
		surface       string
		expectedCount int
	}{
		{"kalem", 1},
		{"kitap", 1},
		{"xyz", 0},
	}

	for _, tt := range tests {
		t.Run(tt.surface, func(t *testing.T) {
			transitions := stemTrans.GetTransitions(tt.surface)
			if len(transitions) != tt.expectedCount {
				t.Errorf("GetTransitions(%q) returned %d, expected %d",
					tt.surface, len(transitions), tt.expectedCount)
			}
		})
	}
}

// Test stem transitions with Turkish-specific attributes
func TestStemTransitions_TurkishAttributes(t *testing.T) {
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("kitap", "kitap", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	matches := stemTrans.GetPrefixMatches("kitap", false)

	if len(matches) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matches))
	}

	match := matches[0]

	// Check phonetic attributes
	if match.PhoneticAttributes == nil {
		t.Error("Phonetic attributes should not be nil")
	}

	// kitap ends with 'p' which is voiceless
	if !match.PhoneticAttributes[turkish.LastLetterVoiceless] {
		t.Error("Expected LastLetterVoiceless attribute for 'kitap'")
	}

	// kitap has 'a' as last vowel which is back
	if !match.PhoneticAttributes[turkish.LastVowelBack] {
		t.Error("Expected LastVowelBack attribute for 'kitap'")
	}
}

// Test stem transitions with different POS types
func TestStemTransitions_DifferentPOS(t *testing.T) {
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("güzel", "güzel", turkish.Adjective, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("yap", "yap", turkish.Verb, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("ev", "ev", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	tests := []struct {
		word        string
		expectedPOS turkish.PrimaryPos
	}{
		{"güzel", turkish.Adjective},
		{"yap", turkish.Verb},
		{"ev", turkish.Noun},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			matches := stemTrans.GetPrefixMatches(tt.word, false)
			if len(matches) != 1 {
				t.Fatalf("Expected 1 match for %q, got %d", tt.word, len(matches))
			}

			if matches[0].Item.PrimaryPos != tt.expectedPOS {
				t.Errorf("Expected POS %v for %q, got %v",
					tt.expectedPOS, tt.word, matches[0].Item.PrimaryPos)
			}
		})
	}
}

// Test edge cases
func TestStemTransitions_EdgeCases(t *testing.T) {
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("a", "a", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("ı", "ı", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{"empty string", "", 0},
		{"single char - a", "a", 1},
		{"single char - ı", "ı", 1},
		{"single char - not in lexicon", "x", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := stemTrans.GetPrefixMatches(tt.input, false)
			if len(matches) != tt.expectedCount {
				t.Errorf("GetPrefixMatches(%q) returned %d matches, expected %d",
					tt.input, len(matches), tt.expectedCount)
			}
		})
	}
}

// Benchmark stem matching performance
func BenchmarkStemTransitions_GetPrefixMatches(b *testing.B) {
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("kalem", "kalem", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("kitap", "kitap", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("ev", "ev", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("masa", "masa", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("sandalye", "sandalye", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	testWords := []string{
		"kalem", "kalemin", "kitap", "kitaplar",
		"ev", "evler", "masa", "masalar", "sandalye",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, word := range testWords {
			stemTrans.GetPrefixMatches(word, false)
		}
	}
}

// Benchmark with larger lexicon
func BenchmarkStemTransitions_LargeLexicon(b *testing.B) {
	// Create a larger lexicon
	items := make([]*lexicon.DictionaryItem, 0, 100)
	words := []string{
		"ev", "el", "su", "iş", "gün", "yıl", "gece", "sabah", "öğle", "akşam",
		"baş", "göz", "kulak", "ağız", "dil", "diş", "saç", "ayak", "kol", "parmak",
		"masa", "sandalye", "kapı", "pencere", "duvar", "tavan", "zemin", "oda", "salon", "mutfak",
		"kalem", "kitap", "defter", "silgi", "kalemlik", "çanta", "tahta", "tebeşir", "sınıf", "okul",
		"öğretmen", "öğrenci", "müdür", "veli", "ders", "sınav", "ödev", "not", "diploma", "belge",
	}

	for i, word := range words {
		items = append(items, lexicon.NewDictionaryItem(word, word, turkish.Noun, turkish.NonePos, nil, "", i))
	}

	lex := lexicon.NewRootLexicon(items)
	tm := NewTurkishMorphotactics(lex)
	stemTrans := tm.GetStemTransitions()

	testWords := []string{
		"evler", "eller", "sular", "işler", "günler",
		"masalar", "kitaplar", "öğretmenler", "öğrenciler",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, word := range testWords {
			stemTrans.GetPrefixMatches(word, false)
		}
	}
}
