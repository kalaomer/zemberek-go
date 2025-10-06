package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology"
	"github.com/kalaomer/zemberek-go/normalization"
	"github.com/kalaomer/zemberek-go/normalization/deasciifier"
	"github.com/kalaomer/zemberek-go/tokenization"
)

func main() {
	fmt.Println("=== Zemberek-Go Examples ===\n")

	// Example 1: Turkish Alphabet Operations
	fmt.Println("1. Turkish Alphabet Operations:")
	alphabet := turkish.Instance

	fmt.Printf("   Is 'ı' a vowel? %v\n", alphabet.IsVowel('ı'))
	fmt.Printf("   Is 'ş' a vowel? %v\n", alphabet.IsVowel('ş'))
	fmt.Printf("   Last letter of 'kitap': %v\n", alphabet.GetLastLetter("kitap").CharValue)
	fmt.Printf("   Last vowel of 'kitap': %v\n", alphabet.GetLastVowel("kitap").CharValue)
	fmt.Println()

	// Example 2: Text Normalization
	fmt.Println("2. Text Normalization:")
	normalized := alphabet.Normalize("Merhaba! Bu bir test.")
	fmt.Printf("   Original: 'Merhaba! Bu bir test.'\n")
	fmt.Printf("   Normalized: '%s'\n", normalized)
	fmt.Println()

	// Example 3: Case Operations
	fmt.Println("3. Turkish Capitalization:")
	word1 := "istanbul"
	word2 := "ışık"
	fmt.Printf("   Capitalize('%s'): %s\n", word1, turkish.Capitalize(word1))
	fmt.Printf("   Capitalize('%s'): %s\n", word2, turkish.Capitalize(word2))
	fmt.Println()

	// Example 4: Sentence Extraction
	fmt.Println("4. Sentence Extraction:")
	paragraph := "Merhaba dünya! Bu bir test cümlesidir. Nasılsınız?"

	extractor, err := tokenization.NewTurkishSentenceExtractor(false, "")
	if err == nil {
		sentences := extractor.FromParagraph(paragraph)
		fmt.Printf("   Paragraph: '%s'\n", paragraph)
		fmt.Printf("   Extracted %d sentences:\n", len(sentences))
		for i, sentence := range sentences {
			fmt.Printf("     %d: '%s'\n", i+1, sentence)
		}
	} else {
		fmt.Printf("   Error: %v\n", err)
	}
	fmt.Println()

	// Example 5: Token Types
	fmt.Println("5. Token Types:")
	token := tokenization.NewToken("Merhaba", tokenization.Word, 0, 7)
	fmt.Printf("   Token: %s\n", token)
	fmt.Printf("   Type: %d (Word)\n", token.Type)
	fmt.Println()

	// Example 6: Span Operations
	fmt.Println("6. Span Operations:")
	text := "Merhaba dünya"
	span, _ := tokenization.NewSpan(0, 7)
	fmt.Printf("   Text: '%s'\n", text)
	fmt.Printf("   Span [%d:%d]: '%s'\n", span.Start, span.End, span.GetSubString(text))
	fmt.Printf("   Span length: %d\n", span.GetLength())
	fmt.Println()

	// Example 7: Turkish Letter Properties
	fmt.Println("7. Turkish Letter Properties:")
	letter := alphabet.GetLetter('ğ')
	fmt.Printf("   Letter: %c\n", letter.CharValue)
	fmt.Printf("   Is vowel: %v\n", letter.IsVowel())
	fmt.Printf("   Is consonant: %v\n", letter.IsConsonant())
	fmt.Printf("   Is frontal: %v\n", letter.IsFrontal())
	fmt.Println()

	// Example 8: Voicing/Devoicing
	fmt.Println("8. Voicing/Devoicing:")
	fmt.Printf("   Voice('k'): %c\n", alphabet.Voice('k'))
	fmt.Printf("   Voice('p'): %c\n", alphabet.Voice('p'))
	fmt.Printf("   Devoice('b'): %c\n", alphabet.Devoice('b'))
	fmt.Printf("   Devoice('d'): %c\n", alphabet.Devoice('d'))
	fmt.Println()

	// Example 9: Deasciifier (ASCII Turkish -> Proper Turkish)
	fmt.Println("9. Deasciifier (ASCII to Turkish):")
	asciiText := "Merhaba dunya! Bugun hava cok guzel."
	d := deasciifier.NewDeasciifier(asciiText)
	turkishText := d.ConvertToTurkish()
	fmt.Printf("   ASCII Input:  '%s'\n", asciiText)
	fmt.Printf("   Turkish Output: '%s'\n", turkishText)
	fmt.Println("   Note: Without pattern table, only toggles existing Turkish chars")
	fmt.Println()

	// Example 10: Character Graph
	fmt.Println("10. Character Graph (for spell checking):")
	graph := normalization.NewCharacterGraph()
	graph.AddWord("kitap", normalization.TypeWord)
	graph.AddWord("kalem", normalization.TypeWord)
	graph.AddWord("defter", normalization.TypeWord)
	fmt.Printf("   Added words: kitap, kalem, defter\n")
	fmt.Printf("   Contains 'kitap': %v\n", graph.ContainsWord("kitap"))
	fmt.Printf("   Contains 'test': %v\n", graph.ContainsWord("test"))
	fmt.Printf("   Total words in graph: %d\n", len(graph.GetAllNodes()))
	fmt.Println()

	// Example 11: Spell Checker with Suggestions
	fmt.Println("11. Spell Checking and Suggestions:")
	decoder := normalization.NewCharacterGraphDecoder(graph)
	matcher := normalization.DiacriticsIgnoringMatcherInstance

	// Test word with typo
	misspelled := "kitab"  // Should suggest "kitap"
	suggestions := decoder.GetSuggestions(misspelled, matcher)
	fmt.Printf("   Misspelled: '%s'\n", misspelled)
	fmt.Printf("   Suggestions: %v\n", suggestions)
	fmt.Println()

	// Example 12: Edit Distance (Levenshtein)
	fmt.Println("12. Edit Distance Calculation:")
	word1_edit := "kitap"
	word2_edit := "kitab"
	// Note: levenshteinDistance is not exported, showing concept
	fmt.Printf("   Word 1: '%s'\n", word1_edit)
	fmt.Printf("   Word 2: '%s'\n", word2_edit)
	fmt.Println("   Edit distance calculation available in spell checker")
	fmt.Println()

	// Example 13: Sentence Normalization (Comprehensive)
	fmt.Println("13. Sentence Normalization (Comprehensive):")

	// Extended word list for better normalization
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
	if err == nil {
		// Test examples from Python version
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

		fmt.Println("   Normalizing informal Turkish sentences:")
		fmt.Println("   ────────────────────────────────────────────────────────────")

		for i, example := range examples {
			normalized := normalizer.Normalize(example)
			fmt.Printf("   %d. Input:  '%s'\n", i+1, example)
			fmt.Printf("      Output: '%s'\n", normalized)
			if i < len(examples)-1 {
				fmt.Println()
			}
		}

		fmt.Println("   ────────────────────────────────────────────────────────────")
		fmt.Println("   Note: Normalization quality depends on vocabulary coverage")
	} else {
		fmt.Printf("   Normalizer creation note: %v\n", err)
	}
	fmt.Println()

	// Example 14: Candidate Generation
	fmt.Println("14. Normalization Candidates:")
	candidate1 := normalization.NewCandidate("merhaba")
	candidate2 := normalization.NewCandidate("selam")
	candidates := normalization.NewCandidates("mrb", []*normalization.Candidate{candidate1, candidate2})
	fmt.Printf("   Original word: '%s'\n", candidates.Word)
	fmt.Printf("   Candidate 1: '%s' (score: %.2f)\n", candidates.Candidates[0].Content, candidates.Candidates[0].Score)
	fmt.Printf("   Candidate 2: '%s' (score: %.2f)\n", candidates.Candidates[1].Content, candidates.Candidates[1].Score)
	fmt.Println()

	// Example 15: Morphological Analysis (Single Word)
	fmt.Println("15. Morphological Analysis (Single Word):")
	morph := morphology.CreateWithDefaults()

	word := "kalemin"
	analysis := morph.Analyze(word)

	fmt.Printf("   Word: '%s'\n", word)
	fmt.Printf("   Analysis results: %d\n", len(analysis.AnalysisResults))

	if len(analysis.AnalysisResults) > 0 {
		fmt.Println("   Analyses:")
		for i, result := range analysis.AnalysisResults {
			fmt.Printf("     %d. %s\n", i+1, result.FormatString())
		}
	} else {
		fmt.Println("   No analysis found (simplified morphology)")
	}
	fmt.Println("   Note: Full morphology requires lexicon resources")
	fmt.Println()

	// Example 16: Sentence Analysis
	fmt.Println("16. Sentence Analysis:")
	sentence := "Yarın kar yağacak"
	sentenceAnalysis := morph.AnalyzeSentence(sentence)

	fmt.Printf("   Sentence: '%s'\n", sentence)
	fmt.Printf("   Word count: %d\n", len(sentenceAnalysis))
	fmt.Println("   Word analyses:")

	for _, wordAnalysis := range sentenceAnalysis {
		fmt.Printf("     Word: '%s' -> %d analysis(es)\n",
			wordAnalysis.Input, len(wordAnalysis.AnalysisResults))

		// Show first analysis if available
		if len(wordAnalysis.AnalysisResults) > 0 {
			fmt.Printf("       Best: %s\n", wordAnalysis.AnalysisResults[0].FormatString())
		}
	}
	fmt.Println()

	// Example 17: Detailed Tokenization
	fmt.Println("17. Detailed Tokenization:")
	tokenText := "Saat 12:00'de buluşalım."

	fmt.Printf("   Text: '%s'\n", tokenText)

	// Simple word tokenization (by spaces for demo)
	words := tokenization.SimpleTokenize(tokenText)
	fmt.Printf("   Tokens found: %d\n", len(words))
	fmt.Println("   Token details:")

	for i, word := range words {
		// Determine token type
		tokenType := tokenization.Word
		if len(word) > 0 {
			if word[0] >= '0' && word[0] <= '9' {
				tokenType = tokenization.Number
			}
			if word == "." || word == "," || word == "!" || word == "?" {
				tokenType = tokenization.Punctuation
			}
		}

		fmt.Printf("     %d. Content='%s', Type=%s\n",
			i+1, word, tokenization.TokenTypeName(tokenType))
	}
	fmt.Println()

	// Example 18: Morphology-based Analysis
	fmt.Println("18. Morphology-based Analysis:")
	testWords := []string{"kitap", "kitaplar", "gidiyorum", "geldi"}

	fmt.Println("   Word analysis results:")
	for _, w := range testWords {
		wa := morph.Analyze(w)
		hasAnalysis := morph.HasAnalysis(w)

		fmt.Printf("     '%s': ", w)
		if hasAnalysis {
			fmt.Printf("✓ Valid (%d analysis)\n", len(wa.AnalysisResults))
		} else {
			fmt.Println("✗ No valid analysis")
		}
	}
	fmt.Println()

	fmt.Println("=== Examples Complete ===")
	fmt.Println("\nImplemented Features:")
	fmt.Println("  ✓ Turkish Alphabet & Letter Operations")
	fmt.Println("  ✓ Text Normalization & Case Conversion")
	fmt.Println("  ✓ Sentence Extraction & Tokenization")
	fmt.Println("  ✓ Deasciifier (ASCII → Turkish)")
	fmt.Println("  ✓ Character Graph & Spell Checking")
	fmt.Println("  ✓ Sentence Normalization (informal → formal)")
	fmt.Println("  ✓ Morphological Analysis (basic)")
	fmt.Println("  ✓ Sentence Analysis")
	fmt.Println("  ✓ Edit Distance & Beam Search")
	fmt.Println("\nNote: Full morphology features require lexicon resources.")
}
