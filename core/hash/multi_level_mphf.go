package hash

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	hashMultiplier  int32 = 16777619
	initialHashSeed int32 = -2128831035
	bitMask21       int32 = 2097151
)

// MultiLevelMphf is a Minimum Perfect Hash Function implementation
type MultiLevelMphf struct {
	hashLevelData []*HashIndexes
}

// HashIndexes represents hash indexes for a level
type HashIndexes struct {
	KeyAmount            int32
	BucketAmount         int32
	BucketHashSeedValues []byte
	FailedIndexes        []int32
}

// NewMultiLevelMphf creates a new MultiLevelMphf
func NewMultiLevelMphf(hashLevelData []*HashIndexes) *MultiLevelMphf {
	return &MultiLevelMphf{
		hashLevelData: hashLevelData,
	}
}

// Deserialize reads a MultiLevelMphf from a binary stream
func DeserializeMultiLevelMphf(r io.Reader) (*MultiLevelMphf, error) {
	var levelCount int32
	if err := binary.Read(r, binary.BigEndian, &levelCount); err != nil {
		return nil, err
	}

	indexes := make([]*HashIndexes, levelCount)

	for i := int32(0); i < levelCount; i++ {
		var keyCount, bucketAmount int32
		if err := binary.Read(r, binary.BigEndian, &keyCount); err != nil {
			return nil, err
		}
		if err := binary.Read(r, binary.BigEndian, &bucketAmount); err != nil {
			return nil, err
		}

		hashSeedValues := make([]byte, bucketAmount)
		if _, err := io.ReadFull(r, hashSeedValues); err != nil {
			return nil, err
		}

		var failedIndexesCount int32
		if err := binary.Read(r, binary.BigEndian, &failedIndexesCount); err != nil {
			return nil, err
		}

		failedIndexes := make([]int32, failedIndexesCount)
		for j := int32(0); j < failedIndexesCount; j++ {
			if err := binary.Read(r, binary.BigEndian, &failedIndexes[j]); err != nil {
				return nil, err
			}
		}

		indexes[i] = &HashIndexes{
			KeyAmount:            keyCount,
			BucketAmount:         bucketAmount,
			BucketHashSeedValues: hashSeedValues,
			FailedIndexes:        failedIndexes,
		}
	}

	return NewMultiLevelMphf(indexes), nil
}

// HashForStr computes hash for a string
func HashForStr(data string, seed int32) int32 {
	d := seed
	if seed <= 0 {
		d = initialHashSeed
	}

	for _, c := range data {
		d = (d ^ int32(c)) * hashMultiplier
	}

	return d & 0x7fffffff
}

// HashForIntSlice computes hash for an int slice
func HashForIntSlice(data []int32, seed int32) int32 {
	d := seed
	if seed <= 0 {
		d = initialHashSeed
	}

	for _, a := range data {
		d = (d ^ a) * hashMultiplier
	}

	return d & 0x7fffffff
}

// Hash computes hash for any supported type
func Hash(data interface{}, seed int32) int32 {
	switch v := data.(type) {
	case string:
		return HashForStr(v, seed)
	case []int32:
		return HashForIntSlice(v, seed)
	default:
		panic(fmt.Sprintf("unsupported data type: %T", data))
	}
}

// GetSeed returns the seed for a fingerprint
func (hi *HashIndexes) GetSeed(fingerprint int32) int32 {
	idx := fingerprint % hi.BucketAmount
	if idx < 0 {
		idx += hi.BucketAmount
	}
	return int32(hi.BucketHashSeedValues[idx]) & 0xFF
}

// Get returns the hash value for the given key
func (m *MultiLevelMphf) Get(key interface{}, initialHash ...int32) int32 {
	switch v := key.(type) {
	case string:
		return m.getForStr(v, initialHash...)
	case []int32:
		var hash int32
		if len(initialHash) > 0 {
			hash = initialHash[0]
		} else {
			hash = HashForIntSlice(v, -1)
		}
		return m.getForIntSlice(v, hash)
	default:
		panic(fmt.Sprintf("unsupported key type: %T", key))
	}
}

func (m *MultiLevelMphf) getForStr(key string, initialHash ...int32) int32 {
	var hash int32
	if len(initialHash) > 0 {
		hash = initialHash[0]
	} else {
		hash = HashForStr(key, -1)
	}

	for i, hd := range m.hashLevelData {
		seed := hd.GetSeed(hash)
		if seed != 0 {
			if i == 0 {
				result := HashForStr(key, seed) % m.hashLevelData[0].KeyAmount
				if result < 0 {
					result += m.hashLevelData[0].KeyAmount
				}
				return result
			} else {
				hashVal := HashForStr(key, seed) % m.hashLevelData[i].KeyAmount
				if hashVal < 0 {
					hashVal += m.hashLevelData[i].KeyAmount
				}
				return m.hashLevelData[i-1].FailedIndexes[hashVal]
			}
		}
	}

	panic("Cannot be here")
}

func (m *MultiLevelMphf) getForIntSlice(key []int32, initialHash int32) int32 {
	for i, hd := range m.hashLevelData {
		seed := hd.GetSeed(initialHash)
		if seed != 0 {
			if i == 0 {
				result := Hash(key, seed) % m.hashLevelData[0].KeyAmount
				if result < 0 {
					result += m.hashLevelData[0].KeyAmount
				}
				return result
			} else {
				hashVal := Hash(key, seed) % m.hashLevelData[i].KeyAmount
				if hashVal < 0 {
					hashVal += m.hashLevelData[i].KeyAmount
				}
				return m.hashLevelData[i-1].FailedIndexes[hashVal]
			}
		}
	}

	panic("Cannot be here")
}
