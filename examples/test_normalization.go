//go:build demo
// +build demo

package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/normalization"
)

func main() {
	fmt.Println("=== Enhanced Turkish Sentence Normalization Test ===\n")

	// Create enhanced normalizer
	normalizer, err := normalization.NewTurkishSentenceNormalizerEnhanced()
	if err != nil {
		fmt.Printf("Error creating normalizer: %v\n", err)
		return
	}

	// Test sentences from Java NormalizeNoisyText.java
	examples := []string{
		"Yrn okua gidicem",
		"Tmm, yarin havuza giricem ve aksama kadar yaticam :)",
		"ah aynen ya annemde fark ettı siz evinizden cıkmayın diyo",
		"gercek mı bu? Yuh! Artık unutulması bile beklenmiyo",
		"Hayır hayat telaşm olmasa alacam buraları gökdelen dikicem.",
		"yok hocam kesınlıkle oyle birşey yok",
		"herseyi soyle hayatında olmaması gerek bence boyle ınsanların falan baskı yapıyosa",
	}

	fmt.Println("Java Test Sentences:")
	fmt.Println("=" + string(make([]byte, 70)) + "=")

	for i, example := range examples {
		normalized := normalizer.Normalize(example)
		fmt.Printf("%d. Input:  %s\n", i+1, example)
		fmt.Printf("   Output: %s\n", normalized)
		fmt.Println()
	}

	// Additional test cases
	fmt.Println("\nAdditional Test Cases:")
	fmt.Println("=" + string(make([]byte, 70)) + "=")

	additionalTests := []string{
		"email adresim zemberek_python@loodos.com",
		"Kredi başvrusu yapmk istiyrum.",
		"Bankanizin hesp blgilerini ogrenmek istyorum.",
		"Bugun hava cok guzel.",
		"Iste bu olay gercek mi?",
	}

	for i, test := range additionalTests {
		normalized := normalizer.Normalize(test)
		fmt.Printf("%d. Input:  %s\n", i+1, test)
		fmt.Printf("   Output: %s\n", normalized)
		fmt.Println()
	}

	fmt.Println("=" + string(make([]byte, 70)) + "=")
	fmt.Println("Test completed!")
}
