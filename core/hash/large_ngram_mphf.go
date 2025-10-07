package hash

import (
	"encoding/binary"
	"fmt"
	"io"
)

const defaultChunkSizeInBits = 22

// LargeNgramMphf represents a large n-gram minimum perfect hash function
type LargeNgramMphf struct {
	MaxBitMask int32
	BucketMask int32
	PageShift  int32
	Mphfs      []*MultiLevelMphf
	Offsets    []int32
}

// NewLargeNgramMphf creates a new LargeNgramMphf
func NewLargeNgramMphf(maxBitMask, bucketMask, pageShift int32, mphfs []*MultiLevelMphf, offsets []int32) *LargeNgramMphf {
	return &LargeNgramMphf{
		MaxBitMask: maxBitMask,
		BucketMask: bucketMask,
		PageShift:  pageShift,
		Mphfs:      mphfs,
		Offsets:    offsets,
	}
}

// DeserializeLargeNgramMphf reads a LargeNgramMphf from a binary stream
func DeserializeLargeNgramMphf(r io.Reader) (*LargeNgramMphf, error) {
	var maxBitMask, bucketMask, pageShift, phfCount int32

	if err := binary.Read(r, binary.BigEndian, &maxBitMask); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &bucketMask); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &pageShift); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &phfCount); err != nil {
		return nil, err
	}

	offsets := make([]int32, phfCount)
	for i := int32(0); i < phfCount; i++ {
		if err := binary.Read(r, binary.BigEndian, &offsets[i]); err != nil {
			return nil, err
		}
	}

	hashes := make([]*MultiLevelMphf, phfCount)
	for i := int32(0); i < phfCount; i++ {
		mphf, err := DeserializeMultiLevelMphf(r)
		if err != nil {
			return nil, err
		}
		hashes[i] = mphf
	}

	return NewLargeNgramMphf(maxBitMask, bucketMask, pageShift, hashes, offsets), nil
}

// Get returns the hash value for the given key and satisfies the Mphf interface
func (l *LargeNgramMphf) Get(key interface{}, initialHash ...int32) int32 {
	switch v := key.(type) {
	case []int32:
		var hashVal int32
		if len(initialHash) > 0 {
			hashVal = initialHash[0]
		} else {
			hashVal = HashForIntSlice(v, -1)
		}
		return l.getForIntSlice(v, hashVal)
	default:
		panic(fmt.Sprintf("unsupported key type for LargeNgramMphf: %T", key))
	}
}

func (l *LargeNgramMphf) getForIntSlice(ngram []int32, hashVal int32) int32 {
	pageIndex := Rshift(hashVal&l.MaxBitMask, uint(l.PageShift))
	return l.Mphfs[pageIndex].Get(ngram, hashVal) + l.Offsets[pageIndex]
}
