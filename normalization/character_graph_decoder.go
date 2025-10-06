package normalization

import (
	"math"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

// Operation represents the type of edit operation
type Operation int

const (
	NoError Operation = iota
	Insertion
	Deletion
	Substitution
	Transposition
	NA
)

// CharMatcher interface for character matching strategies
type CharMatcher interface {
	Matches(c rune) []rune
}

// DiacriticsIgnoringMatcher matches characters ignoring diacritics
type DiacriticsIgnoringMatcher struct {
	matchMap map[rune][]rune
}

// NewDiacriticsIgnoringMatcher creates a new diacritics ignoring matcher
func NewDiacriticsIgnoringMatcher() *DiacriticsIgnoringMatcher {
	matcher := &DiacriticsIgnoringMatcher{
		matchMap: make(map[rune][]rune),
	}

	// Initialize all letters
	allLetters := turkish.Instance.AllLetters + "+.,'-"
	for _, c := range allLetters {
		matcher.matchMap[c] = []rune{c}
	}

	// Turkish character mappings
	matcher.matchMap['c'] = []rune{'c', 'ç'}
	matcher.matchMap['g'] = []rune{'g', 'ğ'}
	matcher.matchMap['ı'] = []rune{'ı', 'i'}
	matcher.matchMap['i'] = []rune{'ı', 'i'}
	matcher.matchMap['o'] = []rune{'o', 'ö'}
	matcher.matchMap['s'] = []rune{'s', 'ş'}
	matcher.matchMap['u'] = []rune{'u', 'ü'}
	matcher.matchMap['a'] = []rune{'a', 'â'}
	matcher.matchMap['C'] = []rune{'C', 'Ç'}
	matcher.matchMap['G'] = []rune{'G', 'Ğ'}
	matcher.matchMap['I'] = []rune{'I', 'İ'}
	matcher.matchMap['İ'] = []rune{'İ', 'I'}
	matcher.matchMap['O'] = []rune{'O', 'Ö'}
	matcher.matchMap['S'] = []rune{'S', 'Ş'}
	matcher.matchMap['U'] = []rune{'U', 'Ü'}
	matcher.matchMap['A'] = []rune{'A', 'Â'}

	return matcher
}

// Matches returns possible character matches for given character
func (dim *DiacriticsIgnoringMatcher) Matches(c rune) []rune {
	if matches, exists := dim.matchMap[c]; exists {
		return matches
	}
	return []rune{c}
}

// CharacterGraphDecoder decodes strings using character graph with error tolerance
type CharacterGraphDecoder struct {
	Graph                      *CharacterGraph
	MaxPenalty                 float64
	CheckNearKeySubstitution   bool
}

// NewCharacterGraphDecoder creates a new decoder
func NewCharacterGraphDecoder(graph *CharacterGraph) *CharacterGraphDecoder {
	return &CharacterGraphDecoder{
		Graph:                      graph,
		MaxPenalty:                 1.0,
		CheckNearKeySubstitution:   false,
	}
}

// GetSuggestions returns suggestions for input string using given matcher
func (cgd *CharacterGraphDecoder) GetSuggestions(input string, matcher CharMatcher) []string {
	decoder := newDecoder(matcher, cgd)
	finished := decoder.decode(input)

	suggestions := make([]string, 0, len(finished))
	for word := range finished {
		suggestions = append(suggestions, word)
	}
	return suggestions
}

// decoder handles the decoding process
type decoder struct {
	finished map[string]float64
	matcher  CharMatcher
	outer    *CharacterGraphDecoder
}

// newDecoder creates a new decoder instance
func newDecoder(matcher CharMatcher, outer *CharacterGraphDecoder) *decoder {
	return &decoder{
		finished: make(map[string]float64),
		matcher:  matcher,
		outer:    outer,
	}
}

// decode performs the decoding
func (d *decoder) decode(input string) map[string]float64 {
	runes := []rune(input)
	hyp := newHypothesis(nil, d.outer.Graph.Root, 0.0, NA, "", "", -1)
	next := d.expand(hyp, runes)

	for {
		newHyps := make(map[*hypothesis]bool)
		for h := range next {
			expanded := d.expand(h, runes)
			for expandedHyp := range expanded {
				newHyps[expandedHyp] = true
			}
		}

		if len(newHyps) == 0 {
			return d.finished
		}
		next = newHyps
	}
}

// expand expands a hypothesis
func (d *decoder) expand(h *hypothesis, input []rune) map[*hypothesis]bool {
	newHypotheses := make(map[*hypothesis]bool)
	nextIndex := h.charIndex + 1

	var nextChar rune
	if nextIndex < len(input) {
		nextChar = input[nextIndex]
	} else {
		nextChar = 0
	}

	if nextIndex < len(input) {
		var cc []rune
		if d.matcher != nil {
			cc = d.matcher.Matches(nextChar)
		}

		if h.node.HasEpsilonConnection() {
			var childList []*Node
			if cc == nil {
				childList = h.node.GetChildList(nextChar)
			} else {
				childList = h.node.GetChildListMulti(cc)
			}

			for _, child := range childList {
				newH := h.getNewMoveForward(child, 0.0, NoError)
				newH.setWord(child)
				newHypotheses[newH] = true
				if nextIndex >= len(input)-1 && newH.node.Word != "" {
					d.addHypothesis(newH)
				}
			}
		} else if cc == nil {
			child := h.node.GetImmediateChild(nextChar)
			if child != nil {
				newH := h.getNewMoveForward(child, 0.0, NoError)
				newH.setWord(child)
				newHypotheses[newH] = true
				if nextIndex >= len(input)-1 && newH.node.Word != "" {
					d.addHypothesis(newH)
				}
			}
		} else {
			for _, c := range cc {
				child := h.node.GetImmediateChild(c)
				if child != nil {
					newH := h.getNewMoveForward(child, 0.0, NoError)
					newH.setWord(child)
					newHypotheses[newH] = true
					if nextIndex >= len(input)-1 && newH.node.Word != "" {
						d.addHypothesis(newH)
					}
				}
			}
		}
	} else if h.node.Word != "" {
		d.addHypothesis(h)
	}

	if h.penalty >= d.outer.MaxPenalty {
		return newHypotheses
	}

	var allChildNodes []*Node
	if h.node.HasEpsilonConnection() {
		allChildNodes = h.node.GetAllChildNodes()
	} else {
		allChildNodes = h.node.GetImmediateChildNodeIterable()
	}

	if nextIndex < len(input) {
		// Substitution
		for _, child := range allChildNodes {
			penalty := 1.0
			if penalty > 0.0 && h.penalty+penalty <= d.outer.MaxPenalty {
				newH := h.getNewMoveForward(child, penalty, Substitution)
				newH.setWord(child)
				if nextIndex == len(input)-1 {
					if newH.node.Word != "" {
						d.addHypothesis(newH)
					}
				} else {
					newHypotheses[newH] = true
				}
			}
		}
	}

	if h.penalty+1.0 > d.outer.MaxPenalty {
		return newHypotheses
	}

	// Deletion
	newHypotheses[h.getNewMoveForward(h.node, 1.0, Deletion)] = true

	// Insertion
	for _, child := range allChildNodes {
		newH := h.getNew(child, 1.0, Insertion, -1)
		newH.setWord(child)
		newHypotheses[newH] = true
	}

	// Transposition
	if len(input) > 2 && nextIndex < len(input)-1 {
		transpose := input[nextIndex+1]
		if d.matcher != nil {
			tt := d.matcher.Matches(transpose)
			cc := d.matcher.Matches(nextChar)
			for _, t := range tt {
				nextNodes := h.node.GetChildList(t)
				for _, nextNode := range nextNodes {
					for _, c := range cc {
						if h.node.HasChild(t) && nextNode.HasChild(c) {
							for _, n := range nextNode.GetChildList(c) {
								newH := h.getNew(n, 1.0, Transposition, nextIndex+1)
								newH.setWord(n)
								if nextIndex == len(input)-1 {
									if newH.node.Word != "" {
										d.addHypothesis(newH)
									}
								} else {
									newHypotheses[newH] = true
								}
							}
						}
					}
				}
			}
		} else {
			nextNodes := h.node.GetChildList(transpose)
			for _, nextNode := range nextNodes {
				if h.node.HasChild(transpose) && nextNode.HasChild(nextChar) {
					for _, n := range nextNode.GetChildList(nextChar) {
						newH := h.getNew(n, 1.0, Transposition, nextIndex+1)
						newH.setWord(n)
						if nextIndex == len(input)-1 {
							if newH.node.Word != "" {
								d.addHypothesis(newH)
							}
						} else {
							newHypotheses[newH] = true
						}
					}
				}
			}
		}
	}

	return newHypotheses
}

// addHypothesis adds hypothesis to finished set
func (d *decoder) addHypothesis(h *hypothesis) {
	hypWord := h.getContent()
	if existingPenalty, exists := d.finished[hypWord]; !exists {
		d.finished[hypWord] = h.penalty
	} else if existingPenalty > h.penalty {
		d.finished[hypWord] = h.penalty
	}
}

// hypothesis represents a decoding hypothesis
type hypothesis struct {
	previous  *hypothesis
	node      *Node
	penalty   float64
	operation Operation
	word      string
	ending    string
	charIndex int
}

// newHypothesis creates a new hypothesis
func newHypothesis(previous *hypothesis, node *Node, penalty float64, op Operation, word string, ending string, charIndex int) *hypothesis {
	return &hypothesis{
		previous:  previous,
		node:      node,
		penalty:   penalty,
		operation: op,
		word:      word,
		ending:    ending,
		charIndex: charIndex,
	}
}

// getNew creates a new hypothesis with same char index
func (h *hypothesis) getNew(node *Node, penaltyToAdd float64, op Operation, index int) *hypothesis {
	charIndex := h.charIndex
	if index != -1 {
		charIndex = index
	}
	return newHypothesis(h, node, h.penalty+penaltyToAdd, op, h.word, h.ending, charIndex)
}

// getNewMoveForward creates a new hypothesis moving forward
func (h *hypothesis) getNewMoveForward(node *Node, penaltyToAdd float64, op Operation) *hypothesis {
	return newHypothesis(h, node, h.penalty+penaltyToAdd, op, h.word, h.ending, h.charIndex+1)
}

// getContent returns the content (word + ending)
func (h *hypothesis) getContent() string {
	w := ""
	if h.word != "" {
		w = h.word
	}
	e := ""
	if h.ending != "" {
		e = h.ending
	}
	return w + e
}

// setWord sets word or ending based on node type
func (h *hypothesis) setWord(node *Node) {
	if node.Word != "" {
		if node.Type == TypeWord {
			h.word = node.Word
		} else if node.Type == TypeEnding {
			h.ending = node.Word
		}
	}
}

// Hash returns hash for hypothesis
func (h *hypothesis) Hash() int {
	result := h.charIndex
	result = 31*result + h.node.Hash()
	result = 31*result + int(math.Float64bits(h.penalty))
	if h.word != "" {
		for _, r := range h.word {
			result = 31*result + int(r)
		}
	}
	if h.ending != "" {
		for _, r := range h.ending {
			result = 31*result + int(r)
		}
	}
	return result
}

// Equals checks equality
func (h *hypothesis) Equals(other *hypothesis) bool {
	if h == other {
		return true
	}
	if other == nil {
		return false
	}
	if h.charIndex != other.charIndex {
		return false
	}
	if h.penalty != other.penalty {
		return false
	}
	if !h.node.Equals(other.node) {
		return false
	}
	if h.word != other.word {
		return false
	}
	return h.ending == other.ending
}

// DiacriticsIgnoringMatcherInstance is the singleton instance
var DiacriticsIgnoringMatcherInstance = NewDiacriticsIgnoringMatcher()
