package morphology

import (
	"fmt"
	"testing"

	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

func TestTracePath(t *testing.T) {
	morphology := CreateWithDefaults()

	// Trace: kutu -> noun_S -> a3sg_S -> pnon_S -> nom_ST -> dim_S
	kutuStem := morphology.Morphotactics.GetStemTransitions().GetPrefixMatches("kutu", false)
	if len(kutuStem) < 2 {
		t.Fatal("Not enough stems for kutu")
	}

	// Find noun_S state
	var nounS *morphotactics.MorphemeState
	for _, stem := range kutuStem {
		fmt.Printf("Stem: %s -> %s\n", stem.Surface, stem.To.ID)
		if stem.To.ID == "noun_S" {
			nounS = stem.To
			break
		}
	}
	if nounS == nil {
		t.Fatal("noun_S not found")
	}

	fmt.Printf("1. noun_S: %s, outgoing: %d\n", nounS.ID, len(nounS.Outgoing))

	// Get a3sg_S
	if len(nounS.Outgoing) == 0 {
		t.Fatal("No outgoing from noun_S")
	}
	a3sgS := nounS.Outgoing[0].GetState()
	fmt.Printf("2. a3sg_S: %s, outgoing: %d\n", a3sgS.ID, len(a3sgS.Outgoing))

	// Get pnon_S
	if len(a3sgS.Outgoing) == 0 {
		t.Fatal("No outgoing from a3sg_S")
	}
	pnonS := a3sgS.Outgoing[0].GetState()
	fmt.Printf("3. pnon_S: %s, outgoing: %d\n", pnonS.ID, len(pnonS.Outgoing))

	// Check if nom_ST is in outgoing
	foundNomST := false
	var nomST *morphotactics.MorphemeState
	for _, trans := range pnonS.Outgoing {
		state := trans.GetState()
		fmt.Printf("   - %s -> %s\n", pnonS.ID, state.ID)
		if state.ID == "nom_ST" {
			foundNomST = true
			nomST = state
		}
	}

	if !foundNomST {
		t.Fatal("nom_ST not found in pnon_S outgoing")
	}

	fmt.Printf("4. nom_ST: %s, outgoing: %d, terminal: %v\n", nomST.ID, len(nomST.Outgoing), nomST.Terminal)

	// Check if dim_S is in nom_ST outgoing
	foundDim := false
	for _, trans := range nomST.Outgoing {
		state := trans.GetState()
		fmt.Printf("   - %s -> %s\n", nomST.ID, state.ID)
		if state.ID == "dim_S" {
			foundDim = true
		}
	}

	if !foundDim {
		t.Error("dim_S not found in nom_ST outgoing")
	}
}
