package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()
	morphotactics := morph.GetMorphotactics()
	
	fmt.Println("=== Checking PnonS transitions ===")
	pnonS := morphotactics.PnonS
	
	if pnonS == nil {
		fmt.Println("PnonS is nil!")
		return
	}
	
	fmt.Printf("PnonS has %d outgoing transitions\n", len(pnonS.GetOutgoing()))
	
	for i, trans := range pnonS.GetOutgoing() {
		fmt.Printf("%d. %s -> %s (template: '%s')\n", 
			i+1, 
			trans.From.GetID(),
			trans.To.GetID(), 
			trans.SurfaceTemplate)
	}
}
