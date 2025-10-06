package turkish

// PhoneticAttribute represents phonetic attributes
type PhoneticAttribute int

const (
	LastLetterVowel PhoneticAttribute = iota
	LastLetterConsonant
	LastVowelFrontal
	LastVowelBack
	LastVowelRounded
	LastVowelUnrounded
	LastLetterVoiceless
	LastLetterVoiced
	LastLetterVoicelessStop
	FirstLetterVowel
	FirstLetterConsonant
	HasNoVowel
	ExpectsVowel
	ExpectsConsonant
	ModifiedPronoun
	UnModifiedPronoun
	LastLetterDropped
	CannotTerminate
)

var phoneticAttributeStrings = map[PhoneticAttribute]string{
	LastLetterVowel:         "LLV",
	LastLetterConsonant:     "LLC",
	LastVowelFrontal:        "LVF",
	LastVowelBack:           "LVB",
	LastVowelRounded:        "LVR",
	LastVowelUnrounded:      "LVuR",
	LastLetterVoiceless:     "LLVless",
	LastLetterVoiced:        "LLVo",
	LastLetterVoicelessStop: "LLVlessStop",
	FirstLetterVowel:        "FLV",
	FirstLetterConsonant:    "FLC",
	HasNoVowel:              "NoVow",
	ExpectsVowel:            "EV",
	ExpectsConsonant:        "EC",
	ModifiedPronoun:         "MP",
	UnModifiedPronoun:       "UMP",
	LastLetterDropped:       "LWD",
	CannotTerminate:         "CNT",
}

// GetStringForm returns the short form string representation
func (p PhoneticAttribute) GetStringForm() string {
	return phoneticAttributeStrings[p]
}
