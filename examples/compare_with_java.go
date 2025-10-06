package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
	"github.com/kalaomer/zemberek-go/normalization"
)

func main() {
	fmt.Println("=== Go Spell Checker Test (Comparable with Java) ===\n")

	// Create word list (similar to Java's lexicon)
	words := []string{
		"kitap", "kitabı", "kitaba", "kitaplar", "kitapçık",
		"yabancı", "yaban",
		"gidiyorum", "gidiyor", "gitti",
		"tam", "tüm", "yarın", "yan", "yön",
		"ve", "ile", "için", "bir", "bu",
	}

	// Build character graph
	graph := normalization.NewCharacterGraph()
	for _, word := range words {
		graph.AddWord(word, normalization.TypeWord)
	}
	decoder := normalization.NewCharacterGraphDecoder(graph)
	matcher := normalization.DiacriticsIgnoringMatcherInstance

	testWords := []string{"kitab", "yabanncı", "gidiyrum", "tmm", "yrn"}

	for _, word := range testWords {
		fmt.Printf("Word: %s\n", word)
		suggestions := decoder.GetSuggestions(word, matcher)
		if len(suggestions) > 0 {
			fmt.Printf("  Suggestions: %v\n", suggestions)
		} else {
			fmt.Println("  No suggestions found")
		}
	}

	fmt.Println("\n=== Morphology Analysis Test ===")
	morph := morphology.CreateWithDefaults()
	morphWords := []string{"kitap", "kalemi", "gidiyorum", "geldi"}
	for _, word := range morphWords {
		analysis := morph.Analyze(word)
		count := len(analysis.AnalysisResults)
		fmt.Printf("Word: %s -> %d analysis\n", word, count)
	}
}
