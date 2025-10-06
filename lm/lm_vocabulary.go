package lm

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

)

const (
	DefaultSentenceBeginMarker = "<s>"
	DefaultSentenceEndMarker   = "</s>"
	DefaultUnknownWord         = "<unk>"
)

// LmVocabulary represents a language model vocabulary
type LmVocabulary struct {
	VocabularyIndexMap  map[string]int
	Vocabulary          []string
	UnknownWord         string
	SentenceStart       string
	SentenceEnd         string
	UnknownWordIndex    int
	SentenceStartIndex  int
	SentenceEndIndex    int
}

// NewLmVocabulary creates a new LmVocabulary from a reader
func NewLmVocabulary(r io.Reader) (*LmVocabulary, error) {
	var vocabularyLength int32
	if err := binary.Read(r, binary.BigEndian, &vocabularyLength); err != nil {
		return nil, err
	}

	vocab := make([]string, vocabularyLength)
	for i := int32(0); i < vocabularyLength; i++ {
		var utfLength uint16
		if err := binary.Read(r, binary.BigEndian, &utfLength); err != nil {
			return nil, err
		}

		bytes := make([]byte, utfLength)
		if _, err := io.ReadFull(r, bytes); err != nil {
			return nil, err
		}
		vocab[i] = string(bytes)
	}

	lv := &LmVocabulary{
		VocabularyIndexMap: make(map[string]int),
		UnknownWordIndex:   -1,
		SentenceStartIndex: -1,
		SentenceEndIndex:   -1,
	}

	lv.generateMap(vocab)
	return lv, nil
}

// IndexOf returns the index of a word
func (lv *LmVocabulary) IndexOf(word string) int {
	if idx, ok := lv.VocabularyIndexMap[word]; ok {
		return idx
	}
	return lv.UnknownWordIndex
}

// Size returns the vocabulary size
func (lv *LmVocabulary) Size() int {
	return len(lv.Vocabulary)
}

// ToIndexes converts words to indexes
func (lv *LmVocabulary) ToIndexes(words []string) []int32 {
	indexes := make([]int32, len(words))
	for i, word := range words {
		if idx, ok := lv.VocabularyIndexMap[word]; ok {
			indexes[i] = int32(idx)
		} else {
			indexes[i] = int32(lv.UnknownWordIndex)
		}
	}
	return indexes
}

func (lv *LmVocabulary) generateMap(inputVocabulary []string) {
	indexCounter := 0
	cleanVocab := make([]string, 0)

	for _, word := range inputVocabulary {
		if _, exists := lv.VocabularyIndexMap[word]; exists {
			fmt.Printf("Warning: Language model vocabulary has duplicate item: %s\n", word)
		} else {
			lower := strings.ToLower(turkishLowerMap(word))

			if lower == "<unk>" {
				if lv.UnknownWordIndex != -1 {
					fmt.Printf("Warning: Unknown word was already defined\n")
				} else {
					lv.UnknownWord = word
					lv.UnknownWordIndex = indexCounter
				}
			} else if lower == "<s>" {
				if lv.SentenceStartIndex != -1 {
					fmt.Printf("Warning: Sentence start was already defined\n")
				} else {
					lv.SentenceStart = word
					lv.SentenceStartIndex = indexCounter
				}
			} else if lower == "</s>" {
				if lv.SentenceEndIndex != -1 {
					fmt.Printf("Warning: Sentence end was already defined\n")
				} else {
					lv.SentenceEnd = word
					lv.SentenceEndIndex = indexCounter
				}
			}

			lv.VocabularyIndexMap[word] = indexCounter
			cleanVocab = append(cleanVocab, word)
			indexCounter++
		}
	}

	if lv.UnknownWordIndex == -1 {
		lv.UnknownWord = "<unk>"
		cleanVocab = append(cleanVocab, lv.UnknownWord)
		lv.VocabularyIndexMap[lv.UnknownWord] = indexCounter
		indexCounter++
	}

	lv.UnknownWordIndex = lv.VocabularyIndexMap[lv.UnknownWord]

	if lv.SentenceStartIndex == -1 {
		lv.SentenceStart = "<s>"
		cleanVocab = append(cleanVocab, lv.SentenceStart)
		lv.VocabularyIndexMap[lv.SentenceStart] = indexCounter
		indexCounter++
	}

	lv.SentenceStartIndex = lv.VocabularyIndexMap[lv.SentenceStart]

	if lv.SentenceEndIndex == -1 {
		lv.SentenceEnd = "</s>"
		cleanVocab = append(cleanVocab, lv.SentenceEnd)
		lv.VocabularyIndexMap[lv.SentenceEnd] = indexCounter
	}

	lv.SentenceEndIndex = lv.VocabularyIndexMap[lv.SentenceEnd]
	lv.Vocabulary = cleanVocab
}

func turkishLowerMap(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r == 'I' {
			result.WriteRune('ı')
		} else if r == 'İ' {
			result.WriteRune('i')
		} else {
			for _, c := range strings.ToLower(string(r)) {
				result.WriteRune(c)
			}
		}
	}
	return result.String()
}
