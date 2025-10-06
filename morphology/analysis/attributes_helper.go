package analysis

import (
	"github.com/kalaomer/zemberek-go/core/turkish"
)

// AttributesHelper helps with phonetic attribute management
type AttributesHelper struct{}

var (
	alphabet = turkish.Instance
	noVowelAttributes = []turkish.PhoneticAttribute{
		turkish.LastLetterConsonant,
		turkish.FirstLetterConsonant,
		turkish.HasNoVowel,
	}
)

// GetMorphemicAttributes calculates phonetic attributes for a morpheme
func GetMorphemicAttributes(seq string, predecessorAttrs map[turkish.PhoneticAttribute]bool) map[turkish.PhoneticAttribute]bool {
	if predecessorAttrs == nil {
		predecessorAttrs = make(map[turkish.PhoneticAttribute]bool)
	}

	if len(seq) == 0 {
		// Copy predecessor attributes
		attrs := make(map[turkish.PhoneticAttribute]bool)
		for k, v := range predecessorAttrs {
			attrs[k] = v
		}
		return attrs
	}

	attrs := make(map[turkish.PhoneticAttribute]bool)

	if alphabet.ContainsVowel(seq) {
		last := alphabet.GetLastLetter(seq)

		if last.IsVowel() {
			attrs[turkish.LastLetterVowel] = true
		} else {
			attrs[turkish.LastLetterConsonant] = true
		}

		lastVowel := last
		if !last.IsVowel() {
			lastVowel = alphabet.GetLastVowel(seq)
		}

		if lastVowel.IsFrontal() {
			attrs[turkish.LastVowelFrontal] = true
		} else {
			attrs[turkish.LastVowelBack] = true
		}

		if lastVowel.IsRounded() {
			attrs[turkish.LastVowelRounded] = true
		} else {
			attrs[turkish.LastVowelUnrounded] = true
		}

		if alphabet.GetFirstLetter(seq).IsVowel() {
			attrs[turkish.FirstLetterVowel] = true
		} else {
			attrs[turkish.FirstLetterConsonant] = true
		}
	} else {
		// Copy predecessor attributes
		for k, v := range predecessorAttrs {
			attrs[k] = v
		}

		// Add no vowel attributes
		for _, attr := range noVowelAttributes {
			attrs[attr] = true
		}

		// Remove conflicting attributes
		delete(attrs, turkish.LastLetterVowel)
		delete(attrs, turkish.ExpectsConsonant)
	}

	last := alphabet.GetLastLetter(seq)
	if last.IsVoiceless() {
		attrs[turkish.LastLetterVoiceless] = true
		if last.IsStopConsonant() {
			attrs[turkish.LastLetterVoicelessStop] = true
		}
	} else {
		attrs[turkish.LastLetterVoiced] = true
	}

	return attrs
}
