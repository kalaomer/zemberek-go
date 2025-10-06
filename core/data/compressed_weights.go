package data

import (
	"os"

	"github.com/kalaomer/zemberek-go/core/compression"
)

// CompressedWeights implements WeightLookup using compressed storage
type CompressedWeights struct {
	Lookup *compression.LossyIntLookup
}

// NewCompressedWeights creates a new CompressedWeights
func NewCompressedWeights(lookup *compression.LossyIntLookup) *CompressedWeights {
	return &CompressedWeights{
		Lookup: lookup,
	}
}

// Size returns the size of the lookup
func (c *CompressedWeights) Size() int {
	return c.Lookup.Size()
}

// Get returns the weight for the given key
func (c *CompressedWeights) Get(key string) float32 {
	return c.Lookup.GetAsFloat(key)
}

// Deserialize reads CompressedWeights from a file
func Deserialize(resource string) (*CompressedWeights, error) {
	f, err := os.Open(resource)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lookup, err := compression.DeserializeLossyIntLookup(f)
	if err != nil {
		return nil, err
	}

	return NewCompressedWeights(lookup), nil
}
