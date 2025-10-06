package normalization

import (
	"testing"

	"github.com/kalaomer/zemberek-go/morphology"
)

// TestNewStemEndingGraphFromMorphology tests creating StemEndingGraph from morphology
func TestNewStemEndingGraphFromMorphology(t *testing.T) {
	// Create morphology with default lexicon (94K+ words)
	morph := morphology.CreateWithDefaults()

	// Debug: Check lexicon
	if morph.Lexicon == nil {
		t.Fatal("Morphology lexicon is nil")
	}
	lexiconSize := morph.Lexicon.Size()
	t.Logf("Lexicon size: %d", lexiconSize)

	if lexiconSize == 0 {
		t.Skip("Skipping test: Lexicon is empty (probably running from test directory without resources)")
	}

	// Create StemEndingGraph from morphology (like Java does)
	stemGraph, err := NewStemEndingGraphFromMorphology(morph, "")
	if err != nil {
		t.Fatalf("Failed to create StemEndingGraph: %v", err)
	}

	// Verify stem graph was created
	if stemGraph == nil {
		t.Fatal("StemEndingGraph is nil")
	}

	if stemGraph.StemGraph == nil {
		t.Fatal("StemGraph is nil")
	}

	// Check that stems were extracted
	// We should have many stems from the lexicon
	allNodes := stemGraph.StemGraph.GetAllNodes()
	if len(allNodes) == 0 {
		t.Error("No stem nodes found in graph")
	}

	t.Logf("✅ Created StemEndingGraph with %d stem nodes", len(allNodes))
}

// TestStemEndingGraph_ContainsCommonWords tests that common Turkish words are in stem graph
func TestStemEndingGraph_ContainsCommonWords(t *testing.T) {
	morph := morphology.CreateWithDefaults()
	stemGraph, err := NewStemEndingGraphFromMorphology(morph, "")
	if err != nil {
		t.Fatalf("Failed to create StemEndingGraph: %v", err)
	}

	// Common Turkish words that should be in the lexicon
	commonWords := []string{
		"kitap",
		"kalem",
		"ev",
		"okul",
		"gitmek",
		"yapmak",
		"güzel",
		"büyük",
		"merhaba",
		"yarın",
	}

	for _, word := range commonWords {
		if !stemGraph.StemGraph.ContainsWord(word) {
			t.Errorf("Common word '%s' not found in stem graph", word)
		} else {
			t.Logf("✅ Found '%s' in stem graph", word)
		}
	}
}

// TestStemEndingGraph_LexiconSize tests that we have a substantial lexicon
func TestStemEndingGraph_LexiconSize(t *testing.T) {
	morph := morphology.CreateWithDefaults()

	if morph.Lexicon == nil {
		t.Fatal("Lexicon is nil")
	}

	itemCount := morph.Lexicon.Size()
	if itemCount == 0 {
		t.Fatal("Lexicon is empty")
	}

	t.Logf("✅ Lexicon contains %d items", itemCount)

	// We should have at least 50,000 items (Java has ~29K stems but we have more)
	if itemCount < 50000 {
		t.Errorf("Expected at least 50,000 lexicon items, got %d", itemCount)
	}
}

// TestStemEndingGraph_ManualStems tests creating graph with manual stem list
func TestStemEndingGraph_ManualStems(t *testing.T) {
	// Manual stem list (old way)
	stems := []string{
		"kitap",
		"kalem",
		"defter",
		"masa",
		"sandalye",
	}

	stemGraph, err := NewStemEndingGraph(stems, "")
	if err != nil {
		t.Fatalf("Failed to create StemEndingGraph with manual stems: %v", err)
	}

	// Verify all manual stems are in graph
	for _, stem := range stems {
		if !stemGraph.StemGraph.ContainsWord(stem) {
			t.Errorf("Manual stem '%s' not found in graph", stem)
		}
	}

	t.Logf("✅ All %d manual stems added to graph", len(stems))
}

// TestStemEndingGraph_CompareWithManualAndMorphology compares manual vs morphology approach
func TestStemEndingGraph_CompareWithManualAndMorphology(t *testing.T) {
	// Manual approach (limited)
	manualStems := []string{"kitap", "kalem", "ev"}
	manualGraph, _ := NewStemEndingGraph(manualStems, "")
	manualCount := len(manualGraph.StemGraph.GetAllNodes())

	// Morphology approach (comprehensive)
	morph := morphology.CreateWithDefaults()
	morphGraph, _ := NewStemEndingGraphFromMorphology(morph, "")
	morphCount := len(morphGraph.StemGraph.GetAllNodes())

	t.Logf("Manual approach: %d nodes", manualCount)
	t.Logf("Morphology approach: %d nodes", morphCount)

	// Morphology should have MUCH more stems
	if morphCount <= manualCount {
		t.Errorf("Expected morphology approach to have more stems than manual. Manual: %d, Morphology: %d",
			manualCount, morphCount)
	}

	// We should have at least 50K stems from morphology
	if morphCount < 50000 {
		t.Errorf("Expected at least 50,000 stems from morphology, got %d", morphCount)
	}

	t.Logf("✅ Morphology approach has %dx more stems than manual", morphCount/manualCount)
}

// BenchmarkStemEndingGraphCreation benchmarks creating graph from morphology
func BenchmarkStemEndingGraphCreation(b *testing.B) {
	morph := morphology.CreateWithDefaults()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewStemEndingGraphFromMorphology(morph, "")
	}
}

// BenchmarkStemLookup benchmarks looking up stems in the graph
func BenchmarkStemLookup(b *testing.B) {
	morph := morphology.CreateWithDefaults()
	stemGraph, _ := NewStemEndingGraphFromMorphology(morph, "")

	testWords := []string{"kitap", "kalem", "ev", "okul", "gitmek"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, word := range testWords {
			stemGraph.StemGraph.ContainsWord(word)
		}
	}
}
