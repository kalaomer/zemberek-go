package lexicon

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

// RootLexicon represents the lexicon dictionary
type RootLexicon struct {
	IDMap   map[string]*DictionaryItem
	ItemSet map[*DictionaryItem]bool
	ItemMap map[string][]*DictionaryItem
}

// NewRootLexicon creates a new RootLexicon
func NewRootLexicon(itemList []*DictionaryItem) *RootLexicon {
	rl := &RootLexicon{
		IDMap:   make(map[string]*DictionaryItem),
		ItemSet: make(map[*DictionaryItem]bool),
		ItemMap: make(map[string][]*DictionaryItem),
	}

	for _, item := range itemList {
		rl.Add(item)
	}

	return rl
}

// Add adds an item to the lexicon
func (rl *RootLexicon) Add(item *DictionaryItem) {
	if rl.ItemSet[item] {
		fmt.Println("Warning: Duplicated item")
		return
	}

	if _, exists := rl.IDMap[item.ID]; exists {
		fmt.Printf("Warning: Duplicated item ID: %s\n", item.ID)
		return
	}

	rl.ItemSet[item] = true
	rl.IDMap[item.ID] = item

	if _, exists := rl.ItemMap[item.Lemma]; exists {
		rl.ItemMap[item.Lemma] = append(rl.ItemMap[item.Lemma], item)
	} else {
		rl.ItemMap[item.Lemma] = []*DictionaryItem{item}
	}
}

// GetItemByID gets an item by ID
func (rl *RootLexicon) GetItemByID(id string) *DictionaryItem {
	return rl.IDMap[id]
}

// GetItems gets items by lemma
func (rl *RootLexicon) GetItems(lemma string) []*DictionaryItem {
	return rl.ItemMap[lemma]
}

// Size returns the number of items
func (rl *RootLexicon) Size() int {
	return len(rl.ItemSet)
}

// GetAllItems returns all items in the lexicon
func (rl *RootLexicon) GetAllItems() []*DictionaryItem {
	items := make([]*DictionaryItem, 0, len(rl.ItemSet))
	for item := range rl.ItemSet {
		items = append(items, item)
	}
	return items
}

// LoadFromResources loads lexicon from a file
func LoadFromResources(resourcePath string) (*RootLexicon, error) {
	file, err := os.Open(resourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'
	reader.LazyQuotes = true

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	items := make([]*DictionaryItem, 0)
	lineMap := make(map[string][]string)

	// Build map for reference lookup
	for _, line := range lines {
		if len(line) > 0 {
			lineMap[line[0]] = line
		}
	}

	for _, line := range lines {
		if len(line) < 9 {
			continue
		}

		item := makeDictItemFromLine(line)

		// Handle reference item
		if line[7] != "null" && line[7] != "" {
			if refLine, ok := lineMap[line[7]]; ok {
				item.SetReferenceItem(makeDictItemFromLine(refLine))
			}
		}

		items = append(items, item)
	}

	return NewRootLexicon(items), nil
}

func makeDictItemFromLine(line []string) *DictionaryItem {
	itemLemma := line[1]
	itemRoot := line[2]
	itemPron := line[5]

	itemPPos := parsePrimaryPos(line[3])
	itemSPos := parseSecondaryPos(line[4])

	itemIndex := 0
	if len(line) > 6 && line[6] != "" {
		itemIndex, _ = strconv.Atoi(line[6])
	}

	var itemAttrs map[turkish.RootAttribute]bool
	if len(line) > 8 && line[8] != "0" && line[8] != "" {
		itemAttrs = parseAttributes(line[8])
	}

	return NewDictionaryItem(itemLemma, itemRoot, itemPPos, itemSPos, itemAttrs, itemPron, itemIndex)
}

func parsePrimaryPos(s string) turkish.PrimaryPos {
	switch s {
	case "Noun":
		return turkish.Noun
	case "Adj":
		return turkish.Adjective
	case "Adv":
		return turkish.Adverb
	case "Conj":
		return turkish.Conjunction
	case "Interj":
		return turkish.Interjection
	case "Verb":
		return turkish.Verb
	case "Pron":
		return turkish.Pronoun
	case "Num":
		return turkish.Numeral
	case "Det":
		return turkish.Determiner
	case "Postp":
		return turkish.PostPositive
	case "Ques":
		return turkish.Question
	case "Dup":
		return turkish.Duplicator
	case "Punc":
		return turkish.Punctuation
	default:
		return turkish.UnknownPos
	}
}

func parseSecondaryPos(s string) turkish.SecondaryPos {
	switch s {
	case "Prop":
		return turkish.ProperNoun
	case "Time":
		return turkish.Time
	case "Pers":
		return turkish.PersonalPron
	// Add more cases as needed
	default:
		return turkish.NonePos
	}
}

func parseAttributes(s string) map[turkish.RootAttribute]bool {
	attrs := make(map[turkish.RootAttribute]bool)
	parts := strings.Fields(s)

	for _, part := range parts {
		attr := parseRootAttribute(part)
		if attr != 0 {
			attrs[attr] = true
		}
	}

	return attrs
}

func parseRootAttribute(s string) turkish.RootAttribute {
	switch s {
	case "Voicing":
		return turkish.Voicing
	case "NoVoicing":
		return turkish.NoVoicing
	case "LastVowelDrop":
		return turkish.LastVowelDrop
	case "Doubling":
		return turkish.Doubling
	// Add more cases as needed
	default:
		return 0
	}
}

// GetDefault loads the default lexicon
func GetDefault(lexiconPath string) (*RootLexicon, error) {
	return LoadFromResources(lexiconPath)
}

// LoadDefaultWithMasterDictionaries loads lexicon from master dictionary files
func LoadDefaultWithMasterDictionaries() (*RootLexicon, error) {
	return LoadDefaultLexicon()
}
