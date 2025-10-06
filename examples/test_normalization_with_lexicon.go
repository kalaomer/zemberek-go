package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kalaomer/zemberek-go/normalization"
)

func main() {
	fmt.Println("=== Turkish Sentence Normalization with Full Lexicon (94K+ words) ===\n")

	startTime := time.Now()
	fmt.Println("Loading lexicon...")

	// Create normalizer with full lexicon
	normalizer, err := normalization.NewTurkishSentenceNormalizerWithLexicon()
	if err != nil {
		fmt.Printf("Error creating normalizer: %v\n", err)
		return
	}

	loadTime := time.Since(startTime)
	fmt.Printf("✅ Lexicon loaded in %v\n", loadTime)
	fmt.Printf("   Total words: ~189,000 (with lowercase variants)\n\n")

	// Test sentences from Java
	examples := []string{
		"Yrn okua gidicem",
		"Tmm, yarin havuza giricem ve aksama kadar yaticam :)",
		"ah aynen ya annemde fark ettı siz evinizden cıkmayın diyo",
		"gercek mı bu? Yuh! Artık unutulması bile beklenmiyo",
		"Hayır hayat telaşm olmasa alacam buraları gökdelen dikicem.",
		"yok hocam kesınlıkle oyle birşey yok",
		"herseyi soyle hayatında olmaması gerek bence boyle ınsanların falan baskı yapıyosa",
	}

	fmt.Println("Java Test Sentences (with Full Lexicon):")
	fmt.Println(strings.Repeat("=", 75))

	for i, example := range examples {
		normalized := normalizer.Normalize(example)
		fmt.Printf("%d. Input:  %s\n", i+1, example)
		fmt.Printf("   Output: %s\n", normalized)
		fmt.Println()
	}

	// Additional test cases
	fmt.Println("\nAdditional Test Cases:")
	fmt.Println(strings.Repeat("=", 75))

	additionalTests := []string{
		"email adresim zemberek_python@loodos.com",
		"Kredi başvrusu yapmk istiyrum.",
		"Bankanizin hesp blgilerini ogrenmek istyorum.",
		"Bugun hava cok guzel.",
		"Iste bu olay gercek mi?",
		"Bilgisayar ve telefon almak istiyorum.",
		"Turkiyenin en buyuk sehri Istanbul.",
	}

	for i, test := range additionalTests {
		normalized := normalizer.Normalize(test)
		fmt.Printf("%d. Input:  %s\n", i+1, test)
		fmt.Printf("   Output: %s\n", normalized)
		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 75))
	fmt.Println("✅ All tests completed!")
}
