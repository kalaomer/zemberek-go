package main

import (
	"fmt"
	"strings"

	"github.com/kalaomer/zemberek-go/normalization"
)

func main() {
	// Extended word list matching Java version
	extendedWords := []string{
		"yarın", "okula", "gideceğim", "tamam", "havuza", "gireceğim", "akşama", "kadar", "yatacağım",
		"anne", "annem", "annesi", "fark", "etti", "ettim", "siz", "sizin", "evinizden", "evimizden",
		"çıkmayın", "çıkmayalım", "diyor", "dedi", "gerçek", "artık", "unutulması", "beklenmiyor",
		"hayır", "hayat", "telaş", "telaşım", "olsa", "olmasaydı", "alacağım", "burayı", "burası", "buraları",
		"gökdelen", "dikeceğim", "yok", "hocam", "kesinlikle", "öyle", "birşey", "bir", "şey",
		"herşey", "herşeyi", "her", "şeyi", "söyle", "hayatında", "olmamalı", "olmak", "bence", "böyle",
		"insan", "insanlar", "insanların", "falan", "baskı", "yapıyorsa", "yapıyor", "email", "adres",
		"adresim", "zemberek", "kredi", "başvuru", "başvrusu", "yapmak", "istiyorum", "banka", "bankanızın",
		"hesap", "bilgi", "bilgiler", "bilgilerini", "öğrenmek", "istyorum",
		"kitap", "kalem", "defter", "masa", "sandalye",
	}

	normalizer, err := normalization.NewTurkishSentenceNormalizer(extendedWords, "")
	if err != nil {
		fmt.Printf("Error creating normalizer: %v\n", err)
		return
	}

	// Test examples
	examples := []string{
		"Yrn okua gidicem",
		"Tmm, yarin havuza giricem ve aksama kadar yaticam :)",
		"ah aynen ya annemde fark ettı siz evinizden cıkmayın diyo",
		"gercek mı bu? Yuh! Artık unutulması bile beklenmiyo",
		"Hayır hayat telaşm olmasa alacam buraları gökdelen dikicem.",
		"yok hocam kesınlıkle oyle birşey yok",
		"herseyi soyle hayatında olmaması gerek bence boyle ınsanların falan baskı yapıyosa",
		"email adresim zemberek_python@loodos.com",
		"Kredi başvrusu yapmk istiyrum.",
		"Bankanizin hesp blgilerini ogrenmek istyorum.",
	}

	fmt.Println("Sentence Normalization Comparison - Go")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	successCount := 0
	totalCount := len(examples)

	for i, input := range examples {
		normalized := normalizer.Normalize(input)

		fmt.Printf("%2d. Input:  '%s'\n", i+1, input)
		fmt.Printf("    Output: '%s'\n", normalized)

		if input != normalized {
			successCount++
		}

		if i < len(examples)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Normalization Stats: %d/%d sentences changed (%.1f%%)\n",
		successCount, totalCount, float64(successCount*100)/float64(totalCount))
}
