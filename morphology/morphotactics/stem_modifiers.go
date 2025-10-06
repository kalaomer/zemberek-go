package morphotactics

import (
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

// Modifier attributes that require stem modification
var modifierAttributes = []turkish.RootAttribute{
	turkish.Voicing,
	turkish.Doubling,
	turkish.LastVowelDrop,
	turkish.ProgressiveVowelDrop,
	turkish.InverseHarmony,
}

// hasModifierAttribute checks if item has any modifier attribute
func (stm *StemTransitionsMapBased) hasModifierAttribute(item *lexicon.DictionaryItem) bool {
	if item.Attributes == nil {
		return false
	}
	for _, attr := range modifierAttributes {
		if item.Attributes[attr] {
			return true
		}
	}
	return false
}

// generateModifiedRootNodes generates both original and modified stem transitions
func (stm *StemTransitionsMapBased) generateModifiedRootNodes(item *lexicon.DictionaryItem) []*StemTransition {
	// Use pronunciation if available, otherwise use root
	baseSeq := item.Root
	if item.Pronunciation != "" {
		baseSeq = item.Pronunciation
	}

	modifiedSeq := []rune(baseSeq)

	// Calculate attributes
	originalAttrs := GetPhoneticAttributes(baseSeq, nil)
	modifiedAttrs := copyAttributes(originalAttrs)

	// Process each modifier attribute
	for attr := range item.Attributes {
		switch attr {
		case turkish.Voicing:
			// Voice the last letter: t→d, k→g, p→b, ç→c
			if len(modifiedSeq) > 0 {
				last := modifiedSeq[len(modifiedSeq)-1]
				voiced := turkish.Instance.Voice(last)

				// Special case: -nk → -ng
				if len(baseSeq) >= 2 && strings.HasSuffix(baseSeq, "nk") {
					voiced = 'g'
				}

				modifiedSeq[len(modifiedSeq)-1] = voiced
				delete(modifiedAttrs, turkish.LastLetterVoicelessStop)
				originalAttrs[turkish.ExpectsConsonant] = true
				modifiedAttrs[turkish.ExpectsVowel] = true
				modifiedAttrs[turkish.CannotTerminate] = true
			}

		case turkish.Doubling:
			// Double the last letter: ek→ekk
			if len(modifiedSeq) > 0 {
				last := modifiedSeq[len(modifiedSeq)-1]
				modifiedSeq = append(modifiedSeq, last)
				originalAttrs[turkish.ExpectsConsonant] = true
				modifiedAttrs[turkish.ExpectsVowel] = true
				modifiedAttrs[turkish.CannotTerminate] = true
			}

		case turkish.LastVowelDrop:
			// Drop last vowel: ara→ar, kabul→kabl
			if len(modifiedSeq) > 0 {
				lastLetter := modifiedSeq[len(modifiedSeq)-1]
				if turkish.Instance.IsVowel(lastLetter) {
					// Last letter is vowel, drop it
					modifiedSeq = modifiedSeq[:len(modifiedSeq)-1]
					modifiedAttrs[turkish.ExpectsConsonant] = true
					modifiedAttrs[turkish.CannotTerminate] = true
				} else if len(modifiedSeq) > 1 {
					// Last letter is consonant, drop second-to-last (vowel)
					modifiedSeq = append(modifiedSeq[:len(modifiedSeq)-2], modifiedSeq[len(modifiedSeq)-1])
					if item.PrimaryPos != turkish.Verb {
						originalAttrs[turkish.ExpectsConsonant] = true
					}
					modifiedAttrs[turkish.ExpectsVowel] = true
					modifiedAttrs[turkish.CannotTerminate] = true
				}
			}

		case turkish.ProgressiveVowelDrop:
			// Drop last letter for progressive: git→gid (when adding -iyor)
			if len(modifiedSeq) > 1 {
				modifiedSeq = modifiedSeq[:len(modifiedSeq)-1]
				modifiedStr := string(modifiedSeq)
				if turkish.Instance.ContainsVowel(modifiedStr) {
					modifiedAttrs = GetPhoneticAttributes(modifiedStr, nil)
				}
				modifiedAttrs[turkish.LastLetterDropped] = true
			}

		case turkish.InverseHarmony:
			// Force frontal vowel harmony
			originalAttrs[turkish.LastVowelFrontal] = true
			delete(originalAttrs, turkish.LastVowelBack)
			modifiedAttrs[turkish.LastVowelFrontal] = true
			delete(modifiedAttrs, turkish.LastVowelBack)
		}
	}

	// Create original stem transition
	unmodifiedRootState := stm.morphotactics.GetRootState(item, originalAttrs)
	original := NewStemTransition(item.Root, item, originalAttrs, unmodifiedRootState)

	// Create modified stem transition
	modifiedRootState := stm.morphotactics.GetRootState(item, modifiedAttrs)
	modifiedSurface := string(modifiedSeq)
	modified := NewStemTransition(modifiedSurface, item, modifiedAttrs, modifiedRootState)

	// If they're the same, return only one
	if item.Root == modifiedSurface {
		return []*StemTransition{original}
	}

	// Return both transitions
	return []*StemTransition{original, modified}
}

// copyAttributes creates a copy of attribute map
func copyAttributes(attrs map[turkish.PhoneticAttribute]bool) map[turkish.PhoneticAttribute]bool {
	copied := make(map[turkish.PhoneticAttribute]bool)
	for k, v := range attrs {
		copied[k] = v
	}
	return copied
}
