package tokenization

// TokenizerBuilder builds a TurkishTokenizer with custom token type filtering
// Matches Java's TurkishTokenizer.Builder pattern
type TokenizerBuilder struct {
	acceptedTypes map[TokenType]bool
}

// NewBuilder creates a new tokenizer builder
func NewBuilder() *TokenizerBuilder {
	return &TokenizerBuilder{
		acceptedTypes: make(map[TokenType]bool),
	}
}

// AcceptAll marks all token types as accepted
func (b *TokenizerBuilder) AcceptAll() *TokenizerBuilder {
	// Add all token types
	for i := SpaceTab; i <= Unknown; i++ {
		b.acceptedTypes[i] = true
	}
	return b
}

// IgnoreAll marks all token types as ignored
func (b *TokenizerBuilder) IgnoreAll() *TokenizerBuilder {
	b.acceptedTypes = make(map[TokenType]bool)
	return b
}

// AcceptTypes marks specific token types as accepted
func (b *TokenizerBuilder) AcceptTypes(types ...TokenType) *TokenizerBuilder {
	for _, t := range types {
		b.acceptedTypes[t] = true
	}
	return b
}

// IgnoreTypes marks specific token types as ignored
func (b *TokenizerBuilder) IgnoreTypes(types ...TokenType) *TokenizerBuilder {
	for _, t := range types {
		delete(b.acceptedTypes, t)
	}
	return b
}

// Build constructs the TurkishTokenizer with configured settings
func (b *TokenizerBuilder) Build() *TurkishTokenizer {
	// Copy the accepted types map to avoid sharing state
	acceptedTypes := make(map[TokenType]bool, len(b.acceptedTypes))
	for k, v := range b.acceptedTypes {
		acceptedTypes[k] = v
	}

	return &TurkishTokenizer{
		acceptedTypes: acceptedTypes,
	}
}
