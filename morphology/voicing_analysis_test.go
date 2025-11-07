package morphology

import (
	"fmt"
	"testing"
)

func TestVoicingAnalysis(t *testing.T) {
	morph := CreateWithDefaults()

	// Test words with voicing
	testWords := []string{
		"kitap",     // base form (no voicing)
		"kitaplar",  // plural (no voicing)
		"kitapları", // plural + accusative (no voicing)
		"kitabı",    // accusative with voicing (p->b)
		"kitaba",    // dative with voicing (p->b)
	}

	fmt.Println("\n=== Voicing Analysis Test ===")
	fmt.Println("Word         | Analyses | Stem(s)                    | Item.Root | Item.Lemma")
	fmt.Println("-------------+----------+----------------------------+-----------+-----------")

	for _, word := range testWords {
		analysis := morph.Analyze(word)

		if len(analysis.AnalysisResults) == 0 {
			fmt.Printf("%-12s | %8d | %-26s | %9s | %s\n",
				word, 0, "NO ANALYSIS", "-", "-")
			continue
		}

		// Show first analysis
		firstAnalysis := analysis.AnalysisResults[0]

		// Get dictionary item info
		itemRoot := "-"
		itemLemma := "-"
		if firstAnalysis.Item != nil {
			itemRoot = firstAnalysis.Item.Root
			itemLemma = firstAnalysis.Item.Lemma
		}

		// Collect all unique stems
		stems := make(map[string]bool)
		for _, a := range analysis.AnalysisResults {
			stems[a.GetStem()] = true
		}

		// Format stems list
		stemsList := ""
		for s := range stems {
			if stemsList != "" {
				stemsList += ", "
			}
			stemsList += s
		}

		fmt.Printf("%-12s | %8d | %-26s | %9s | %s\n",
			word, len(analysis.AnalysisResults), stemsList, itemRoot, itemLemma)

		// Show detailed morpheme breakdown for first analysis
		if testing.Verbose() {
			fmt.Printf("  Morphemes: ")
			for i, md := range firstAnalysis.MorphemeDataList {
				if i > 0 {
					fmt.Print(" + ")
				}
				if md.Surface != "" {
					fmt.Printf("%s:%s", md.Surface, md.Morpheme.ID)
				} else {
					fmt.Printf(":%s", md.Morpheme.ID)
				}
			}
			fmt.Println()
		}
	}

	fmt.Println()

	// Key finding test
	t.Run("Voicing_Behavior", func(t *testing.T) {
		// Test that kitabı gets analyzed
		analysis := morph.Analyze("kitabı")

		if len(analysis.AnalysisResults) == 0 {
			t.Fatalf("Expected analysis for 'kitabı', got none")
		}

		stem := analysis.AnalysisResults[0].GetStem()
		item := analysis.AnalysisResults[0].Item

		t.Logf("Word: 'kitabı'")
		t.Logf("  Stem: '%s'", stem)
		t.Logf("  Item.Root: '%s'", item.Root)
		t.Logf("  Item.Lemma: '%s'", item.Lemma)

		// Document current behavior
		if stem == "kitap" {
			t.Log("✅ Stem is 'kitap' (voicing corrected)")
		} else if stem == "kitab" {
			t.Log("⚠️  Stem is 'kitab' (voicing preserved in surface)")
		} else {
			t.Logf("❓ Unexpected stem: '%s'", stem)
		}

		// Check if item has voicing attribute
		if item != nil {
			t.Logf("  Item has Voicing attribute: %v", item.HasAttribute(0x1)) // Voicing = 0x1
		}
	})
}

func TestDictionaryItemLookup(t *testing.T) {
	morph := CreateWithDefaults()

	// Check if "kitap" is in the dictionary
	analysis := morph.Analyze("kitap")

	if len(analysis.AnalysisResults) == 0 {
		t.Fatal("'kitap' not found in dictionary")
	}

	firstItem := analysis.AnalysisResults[0].Item
	t.Logf("Dictionary item for 'kitap':")
	t.Logf("  Lemma: %s", firstItem.Lemma)
	t.Logf("  Root: %s", firstItem.Root)
	t.Logf("  PrimaryPos: %v", firstItem.PrimaryPos)
	t.Logf("  SecondaryPos: %v", firstItem.SecondaryPos)

	// Check attributes
	t.Log("  Attributes:")
	for attr := range firstItem.Attributes {
		t.Logf("    - %v", attr)
	}
}
