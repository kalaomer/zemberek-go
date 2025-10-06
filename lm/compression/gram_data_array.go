package compression

import (
	"encoding/binary"
	"io"
)

const maxBuf int32 = 0x3fffffff

// GramDataArray represents compressed n-gram data
type GramDataArray struct {
	Count       int32
	FpSize      int32
	ProbSize    int32
	BackoffSize int32
	FpMask      int32
	BlockSize   int32
	PageShift   uint
	IndexMask   int32
	Data        [][]byte
}

// NewGramDataArray creates a new GramDataArray from a reader
func NewGramDataArray(r io.Reader) (*GramDataArray, error) {
	gda := &GramDataArray{}

	var values [4]int32
	if err := binary.Read(r, binary.BigEndian, &values); err != nil {
		return nil, err
	}

	gda.Count = values[0]
	gda.FpSize = values[1]
	gda.ProbSize = values[2]
	gda.BackoffSize = values[3]

	if gda.FpSize == 4 {
		gda.FpMask = -1
	} else {
		gda.FpMask = (1 << uint(gda.FpSize*8)) - 1
	}

	gda.BlockSize = gda.FpSize + gda.ProbSize + gda.BackoffSize

	// Calculate page parameters based on block size
	var pageLength int32
	switch gda.BlockSize {
	case 2:
		pageLength = 268435456
		gda.PageShift = 28
		gda.IndexMask = 134217727
	case 4:
		pageLength = 134217728
		gda.PageShift = 27
		gda.IndexMask = 67108863
	default:
		pageLength = 536870912
		gda.PageShift = 29
		gda.IndexMask = 268435455
	}

	pageCounter := int32(1)
	gda.Data = make([][]byte, pageCounter)

	total := int32(0)
	for i := int32(0); i < pageCounter; i++ {
		var readCount int32
		if i < (pageCounter - 1) {
			readCount = pageLength * gda.BlockSize
		} else {
			readCount = (gda.Count * gda.BlockSize) - total
		}

		gda.Data[i] = make([]byte, readCount)
		if _, err := io.ReadFull(r, gda.Data[i]); err != nil {
			return nil, err
		}
		total += readCount
	}

	return gda, nil
}

// GetProbabilityRank returns the probability rank for an index
func (gda *GramDataArray) GetProbabilityRank(index int32) int32 {
	pageID := rshift(index, gda.PageShift)
	pageIndex := (index&gda.IndexMask)*gda.BlockSize + gda.FpSize

	d := gda.Data[pageID]

	switch gda.ProbSize {
	case 1:
		return int32(d[pageIndex]) & 255
	case 2:
		return (int32(d[pageIndex])&255)<<8 | int32(d[pageIndex+1])&255
	case 3:
		return (int32(d[pageIndex])&255)<<16 | (int32(d[pageIndex+1])&255)<<8 | int32(d[pageIndex+2])&255
	default:
		return -1
	}
}

// GetBackOffRank returns the backoff rank for an index
func (gda *GramDataArray) GetBackOffRank(index int32) int32 {
	pageID := rshift(index, gda.PageShift)
	pageIndex := (index&gda.IndexMask)*gda.BlockSize + gda.FpSize + gda.ProbSize

	d := gda.Data[pageID]

	switch gda.BackoffSize {
	case 1:
		return int32(d[pageIndex]) & 255
	case 2:
		return (int32(d[pageIndex])&255)<<8 | int32(d[pageIndex+1])&255
	case 3:
		return (int32(d[pageIndex])&255)<<16 | (int32(d[pageIndex+1])&255)<<8 | int32(d[pageIndex+2])&255
	default:
		return -1
	}
}

// CheckFingerPrint checks if a fingerprint matches
func (gda *GramDataArray) CheckFingerPrint(fpToCheck int32, globalIndex int32) bool {
	fp := fpToCheck & gda.FpMask
	pageIndex := (globalIndex & gda.IndexMask) * gda.BlockSize
	d := gda.Data[rshift(globalIndex, gda.PageShift)]

	switch gda.FpSize {
	case 1:
		return fp == (int32(d[pageIndex]) & 0xFF)
	case 2:
		return (rshift(fp, 8) == int32(d[pageIndex])&0xFF) &&
			((fp & 0xFF) == int32(d[pageIndex+1])&0xFF)
	case 3:
		return (rshift(fp, 16) == int32(d[pageIndex])&0xFF) &&
			((rshift(fp, 8) & 0xFF) == int32(d[pageIndex+1])&0xFF) &&
			((fp & 0xFF) == int32(d[pageIndex+2])&0xFF)
	case 4:
		return (rshift(fp, 24) == int32(d[pageIndex])&0xFF) &&
			((rshift(fp, 16) & 0xFF) == int32(d[pageIndex+1])&0xFF) &&
			((rshift(fp, 8) & 0xFF) == int32(d[pageIndex+2])&0xFF) &&
			((fp & 0xFF) == int32(d[pageIndex+3])&0xFF)
	default:
		return false
	}
}

func rshift(val int32, n uint) int32 {
	return int32(uint32(val) >> n)
}
