package normalization

import (
	"bufio"
	"os"
	"strings"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology"
)

// StemEndingGraph creates a character graph from stems and endings
type StemEndingGraph struct {
	EndingGraph  *CharacterGraph
	StemGraph    *CharacterGraph
	stemWords    []string // Store stem words directly
}

// NewStemEndingGraph creates a new stem-ending graph
func NewStemEndingGraph(stemWords []string, endingsPath string) (*StemEndingGraph, error) {
	seg := &StemEndingGraph{
		stemWords: stemWords,
	}

	// Load endings (optional)
	endings := []string{}
	if endingsPath != "" {
		loadedEndings, err := seg.loadLinesFromResource(endingsPath)
		if err == nil {
			endings = loadedEndings
		}
		// If error, continue with empty endings (stems only)
	}

	seg.EndingGraph = seg.generateEndingGraph(endings)
	seg.StemGraph = seg.generateStemGraph(stemWords)

	// Connect stem nodes to ending graph via epsilon transitions (if endings exist)
	if len(endings) > 0 {
		stemWordNodes := seg.StemGraph.GetAllNodes()
		for _, node := range stemWordNodes {
			node.ConnectEpsilon(seg.EndingGraph.Root)
		}
	}

	return seg, nil
}

// loadLinesFromResource loads lines from a file
func (seg *StemEndingGraph) loadLinesFromResource(path string) ([]string, error) {
	// If no path provided, use default
	if path == "" {
		// Default path would be set here
		path = "resources/normalization/endings.txt"
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// generateStemGraph creates graph from stem words
func (seg *StemEndingGraph) generateStemGraph(stemWords []string) *CharacterGraph {
	stemGraph := NewCharacterGraph()

	for _, word := range stemWords {
		if word != "" {
			stemGraph.AddWord(word, TypeWord)
		}
	}

	return stemGraph
}

// generateEndingGraph creates graph from endings
func (seg *StemEndingGraph) generateEndingGraph(endings []string) *CharacterGraph {
	graph := NewCharacterGraph()

	for _, ending := range endings {
		if ending != "" {
			graph.AddWord(ending, TypeEnding)
		}
	}

	return graph
}

// GetStemGraph returns the stem graph
func (seg *StemEndingGraph) GetStemGraph() *CharacterGraph {
	return seg.StemGraph
}

// GetEndingGraph returns the ending graph
func (seg *StemEndingGraph) GetEndingGraph() *CharacterGraph {
	return seg.EndingGraph
}

// NewStemEndingGraphFromMorphology creates a stem-ending graph from TurkishMorphology
// This matches Java's implementation: extracting stems from morphology
func NewStemEndingGraphFromMorphology(morph *morphology.TurkishMorphology, endingsPath string) (*StemEndingGraph, error) {
	// Extract stems from lexicon (like Java does from StemTransitions)
	stems := make([]string, 0)

	if morph.Lexicon != nil {
		for _, item := range morph.Lexicon.GetAllItems() {
			if item.Root != "" && item.PrimaryPos != turkish.Punctuation {
				stems = append(stems, item.Root)
				// Also add lemma if different from root
				if item.Lemma != item.Root && item.Lemma != "" {
					stems = append(stems, item.Lemma)
				}
			}
		}
	}

	// Create StemEndingGraph with extracted stems
	return NewStemEndingGraph(stems, endingsPath)
}
