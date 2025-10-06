package morphology

import (
	"fmt"
	"testing"
)

func TestDimStateTransitions(t *testing.T) {
	morphology := CreateWithDefaults()

	// Get dim_S state
	dimS := morphology.Morphotactics.DimS

	fmt.Printf("dim_S outgoing transitions: %d\n", len(dimS.Outgoing))
	for i, trans := range dimS.Outgoing {
		fmt.Printf("  [%d] %s -> %s\n", i, dimS.ID, trans.GetState().ID)
	}

	// dim_S should connect to noun_S
	foundNoun := false
	for _, trans := range dimS.Outgoing {
		if trans.GetState().ID == "noun_S" {
			foundNoun = true
		}
	}

	if !foundNoun {
		t.Error("dim_S doesn't connect to noun_S")
	}
}
