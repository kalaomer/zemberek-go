//go:build demo
// +build demo

package main

import (
	"fmt"

	pb "github.com/kalaomer/zemberek-go/morphology/lexicon/proto"
	"google.golang.org/protobuf/proto"
	"os"
)

func main() {
	data, err := os.ReadFile("/Users/kalaomer/Projects/zemberek/zemberek-go/morphology/lexicon/data/lexicon.bin")
	if err != nil {
		panic(err)
	}

	dictionary := &pb.Dictionary{}
	if err := proto.Unmarshal(data, dictionary); err != nil {
		panic(err)
	}

	fmt.Printf("Total items: %d\n", len(dictionary.Items))

	// Print first 5 items
	for i := 0; i < 5 && i < len(dictionary.Items); i++ {
		item := dictionary.Items[i]
		fmt.Printf("\nItem %d:\n", i)
		fmt.Printf("  Lemma: '%s'\n", item.Lemma)
		fmt.Printf("  Root: '%s'\n", item.Root)
		fmt.Printf("  Pronunciation: '%s'\n", item.Pronunciation)
		fmt.Printf("  PrimaryPos: %v\n", item.PrimaryPos)
		fmt.Printf("  SecondaryPos: %v\n", item.SecondaryPos)
	}

	// Find kitap
	fmt.Println("\n=== Searching for kitap ===")
	found := 0
	for i, item := range dictionary.Items {
		if item.Lemma == "kitap" {
			fmt.Printf("\nFound at index %d:\n", i)
			fmt.Printf("  Lemma: '%s'\n", item.Lemma)
			fmt.Printf("  Root: '%s'\n", item.Root)
			fmt.Printf("  Pronunciation: '%s'\n", item.Pronunciation)
			fmt.Printf("  PrimaryPos: %v\n", item.PrimaryPos)
			fmt.Printf("  SecondaryPos: %v\n", item.SecondaryPos)
			found++
		}
	}
	fmt.Printf("\nTotal 'kitap' items found: %d\n", found)
}
