package turkish

// TurkicLetter represents a Letter which contains Turkic language specific attributes,
// such as vowel type, English equivalent characters.
type TurkicLetter struct {
	CharValue  rune
	Vowel      bool
	Frontal    bool
	Rounded    bool
	Voiceless  bool
	Continuant bool
}

// Undefined letter constant
var Undefined = &TurkicLetter{CharValue: '\u0000'}

// NewTurkicLetter creates a new TurkicLetter instance
func NewTurkicLetter(charValue rune, vowel, frontal, rounded, voiceless, continuant bool) *TurkicLetter {
	return &TurkicLetter{
		CharValue:  charValue,
		Vowel:      vowel,
		Frontal:    frontal,
		Rounded:    rounded,
		Voiceless:  voiceless,
		Continuant: continuant,
	}
}

// IsVowel returns true if this is a vowel
func (t *TurkicLetter) IsVowel() bool {
	return t.Vowel
}

// IsConsonant returns true if this is a consonant
func (t *TurkicLetter) IsConsonant() bool {
	return !t.Vowel
}

// IsFrontal returns true if this is a frontal letter
func (t *TurkicLetter) IsFrontal() bool {
	return t.Frontal
}

// IsRounded returns true if this is a rounded letter
func (t *TurkicLetter) IsRounded() bool {
	return t.Rounded
}

// IsVoiceless returns true if this is a voiceless letter
func (t *TurkicLetter) IsVoiceless() bool {
	return t.Voiceless
}

// IsStopConsonant returns true if this is a stop consonant
func (t *TurkicLetter) IsStopConsonant() bool {
	return t.Voiceless && !t.Continuant
}

// CopyFor creates a copy of this letter with a different character value
func (t *TurkicLetter) CopyFor(c rune) *TurkicLetter {
	return NewTurkicLetter(c, t.Vowel, t.Frontal, t.Rounded, t.Voiceless, t.Continuant)
}
