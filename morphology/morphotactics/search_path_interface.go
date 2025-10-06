package morphotactics

import (
	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

// SearchPathInterface defines the interface for search paths used in conditions
// This avoids circular dependency between morphotactics and analysis packages
type SearchPathInterface interface {
	// GetDictionaryItem returns the dictionary item for this path
	GetDictionaryItem() *lexicon.DictionaryItem

	// HasDictionaryItem checks if path has the given item
	HasDictionaryItem(item *lexicon.DictionaryItem) bool

	// GetPreviousState returns the previous morpheme state
	GetPreviousState() *MorphemeState

	// GetStemTransition returns the stem transition
	GetStemTransition() *StemTransition

	// ContainsSuffixWithSurface checks if path has suffix with surface form
	GetContainsSuffixWithSurface() bool

	// GetPhoneticAttributes returns phonetic attributes
	GetPhoneticAttributes() map[turkish.PhoneticAttribute]bool

	// GetTail returns the remaining tail string
	GetTail() string

	// GetTransitions returns the transitions
	GetTransitions() []TransitionInterface
}

// TransitionInterface defines interface for transitions
type TransitionInterface interface {
	GetState() *MorphemeState
	GetMorpheme() *Morpheme
	GetSurface() string
}
