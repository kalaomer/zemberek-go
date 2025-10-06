package normalization

import (
	"fmt"
	"sort"
	"strings"
)

// NodeType represents the type of node in the graph
type NodeType int

const (
	TypeEmpty      NodeType = 0
	TypeWord       NodeType = 1
	TypeEnding     NodeType = 2
	TypeGraphRoot  NodeType = 3
)

// Node represents a node in the character graph
type Node struct {
	Index        int
	Char         rune
	Type         NodeType
	Word         string
	EpsilonNodes []*Node
	Nodes        map[rune]*Node
}

// NewNode creates a new node
func NewNode(index int, char rune, nodeType NodeType, word string) *Node {
	return &Node{
		Index:        index,
		Char:         char,
		Type:         nodeType,
		Word:         word,
		EpsilonNodes: nil,
		Nodes:        make(map[rune]*Node),
	}
}

// String returns string representation of the node
func (n *Node) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	sb.WriteRune(n.Char)

	if len(n.Nodes) > 0 {
		chars := make([]rune, 0, len(n.Nodes))
		for c := range n.Nodes {
			chars = append(chars, c)
		}
		sort.Slice(chars, func(i, j int) bool {
			return chars[i] < chars[j]
		})

		sb.WriteString(" children=")
		charStrings := make([]string, len(chars))
		for i, c := range chars {
			charStrings[i] = string(c)
		}
		sb.WriteString(strings.Join(charStrings, ", "))
	}

	if n.Word != "" {
		sb.WriteString(" word=")
		sb.WriteString(n.Word)
	}
	sb.WriteString("]")
	return sb.String()
}

// HasEpsilonConnection checks if node has epsilon connections
func (n *Node) HasEpsilonConnection() bool {
	return n.EpsilonNodes != nil && len(n.EpsilonNodes) > 0
}

// HasChild checks if node has a child with given character (including epsilon nodes)
func (n *Node) HasChild(c rune) bool {
	if n.HasImmediateChild(c) {
		return true
	}
	if n.EpsilonNodes == nil {
		return false
	}
	for _, node := range n.EpsilonNodes {
		if node.HasImmediateChild(c) {
			return true
		}
	}
	return false
}

// HasImmediateChild checks if node has an immediate child with given character
func (n *Node) HasImmediateChild(c rune) bool {
	_, exists := n.Nodes[c]
	return exists
}

// GetImmediateChild returns immediate child node for given character
func (n *Node) GetImmediateChild(c rune) *Node {
	return n.Nodes[c]
}

// GetImmediateChildNodes returns all immediate child nodes
func (n *Node) GetImmediateChildNodes() []*Node {
	nodes := make([]*Node, 0, len(n.Nodes))
	for _, node := range n.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetImmediateChildNodeIterable returns all immediate child nodes (same as GetImmediateChildNodes)
func (n *Node) GetImmediateChildNodeIterable() []*Node {
	return n.GetImmediateChildNodes()
}

// GetAllChildNodes returns all child nodes including epsilon-connected nodes
func (n *Node) GetAllChildNodes() []*Node {
	if n.EpsilonNodes == nil || len(n.EpsilonNodes) == 0 {
		return n.GetImmediateChildNodes()
	}

	nodeList := n.GetImmediateChildNodes()
	for _, emptyNode := range n.EpsilonNodes {
		for _, node := range emptyNode.Nodes {
			nodeList = append(nodeList, node)
		}
	}
	return nodeList
}

// GetChildList returns list of children matching given character
func (n *Node) GetChildList(c rune) []*Node {
	children := make([]*Node, 0)
	n.addIfChildExists(c, &children)

	if n.EpsilonNodes != nil {
		for _, emptyNode := range n.EpsilonNodes {
			emptyNode.addIfChildExists(c, &children)
		}
	}
	return children
}

// GetChildListMulti returns list of children matching any of given characters
func (n *Node) GetChildListMulti(chars []rune) []*Node {
	children := make([]*Node, 0)
	for _, c := range chars {
		n.addIfChildExists(c, &children)
		if n.EpsilonNodes != nil {
			for _, emptyNode := range n.EpsilonNodes {
				emptyNode.addIfChildExists(c, &children)
			}
		}
	}
	return children
}

// ConnectEpsilon connects this node to another node via epsilon transition
func (n *Node) ConnectEpsilon(node *Node) bool {
	if n.EpsilonNodes == nil {
		n.EpsilonNodes = []*Node{node}
		return true
	}

	for _, existingNode := range n.EpsilonNodes {
		if existingNode == node {
			return false
		}
	}
	n.EpsilonNodes = append(n.EpsilonNodes, node)
	return true
}

// addIfChildExists adds child to list if it exists
func (n *Node) addIfChildExists(c rune, nodeList *[]*Node) {
	if child, exists := n.Nodes[c]; exists {
		*nodeList = append(*nodeList, child)
	}
}

// AddChild adds or updates a child node
func (n *Node) AddChild(index int, c rune, nodeType NodeType, word string) *Node {
	node, exists := n.Nodes[c]
	if word != "" {
		if !exists {
			node = NewNode(index, c, nodeType, word)
			n.Nodes[c] = node
		} else {
			node.Word = word
			node.Type = nodeType
		}
	} else {
		if !exists {
			node = NewNode(index, c, nodeType, "")
			n.Nodes[c] = node
		}
	}
	return node
}

// Hash returns hash code for the node (based on index)
func (n *Node) Hash() int {
	return n.Index
}

// Equals checks equality based on index
func (n *Node) Equals(other *Node) bool {
	if n == other {
		return true
	}
	if other == nil {
		return false
	}
	return n.Index == other.Index
}

// Debug returns detailed debug string
func (n *Node) Debug() string {
	return fmt.Sprintf("Node{Index: %d, Char: %c, Type: %d, Word: %q, Children: %d, Epsilon: %d}",
		n.Index, n.Char, n.Type, n.Word, len(n.Nodes), len(n.EpsilonNodes))
}
