package lexicon

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

// DictionaryItem represents a word and its properties in the lexicon dictionary
type DictionaryItem struct {
	Lemma         string
	Root          string
	PrimaryPos    turkish.PrimaryPos
	SecondaryPos  turkish.SecondaryPos
	Attributes    map[turkish.RootAttribute]bool
	Pronunciation string
	Index         int
	ReferenceItem *DictionaryItem
	ID            string
}

// Unknown dictionary item
var Unknown *DictionaryItem

func init() {
	Unknown = &DictionaryItem{
		Lemma:        "UNK",
		Root:         "UNK",
		Pronunciation: "UNK",
		PrimaryPos:   turkish.UnknownPos,
		SecondaryPos: turkish.UnknownSec,
		ID:           "UNK_Unk",
	}
}

// NewDictionaryItem creates a new DictionaryItem
func NewDictionaryItem(lemma, root string, primaryPos turkish.PrimaryPos, secondaryPos turkish.SecondaryPos,
	attributes map[turkish.RootAttribute]bool, pronunciation string, index int) *DictionaryItem {

	if pronunciation == "" {
		pronunciation = root
	}

	if attributes == nil {
		attributes = make(map[turkish.RootAttribute]bool)
	}

	id := GenerateID(lemma, primaryPos, secondaryPos, index)

	return &DictionaryItem{
		Lemma:         lemma,
		Root:          root,
		PrimaryPos:    primaryPos,
		SecondaryPos:  secondaryPos,
		Attributes:    attributes,
		Pronunciation: pronunciation,
		Index:         index,
		ID:            id,
	}
}

// GenerateID generates an ID for a word
func GenerateID(lemma string, pos turkish.PrimaryPos, spos turkish.SecondaryPos, index int) string {
	itemID := fmt.Sprintf("%s_%s", lemma, pos.GetStringForm())

	if spos != turkish.NonePos {
		itemID = fmt.Sprintf("%s_%s", itemID, spos.GetStringForm())
	}
	if index > 0 {
		itemID = fmt.Sprintf("%s_%d", itemID, index)
	}
	return itemID
}

// HasAttribute checks if the item has a specific attribute
func (d *DictionaryItem) HasAttribute(attribute turkish.RootAttribute) bool {
	return d.Attributes[attribute]
}

// SetReferenceItem sets a reference item
func (d *DictionaryItem) SetReferenceItem(referenceItem *DictionaryItem) {
	d.ReferenceItem = referenceItem
}

// String returns string representation
func (d *DictionaryItem) String() string {
	str := d.Lemma + " [P:" + d.PrimaryPos.GetStringForm()

	if d.SecondaryPos != turkish.NonePos {
		str += ", " + d.SecondaryPos.GetStringForm()
	}

	if len(d.Attributes) == 0 {
		str += "]"
	} else {
		str = printAttributes(str, d.Attributes)
	}

	return str
}

// NormalizedLemma returns normalized lemma
func (d *DictionaryItem) NormalizedLemma() string {
	if d.PrimaryPos == turkish.Verb && len(d.Lemma) >= 3 {
		return d.Lemma[:len(d.Lemma)-3]
	}
	return d.Lemma
}

func printAttributes(str string, attrs map[turkish.RootAttribute]bool) string {
	if len(attrs) > 0 {
		str += "; A:"
		i := 0
		for attr := range attrs {
			str += attr.GetStringForm()
			if i < len(attrs)-1 {
				str += ", "
			}
			i++
		}
		str += "]"
	}
	return str
}

// IsUnknown checks if this is the unknown item
func (d *DictionaryItem) IsUnknown() bool {
	return d == Unknown
}
