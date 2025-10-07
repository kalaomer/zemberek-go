//go:build demo
// +build demo

package timing

import (
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/kalaomer/zemberek-go/morphology"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var trLowerTiming = cases.Lower(language.Turkish)

func TestMorphologyStemmingTimingStress(t *testing.T) {
	text := `Yarın okula gideceğim. Tüketici mahkemesindeki dosyada yer alan kararın dayandığı deliller incelendiğinde temyiz edilen hükmün yasaya uygun olduğu sonucuna varıldı. Kârınızın ne kadar olduğunu nasıl anlarsınız? Hukuk dairesinin incelediği dosyada borçlunun reddi ve temyiz talebi ayrıntılı biçimde değerlendirildi. Gecikmiş ödeme, faizi ve cezai şartların nasıl uygulanacağı raporda açıkça belirtiliyor. Öğleden sonra yapılacak olan toplantıda yeni düzenlemeler konuşulacak. Bankanızın hesap bilgilerini öğrenmek istiyorum; ayrıca yarın havuza gireceğim ve akşama kadar yatacağım.`

	tokens := tokenizeTiming(text)
	if len(tokens) == 0 {
		t.Fatalf("tokenization produced zero tokens")
	}

	iterations := 500

	t.Log("=== Morphology Timing Stress ===")

	startInit := time.Now()
	morph := morphology.CreateWithDefaults()
	initDuration := time.Since(startInit)
	t.Logf("Initialization: %v", initDuration)

	startStem := time.Now()
	analyzed := 0
	unknown := 0
	for i := 0; i < iterations; i++ {
		for _, token := range tokens {
			result := morph.Analyze(token)
			if len(result.AnalysisResults) > 0 {
				analyzed++
			} else {
				unknown++
			}
		}
	}
	stemDuration := time.Since(startStem)

	totalProcessed := analyzed + unknown
	expectedTotal := len(tokens) * iterations
	if totalProcessed != expectedTotal {
		t.Logf("warning: processed count mismatch (expected %d, got %d)", expectedTotal, totalProcessed)
	}

	var perToken time.Duration
	if totalProcessed > 0 {
		perToken = stemDuration / time.Duration(totalProcessed)
	}

	t.Logf("Iterations: %d", iterations)
	t.Logf("Tokens per iteration: %d", len(tokens))
	t.Logf("Total tokens processed: %d (analyzed %d, unknown %d)", totalProcessed, analyzed, unknown)
	t.Logf("Stemming time: %v", stemDuration)
	t.Logf("Average per token: %v", perToken)
	t.Logf("Total time (init + stemming): %v", initDuration+stemDuration)
}

func tokenizeTiming(text string) []string {
	lower := trLowerTiming.String(text)
	raw := strings.Fields(lower)
	cleaned := make([]string, 0, len(raw))
	for _, token := range raw {
		trimmed := strings.TrimFunc(token, func(r rune) bool {
			if unicode.IsLetter(r) || r == '\'' || r == '’' {
				return false
			}
			return true
		})
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, trimmed)
	}
	return cleaned
}
