package compression

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// FloatLookup stores dequantized float values for probability/backoff tables.
type FloatLookup struct {
	values []float32
}

// DeserializeFloatLookupDouble reads a float lookup table serialized as double precision values.
func DeserializeFloatLookupDouble(r io.Reader) (*FloatLookup, error) {
	var length int32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	if length < 0 || length > math.MaxInt32/8 {
		return nil, fmt.Errorf("invalid float lookup length: %d", length)
	}
	vals := make([]float32, length)
	for i := int32(0); i < length; i++ {
		var v float64
		if err := binary.Read(r, binary.BigEndian, &v); err != nil {
			return nil, err
		}
		vals[i] = float32(v)
	}
	return &FloatLookup{values: vals}, nil
}

// Get returns the value for the given rank. Out-of-range ranks return negative infinity.
func (fl *FloatLookup) Get(rank int32) float32 {
	if rank < 0 || int(rank) >= len(fl.values) {
		return float32(math.Inf(-1))
	}
	return fl.values[rank]
}

// ChangeBase converts values from source log base to target log base.
func (fl *FloatLookup) ChangeBase(source, target float32) {
	if fl == nil || len(fl.values) == 0 {
		return
	}
	multiplier := float32(math.Log(float64(source)) / math.Log(float64(target)))
	for i, v := range fl.values {
		fl.values[i] = v * multiplier
	}
}

// Range returns the number of entries in the lookup.
func (fl *FloatLookup) Range() int {
	if fl == nil {
		return 0
	}
	return len(fl.values)
}
