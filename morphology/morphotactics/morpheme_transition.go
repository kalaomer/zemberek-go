package morphotactics

// MorphemeTransition is the base interface for transitions
type MorphemeTransition interface {
	GetCopy() MorphemeTransition
	HasSurfaceForm() bool
	GetState() *MorphemeState
	GetMorpheme() *Morpheme
	String() string
}

// BaseMorphemeTransition provides common fields for transitions
type BaseMorphemeTransition struct {
	From *MorphemeState
}
