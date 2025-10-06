package normalization

import (
	"strings"
	"testing"
)

// TestNewTurkishSentenceNormalizer tests normalizer creation
func TestNewTurkishSentenceNormalizer(t *testing.T) {
	stemWords := []string{"git", "gel", "al", "ver", "yap"}

	// Test with empty resource path (will fail to load files, but struct should be created)
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Expected normalizer to be created even without resources, got error: %v", err)
	}

	if normalizer == nil {
		t.Fatal("Expected non-nil normalizer")
	}

	if normalizer.Replacements == nil {
		t.Error("Expected non-nil Replacements map")
	}
	if normalizer.NoSplitWords == nil {
		t.Error("Expected non-nil NoSplitWords map")
	}
	if normalizer.CommonSplits == nil {
		t.Error("Expected non-nil CommonSplits map")
	}
	if normalizer.LookupManual == nil {
		t.Error("Expected non-nil LookupManual map")
	}
}

// TestCandidate tests Candidate structure
func TestCandidate(t *testing.T) {
	candidate := NewCandidate("test")

	if candidate == nil {
		t.Fatal("Expected non-nil candidate")
	}
	if candidate.Content != "test" {
		t.Errorf("Expected content 'test', got '%s'", candidate.Content)
	}
	if candidate.Score != 1.0 {
		t.Errorf("Expected score 1.0, got %f", candidate.Score)
	}
}

// TestCandidates tests Candidates structure
func TestCandidates(t *testing.T) {
	c1 := NewCandidate("test1")
	c2 := NewCandidate("test2")
	candidates := NewCandidates("original", []*Candidate{c1, c2})

	if candidates == nil {
		t.Fatal("Expected non-nil candidates")
	}
	if candidates.Word != "original" {
		t.Errorf("Expected word 'original', got '%s'", candidates.Word)
	}
	if len(candidates.Candidates) != 2 {
		t.Errorf("Expected 2 candidates, got %d", len(candidates.Candidates))
	}
}

// TestHypothesis tests Hypothesis structure
func TestHypothesis(t *testing.T) {
	hyp := NewHypothesis()

	if hyp == nil {
		t.Fatal("Expected non-nil hypothesis")
	}
	if hyp.Score != 0.0 {
		t.Errorf("Expected score 0.0, got %f", hyp.Score)
	}
	if hyp.History != nil {
		t.Error("Expected nil History initially")
	}
	if hyp.Current != nil {
		t.Error("Expected nil Current initially")
	}
}

// TestHypothesis_Equals tests hypothesis equality
func TestHypothesis_Equals(t *testing.T) {
	hyp1 := NewHypothesis()
	hyp2 := NewHypothesis()

	if !hyp1.Equals(hyp1) {
		t.Error("Hypothesis should equal itself")
	}

	if !hyp1.Equals(hyp2) {
		t.Error("Two empty hypotheses should be equal")
	}

	hyp1.Current = NewCandidate("test")
	if hyp1.Equals(hyp2) {
		t.Error("Hypotheses with different current should not be equal")
	}

	// Use same candidate instance for equality
	sameCandidate := NewCandidate("test")
	hyp1.Current = sameCandidate
	hyp2.Current = sameCandidate
	if !hyp1.Equals(hyp2) {
		t.Error("Hypotheses with same current candidate instance should be equal")
	}
}

// TestHypothesis_Hash tests hypothesis hashing
func TestHypothesis_Hash(t *testing.T) {
	hyp1 := NewHypothesis()
	hyp2 := NewHypothesis()

	hash1 := hyp1.Hash()
	hash2 := hyp2.Hash()

	if hash1 != hash2 {
		t.Error("Equal hypotheses should have same hash")
	}

	hyp1.Current = NewCandidate("test")
	hash1 = hyp1.Hash()
	if hash1 == hash2 {
		t.Error("Different hypotheses should have different hash (usually)")
	}
}

// TestGetStartCandidate tests START sentinel
func TestGetStartCandidate(t *testing.T) {
	start := GetStartCandidate()

	if start == nil {
		t.Fatal("Expected non-nil START candidate")
	}
	if start.Content != "<s>" {
		t.Errorf("Expected START content '<s>', got '%s'", start.Content)
	}
}

// TestGetEndCandidate tests END sentinel
func TestGetEndCandidate(t *testing.T) {
	end := GetEndCandidate()

	if end == nil {
		t.Fatal("Expected non-nil END candidate")
	}
	if end.Content != "</s>" {
		t.Errorf("Expected END content '</s>', got '%s'", end.Content)
	}
}

// TestGetEndCandidates tests END candidates structure
func TestGetEndCandidates(t *testing.T) {
	endCandidates := GetEndCandidates()

	if endCandidates == nil {
		t.Fatal("Expected non-nil END candidates")
	}
	if endCandidates.Word != "</s>" {
		t.Errorf("Expected word '</s>', got '%s'", endCandidates.Word)
	}
	if len(endCandidates.Candidates) != 1 {
		t.Errorf("Expected 1 candidate, got %d", len(endCandidates.Candidates))
	}
}

// TestGetBestHypothesis tests best hypothesis selection
func TestGetBestHypothesis(t *testing.T) {
	// Test with empty list
	best := GetBestHypothesis([]*Hypothesis{})
	if best != nil {
		t.Error("Expected nil for empty list")
	}

	// Test with single hypothesis
	hyp1 := NewHypothesis()
	hyp1.Score = 5.0
	best = GetBestHypothesis([]*Hypothesis{hyp1})
	if best != hyp1 {
		t.Error("Expected hyp1 as best")
	}

	// Test with multiple hypotheses
	hyp2 := NewHypothesis()
	hyp2.Score = 10.0
	hyp3 := NewHypothesis()
	hyp3.Score = 3.0

	best = GetBestHypothesis([]*Hypothesis{hyp1, hyp2, hyp3})
	if best != hyp2 {
		t.Errorf("Expected hyp2 (score 10.0) as best, got score %f", best.Score)
	}
}

// TestTurkishSentenceNormalizer_GetCandidates tests candidate generation
func TestTurkishSentenceNormalizer_GetCandidates(t *testing.T) {
	stemWords := []string{"git", "gel", "al", "ver"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	// Test with word
	candidates := normalizer.getCandidates("test")
	if len(candidates) == 0 {
		t.Error("Expected at least one candidate (the word itself)")
	}

	// Should include original word
	found := false
	for _, c := range candidates {
		if c == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Candidates should include original word")
	}
}

// TestTurkishSentenceNormalizer_Normalize tests basic normalization
func TestTurkishSentenceNormalizer_Normalize(t *testing.T) {
	stemWords := []string{"git", "gel", "al", "ver", "yap", "gör"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		contains string // What the output should contain
	}{
		{
			name:     "simple sentence",
			input:    "merhaba dünya",
			contains: "merhaba",
		},
		{
			name:     "single word",
			input:    "test",
			contains: "test",
		},
		{
			name:     "empty string",
			input:    "",
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.Normalize(tt.input)
			if tt.contains != "" && result != "" && !contains(result, tt.contains) {
				t.Errorf("Expected output to contain '%s', got '%s'", tt.contains, result)
			}
		})
	}
}

// TestTurkishSentenceNormalizer_PreProcess tests preprocessing
func TestTurkishSentenceNormalizer_PreProcess(t *testing.T) {
	stemWords := []string{"git", "gel"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "uppercase conversion",
			input: "MERHABA DÜNYA",
		},
		{
			name:  "mixed case",
			input: "Merhaba Dünya",
		},
		{
			name:  "already lowercase",
			input: "merhaba dünya",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.preProcess(tt.input)
			// Result should be lowercase
			if result != strings.ToLower(result) {
				t.Errorf("Expected lowercase output, got '%s'", result)
			}
		})
	}
}

// TestTurkishSentenceNormalizer_ReplaceCommon tests common replacements
func TestTurkishSentenceNormalizer_ReplaceCommon(t *testing.T) {
	stemWords := []string{}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	// Add test replacement
	normalizer.Replacements["tmm"] = "tamam"

	tokens := []string{"tmm", "yarın", "geliyorum"}
	result := normalizer.replaceCommon(tokens)

	if !contains(result, "tamam") {
		t.Errorf("Expected 'tamam' in result, got '%s'", result)
	}
}

// TestTurkishSentenceNormalizer_CombineCommon tests word combining
func TestTurkishSentenceNormalizer_CombineCommon(t *testing.T) {
	stemWords := []string{"kitap", "bilgisayar"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	tests := []struct {
		name     string
		w1       string
		w2       string
		expected string
	}{
		{
			name:     "no combination",
			w1:       "test",
			w2:       "word",
			expected: "",
		},
		{
			name:     "apostrophe combination",
			w1:       "kitap",
			w2:       "'ı",
			expected: "kitap'ı",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.combineCommon(tt.w1, tt.w2)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestTurkishSentenceNormalizer_SeparateCommon tests word separation
func TestTurkishSentenceNormalizer_SeparateCommon(t *testing.T) {
	stemWords := []string{"git", "gel"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	// Add no-split word
	normalizer.NoSplitWords["gitme"] = true

	result := normalizer.separateCommon("gitme", false)
	if result != "gitme" {
		t.Errorf("No-split word should not be separated, got '%s'", result)
	}
}

// TestTurkishSentenceNormalizer_NormalizeWithBeamSearch tests beam search normalization
func TestTurkishSentenceNormalizer_NormalizeWithBeamSearch(t *testing.T) {
	stemWords := []string{"git", "gel", "al", "ver"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "simple sentence",
			input: "merhaba dünya",
		},
		{
			name:  "single word",
			input: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.NormalizeWithBeamSearch(tt.input)
			if result == "" && tt.input != "" {
				t.Error("Expected non-empty result for non-empty input")
			}
		})
	}
}

// TestTurkishSentenceNormalizer_DecodeSimple tests simple decoding
func TestTurkishSentenceNormalizer_DecodeSimple(t *testing.T) {
	stemWords := []string{}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		t.Fatalf("Failed to create normalizer: %v", err)
	}

	candidates1 := NewCandidates("word1", []*Candidate{
		NewCandidate("option1"),
		NewCandidate("option2"),
	})
	candidates2 := NewCandidates("word2", []*Candidate{
		NewCandidate("option3"),
	})

	candidatesList := []*Candidates{candidates1, candidates2}
	result := normalizer.decodeSimple(candidatesList)

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}
	if result[0] != "option1" {
		t.Errorf("Expected first result to be 'option1', got '%s'", result[0])
	}
	if result[1] != "option3" {
		t.Errorf("Expected second result to be 'option3', got '%s'", result[1])
	}
}

// Helper function
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// BenchmarkTurkishSentenceNormalizer_Normalize benchmarks normalization
func BenchmarkTurkishSentenceNormalizer_Normalize(b *testing.B) {
	stemWords := []string{"git", "gel", "al", "ver", "yap", "gör", "bil", "kal"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		b.Fatalf("Failed to create normalizer: %v", err)
	}

	sentence := "merhaba dünya nasılsın"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = normalizer.Normalize(sentence)
	}
}

// BenchmarkTurkishSentenceNormalizer_GetCandidates benchmarks candidate generation
func BenchmarkTurkishSentenceNormalizer_GetCandidates(b *testing.B) {
	stemWords := []string{"git", "gel", "al", "ver", "yap", "gör", "bil", "kal"}
	normalizer, err := NewTurkishSentenceNormalizer(stemWords, "")
	if err != nil {
		b.Fatalf("Failed to create normalizer: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = normalizer.getCandidates("test")
	}
}
