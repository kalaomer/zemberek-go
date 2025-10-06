package lexicon

import (
	_ "embed"
	
	"github.com/kalaomer/zemberek-go/core/turkish"
	pb "github.com/kalaomer/zemberek-go/morphology/lexicon/proto"
	"google.golang.org/protobuf/proto"
)

//go:embed data/lexicon.bin
var lexiconBinData []byte

// GetLexiconBinData returns the embedded binary data for debugging
func GetLexiconBinData() []byte {
	return lexiconBinData
}

// LoadBinaryLexicon loads the binary lexicon file (lexicon.bin)
func LoadBinaryLexicon() ([]*DictionaryItem, error) {
	// Parse protobuf
	dictionary := &pb.Dictionary{}
	if err := proto.Unmarshal(lexiconBinData, dictionary); err != nil {
		return nil, err
	}

	// Convert protobuf items to DictionaryItem
	items := make([]*DictionaryItem, 0, len(dictionary.Items))
	itemMap := make(map[string]*DictionaryItem) // For reference resolution

	for _, pbItem := range dictionary.Items {
		item := convertProtoToDictionaryItem(pbItem)
		items = append(items, item)
		itemMap[pbItem.Lemma] = item
	}

	// Resolve references
	for i, pbItem := range dictionary.Items {
		if pbItem.Reference != "" {
			if refItem, ok := itemMap[pbItem.Reference]; ok {
				items[i].SetReferenceItem(refItem)
			}
		}
	}

	return items, nil
}

// convertProtoToDictionaryItem converts protobuf DictionaryItem to Go DictionaryItem
func convertProtoToDictionaryItem(pbItem *pb.DictionaryItem) *DictionaryItem {
	// Convert PrimaryPos
	primaryPos := convertProtoPrimaryPos(pbItem.PrimaryPos)

	// Convert SecondaryPos
	secondaryPos := convertProtoSecondaryPos(pbItem.SecondaryPos)

	// Convert RootAttributes
	attributes := make(map[turkish.RootAttribute]bool)
	for _, pbAttr := range pbItem.RootAttributes {
		attr := convertProtoRootAttribute(pbAttr)
		if attr != 0 {
			attributes[attr] = true
		}
	}

	// If root is empty, use lemma as root (like Java does)
	root := pbItem.Root
	if root == "" {
		root = pbItem.Lemma
	}

	return NewDictionaryItem(
		pbItem.Lemma,
		root,
		primaryPos,
		secondaryPos,
		attributes,
		pbItem.Pronunciation,
		int(pbItem.Index),
	)
}

// convertProtoPrimaryPos converts protobuf PrimaryPos to turkish.PrimaryPos
func convertProtoPrimaryPos(pbPos pb.PrimaryPos) turkish.PrimaryPos {
	switch pbPos {
	case pb.PrimaryPos_Noun:
		return turkish.Noun
	case pb.PrimaryPos_Adjective:
		return turkish.Adjective
	case pb.PrimaryPos_Adverb:
		return turkish.Adverb
	case pb.PrimaryPos_Conjunction:
		return turkish.Conjunction
	case pb.PrimaryPos_Interjection:
		return turkish.Interjection
	case pb.PrimaryPos_Verb:
		return turkish.Verb
	case pb.PrimaryPos_Pronoun:
		return turkish.Pronoun
	case pb.PrimaryPos_Numeral:
		return turkish.Numeral
	case pb.PrimaryPos_Determiner:
		return turkish.Determiner
	case pb.PrimaryPos_PostPositive:
		return turkish.PostPositive
	case pb.PrimaryPos_Question:
		return turkish.Question
	case pb.PrimaryPos_Duplicator:
		return turkish.Duplicator
	case pb.PrimaryPos_Punctuation:
		return turkish.Punctuation
	default:
		return turkish.Noun // Default
	}
}

// convertProtoSecondaryPos converts protobuf SecondaryPos to turkish.SecondaryPos
func convertProtoSecondaryPos(pbPos pb.SecondaryPos) turkish.SecondaryPos {
	switch pbPos {
	case pb.SecondaryPos_ProperNoun:
		return turkish.ProperNoun
	case pb.SecondaryPos_Time:
		return turkish.Time
	case pb.SecondaryPos_Abbreviation:
		return turkish.Abbreviation
	default:
		return turkish.NonePos
	}
}

// convertProtoRootAttribute converts protobuf RootAttribute to turkish.RootAttribute
func convertProtoRootAttribute(pbAttr pb.RootAttribute) turkish.RootAttribute {
	switch pbAttr {
	case pb.RootAttribute_Voicing:
		return turkish.Voicing
	case pb.RootAttribute_NoVoicing:
		return turkish.NoVoicing
	case pb.RootAttribute_LastVowelDrop:
		return turkish.LastVowelDrop
	case pb.RootAttribute_InverseHarmony:
		return turkish.InverseHarmony
	case pb.RootAttribute_Doubling:
		return turkish.Doubling
	// Add more as needed
	default:
		return 0
	}
}
