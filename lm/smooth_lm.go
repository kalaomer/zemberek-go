package lm

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/kalaomer/zemberek-go/core/hash"
	"github.com/kalaomer/zemberek-go/lm/compression"
)

// SmoothLM is a compressed bigram language model reader (subset of Java SmoothLm).
type SmoothLM struct {
	order                 int
	logBase               float32
	unknownBackoff        float32
	vocab                 *Vocabulary
	probabilityLookups    []*compression.FloatLookup
	backoffLookups        []*compression.FloatLookup
	ngramData             []*compression.GramDataArray
	mphfs                 []hash.Mphf
	unigramProbs          []float32
	unigramBackoffs       []float32
	unknownLogProbability float32
}

func LoadSmoothLM(path string) (*SmoothLM, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		path = filepath.Join(path, "lm.2gram.slm")
	}
	if filepath.Ext(path) == "" {
		path = filepath.Join(path, "lm.2gram.slm")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var version int32
	if err := binary.Read(reader, binary.BigEndian, &version); err != nil {
		return nil, err
	}
	var typeInt int32
	if err := binary.Read(reader, binary.BigEndian, &typeInt); err != nil {
		return nil, err
	}
	useLarge := typeInt != 0
	var logBase float64
	if err := binary.Read(reader, binary.BigEndian, &logBase); err != nil {
		return nil, err
	}
	var order int32
	if err := binary.Read(reader, binary.BigEndian, &order); err != nil {
		return nil, err
	}
	if order < 1 {
		return nil, fmt.Errorf("unsupported lm order: %d", order)
	}

	counts := make([]int32, order+1)
	for i := int32(1); i <= order; i++ {
		if err := binary.Read(reader, binary.BigEndian, &counts[i]); err != nil {
			return nil, err
		}
	}

	probabilityLookups := make([]*compression.FloatLookup, order+1)
	for i := int32(1); i <= order; i++ {
		lookup, err := compression.DeserializeFloatLookupDouble(reader)
		if err != nil {
			return nil, fmt.Errorf("reading probability lookup %d: %w", i, err)
		}
		probabilityLookups[i] = lookup
	}

	backoffLookups := make([]*compression.FloatLookup, order)
	for i := int32(1); i < order; i++ {
		lookup, err := compression.DeserializeFloatLookupDouble(reader)
		if err != nil {
			return nil, fmt.Errorf("reading backoff lookup %d: %w", i, err)
		}
		backoffLookups[i] = lookup
	}

	ngramData := make([]*compression.GramDataArray, order+1)
	for i := int32(1); i <= order; i++ {
		gda, err := compression.NewGramDataArray(reader)
		if err != nil {
			return nil, fmt.Errorf("reading gram data %d: %w", i, err)
		}
		ngramData[i] = gda
	}

	unigramCount := ngramData[1].Count
	unigramProbs := make([]float32, unigramCount)
	unigramBackoffs := make([]float32, unigramCount)
	for i := int32(0); i < unigramCount; i++ {
		rank := ngramData[1].GetProbabilityRank(i)
		unigramProbs[i] = probabilityLookups[1].Get(rank)
		if order > 1 {
			backRank := ngramData[1].GetBackOffRank(i)
			unigramBackoffs[i] = backoffLookups[1].Get(backRank)
		}
	}

	mphfs := make([]hash.Mphf, order+1)
	if useLarge {
		for i := int32(2); i <= order; i++ {
			mphf, err := hash.DeserializeLargeNgramMphf(reader)
			if err != nil {
				return nil, fmt.Errorf("reading large mphf level %d: %w", i, err)
			}
			mphfs[i] = mphf
		}
	} else {
		for i := int32(2); i <= order; i++ {
			mphf, err := hash.DeserializeMultiLevelMphf(reader)
			if err != nil {
				return nil, fmt.Errorf("reading mphf level %d: %w", i, err)
			}
			mphfs[i] = mphf
		}
	}

	vocab, err := NewLmVocabulary(reader)
	if err != nil {
		return nil, fmt.Errorf("reading vocabulary: %w", err)
	}

	lm := &SmoothLM{
		order:                 int(order),
		logBase:               float32(logBase),
		unknownBackoff:        0,
		vocab:                 vocab,
		probabilityLookups:    probabilityLookups,
		backoffLookups:        backoffLookups,
		ngramData:             ngramData,
		mphfs:                 mphfs,
		unigramProbs:          unigramProbs,
		unigramBackoffs:       unigramBackoffs,
		unknownLogProbability: float32(math.Log(1e-8)),
	}

	// Adjust log base if necessary
	if math.Abs(logBase-math.E) > 1e-5 {
		lm.changeLogBase(float32(logBase), float32(math.E))
	}

	return lm, nil
}

func (lm *SmoothLM) changeLogBase(source, target float32) {
	if source == target {
		return
	}
	multiplier := float32(math.Log(float64(source)) / math.Log(float64(target)))
	for i, v := range lm.unigramProbs {
		lm.unigramProbs[i] = v * multiplier
	}
	for i, v := range lm.unigramBackoffs {
		lm.unigramBackoffs[i] = v * multiplier
	}
	for i := 2; i < len(lm.probabilityLookups); i++ {
		if lm.probabilityLookups[i] != nil {
			lm.probabilityLookups[i].ChangeBase(source, target)
		}
		if i < len(lm.backoffLookups) && lm.backoffLookups[i] != nil {
			lm.backoffLookups[i].ChangeBase(source, target)
		}
	}
	lm.logBase = target
}

func (lm *SmoothLM) GetOrder() int {
	return lm.order
}

func (lm *SmoothLM) GetVocabulary() *Vocabulary {
	return lm.vocab
}

func (lm *SmoothLM) NgramExists(indexes []int) bool {
	n := len(indexes)
	if n == 0 || n > lm.order {
		return false
	}
	if n == 1 {
		idx := indexes[0]
		return idx >= 0 && idx < len(lm.unigramProbs)
	}
	ints := toInt32Slice(indexes)
	hashVal := hash.HashForIntSlice(ints, -1)
	mphf := lm.mphfs[n]
	if mphf == nil {
		return false
	}
	idx := mphf.Get(ints, hashVal)
	return lm.ngramData[n].CheckFingerPrint(hashVal, idx)
}

func (lm *SmoothLM) GetProbability(indexes []int) float32 {
	n := len(indexes)
	if n == 0 {
		return 0
	}
	if n == 1 {
		idx := indexes[0]
		if idx >= 0 && idx < len(lm.unigramProbs) {
			return lm.unigramProbs[idx]
		}
		return lm.unknownLogProbability
	}
	if n > lm.order {
		// back off recursively
		return lm.GetProbability(indexes[1:])
	}

	ints := toInt32Slice(indexes)
	hashVal := hash.HashForIntSlice(ints, -1)
	mphf := lm.mphfs[n]
	if mphf != nil {
		idx := mphf.Get(ints, hashVal)
		if lm.ngramData[n].CheckFingerPrint(hashVal, idx) {
			rank := lm.ngramData[n].GetProbabilityRank(idx)
			if lm.probabilityLookups[n] != nil {
				val := lm.probabilityLookups[n].Get(rank)
				if !math.IsInf(float64(val), -1) {
					return val
				}
			}
		}
	}

	if n == 2 {
		prev := indexes[0]
		back := lm.unknownBackoff
		if prev >= 0 && prev < len(lm.unigramBackoffs) {
			back = lm.unigramBackoffs[prev]
		}
		return back + lm.GetProbability(indexes[1:])
	}
	return lm.GetProbability(indexes[1:])
}

func toInt32Slice(in []int) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}
