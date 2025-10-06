package turkish

// RootAttribute represents attributes of a root
type RootAttribute int

const (
	AoristI RootAttribute = iota + 1
	AoristA
	ProgressiveVowelDrop
	PassiveIn
	CausativeT
	Voicing
	NoVoicing
	InverseHarmony
	Doubling
	LastVowelDrop
	CompoundP3sg
	NoSuffix
	NounConsInsertN
	NoQuote
	CompoundP3sgRoot
	Reflexive
	Reciprocal
	NonReciprocal
	Ext
	Runtime
	Dummy
	ImplicitDative
	ImplicitPlural
	ImplicitP1sg
	ImplicitP2sg
	FamilyMember
	PronunciationGuessed
	Informal
	LocaleEn
	Unknown
)

var rootAttributeNames = map[RootAttribute]string{
	AoristI:              "Aorist_I",
	AoristA:              "Aorist_A",
	ProgressiveVowelDrop: "ProgressiveVowelDrop",
	PassiveIn:            "Passive_In",
	CausativeT:           "Causative_t",
	Voicing:              "Voicing",
	NoVoicing:            "NoVoicing",
	InverseHarmony:       "InverseHarmony",
	Doubling:             "Doubling",
	LastVowelDrop:        "LastVowelDrop",
	CompoundP3sg:         "CompoundP3sg",
	NoSuffix:             "NoSuffix",
	NounConsInsertN:      "NounConsInsert_n",
	NoQuote:              "NoQuote",
	CompoundP3sgRoot:     "CompoundP3sgRoot",
	Reflexive:            "Reflexive",
	Reciprocal:           "Reciprocal",
	NonReciprocal:        "NonReciprocal",
	Ext:                  "Ext",
	Runtime:              "Runtime",
	Dummy:                "Dummy",
	ImplicitDative:       "ImplicitDative",
	ImplicitPlural:       "ImplicitPlural",
	ImplicitP1sg:         "ImplicitP1sg",
	ImplicitP2sg:         "ImplicitP2sg",
	FamilyMember:         "FamilyMember",
	PronunciationGuessed: "PronunciationGuessed",
	Informal:             "Informal",
	LocaleEn:             "LocaleEn",
	Unknown:              "Unknown",
}

// GetStringForm returns the name of the attribute
func (r RootAttribute) GetStringForm() string {
	return rootAttributeNames[r]
}
