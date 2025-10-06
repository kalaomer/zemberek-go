package morphology

import (
	"fmt"
	"testing"
)

func TestDebugDiminutive(t *testing.T) {
	morphology := CreateWithDefaults()

	word := "kutucuğ"

	// Check stem candidates
	stemTrans := morphology.Morphotactics.GetStemTransitions()
	candidates := stemTrans.GetPrefixMatches(word, false)

	fmt.Printf("Stem candidates for '%s': %d\n", word, len(candidates))
	for _, cand := range candidates {
		fmt.Printf("  - Surface: '%s', Item: %s, To: %s\n", cand.Surface, cand.Item.Root, cand.To.ID)
	}

	// Check if "kutu" has dim transitions
	kutuCandidates := stemTrans.GetPrefixMatches("kutu", false)
	if len(kutuCandidates) > 0 {
		fmt.Printf("\n'kutu' candidates: %d\n", len(kutuCandidates))
		kutuState := kutuCandidates[0].To
		fmt.Printf("State: %s, Outgoing transitions: %d\n", kutuState.ID, len(kutuState.Outgoing))

		for i, trans := range kutuState.Outgoing {
			if i < 10 { // Show first 10
				fmt.Printf("  [%d] %s -> %s\n", i, kutuState.ID, trans.GetState().ID)
			}
		}
	}
}
