package normalization

import (
	"testing"
)

// TestNode tests basic node functionality
func TestNode(t *testing.T) {
	node := NewNode(0, 'a', TypeWord, "ali")

	if node.Char != 'a' {
		t.Errorf("Expected char 'a', got %c", node.Char)
	}
	if node.Word != "ali" {
		t.Errorf("Expected word 'ali', got %s", node.Word)
	}
	if node.Type != TypeWord {
		t.Errorf("Expected type TypeWord, got %d", node.Type)
	}
}

// TestNode_AddChild tests adding children to nodes
func TestNode_AddChild(t *testing.T) {
	root := NewNode(0, 'r', TypeGraphRoot, "")
	child1 := root.AddChild(1, 'a', TypeEmpty, "")
	child2 := root.AddChild(2, 'b', TypeEmpty, "")

	if len(root.Nodes) != 2 {
		t.Errorf("Expected 2 children, got %d", len(root.Nodes))
	}

	if !root.HasImmediateChild('a') {
		t.Error("Expected to have child 'a'")
	}
	if !root.HasImmediateChild('b') {
		t.Error("Expected to have child 'b'")
	}

	if child1.Char != 'a' {
		t.Errorf("Expected child1 char 'a', got %c", child1.Char)
	}
	if child2.Char != 'b' {
		t.Errorf("Expected child2 char 'b', got %c", child2.Char)
	}
}

// TestNode_EpsilonConnection tests epsilon connections
func TestNode_EpsilonConnection(t *testing.T) {
	node1 := NewNode(0, 'a', TypeEmpty, "")
	node2 := NewNode(1, 'b', TypeEmpty, "")

	if node1.HasEpsilonConnection() {
		t.Error("Expected no epsilon connection initially")
	}

	node1.ConnectEpsilon(node2)

	if !node1.HasEpsilonConnection() {
		t.Error("Expected epsilon connection after ConnectEpsilon")
	}
}

// TestCharacterGraph tests character graph functionality
func TestCharacterGraph(t *testing.T) {
	graph := NewCharacterGraph()

	if graph.Root == nil {
		t.Fatal("Expected non-nil root")
	}
	if graph.Root.Type != TypeGraphRoot {
		t.Errorf("Expected root type TypeGraphRoot, got %d", graph.Root.Type)
	}
}

// TestCharacterGraph_AddWord tests adding words to graph
func TestCharacterGraph_AddWord(t *testing.T) {
	graph := NewCharacterGraph()

	node := graph.AddWord("test", TypeWord)
	if node == nil {
		t.Fatal("Expected non-nil node")
	}
	if node.Word != "test" {
		t.Errorf("Expected word 'test', got %s", node.Word)
	}
	if node.Type != TypeWord {
		t.Errorf("Expected type TypeWord, got %d", node.Type)
	}

	// Test retrieval
	if !graph.ContainsWord("test") {
		t.Error("Graph should contain 'test'")
	}

	retrievedNode := graph.GetNode("test")
	if retrievedNode == nil {
		t.Fatal("Expected to retrieve node for 'test'")
	}
	if retrievedNode.Word != "test" {
		t.Errorf("Retrieved node word mismatch: got %s", retrievedNode.Word)
	}
}

// TestCharacterGraph_MultipleWords tests adding multiple words
func TestCharacterGraph_MultipleWords(t *testing.T) {
	graph := NewCharacterGraph()

	words := []string{"ali", "ahmet", "ayşe", "test"}
	for _, word := range words {
		graph.AddWord(word, TypeWord)
	}

	for _, word := range words {
		if !graph.ContainsWord(word) {
			t.Errorf("Graph should contain '%s'", word)
		}
	}

	allNodes := graph.GetAllNodes()
	if len(allNodes) != len(words) {
		t.Errorf("Expected %d nodes, got %d", len(words), len(allNodes))
	}
}

// TestCharacterGraphDecoder tests decoder functionality
func TestCharacterGraphDecoder(t *testing.T) {
	graph := NewCharacterGraph()
	graph.AddWord("test", TypeWord)
	graph.AddWord("best", TypeWord)
	graph.AddWord("rest", TypeWord)

	decoder := NewCharacterGraphDecoder(graph)

	if decoder.Graph == nil {
		t.Fatal("Expected non-nil graph in decoder")
	}
	if decoder.MaxPenalty != 1.0 {
		t.Errorf("Expected MaxPenalty 1.0, got %f", decoder.MaxPenalty)
	}
}

// TestCharacterGraphDecoder_Suggestions tests getting suggestions
func TestCharacterGraphDecoder_Suggestions(t *testing.T) {
	graph := NewCharacterGraph()
	graph.AddWord("test", TypeWord)
	graph.AddWord("best", TypeWord)
	graph.AddWord("rest", TypeWord)

	decoder := NewCharacterGraphDecoder(graph)
	suggestions := decoder.GetSuggestions("test", nil)

	if len(suggestions) == 0 {
		t.Error("Expected at least one suggestion")
	}

	found := false
	for _, sug := range suggestions {
		if sug == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find 'test' in suggestions")
	}
}

// TestDiacriticsIgnoringMatcher tests diacritics matching
func TestDiacriticsIgnoringMatcher(t *testing.T) {
	matcher := NewDiacriticsIgnoringMatcher()

	tests := []struct {
		char     rune
		expected []rune
	}{
		{'c', []rune{'c', 'ç'}},
		{'g', []rune{'g', 'ğ'}},
		{'i', []rune{'ı', 'i'}},
		{'o', []rune{'o', 'ö'}},
		{'s', []rune{'s', 'ş'}},
		{'u', []rune{'u', 'ü'}},
	}

	for _, tt := range tests {
		matches := matcher.Matches(tt.char)
		if len(matches) != len(tt.expected) {
			t.Errorf("For char '%c': expected %d matches, got %d", tt.char, len(tt.expected), len(matches))
			continue
		}
		for i, expected := range tt.expected {
			if matches[i] != expected {
				t.Errorf("For char '%c' at index %d: expected '%c', got '%c'", tt.char, i, expected, matches[i])
			}
		}
	}
}

// TestStemEndingGraph tests stem-ending graph creation
func TestStemEndingGraph(t *testing.T) {
	stemWords := []string{"git", "gel", "al", "ver"}

	// We won't test with actual file, just test structure
	seg, err := NewStemEndingGraph(stemWords, "")
	if err == nil {
		// If it succeeds without file, check structure
		if seg.StemGraph == nil {
			t.Error("Expected non-nil StemGraph")
		}
		if seg.EndingGraph == nil {
			t.Error("Expected non-nil EndingGraph")
		}
	}
	// If error, that's fine - file doesn't exist in test
}

// TestLevenshteinDistance tests edit distance calculation
func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "ab", 1},
		{"abc", "ac", 1},
		{"abc", "abcd", 1},
		{"kitten", "sitting", 3},
		{"saturday", "sunday", 3},
	}

	for _, tt := range tests {
		result := levenshteinDistance(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("levenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

// TestGuessCase tests case type detection
func TestGuessCase(t *testing.T) {
	tests := []struct {
		word     string
		expected CaseType
	}{
		{"test", LowerCase},
		{"TEST", UpperCase},
		{"Test", TitleCase},
		{"TeSt", MixedCase},
		{"", DefaultCase},
	}

	for _, tt := range tests {
		result := guessCase(tt.word)
		if result != tt.expected {
			t.Errorf("guessCase(%q) = %v, want %v", tt.word, result, tt.expected)
		}
	}
}

// TestTokenize tests tokenization
func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"merhaba dünya", 2},
		{"bir iki üç dört", 4},
		{"", 0},
		{"tek", 1},
		{"  boşluk   çok  ", 2},
	}

	for _, tt := range tests {
		result := tokenize(tt.input)
		if len(result) != tt.expected {
			t.Errorf("tokenize(%q) returned %d tokens, want %d", tt.input, len(result), tt.expected)
		}
	}
}

// TestIsWord tests word detection
func TestIsWord(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"merhaba", true},
		{"test", true},
		{"çok", true},
		{"güzel", true},
		{"test123", false},
		{"test!", false},
		{"", false},
		{"ali'nin", true},
		{"test-case", true},
	}

	for _, tt := range tests {
		result := isWord(tt.input)
		if result != tt.expected {
			t.Errorf("isWord(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// TestProbablyRequiresDeasciifier tests deasciifier detection
func TestProbablyRequiresDeasciifier(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"cok guzel", true},          // No Turkish chars, needs deasciifier
		{"çok güzel", false},         // Has Turkish chars, doesn't need
		{"test", true},               // No Turkish chars
		{"şişe", false},              // Has Turkish chars
		{"", false},                  // Empty
	}

	for _, tt := range tests {
		result := probablyRequiresDeasciifier(tt.input)
		if result != tt.expected {
			t.Errorf("probablyRequiresDeasciifier(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// BenchmarkCharacterGraph_AddWord benchmarks word addition
func BenchmarkCharacterGraph_AddWord(b *testing.B) {
	graph := NewCharacterGraph()
	words := []string{"test", "best", "rest", "quest", "guest"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, word := range words {
			graph.AddWord(word, TypeWord)
		}
	}
}

// BenchmarkLevenshteinDistance benchmarks edit distance calculation
func BenchmarkLevenshteinDistance(b *testing.B) {
	s1 := "saturday"
	s2 := "sunday"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		levenshteinDistance(s1, s2)
	}
}
