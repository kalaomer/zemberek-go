package analysis

import (
	"fmt"
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

// SearchPath represents a path in morphological analysis
type SearchPath struct {
	Tail                       string
	CurrentState               *morphotactics.MorphemeState
	Transitions                []*SurfaceTransition
	PhoneticAttributes         map[turkish.PhoneticAttribute]bool
	Terminal                   bool
	ContainsDerivation         bool
	ContainsSuffixWithSurface  bool
}

// NewSearchPath creates a new SearchPath
func NewSearchPath(tail string, currentState *morphotactics.MorphemeState,
	transitions []*SurfaceTransition, phoneticAttributes map[turkish.PhoneticAttribute]bool,
	terminal bool) *SearchPath {

	return &SearchPath{
		Tail:               tail,
		CurrentState:       currentState,
		Transitions:        transitions,
		PhoneticAttributes: phoneticAttributes,
		Terminal:           terminal,
	}
}

// HasDictionaryItem checks if the path has the given dictionary item
func (sp *SearchPath) HasDictionaryItem(item *lexicon.DictionaryItem) bool {
	stemTrans := sp.GetStemTransition()
	if stemTrans == nil {
		return false
	}
	return stemTrans.Item == item
}

// GetStemTransition returns the stem transition
func (sp *SearchPath) GetStemTransition() *morphotactics.StemTransition {
	if len(sp.Transitions) == 0 {
		return nil
	}
	return sp.Transitions[0].LexicalTransition.(*morphotactics.StemTransition)
}

// GetLastTransition returns the last transition
func (sp *SearchPath) GetLastTransition() *SurfaceTransition {
	if len(sp.Transitions) == 0 {
		return nil
	}
	return sp.Transitions[len(sp.Transitions)-1]
}

// GetDictionaryItem returns the dictionary item
func (sp *SearchPath) GetDictionaryItem() *lexicon.DictionaryItem {
	stemTrans := sp.GetStemTransition()
	if stemTrans == nil {
		return nil
	}
	return stemTrans.Item
}

// GetPhoneticAttributes returns phonetic attributes (for interface)
func (sp *SearchPath) GetPhoneticAttributes() map[turkish.PhoneticAttribute]bool {
	return sp.PhoneticAttributes
}

// GetTail returns the tail string (for interface)
func (sp *SearchPath) GetTail() string {
	return sp.Tail
}

// GetTransitions returns transitions as interface slice
func (sp *SearchPath) GetTransitions() []morphotactics.TransitionInterface {
	result := make([]morphotactics.TransitionInterface, len(sp.Transitions))
	for i, t := range sp.Transitions {
		result[i] = t
	}
	return result
}

// ContainsSuffixWithSurface returns true if path contains suffix with surface
func (sp *SearchPath) GetContainsSuffixWithSurface() bool {
	return sp.ContainsSuffixWithSurface
}

// GetPreviousState returns the previous state
func (sp *SearchPath) GetPreviousState() *morphotactics.MorphemeState {
	if len(sp.Transitions) < 2 {
		return nil
	}
	return sp.Transitions[len(sp.Transitions)-2].GetState()
}

// GetCopy creates a copy of the path with a new transition
func (sp *SearchPath) GetCopy(surfaceNode *SurfaceTransition,
	phoneticAttributes map[turkish.PhoneticAttribute]bool) *SearchPath {

	isTerminal := surfaceNode.GetState().Terminal
	newTransitions := append([]*SurfaceTransition{}, sp.Transitions...)
	newTransitions = append(newTransitions, surfaceNode)

	newTail := ""
	if len(surfaceNode.Surface) < len(sp.Tail) {
		newTail = sp.Tail[len(surfaceNode.Surface):]
	}

	path := NewSearchPath(newTail, surfaceNode.GetState(), newTransitions, phoneticAttributes, isTerminal)
	path.ContainsSuffixWithSurface = sp.ContainsSuffixWithSurface || len(surfaceNode.Surface) != 0
	path.ContainsDerivation = sp.ContainsDerivation || surfaceNode.GetState().Derivative
	return path
}

// GetCopyForGeneration creates a copy for word generation
func (sp *SearchPath) GetCopyForGeneration(surfaceNode *SurfaceTransition,
	phoneticAttributes map[turkish.PhoneticAttribute]bool) *SearchPath {

	isTerminal := surfaceNode.GetState().Terminal
	newTransitions := append([]*SurfaceTransition{}, sp.Transitions...)
	newTransitions = append(newTransitions, surfaceNode)

	path := NewSearchPath(sp.Tail, surfaceNode.GetState(), newTransitions, phoneticAttributes, isTerminal)
	path.ContainsSuffixWithSurface = sp.ContainsSuffixWithSurface || len(surfaceNode.Surface) != 0
	path.ContainsDerivation = sp.ContainsDerivation || surfaceNode.GetState().Derivative
	return path
}

// InitialPath creates an initial search path from a stem transition
func InitialPath(stemTransition *morphotactics.StemTransition, tail string) *SearchPath {
	root := NewSurfaceTransition(stemTransition.Surface, stemTransition)

	// Copy phonetic attributes
	attrs := make(map[turkish.PhoneticAttribute]bool)
	for k, v := range stemTransition.PhoneticAttributes {
		attrs[k] = v
	}

	transitions := []*SurfaceTransition{root}
	return NewSearchPath(tail, stemTransition.To, transitions, attrs, stemTransition.To.Terminal)
}

// String returns string representation
func (sp *SearchPath) String() string {
	st := sp.GetStemTransition()
	if st == nil {
		return "[empty path]"
	}

	var parts []string
	for _, s := range sp.Transitions {
		parts = append(parts, s.String())
	}
	morphemeStr := strings.Join(parts, " + ")

	return fmt.Sprintf("[(%s)(-%s) %s]", st.Item, sp.Tail, morphemeStr)
}
