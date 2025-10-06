package morphology

import (
	"fmt"
	"testing"

	"github.com/kalaomer/zemberek-go/morphology/analysis"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

func TestTokenParsing(t *testing.T) {
	templates := []string{">cI~k", ">cI!ğ"}

	for _, template := range templates {
		fmt.Printf("\nTemplate: '%s'\n", template)

		// Parse with morphotactics tokenizer
		morphTokens := morphotactics.TokenizeSuffixTemplate(template)
		fmt.Printf("  Morphotactics tokens: %d\n", len(morphTokens))
		for i, tok := range morphTokens {
			fmt.Printf("    [%d] Type=%v, Value='%c' (%d)\n", i, tok.Type, tok.Value, tok.Value)
		}

		// Parse with analysis tokenizer
		analysisTokenizer := analysis.NewSuffixTemplateTokenizer(template)
		analysisTokens := make([]*analysis.SuffixTemplateToken, 0)
		for analysisTokenizer.HasNext() {
			tok := analysisTokenizer.Next()
			analysisTokens = append(analysisTokens, tok)
		}
		fmt.Printf("  Analysis tokens: %d\n", len(analysisTokens))
		for i, tok := range analysisTokens {
			fmt.Printf("    [%d] Type=%v, Value='%c' (%d)\n", i, tok.Type, tok.Value, tok.Value)
		}

		// Check if same count
		if len(morphTokens) != len(analysisTokens) {
			t.Errorf("Token count mismatch for '%s': morph=%d, analysis=%d",
				template, len(morphTokens), len(analysisTokens))
		}
	}
}
