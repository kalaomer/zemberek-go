package turkish

// PrimaryPos represents primary part-of-speech tags
type PrimaryPos int

const (
	Noun PrimaryPos = iota
	Adjective
	Adverb
	Conjunction
	Interjection
	Verb
	Pronoun
	Numeral
	Determiner
	PostPositive
	Question
	Duplicator
	Punctuation
	UnknownPos
)

var primaryPosStrings = map[PrimaryPos]string{
	Noun:         "Noun",
	Adjective:    "Adj",
	Adverb:       "Adv",
	Conjunction:  "Conj",
	Interjection: "Interj",
	Verb:         "Verb",
	Pronoun:      "Pron",
	Numeral:      "Num",
	Determiner:   "Det",
	PostPositive: "Postp",
	Question:     "Ques",
	Duplicator:   "Dup",
	Punctuation:  "Punc",
	UnknownPos:   "Unk",
}

// GetStringForm returns the short form of the POS tag
func (p PrimaryPos) GetStringForm() string {
	return primaryPosStrings[p]
}
