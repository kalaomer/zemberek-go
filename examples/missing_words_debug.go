package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	fmt.Println("=== DEBUGGING MISSING WORD ANALYSIS ===\n")

	morph := morphology.CreateWithDefaults()

	// Words that Java finds but Go doesn't
	testWords := []struct {
		word         string
		javaRoot     string
		javaAnalysis string
	}{
		{"dairesi", "daire", "daire:Noun+A3sg+P3sg"},
		{"mahkemesi", "mahkeme", "mahkeme:Noun+A3sg+P3sg"},
		{"arasındaki", "ara", "ara:Noun+A3sg+P3sg+Loc+Rel+Adj"},
		{"yapılan", "yapılanmak", "yapılan:Verb+Imp+A2sg"},
		{"yargılaması", "yargılamak", "yargıla:Verb+Inf2+A3sg+P3sg"},
		{"ilamda", "ilâm", "ilam:Noun+A3sg+Loc"},
		{"nedenlerden", "neden", "neden:Noun+A3pl+Abl"},
		{"olarak", "olmak", "ol:Verb+ByDoingSo+Adv"},
		{"verilen", "vermek", "ver:Verb+Pass+PresPart+Adj"},
		{"avukatınca", "avukat", "avukat:Noun+A3sg+P2sg+Equ"},
		{"edilmesi", "etmek", "ed:Verb+Pass+Inf2+A3sg+P3sg"},
		{"konuşulup", "konuşmak", "konuş:Verb+Pass+AfterDoingSo+Adv"},
		{"dosyadaki", "dosya", "dosya:Noun+A3sg+Loc+Rel+Adj"},
		{"dayandığı", "dayanmak", "dayan:Verb+PastPart+Adj+P3sg"},
		{"delillerle", "delil", "delil:Noun+A3pl+Ins"},
		{"gerektirici", "gerekmek", "gerek:Verb+Caus+Agt+Adj"},
		{"isabetsizlik", "isabet", "isabet:Noun+A3sg+Without+Adj+Ness+A3sg"},
		{"bulunmamasına", "bulunmak", "bulun:Verb+Neg+Inf2+A3sg+P3sg+Dat"},
		{"olmayan", "olmak", "ol:Verb+Neg+PresPart+Adj"},
		{"itirazlarının", "itiraz", "itiraz:Noun+A3sg+P3pl+Gen"},
		{"olan", "olmak", "ol:Verb+PresPart+Adj"},
		{"aşağıda", "aşağı", "aşağı:Noun+A3sg+Loc"},
		{"kalan", "kalmak", "kal:Verb+PresPart+Adj"},
		{"gününde", "gün", "gün:Noun+A3sg+P2sg+Loc"},
	}

	fmt.Println("Checking if roots exist in lexicon:\n")

	foundInLexicon := 0
	notFoundInLexicon := 0
	hasAnalysis := 0
	noAnalysis := 0

	for i, test := range testWords {
		// Check if the root word exists in lexicon
		rootAnalysis := morph.Analyze(test.javaRoot)
		rootExists := len(rootAnalysis.AnalysisResults) > 0

		// Check if the inflected word can be analyzed
		wordAnalysis := morph.Analyze(test.word)
		wordExists := len(wordAnalysis.AnalysisResults) > 0

		status := "❌"
		if wordExists {
			status = "✅"
			hasAnalysis++
		} else {
			noAnalysis++
		}

		rootStatus := "❌"
		if rootExists {
			rootStatus = "✅"
			foundInLexicon++
		} else {
			notFoundInLexicon++
		}

		fmt.Printf("%2d. %-20s (root: %-15s) Root in lexicon: %s  Word analyzed: %s\n",
			i+1, test.word, test.javaRoot, rootStatus, status)

		// If word exists, show Go's analysis
		if wordExists {
			fmt.Printf("    Go:   %s\n", wordAnalysis.AnalysisResults[0].FormatString())
		}
		fmt.Printf("    Java: %s\n", test.javaAnalysis)

		// If root exists but word doesn't, this is a morphotactics issue
		if rootExists && !wordExists {
			fmt.Printf("    ⚠️  ROOT EXISTS but MORPHOTACTICS MISSING!\n")
		}

		fmt.Println()
	}

	fmt.Println("=== SUMMARY ===")
	fmt.Printf("Total test words: %d\n", len(testWords))
	fmt.Printf("Roots found in lexicon: %d/%d (%.1f%%)\n", foundInLexicon, len(testWords), float64(foundInLexicon)/float64(len(testWords))*100)
	fmt.Printf("Roots NOT in lexicon: %d/%d (%.1f%%)\n", notFoundInLexicon, len(testWords), float64(notFoundInLexicon)/float64(len(testWords))*100)
	fmt.Printf("Words analyzed by Go: %d/%d (%.1f%%)\n", hasAnalysis, len(testWords), float64(hasAnalysis)/float64(len(testWords))*100)
	fmt.Printf("Words NOT analyzed by Go: %d/%d (%.1f%%)\n", noAnalysis, len(testWords), float64(noAnalysis)/float64(len(testWords))*100)

	fmt.Println("\n=== DIAGNOSIS ===")
	morphotacticsIssues := 0
	for _, test := range testWords {
		rootAnalysis := morph.Analyze(test.javaRoot)
		wordAnalysis := morph.Analyze(test.word)
		if len(rootAnalysis.AnalysisResults) > 0 && len(wordAnalysis.AnalysisResults) == 0 {
			morphotacticsIssues++
		}
	}
	fmt.Printf("Words with morphotactics issues: %d\n", morphotacticsIssues)
	fmt.Printf("Words with lexicon issues: %d\n", notFoundInLexicon)

	if morphotacticsIssues > 0 {
		fmt.Println("\n⚠️  MAIN ISSUE: MORPHOTACTICS - Missing transitions for:")
		for _, test := range testWords {
			rootAnalysis := morph.Analyze(test.javaRoot)
			wordAnalysis := morph.Analyze(test.word)
			if len(rootAnalysis.AnalysisResults) > 0 && len(wordAnalysis.AnalysisResults) == 0 {
				fmt.Printf("  - %s (%s)\n", test.word, test.javaAnalysis)
			}
		}
	}
}
