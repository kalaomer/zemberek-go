package analysis

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

// RuleBasedAnalyzer performs morphological analysis using rules
type RuleBasedAnalyzer struct {
	Lexicon         *lexicon.RootLexicon
	StemTransitions *morphotactics.StemTransitionsMapBased
	Morphotactics   *morphotactics.TurkishMorphotactics
	DebugMode       bool
	ASCIITolerant   bool
}

// NewRuleBasedAnalyzer creates a new analyzer
func NewRuleBasedAnalyzer(morph *morphotactics.TurkishMorphotactics) *RuleBasedAnalyzer {
	return &RuleBasedAnalyzer{
		Lexicon:         morph.GetRootLexicon(),
		StemTransitions: morph.GetStemTransitions(),
		Morphotactics:   morph,
		DebugMode:       false,
		ASCIITolerant:   false,
	}
}

// NewIgnoreDiacriticsAnalyzer creates an analyzer that ignores diacritics
func NewIgnoreDiacriticsAnalyzer(morph *morphotactics.TurkishMorphotactics) *RuleBasedAnalyzer {
	analyzer := NewRuleBasedAnalyzer(morph)
	analyzer.ASCIITolerant = true
	return analyzer
}

// Analyze analyzes a word
func (rba *RuleBasedAnalyzer) Analyze(input string) []*SingleAnalysis {
	// Get stem candidates
	candidates := rba.StemTransitions.GetPrefixMatches(input, rba.ASCIITolerant)

	// Generate initial search paths
	paths := make([]*SearchPath, 0, len(candidates))
	for _, candidate := range candidates {
		length := len(candidate.Surface)
		tail := ""
		if length < len(input) {
			tail = input[length:]
		}
		paths = append(paths, InitialPath(candidate, tail))
	}

	// Search graph
	resultPaths := rba.Search(paths)

	// Generate results from successful paths
	result := make([]*SingleAnalysis, 0, len(resultPaths))
	for _, path := range resultPaths {
		analysis := FromSearchPath(path)
		result = append(result, analysis)
	}

	return result
}

// Search performs graph search for analysis
func (rba *RuleBasedAnalyzer) Search(currentPaths []*SearchPath) []*SearchPath {
	if len(currentPaths) > 30 {
		currentPaths = rba.PruneCyclicPaths(currentPaths)
	}

	result := make([]*SearchPath, 0)

	for len(currentPaths) > 0 {
		allNewPaths := make([]*SearchPath, 0)

		for _, path := range currentPaths {
			// If tail is empty and path is terminal, add to results
			if len(path.Tail) == 0 {
				if path.Terminal && !path.PhoneticAttributes[turkish.CannotTerminate] {
					result = append(result, path)
					continue
				}
			}

			newPaths := rba.Advance(path)
			allNewPaths = append(allNewPaths, newPaths...)
		}

		currentPaths = allNewPaths
	}

	return result
}

// Advance advances a search path
func (rba *RuleBasedAnalyzer) Advance(path *SearchPath) []*SearchPath {
	newPaths := make([]*SearchPath, 0, 2)

	if rba.DebugMode {
		fmt.Printf("\n[ADVANCE] State=%s, Tail='%s', Outgoing=%d\n",
			path.CurrentState.ID, path.Tail, len(path.CurrentState.Outgoing))
	}

	// For all outgoing transitions
	for _, transition := range path.CurrentState.Outgoing {
		suffixTransition, ok := transition.(*morphotactics.SuffixTransition)
		if !ok {
			continue
		}

		if rba.DebugMode {
			fmt.Printf("  [TRANSITION] %s -> %s\n",
				path.CurrentState.ID, suffixTransition.To.ID)
		}

		// If tail is empty and this transition has surface, skip
		if len(path.Tail) == 0 && suffixTransition.HasSurfaceForm() {
			if rba.DebugMode {
				fmt.Printf("    [SKIP] Tail empty but has surface form\n")
			}
			continue
		}

		// Generate surface form
		surface := GenerateSurface(suffixTransition, path.PhoneticAttributes)

		if rba.DebugMode {
			fmt.Printf("    [SURFACE] '%s'\n", surface)
		}

		// Check if tail starts with surface
		tailStartsWith := false
		if rba.ASCIITolerant {
			tailStartsWith = turkish.Instance.StartsWithIgnoreDiacritics(path.Tail, surface)
		} else {
			tailStartsWith = len(path.Tail) >= len(surface) && path.Tail[:len(surface)] == surface
		}

		if !tailStartsWith {
			if rba.DebugMode {
				fmt.Printf("    [SKIP] Tail doesn't start with surface\n")
			}
			continue
		}

		if rba.DebugMode {
			fmt.Printf("    [MATCH] Tail starts with surface\n")
		}

		// Check conditions
		canPass := suffixTransition.CanPass(path)
		if rba.DebugMode {
			if suffixTransition.Condition != nil {
				fmt.Printf("    [CONDITION] Type=%T Ptr=%p\n", suffixTransition.Condition, suffixTransition.Condition)
				// Check if it's ExpectsConsonant/ExpectsVowel check
				ec := turkish.ExpectsConsonant
				ev := turkish.ExpectsVowel
				hasEC := path.PhoneticAttributes[ec]
				hasEV := path.PhoneticAttributes[ev]
				fmt.Printf("    [CONDITION_DEBUG] ExpectsConsonant=%v ExpectsVowel=%v\n", hasEC, hasEV)

				// Manual condition test
				manualResult := suffixTransition.Condition.Accept(path)
				fmt.Printf("    [CONDITION_MANUAL] Accept()=%v\n", manualResult)
			}
			fmt.Printf("    [ATTRIBUTES] ")
			for attr := range path.PhoneticAttributes {
				fmt.Printf("%s(%d) ", attr.GetStringForm(), attr)
			}
			fmt.Printf("→CanPass=%v\n", canPass)
		}
		if !canPass {
			if rba.DebugMode {
				fmt.Printf("    [SKIP] Condition check failed\n")
			}
			continue
		}

		if rba.DebugMode {
			fmt.Printf("    [PASS] Condition check passed\n")
		}

		// Epsilon (empty) transition - use existing attributes
		if !suffixTransition.HasSurfaceForm() {
			if rba.DebugMode {
				fmt.Printf("    [EPSILON] Preserving attributes: ")
				for attr := range path.PhoneticAttributes {
					fmt.Printf("%v ", attr)
				}
				fmt.Println()
			}
			newPaths = append(newPaths, path.GetCopy(
				NewSurfaceTransition("", suffixTransition),
				path.PhoneticAttributes))
			continue
		}

		surfaceTransition := NewSurfaceTransition(surface, suffixTransition)

		// If tail equals surface, no need to recalculate attributes
		var attributes map[turkish.PhoneticAttribute]bool
		tailEqualsSurface := false
		if rba.ASCIITolerant {
			tailEqualsSurface = turkish.Instance.EqualsIgnoreDiacritics(path.Tail, surface)
		} else {
			tailEqualsSurface = path.Tail == surface
		}

		if tailEqualsSurface {
			// Copy attributes
			attributes = make(map[turkish.PhoneticAttribute]bool)
			for k, v := range path.PhoneticAttributes {
				attributes[k] = v
			}
		} else {
			attributes = GetMorphemicAttributes(surface, path.PhoneticAttributes)
		}

		// Remove CannotTerminate
		delete(attributes, turkish.CannotTerminate)

		// Handle last token types
		lastToken := suffixTransition.GetLastTemplateToken()
		if lastToken != nil {
			if lastToken.Type == morphotactics.LAST_VOICED {
				attributes[turkish.ExpectsConsonant] = true
			} else if lastToken.Type == morphotactics.LAST_NOT_VOICED {
				attributes[turkish.ExpectsVowel] = true
				attributes[turkish.CannotTerminate] = true
			}
		}

		p := path.GetCopy(surfaceTransition, attributes)
		newPaths = append(newPaths, p)
	}

	return newPaths
}

// PruneCyclicPaths removes paths with too many repetitions
func (rba *RuleBasedAnalyzer) PruneCyclicPaths(paths []*SearchPath) []*SearchPath {
	result := make([]*SearchPath, 0)

	for _, path := range paths {
		remove := false
		typeCounts := make(map[string]int)

		for _, transition := range path.Transitions {
			state := transition.GetState()
			if state == nil {
				continue
			}

			typeCounts[state.ID]++
			if typeCounts[state.ID] > 3 {
				remove = true
				break
			}
		}

		if !remove {
			result = append(result, path)
		}
	}

	return result
}
