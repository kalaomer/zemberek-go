//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology/lexicon"
	pb "github.com/kalaomer/zemberek-go/morphology/lexicon/proto"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Test raw protobuf parsing
	items, err := lexicon.LoadBinaryLexicon()
	if err != nil {
		panic(err)
	}

	// Find kitap
	for i, item := range items {
		if item.Lemma == "kitap" {
			fmt.Printf("Item %d:\n", i)
			fmt.Printf("  Lemma: %s\n", item.Lemma)
			fmt.Printf("  Root: %s\n", item.Root)
			fmt.Printf("  PrimaryPos: %v\n", item.PrimaryPos)
			fmt.Printf("  SecondaryPos: %v\n", item.SecondaryPos)
			fmt.Printf("  Attributes: %v\n", item.Attributes)
			fmt.Println()
		}
		if i > 10 {
			break
		}
	}

	// Also check raw protobuf
	fmt.Println("=== Checking raw protobuf for 'kitap' ===")
	dictionary := &pb.Dictionary{}
	lexiconData := lexicon.GetLexiconBinData()
	if err := proto.Unmarshal(lexiconData, dictionary); err != nil {
		panic(err)
	}

	for i, pbItem := range dictionary.Items {
		if pbItem.Lemma == "kitap" {
			fmt.Printf("Proto item %d:\n", i)
			fmt.Printf("  Lemma: %s\n", pbItem.Lemma)
			fmt.Printf("  Root: %s\n", pbItem.Root)
			fmt.Printf("  PrimaryPos: %v\n", pbItem.PrimaryPos)
			fmt.Printf("  SecondaryPos: %v\n", pbItem.SecondaryPos)
			fmt.Printf("  RootAttributes: %v\n", pbItem.RootAttributes)
			fmt.Println()
		}
		if i > 100 {
			break
		}
	}
}
