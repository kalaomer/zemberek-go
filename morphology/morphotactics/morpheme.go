package morphotactics

import (
	"github.com/kalaomer/zemberek-go/core/turkish"
)

// Morpheme represents a morpheme
type Morpheme struct {
	Name            string
	ID              string
	Informal        bool
	Derivational    bool
	Pos             turkish.PrimaryPos
	MappedMorpheme  *Morpheme
}

// UnknownMorpheme represents an unknown morpheme
var UnknownMorpheme *Morpheme

func init() {
	UnknownMorpheme = &Morpheme{
		Name: "Unknown",
		ID:   "Unknown",
	}
}

// NewMorpheme creates a new Morpheme
func NewMorpheme(name, id string) *Morpheme {
	return &Morpheme{
		Name: name,
		ID:   id,
	}
}

// NewMorphemeWithPos creates a new Morpheme with POS
func NewMorphemeWithPos(name, id string, pos turkish.PrimaryPos) *Morpheme {
	return &Morpheme{
		Name: name,
		ID:   id,
		Pos:  pos,
	}
}

// NewDerivationalMorpheme creates a new derivational morpheme
func NewDerivationalMorpheme(name, id string) *Morpheme {
	return &Morpheme{
		Name:         name,
		ID:           id,
		Derivational: true,
	}
}

// String returns string representation
func (m *Morpheme) String() string {
	return m.Name + ":" + m.ID
}

// MorphemeBuilder helps build a Morpheme
type MorphemeBuilder struct {
	name           string
	id             string
	derivational   bool
	informal       bool
	pos            turkish.PrimaryPos
	mappedMorpheme *Morpheme
}

// NewMorphemeBuilder creates a new builder
func NewMorphemeBuilder(name, id string) *MorphemeBuilder {
	return &MorphemeBuilder{
		name: name,
		id:   id,
	}
}

// Informal marks the morpheme as informal
func (mb *MorphemeBuilder) Informal() *MorphemeBuilder {
	mb.informal = true
	return mb
}

// Derivational marks the morpheme as derivational
func (mb *MorphemeBuilder) Derivational() *MorphemeBuilder {
	mb.derivational = true
	return mb
}

// WithPos sets the POS
func (mb *MorphemeBuilder) WithPos(pos turkish.PrimaryPos) *MorphemeBuilder {
	mb.pos = pos
	return mb
}

// MappedMorpheme sets the mapped morpheme
func (mb *MorphemeBuilder) MappedMorpheme(morpheme *Morpheme) *MorphemeBuilder {
	mb.mappedMorpheme = morpheme
	return mb
}

// Build builds the morpheme
func (mb *MorphemeBuilder) Build() *Morpheme {
	return &Morpheme{
		Name:           mb.name,
		ID:             mb.id,
		Informal:       mb.informal,
		Derivational:   mb.derivational,
		Pos:            mb.pos,
		MappedMorpheme: mb.mappedMorpheme,
	}
}
