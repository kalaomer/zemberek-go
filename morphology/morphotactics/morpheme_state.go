package morphotactics

import "fmt"

// MorphemeState represents a state in the morphotactic graph
type MorphemeState struct {
	ID         string
	Morpheme   *Morpheme
	Terminal   bool
	Derivative bool
	PosRoot    bool
	Outgoing   []MorphemeTransition
	Incoming   []MorphemeTransition
}

// NewMorphemeState creates a new MorphemeState
func NewMorphemeState(id string, morpheme *Morpheme, terminal, derivative, posRoot bool) *MorphemeState {
	return &MorphemeState{
		ID:         id,
		Morpheme:   morpheme,
		Terminal:   terminal,
		Derivative: derivative,
		PosRoot:    posRoot,
		Outgoing:   make([]MorphemeTransition, 0),
		Incoming:   make([]MorphemeTransition, 0),
	}
}

// Terminal creates a terminal state
func Terminal(id string, morpheme *Morpheme, posRoot bool) *MorphemeState {
	return NewMorphemeState(id, morpheme, true, false, posRoot)
}

// NonTerminal creates a non-terminal state
func NonTerminal(id string, morpheme *Morpheme, posRoot bool) *MorphemeState {
	return NewMorphemeState(id, morpheme, false, false, posRoot)
}

// TerminalDerivative creates a terminal derivative state
func TerminalDerivative(id string, morpheme *Morpheme, posRoot bool) *MorphemeState {
	return NewMorphemeState(id, morpheme, true, true, posRoot)
}

// NonTerminalDerivative creates a non-terminal derivative state
func NonTerminalDerivative(id string, morpheme *Morpheme, posRoot bool) *MorphemeState {
	return NewMorphemeState(id, morpheme, false, true, posRoot)
}

// String returns string representation
func (ms *MorphemeState) String() string {
	return fmt.Sprintf("[%s:%s]", ms.ID, ms.Morpheme.ID)
}

// AddOutgoing adds outgoing transitions
func (ms *MorphemeState) AddOutgoing(transitions ...MorphemeTransition) *MorphemeState {
	ms.Outgoing = append(ms.Outgoing, transitions...)
	return ms
}

// AddIncoming adds incoming transitions
func (ms *MorphemeState) AddIncoming(transitions ...MorphemeTransition) *MorphemeState {
	ms.Incoming = append(ms.Incoming, transitions...)
	return ms
}

// CopyOutgoingTransitionsFrom copies outgoing transitions from another state
func (ms *MorphemeState) CopyOutgoingTransitionsFrom(state *MorphemeState) {
	for _, transition := range state.Outgoing {
		copy := transition.GetCopy()
		// Set the from state to this state
		// This would need to be handled by the specific transition type
		ms.AddOutgoing(copy)
	}
}

// RemoveTransitionsTo removes transitions to a specific morpheme
func (ms *MorphemeState) RemoveTransitionsTo(morpheme *Morpheme) {
	newOutgoing := make([]MorphemeTransition, 0)
	for _, transition := range ms.Outgoing {
		// This needs to be implemented properly based on transition type
		// For now, keep all transitions
		newOutgoing = append(newOutgoing, transition)
	}
	ms.Outgoing = newOutgoing
}

// Helper functions for creating states
func NewMorphemeStateTerminal(id string, morpheme *Morpheme) *MorphemeState {
	return Terminal(id, morpheme, false)
}

func NewMorphemeStateNonTerminal(id string, morpheme *Morpheme) *MorphemeState {
	return NonTerminal(id, morpheme, false)
}

func NewMorphemeStateTerminalDerivative(id string, morpheme *Morpheme) *MorphemeState {
	return TerminalDerivative(id, morpheme, false)
}

func NewMorphemeStateNonTerminalDerivative(id string, morpheme *Morpheme) *MorphemeState {
	return NonTerminalDerivative(id, morpheme, false)
}

// MorphemeStateBuilder helps build MorphemeState
type MorphemeStateBuilder struct {
	id         string
	morpheme   *Morpheme
	terminal   bool
	derivative bool
	posRoot    bool
}

// NewMorphemeStateBuilder creates a builder
func NewMorphemeStateBuilder(id string, morpheme *Morpheme) *MorphemeStateBuilder {
	return &MorphemeStateBuilder{
		id:       id,
		morpheme: morpheme,
		posRoot:  false,
	}
}

// SetTerminal sets terminal flag
func (b *MorphemeStateBuilder) SetTerminal(terminal bool) *MorphemeStateBuilder {
	b.terminal = terminal
	return b
}

// SetDerivative sets derivative flag
func (b *MorphemeStateBuilder) SetDerivative(derivative bool) *MorphemeStateBuilder {
	b.derivative = derivative
	return b
}

// SetPosRoot sets posRoot flag
func (b *MorphemeStateBuilder) SetPosRoot(posRoot bool) *MorphemeStateBuilder {
	b.posRoot = posRoot
	return b
}

// Build builds the state
func (b *MorphemeStateBuilder) Build() *MorphemeState {
	return NewMorphemeState(b.id, b.morpheme, b.terminal, b.derivative, b.posRoot)
}
