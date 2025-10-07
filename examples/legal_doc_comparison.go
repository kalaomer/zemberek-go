//go:build demo
// +build demo

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	// Legal document text (same for Java and Go comparison)
	text := `13. Hukuk Dairesi 2014/4087 E., 2014/3970 K.

İçtihat Metni

MAHKEMESİ: Tüketici Mahkemesi

Taraflar arasındaki alacak davasının yapılan yargılaması sonunda ilamda yazılı nedenlerden dolayı davanın kabulüne yönelik olarak verilen hükmün süresi içinde davalı avukatınca temyiz edilmesi üzerine dosya incelendi gereği konuşulup düşünüldü.

KARAR

Dosyadaki yazılara, kararın dayandığı delillerle yasaya uygun gerektirici nedenlere ve özellikle delillerin takdirinde bir isabetsizlik bulunmamasına göre yerinde olmayan bütün temyiz itirazlarının reddi ile usul ve yasaya uygun olan hükmün ONANMASINA, aşağıda dökümü yazılı 24,30 TL. kalan harcın temyiz edene iadesine, 17.02.2014 gününde oybirliğiyle karar verildi.`

	fmt.Println("=== ZEMBEREK GO - LEGAL DOCUMENT ANALYSIS ===")
	fmt.Println()

	// Initialize morphology analyzer
	startInit := time.Now()
	morph := morphology.CreateWithDefaults()
	initDuration := time.Since(startInit)

	fmt.Printf("Initialization time: %v\n\n", initDuration)

	// Normalize text (simple lowercase for comparison)
	normalized := strings.ToLower(text)

	// Analyze words
	startAnalysis := time.Now()
	words := strings.Fields(normalized)

	analyzedCount := 0
	unknownCount := 0

	fmt.Println("--- Morphological Analysis Results ---\n")

	wordNum := 0
	for _, word := range words {
		// Skip non-alphabetic
		if !isAlphabetic(word) {
			continue
		}

		wordNum++
		wa := morph.Analyze(word)

		if len(wa.AnalysisResults) > 0 {
			analysis := wa.AnalysisResults[0]
			root := analysis.Item.Lemma
			pos := analysis.Item.PrimaryPos.GetStringForm()

			// Format morphemes
			morphStr := analysis.FormatString()

			fmt.Printf("%3d. %-20s → %-15s [%s] %s\n",
				wordNum, word, root, pos, morphStr)
			analyzedCount++
		} else {
			fmt.Printf("%3d. %-20s → NO ANALYSIS\n", wordNum, word)
			unknownCount++
		}
	}

	analysisDuration := time.Since(startAnalysis)

	// Statistics
	fmt.Println("\n--- Statistics ---")
	fmt.Printf("Total words: %d\n", wordNum)
	fmt.Printf("Analyzed: %d (%.1f%%)\n", analyzedCount, float64(analyzedCount)/float64(wordNum)*100)
	fmt.Printf("Unknown: %d (%.1f%%)\n", unknownCount, float64(unknownCount)/float64(wordNum)*100)
	fmt.Printf("Analysis time: %v\n", analysisDuration)
	fmt.Printf("Total time: %v\n", initDuration+analysisDuration)
}

func isAlphabetic(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			r == 'ç' || r == 'ğ' || r == 'ı' || r == 'İ' || r == 'ö' || r == 'ş' || r == 'ü' ||
			r == 'Ç' || r == 'Ğ' || r == 'Ö' || r == 'Ş' || r == 'Ü') {
			return false
		}
	}
	return true
}
