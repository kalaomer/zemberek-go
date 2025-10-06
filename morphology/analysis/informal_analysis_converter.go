package analysis

import (
	"github.com/kalaomer/zemberek-go/morphology/generator"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

// InformalAnalysisConverter converts informal morphemes to formal ones
type InformalAnalysisConverter struct {
	Generator *generator.WordGenerator
}

// NewInformalAnalysisConverter creates a new converter
func NewInformalAnalysisConverter(gen *generator.WordGenerator) *InformalAnalysisConverter {
	return &InformalAnalysisConverter{
		Generator: gen,
	}
}

// Convert converts informal analysis to formal
func (iac *InformalAnalysisConverter) Convert(input string, analysis *SingleAnalysis) *generator.Result {
	if !analysis.ContainsInformalMorpheme() {
		return &generator.Result{
			Surface:  input,
			Analysis: analysis,
		}
	}

	formalMorphemes := iac.ToFormalMorphemeNames(analysis)
	generations := iac.Generator.Generate(analysis.Item, formalMorphemes)

	if len(generations) > 0 {
		return generations[0]
	}

	return nil
}

// ToFormalMorphemeNames converts informal morphemes to formal
func (iac *InformalAnalysisConverter) ToFormalMorphemeNames(analysis *SingleAnalysis) []*morphotactics.Morpheme {
	transform := make([]*morphotactics.Morpheme, 0)

	for _, m := range analysis.GetMorphemes() {
		if m.Informal && m.MappedMorpheme != nil {
			transform = append(transform, m.MappedMorpheme)
		} else {
			transform = append(transform, m)
		}
	}

	return transform
}
