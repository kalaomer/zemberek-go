package morphotactics

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

// StemTransition represents a transition from root/stem
type StemTransition struct {
	BaseMorphemeTransition
	Surface             string
	Item                *lexicon.DictionaryItem
	PhoneticAttributes  map[turkish.PhoneticAttribute]bool
	To                  *MorphemeState
	cachedHash          int
}

// NewStemTransition creates a new StemTransition
func NewStemTransition(surface string, item *lexicon.DictionaryItem,
	phoneticAttributes map[turkish.PhoneticAttribute]bool, toState *MorphemeState) *StemTransition {

	st := &StemTransition{
		Surface:            surface,
		Item:               item,
		PhoneticAttributes: phoneticAttributes,
		To:                 toState,
	}
	st.cachedHash = st.computeHash()
	return st
}

// GetCopy returns a copy of this transition
func (st *StemTransition) GetCopy() MorphemeTransition {
	return &StemTransition{
		Surface:            st.Surface,
		Item:               st.Item,
		PhoneticAttributes: st.PhoneticAttributes,
		To:                 st.To,
		BaseMorphemeTransition: BaseMorphemeTransition{
			From: st.From,
		},
	}
}

// HasSurfaceForm returns true (stem transitions have surface forms)
func (st *StemTransition) HasSurfaceForm() bool {
	return true
}

// GetState returns the target state
func (st *StemTransition) GetState() *MorphemeState {
	return st.To
}

// GetMorpheme returns the morpheme
func (st *StemTransition) GetMorpheme() *Morpheme {
	if st.To == nil {
		return nil
	}
	return st.To.Morpheme
}

// String returns string representation
func (st *StemTransition) String() string {
	return fmt.Sprintf("[(Dict:%s):%s → %s]", st.Item, st.Surface, st.To)
}

func (st *StemTransition) computeHash() int {
	result := hashString(st.Surface)
	result = 31*result + hashPointer(st.Item)
	for attr := range st.PhoneticAttributes {
		result = 31*result + int(attr)
	}
	return result
}

func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	return h
}

func hashPointer(p interface{}) int {
	// Simple pointer-based hash
	s := fmt.Sprintf("%p", p)
	if len(s) > 2 {
		s = s[2:]
	}
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	return h
}
