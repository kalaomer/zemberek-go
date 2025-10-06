package compression

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/kalaomer/zemberek-go/core/hash"
)

const magic int32 = -889274641

// LossyIntLookup represents a lossy integer lookup table
type LossyIntLookup struct {
	Mphf hash.Mphf
	Data []int32
}

// NewLossyIntLookup creates a new LossyIntLookup
func NewLossyIntLookup(mphf hash.Mphf, data []int32) *LossyIntLookup {
	return &LossyIntLookup{
		Mphf: mphf,
		Data: data,
	}
}

// Get returns the value for the given key
func (l *LossyIntLookup) Get(s string) int32 {
	index := l.Mphf.Get(s) * 2
	fingerprint := GetFingerprint(s)

	if fingerprint == l.Data[index] {
		return l.Data[index+1]
	}
	return 0
}

// Size returns the size of the lookup
func (l *LossyIntLookup) Size() int {
	return len(l.Data) / 2
}

// GetAsFloat returns the value as a float32
func (l *LossyIntLookup) GetAsFloat(s string) float32 {
	return JavaIntBitsToFloat(l.Get(s))
}

// GetFingerprint computes the fingerprint for a string using Java's hashCode
func GetFingerprint(s string) int32 {
	return JavaHashCode(s) & 0x7ffffff
}

// JavaIntBitsToFloat converts an int32 to float32 using Java's intBitsToFloat semantics
func JavaIntBitsToFloat(b int32) float32 {
	return math.Float32frombits(uint32(b))
}

// JavaHashCode computes Java's String.hashCode()
func JavaHashCode(s string) int32 {
	var hash int32 = 0
	for _, c := range s {
		hash = 31*hash + int32(c)
	}
	return hash
}

// Deserialize reads a LossyIntLookup from a binary stream
func DeserializeLossyIntLookup(r io.Reader) (*LossyIntLookup, error) {
	var mag int32
	if err := binary.Read(r, binary.BigEndian, &mag); err != nil {
		return nil, err
	}

	if mag != magic {
		return nil, fmt.Errorf("file does not carry expected magic value: got %d, want %d", mag, magic)
	}

	var length int32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	data := make([]int32, length)
	for i := int32(0); i < length; i++ {
		if err := binary.Read(r, binary.BigEndian, &data[i]); err != nil {
			return nil, err
		}
	}

	mphf, err := hash.DeserializeMultiLevelMphf(r)
	if err != nil {
		return nil, err
	}

	return NewLossyIntLookup(mphf, data), nil
}
