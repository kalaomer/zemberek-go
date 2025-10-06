package normalization

import (
	"sync"
	"sync/atomic"
)

// AtomicCounter provides thread-safe counter
type AtomicCounter struct {
	count int64
}

// NewAtomicCounter creates a new atomic counter
func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{count: 0}
}

// GetAndIncrement atomically gets current value and increments
func (ac *AtomicCounter) GetAndIncrement() int {
	return int(atomic.AddInt64(&ac.count, 1) - 1)
}

// CharacterGraph represents a graph structure for character-based operations
type CharacterGraph struct {
	Root             *Node
	nodeIndexCounter *AtomicCounter
	mu               sync.RWMutex
}

// NewCharacterGraph creates a new character graph
func NewCharacterGraph() *CharacterGraph {
	counter := NewAtomicCounter()
	return &CharacterGraph{
		Root:             NewNode(counter.GetAndIncrement(), '\u0000', TypeGraphRoot, ""),
		nodeIndexCounter: counter,
	}
}

// AddWord adds a word to the graph with given type
func (cg *CharacterGraph) AddWord(word string, nodeType NodeType) *Node {
	if word == "" {
		return nil
	}
	runes := []rune(word)
	return cg.add(cg.Root, 0, runes, nodeType)
}

// add recursively adds word to graph
func (cg *CharacterGraph) add(currentNode *Node, index int, word []rune, nodeType NodeType) *Node {
	c := word[index]

	if index == len(word)-1 {
		// Last character - create terminal node with word
		return currentNode.AddChild(cg.nodeIndexCounter.GetAndIncrement(), c, nodeType, string(word))
	}

	// Intermediate character - create non-terminal node
	child := currentNode.AddChild(cg.nodeIndexCounter.GetAndIncrement(), c, TypeEmpty, "")
	return cg.add(child, index+1, word, nodeType)
}

// GetAllNodes returns all nodes in the graph that have words
func (cg *CharacterGraph) GetAllNodes() []*Node {
	nodes := make([]*Node, 0)
	visited := make(map[*Node]bool)
	cg.walk(cg.Root, &nodes, visited)
	return nodes
}

// walk recursively walks the graph collecting nodes with words
func (cg *CharacterGraph) walk(current *Node, nodes *[]*Node, visited map[*Node]bool) {
	if visited[current] {
		return
	}
	visited[current] = true

	if current.Word != "" {
		*nodes = append(*nodes, current)
	}

	for _, node := range current.GetImmediateChildNodes() {
		cg.walk(node, nodes, visited)
	}
}

// GetNodeCount returns total number of nodes in graph
func (cg *CharacterGraph) GetNodeCount() int {
	visited := make(map[*Node]bool)
	return cg.countNodes(cg.Root, visited)
}

// countNodes recursively counts all nodes
func (cg *CharacterGraph) countNodes(current *Node, visited map[*Node]bool) int {
	if visited[current] {
		return 0
	}
	visited[current] = true

	count := 1
	for _, node := range current.GetImmediateChildNodes() {
		count += cg.countNodes(node, visited)
	}
	return count
}

// ContainsWord checks if graph contains a word
func (cg *CharacterGraph) ContainsWord(word string) bool {
	if word == "" {
		return false
	}

	current := cg.Root
	for _, c := range word {
		child := current.GetImmediateChild(c)
		if child == nil {
			return false
		}
		current = child
	}

	return current.Word == word
}

// GetNode returns the node corresponding to the given word, or nil if not found
func (cg *CharacterGraph) GetNode(word string) *Node {
	if word == "" {
		return nil
	}

	current := cg.Root
	for _, c := range word {
		child := current.GetImmediateChild(c)
		if child == nil {
			return nil
		}
		current = child
	}

	if current.Word == word {
		return current
	}
	return nil
}
