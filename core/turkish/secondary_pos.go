package turkish

// SecondaryPos represents secondary part-of-speech tags
type SecondaryPos int

const (
	UnknownSec SecondaryPos = iota
	DemonstrativePron
	Time
	QuantitivePron
	QuestionPron
	ProperNoun
	PersonalPron
	ReflexivePron
	NonePos
	Ordinal
	Cardinal
	Percentage
	Ratio
	Range
	Real
	Distribution
	Clock
	Date
	Email
	Url
	Mention
	HashTag
	Emoticon
	RomanNumeral
	RegularAbbreviation
	Abbreviation
	PCDat
	PCAcc
	PCIns
	PCNom
	PCGen
	PCAbl
)

var secondaryPosStrings = map[SecondaryPos]string{
	UnknownSec:          "Unk",
	DemonstrativePron:   "Demons",
	Time:                "Time",
	QuantitivePron:      "Quant",
	QuestionPron:        "Ques",
	ProperNoun:          "Prop",
	PersonalPron:        "Pers",
	ReflexivePron:       "Reflex",
	NonePos:             "None",
	Ordinal:             "Ord",
	Cardinal:            "Card",
	Percentage:          "Percent",
	Ratio:               "Ratio",
	Range:               "Range",
	Real:                "Real",
	Distribution:        "Dist",
	Clock:               "Clock",
	Date:                "Date",
	Email:               "Email",
	Url:                 "Url",
	Mention:             "Mention",
	HashTag:             "HashTag",
	Emoticon:            "Emoticon",
	RomanNumeral:        "RomanNumeral",
	RegularAbbreviation: "RegAbbrv",
	Abbreviation:        "Abbrv",
	PCDat:               "PCDat",
	PCAcc:               "PCAcc",
	PCIns:               "PCIns",
	PCNom:               "PCNom",
	PCGen:               "PCGen",
	PCAbl:               "PCAbl",
}

// GetStringForm returns the short form of the POS tag
func (s SecondaryPos) GetStringForm() string {
	return secondaryPosStrings[s]
}
