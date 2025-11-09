package lexicon

import (
	"bufio"
	_ "embed"
	"os"
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

//go:embed data/master-dictionary.dict
var masterDictData string

//go:embed data/non-tdk.dict
var nonTdkData string

//go:embed data/proper.dict
var properData string

//go:embed data/proper-from-corpus.dict
var properCorpusData string

//go:embed data/person-names.dict
var personNamesData string

//go:embed data/locations-tr.dict
var locationsData string

//go:embed data/abbreviations.dict
var abbreviationsData string

// LoadMasterDictionary loads the master dictionary file in Zemberek format
// Format: word [P:POS; A:Attr1, Attr2]
// Examples:
//   kitap
//   gitmek [P:Verb]
//   güzel [P:Adj]
//   abdest [A:NoVoicing]
func LoadMasterDictionary(filePath string) ([]*DictionaryItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []*DictionaryItem
	scanner := bufio.NewScanner(file)

	itemIndex := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		item := ParseDictionaryLine(line, itemIndex)
		if item != nil && item.PrimaryPos != turkish.Punctuation {
			items = append(items, item)
			itemIndex++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// ParseDictionaryLine parses a single line from Zemberek dictionary format
// Examples:
//   "kitap" → Noun (default)
//   "gitmek [P:Verb]" → Verb
//   "güzel [P:Adj]" → Adjective
//   "abdest [A:NoVoicing]" → Noun with NoVoicing attribute
func ParseDictionaryLine(line string, index int) *DictionaryItem {
	var lemma string
	var primaryPos turkish.PrimaryPos = turkish.Noun // Default
	var secondaryPos turkish.SecondaryPos = turkish.NonePos
	attributes := make(map[turkish.RootAttribute]bool)

	// Check if line has attributes [...]
	if strings.Contains(line, "[") {
		parts := strings.Split(line, "[")
		lemma = strings.TrimSpace(parts[0])

		// Parse attributes section
		if len(parts) > 1 {
			attrPart := strings.TrimSuffix(parts[1], "]")
			primaryPos, secondaryPos, attributes = parseAttributeSection(attrPart)
		}
	} else {
		// Simple word without attributes
		lemma = line
		// Infer POS from lemma
		if strings.HasSuffix(lemma, "mek") || strings.HasSuffix(lemma, "mak") {
			primaryPos = turkish.Verb
		}
	}

	// Create dictionary item
	root := lemma
	pronunciation := lemma

	// For verbs, remove -mak/-mek suffix to get the root
	if primaryPos == turkish.Verb {
		if strings.HasSuffix(lemma, "mak") {
			root = strings.TrimSuffix(lemma, "mak")
		} else if strings.HasSuffix(lemma, "mek") {
			root = strings.TrimSuffix(lemma, "mek")
		}
	}

	return NewDictionaryItem(lemma, root, primaryPos, secondaryPos, attributes, pronunciation, index)
}

// parseAttributeSection parses the attribute section of a dictionary line
// Format: "P:Noun; A:NoVoicing, Voicing"
func parseAttributeSection(attrPart string) (turkish.PrimaryPos, turkish.SecondaryPos, map[turkish.RootAttribute]bool) {
	primaryPos := turkish.Noun // Default
	secondaryPos := turkish.NonePos
	attributes := make(map[turkish.RootAttribute]bool)

	// Split by semicolon for different attribute types
	sections := strings.Split(attrPart, ";")

	for _, section := range sections {
		section = strings.TrimSpace(section)

		if strings.HasPrefix(section, "P:") {
			// Primary POS
			posStr := strings.TrimPrefix(section, "P:")
			primaryPos = parsePrimaryPosFromString(posStr)
		} else if strings.HasPrefix(section, "A:") {
			// Attributes
			attrStr := strings.TrimPrefix(section, "A:")
			attrs := strings.Split(attrStr, ",")
			for _, attr := range attrs {
				attrVal := parseRootAttributeFromString(strings.TrimSpace(attr))
				if attrVal != 0 {
					attributes[attrVal] = true
				}
			}
		} else if strings.HasPrefix(section, "Pr:") {
			// Pronunciation (skip for now)
			continue
		}
	}

	return primaryPos, secondaryPos, attributes
}

// parsePrimaryPosFromString parses PrimaryPos from string
func parsePrimaryPosFromString(s string) turkish.PrimaryPos {
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
		return turkish.Noun // Default to Noun
	}
}

// parseRootAttributeFromString parses RootAttribute from string
func parseRootAttributeFromString(s string) turkish.RootAttribute {
	switch s {
	case "Voicing":
		return turkish.Voicing
	case "NoVoicing":
		return turkish.NoVoicing
	case "LastVowelDrop":
		return turkish.LastVowelDrop
	case "Doubling":
		return turkish.Doubling
	case "InverseHarmony":
		return turkish.InverseHarmony
	case "ProgressiveVowelDrop":
		return turkish.ProgressiveVowelDrop
	default:
		return 0
	}
}

// GetDefaultDictionaryPaths returns the default dictionary file paths
func GetDefaultDictionaryPaths() []string {
	return []string{
		"resources/tr/master-dictionary.dict",
		"resources/tr/non-tdk.dict",
		"resources/tr/proper.dict",
		"resources/tr/proper-from-corpus.dict",
		"resources/tr/person-names.dict",
		"resources/tr/locations-tr.dict",
		"resources/tr/abbreviations.dict",
	}
}

// LoadAllDefaultDictionaries loads all default dictionaries from embedded data
func LoadAllDefaultDictionaries() ([]*DictionaryItem, error) {
	var allItems []*DictionaryItem

	// Load from embedded strings (same as Java's DEFAULT_DICTIONARY_RESOURCES)
	dictData := []string{
		masterDictData,
		nonTdkData,
		properData,
		properCorpusData,
		abbreviationsData,
		personNamesData,
		// Note: locations-tr.dict NOT included in Java's default
	}

	for _, data := range dictData {
		items := ParseDictionaryData(data)
		allItems = append(allItems, items...)
	}

	return allItems, nil
}

// ParseDictionaryData parses dictionary data from string
func ParseDictionaryData(data string) []*DictionaryItem {
	var items []*DictionaryItem
	scanner := bufio.NewScanner(strings.NewReader(data))
	index := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		item := ParseDictionaryLine(line, index)
		if item != nil {
			items = append(items, item)
			index++
		}
	}

	return items
}

// BuildWordSet creates a set of all lemmas from dictionary items
func BuildWordSet(items []*DictionaryItem) map[string]bool {
	wordSet := make(map[string]bool, len(items)*2)
	for _, item := range items {
		if item.Lemma != "" {
			wordSet[item.Lemma] = true
			// Also add lowercase version
			wordSet[turkish.Instance.ToLower(item.Lemma)] = true
		}
	}
	return wordSet
}

// GetLemmas extracts all lemmas from dictionary items (excluding punctuation)
func GetLemmas(items []*DictionaryItem) []string {
	lemmas := make([]string, 0, len(items))
	for _, item := range items {
		if item.Lemma != "" && item.PrimaryPos != turkish.Punctuation {
			lemmas = append(lemmas, item.Lemma)
		}
	}
	return lemmas
}

// LoadDefaultLexicon loads the default lexicon with all dictionaries
func LoadDefaultLexicon() (*RootLexicon, error) {
	items, err := LoadAllDefaultDictionaries()
	if err != nil {
		return nil, err
	}
	return NewRootLexicon(items), nil
}
