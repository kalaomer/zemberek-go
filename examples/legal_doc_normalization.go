//go:build demo
// +build demo

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kalaomer/zemberek-go/morphology"
	"github.com/kalaomer/zemberek-go/normalization"
)

func main() {
	// Legal document text
	text := `13. Hukuk Dairesi         2014/4087 E.  ,  2014/3970 K.

•

"İçtihat Metni"

MAHKEMESİ :Tüketici Mahkemesi

Taraflar arasındaki alacak davasının yapılan yargılaması sonunda ilamda yazılı nedenlerden dolayı davanın kabulüne yönelik olarak verilen hükmün süresi içinde davalı avukatınca temyiz edilmesi üzerine dosya incelendi gereği konuşulup düşünüldü.

K A R A R

Dosyadaki yazılara, kararın dayandığı delillerle yasaya uygun gerektirici nedenlere ve özellikle delillerin takdirinde bir isabetsizlik bulunmamasına göre yerinde olmayan bütün temyiz itirazlarının reddi ile usul ve yasaya uygun olan hükmün ONANMASINA, aşağıda dökümü yazılı 24,30 TL. kalan harcın temyiz edene iadesine, 17.02.2014 gününde oybirliğiyle karar verildi.`

	fmt.Println("=== LEGAL DOCUMENT NORMALIZATION ===")
	fmt.Println("\n--- Original Text ---")
	fmt.Println(text)
	fmt.Println("\n--- Normalization Process ---")

	// Word list for legal documents
	legalWords := []string{
		"hukuk", "dairesi", "mahkemesi", "tüketici", "taraf", "taraflar", "alacak",
		"dava", "davacı", "davalı", "yargılama", "ilam", "hüküm", "karar", "temyiz",
		"avukat", "dosya", "delil", "yasa", "yasaya", "usul", "itiraz", "red", "reddi",
		"onama", "onanması", "harç", "iade", "oybirliği", "oybirliğiyle", "içtihat",
		"içtihad", "metin", "metni", "neden", "nedenler", "nedenlerden", "dolayı",
		"üzerine", "göre", "yerinde", "bütün", "aşağıda", "döküm", "dökümü", "yazılı",
		"gün", "gününde", "verildi",
	}

	// Measure time for normalization initialization
	startInit := time.Now()
	normalizer, err := normalization.NewTurkishSentenceNormalizer(legalWords, "")
	initDuration := time.Since(startInit)

	if err != nil {
		fmt.Printf("⚠️  Normalizer initialization warning: %v\n", err)
		fmt.Println("Continuing with limited functionality...")
	}

	fmt.Printf("\n✓ Normalizer initialized in: %v\n", initDuration)

	// Measure time for normalization
	startNorm := time.Now()
	normalized := normalizer.Normalize(text)
	normDuration := time.Since(startNorm)

	fmt.Println("\n--- Normalized Text ---")
	fmt.Println(normalized)

	// Statistics
	fmt.Println("\n--- Statistics ---")
	fmt.Printf("Original length: %d chars\n", len(text))
	fmt.Printf("Normalized length: %d chars\n", len(normalized))
	fmt.Printf("Initialization time: %v\n", initDuration)
	fmt.Printf("Normalization time: %v\n", normDuration)
	fmt.Printf("Total time: %v\n", initDuration+normDuration)

	// Check if text changed
	if text == normalized {
		fmt.Println("\nℹ️  Text unchanged (already normalized)")
	} else {
		fmt.Println("\n✓ Text was normalized")
	}

	// === MORPHOLOGICAL ANALYSIS ===
	fmt.Println("\n\n=== MORPHOLOGICAL ANALYSIS (ROOT EXTRACTION) ===")

	// Initialize morphology analyzer
	startMorphInit := time.Now()
	morph := morphology.CreateWithDefaults()
	morphInitDuration := time.Since(startMorphInit)
	fmt.Printf("\n✓ Morphology analyzer initialized in: %v\n", morphInitDuration)

	// Analyze normalized text
	startAnalysis := time.Now()
	words := strings.Fields(normalized)
	analysisDuration := time.Since(startAnalysis)

	fmt.Println("\n--- Word Roots Analysis ---")
	fmt.Printf("Total words: %d\n\n", len(words))

	// Analyze each word and extract roots
	wordCount := 0
	for _, word := range words {
		// Skip punctuation and numbers
		if !isAlphabetic(word) {
			continue
		}

		wordCount++
		wa := morph.Analyze(word)

		if len(wa.AnalysisResults) > 0 {
			// Take first (best) analysis
			analysis := wa.AnalysisResults[0]

			// Extract root/stem
			root := analysis.Item.Lemma
			stem := analysis.GetStem()
			pos := analysis.Item.PrimaryPos.GetStringForm()

			// Get morphemes
			morphemes := make([]string, 0)
			for _, md := range analysis.MorphemeDataList {
				if md.Morpheme.ID != "Noun" && md.Morpheme.ID != "Verb" && md.Morpheme.ID != root {
					morphemes = append(morphemes, md.Morpheme.ID)
				}
			}

			fmt.Printf("%3d. %-20s → Root: %-15s Stem: %-15s POS: %-8s",
				wordCount, word, root, stem, pos)

			if len(morphemes) > 0 {
				fmt.Printf(" [%s]", strings.Join(morphemes, "+"))
			}
			fmt.Println()

		} else {
			fmt.Printf("%3d. %-20s → [No analysis - Unknown word]\n", wordCount, word)
		}
	}

	// Final statistics
	fmt.Println("\n--- Analysis Statistics ---")
	fmt.Printf("Morphology initialization: %v\n", morphInitDuration)
	fmt.Printf("Word analysis time: %v\n", analysisDuration)
	fmt.Printf("Total processing time: %v\n", initDuration+normDuration+morphInitDuration+analysisDuration)
}

// isAlphabetic checks if string contains only letters
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
