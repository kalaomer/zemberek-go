package morphology

import (
	"fmt"
	"testing"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/analysis"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

func TestTemplateGeneration(t *testing.T) {
	// Test template ">cI!ğ" (diminutive -cuğ)
	template := ">cI!ğ"

	// Parse tokens
	tokens := morphotactics.TokenizeSuffixTemplate(template)
	fmt.Printf("Template: %s\n", template)
	fmt.Printf("Tokens: %d\n", len(tokens))
	for i, token := range tokens {
		fmt.Printf("  [%d] Type: %v, Value: %c\n", i, token.Type, token.Value)
	}

	// Create a fake suffix transition to test surface generation
	dummyFrom := morphotactics.NewMorphemeStateNonTerminal("test_from", morphotactics.Noun)
	dummyTo := morphotactics.NewMorphemeStateNonTerminal("test_to", morphotactics.Dim)
	st := morphotactics.NewSuffixTransition(dummyFrom, dummyTo, template, nil)

	// Test with "kutu" phonetic attributes
	// 'kutu' → last vowel is 'u' (back, rounded)
	attrs := make(map[turkish.PhoneticAttribute]bool)
	attrs[turkish.LastLetterVowel] = true
	attrs[turkish.LastVowelBack] = true
	attrs[turkish.LastVowelRounded] = true

	surface := analysis.GenerateSurface(st, attrs)
	fmt.Printf("\nGenerated surface from 'kutu' + '%s': '%s'\n", template, surface)

	if surface != "cuğ" {
		t.Errorf("Expected 'cuğ', got '%s'", surface)
	}
}
