//go:build demo
// +build demo

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kalaomer/zemberek-go/morphology"
	"github.com/kalaomer/zemberek-go/normalization"
)

type normalizerIface interface{ Normalize(string) string }

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

	// Resolve normalization resources path for Go to align with Java
	dataRoot := resolveDataRoot()
	normRoot := filepath.Join(dataRoot, "normalization")
	fmt.Printf("Config (Go):\n  dataRoot          = %s\n  normalizationRoot = %s\n\n", dataRoot, normRoot)

	// Prefer Advanced normalizer (Java-parity: morphology + LM + lookups)
	var normalizer normalizerIface
	morph := morphology.CreateWithDefaults()
	if morph != nil {
		if adv, err := normalization.NewTurkishSentenceNormalizerAdvanced(morph, dataRoot); err == nil {
			normalizer = adv
		}
	}
	// Fallback to basic normalizer
	if normalizer == nil {
		if basic, err := normalization.NewTurkishSentenceNormalizer(extendedWords, normRoot); err == nil {
			normalizer = basic
		}
	}
	if normalizer == nil {
		fmt.Printf("Error creating normalizer: could not initialize advanced or basic normalizer\n")
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
		"Kârınızın ne kadar oldugunu nasil anlarsınız?",
	}

	// Expected (gold) outputs aligned with Java test
	expected := []string{
		"yarın okula gideceğim",
		"tamam, yarın havuza gireceğim ve akşama kadar yatacağım :)",
		"ah aynen ya annemde fark etti siz evinizden çıkmayın diyor",
		"gerçek mi bu? yuh! artık unutulması bile beklenmiyor",
		"hayır hayat telaşım olmasa alacağım buraları gökdelen dikeceğim.",
		"yok hocam kesinlikle öyle bir şey yok",
		"her şeyi söyle hayatında olmaması gerek bence böyle insanların falan baskı yapıyorsa",
		"email adresim zemberek_python@loodos.com",
		"kredi başvurusu yapmak istiyorum.",
		"bankanızın hesap bilgilerini öğrenmek istiyorum.",
		"kârınızın ne kadar olduğunu nasıl anlarsınız?",
	}

	fmt.Println("Sentence Normalization Comparison - Go")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	successCount := 0
	totalCount := len(examples)
	exactMatches := 0
	var totalSim float64

	for i, input := range examples {
		normalized := normalizer.Normalize(input)
		gold := expected[i]

		fmt.Printf("%2d. Input:    '%s'\n", i+1, input)
		fmt.Printf("    Expected: '%s'\n", gold)
		fmt.Printf("    Output:   '%s'\n", normalized)

		if normalized == gold {
			exactMatches++
		}

		sim := similarity(gold, normalized)
		totalSim += sim
		fmt.Printf("    Match: %s  Similarity: %.1f%%\n", tern(normalized == gold), sim*100)

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
	fmt.Printf("Exact Match: %d/%d (%.1f%%)  Avg Similarity: %.1f%%\n",
		exactMatches, totalCount, float64(exactMatches*100)/float64(totalCount),
		(totalSim*100)/float64(totalCount))
}

// Levenshtein-based similarity 0..1
func similarity(a, b string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	d := levenshtein(a, b)
	max := len(a)
	if len(b) > max {
		max = len(b)
	}
	if max == 0 {
		return 1.0
	}
	return 1.0 - float64(d)/float64(max)
}

func levenshtein(a, b string) int {
	n := len(a)
	m := len(b)
	dp := make([]int, m+1)
	for j := 0; j <= m; j++ {
		dp[j] = j
	}
	for i := 1; i <= n; i++ {
		prev := dp[0]
		dp[0] = i
		for j := 1; j <= m; j++ {
			temp := dp[j]
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			del := dp[j] + 1
			ins := dp[j-1] + 1
			sub := prev + cost
			dp[j] = min(del, min(ins, sub))
			prev = temp
		}
	}
	return dp[m]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func tern(ok bool) string {
	if ok {
		return "✓"
	}
	return "✗"
}

// resolveNormalizationRoot tries to find the same data root as Java side
// Priority: ZEMBEREK_DATA_ROOT/normalization -> ../resources/normalization -> ""
// resolveDataRoot tries ZEMBEREK_DATA_ROOT -> ../resources -> "resources"
func resolveDataRoot() string {
	// 1) ZEMBEREK_DATA_ROOT env
	if root := os.Getenv("ZEMBEREK_DATA_ROOT"); root != "" {
		if st, err := os.Stat(root); err == nil && st.IsDir() {
			return root
		}
	}
	// 2) ../resources (repo default)
	if p := filepath.Clean(filepath.Join("..", "resources")); existsDir(p) {
		return p
	}
	// 3) fallback
	return "resources"
}

func existsDir(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}
