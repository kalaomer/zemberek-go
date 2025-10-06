package analysis

import (
	"fmt"
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

// SingleAnalysis represents a single morphological analysis
type SingleAnalysis struct {
	Item             *lexicon.DictionaryItem
	MorphemeDataList []*MorphemeData
	GroupBoundaries  []int
	hash             int
}

// MorphemeData represents morpheme and surface data
type MorphemeData struct {
	Morpheme *morphotactics.Morpheme
	Surface  string
}

// NewMorphemeData creates a new MorphemeData
func NewMorphemeData(morpheme *morphotactics.Morpheme, surface string) *MorphemeData {
	return &MorphemeData{
		Morpheme: morpheme,
		Surface:  surface,
	}
}

// String returns string representation
func (md *MorphemeData) String() string {
	if len(md.Surface) == 0 {
		return md.Morpheme.ID
	}
	return md.Surface + ":" + md.Morpheme.ID
}

// NewSingleAnalysis creates a new SingleAnalysis
func NewSingleAnalysis(item *lexicon.DictionaryItem, morphemeDataList []*MorphemeData, groupBoundaries []int) *SingleAnalysis {
	return &SingleAnalysis{
		Item:             item,
		MorphemeDataList: morphemeDataList,
		GroupBoundaries:  groupBoundaries,
	}
}

// Unknown creates an unknown analysis
func Unknown(input string) *SingleAnalysis {
	morphemeData := NewMorphemeData(morphotactics.UnknownMorpheme, input)
	return NewSingleAnalysis(lexicon.Unknown, []*MorphemeData{morphemeData}, []int{0})
}

// FromSearchPath creates a SingleAnalysis from a SearchPath
func FromSearchPath(searchPath *SearchPath) *SingleAnalysis {
	morphemes := make([]*MorphemeData, 0)
	derivationCount := 0

	for _, transition := range searchPath.Transitions {
		morpheme := transition.GetMorpheme()
		if morpheme == nil {
			continue
		}

		// Skip Pnon and Nom morphemes (visual noise, like Java does)
		// See: SingleAnalysis.fromSearchPath() in Java (lines 15-17)
		if morpheme.ID == "Pnon" || morpheme.ID == "Nom" {
			continue
		}

		// Count derivations AFTER skipping
		if transition.IsDerivative() {
			derivationCount++
		}

		if len(transition.Surface) == 0 {
			morphemeData := NewMorphemeData(morpheme, "")
			morphemes = append(morphemes, morphemeData)
		} else {
			morphemeData := NewMorphemeData(morpheme, transition.Surface)
			morphemes = append(morphemes, morphemeData)
		}
	}

	groupBoundaries := make([]int, derivationCount+1)
	morphemeCounter := 0
	derivationCounter := 1

	for _, morphemeData := range morphemes {
		if morphemeData.Morpheme.Derivational {
			groupBoundaries[derivationCounter] = morphemeCounter
			derivationCounter++
		}
		morphemeCounter++
	}

	item := searchPath.GetDictionaryItem()
	if item != nil && item.HasAttribute(turkish.Dummy) && item.ReferenceItem != nil {
		item = item.ReferenceItem
	}

	return NewSingleAnalysis(item, morphemes, groupBoundaries)
}

// IsUnknown returns true if this is an unknown analysis
func (sa *SingleAnalysis) IsUnknown() bool {
	return sa.Item.IsUnknown()
}

// GetStem returns the stem
func (sa *SingleAnalysis) GetStem() string {
	if len(sa.MorphemeDataList) == 0 {
		return ""
	}
	return sa.MorphemeDataList[0].Surface
}

// GetEnding returns the ending
func (sa *SingleAnalysis) GetEnding() string {
	if len(sa.MorphemeDataList) <= 1 {
		return ""
	}

	var sb strings.Builder
	for _, md := range sa.MorphemeDataList[1:] {
		sb.WriteString(md.Surface)
	}
	return sb.String()
}

// SurfaceForm returns the complete surface form
func (sa *SingleAnalysis) SurfaceForm() string {
	return sa.GetStem() + sa.GetEnding()
}

// String returns string representation
func (sa *SingleAnalysis) String() string {
	return sa.FormatString()
}

// FormatString formats the analysis as a string
func (sa *SingleAnalysis) FormatString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s:%s", sa.Item.Lemma, sa.Item.PrimaryPos.GetStringForm()))

	if sa.Item.SecondaryPos != turkish.NonePos {
		sb.WriteString(", " + sa.Item.SecondaryPos.GetStringForm())
	}

	sb.WriteString("] ")

	// Format morphemes
	surfaces := sa.MorphemeDataList
	sb.WriteString(sa.GetStem() + ":" + surfaces[0].Morpheme.ID)

	if len(surfaces) > 1 && !surfaces[1].Morpheme.Derivational {
		sb.WriteString("+")
	}

	for i := 1; i < len(surfaces); i++ {
		s := surfaces[i]
		morpheme := s.Morpheme

		if morpheme.Derivational {
			sb.WriteString("|")
		}

		if len(s.Surface) > 0 {
			sb.WriteString(s.Surface + ":")
		}

		sb.WriteString(s.Morpheme.ID)

		if morpheme.Derivational {
			sb.WriteString("→")
		} else if i < len(surfaces)-1 && !surfaces[i+1].Morpheme.Derivational {
			sb.WriteString("+")
		}
	}

	return sb.String()
}

// ContainsMorpheme checks if the analysis contains a specific morpheme
func (sa *SingleAnalysis) ContainsMorpheme(morpheme *morphotactics.Morpheme) bool {
	for _, md := range sa.MorphemeDataList {
		if md.Morpheme == morpheme {
			return true
		}
	}
	return false
}

// ContainsInformalMorpheme checks if analysis contains informal morphemes
func (sa *SingleAnalysis) ContainsInformalMorpheme() bool {
	for _, md := range sa.MorphemeDataList {
		if md.Morpheme.Informal {
			return true
		}
	}
	return false
}

// GetMorphemes returns all morphemes in the analysis
func (sa *SingleAnalysis) GetMorphemes() []*morphotactics.Morpheme {
	morphemes := make([]*morphotactics.Morpheme, len(sa.MorphemeDataList))
	for i, md := range sa.MorphemeDataList {
		morphemes[i] = md.Morpheme
	}
	return morphemes
}

// IsRuntime checks if this is a runtime-generated analysis
func (sa *SingleAnalysis) IsRuntime() bool {
	// Runtime analyses are typically generated on-the-fly
	// Would check: sa.Item.HasAttribute(turkish.Runtime)
	return false
}
