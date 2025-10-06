package quantization

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// FloatLookup represents a float lookup table
type FloatLookup struct {
	Data   []float32
	Range_ int32
}

// NewFloatLookup creates a new FloatLookup
func NewFloatLookup(data []float32) *FloatLookup {
	return &FloatLookup{
		Data:   data,
		Range_: int32(len(data)),
	}
}

// ChangeSelfBase changes the base of the lookup data
func (f *FloatLookup) ChangeSelfBase(source, target float64) {
	multiplier := float32(math.Log(source) / math.Log(target))
	for i := range f.Data {
		f.Data[i] *= multiplier
	}
}

// ChangeBase changes the base of the given data
func ChangeBase(data []float32, source, target float64) {
	multiplier := float32(math.Log(source) / math.Log(target))
	for i := range data {
		data[i] *= multiplier
	}
}

// GetLookupFromDouble reads a FloatLookup from a binary stream with double values
func GetLookupFromDouble(r io.Reader) (*FloatLookup, error) {
	var rangeVal int32
	if err := binary.Read(r, binary.BigEndian, &rangeVal); err != nil {
		return nil, err
	}

	values := make([]float32, rangeVal)
	for i := int32(0); i < rangeVal; i++ {
		var val float64
		if err := binary.Read(r, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		values[i] = float32(val)
	}

	return NewFloatLookup(values), nil
}

// Get returns the value at index n
func (f *FloatLookup) Get(n int32) (float32, error) {
	if n >= 0 && n < f.Range_ {
		return f.Data[n], nil
	}
	return 0, fmt.Errorf("value is out of range: %d not in [0, %d)", n, f.Range_)
}
