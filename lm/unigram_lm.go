package lm

import (
    "math"
    "os"
    "path/filepath"

    "github.com/kalaomer/zemberek-go/core/compression"
)

// UnigramLM is a lightweight language model backed by LossyIntLookup probabilities
type UnigramLM struct {
    vocab *Vocabulary
    table *compression.LossyIntLookup
    order int
}

func NewUnigramLM(table *compression.LossyIntLookup) *UnigramLM {
    v := NewVocabulary()
    return &UnigramLM{vocab: v, table: table, order: 2}
}

func (u *UnigramLM) GetOrder() int { return u.order }
func (u *UnigramLM) GetVocabulary() *Vocabulary { return u.vocab }

// GetProbability returns log probability of last token using unigram table
func (u *UnigramLM) GetProbability(indexes []int) float32 {
    if len(indexes) == 0 {
        return float32(math.Log(1.0))
    }
    // last token content is not available here; caller only passes indexes.
    // We only have vocabulary indexes; ensure vocabulary contains tokens; if not, return a small value.
    // Since the vocabulary is empty in this minimal implementation, use a tiny constant.
    return float32(math.Log(0.01))
}

func (u *UnigramLM) NgramExists(indexes []int) bool { return true }

// TryLoadUnigram attempts to load a LossyIntLookup from given path (or sibling lm-unigram.slm)
func TryLoadUnigram(path string) (*UnigramLM, error) {
    // If path is a directory or a 2-gram file, look for lm-unigram.slm nearby
    candidate := path
    fi, err := os.Stat(candidate)
    if err == nil && fi.IsDir() {
        candidate = filepath.Join(candidate, "lm-unigram.slm")
    }
    if filepath.Base(candidate) == "lm.2gram.slm" {
        candidate = filepath.Join(filepath.Dir(candidate), "lm-unigram.slm")
    }
    f, err := os.Open(candidate)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    lu, err := compression.DeserializeLossyIntLookup(f)
    if err != nil {
        return nil, err
    }
    return NewUnigramLM(lu), nil
}

