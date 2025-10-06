package morphotactics

// TemplateTokenType represents token types in suffix templates
type TemplateTokenType int

const (
	LETTER TemplateTokenType = iota
	A_WOVEL
	I_WOVEL
	LAST_VOICED
	LAST_NOT_VOICED
	DEVOICE_FIRST
	APPEND
	BUFFER_LETTER
)

// SuffixTemplateToken represents a token in suffix template
type SuffixTemplateToken struct {
	Type     TemplateTokenType
	Value    rune
	Optional bool
}

// SuffixTemplateTokenizer tokenizes suffix templates
type SuffixTemplateTokenizer struct {
	template []rune
	pos      int
}

// NewSuffixTemplateTokenizer creates a new tokenizer
func NewSuffixTemplateTokenizer(template string) *SuffixTemplateTokenizer {
	return &SuffixTemplateTokenizer{
		template: []rune(template),
		pos:      0,
	}
}

// HasNext checks if there are more tokens
func (stt *SuffixTemplateTokenizer) HasNext() bool {
	return stt.pos < len(stt.template)
}

// Next returns the next token
func (stt *SuffixTemplateTokenizer) Next() *SuffixTemplateToken {
	if !stt.HasNext() {
		return nil
	}

	ch := stt.template[stt.pos]
	stt.pos++

	var nextCh rune
	if stt.pos < len(stt.template) {
		nextCh = stt.template[stt.pos]
	}

	switch ch {
	case '!':
		// Last letter devoiced
		stt.pos++
		return &SuffixTemplateToken{
			Type:  LAST_NOT_VOICED,
			Value: nextCh,
		}
	case '+':
		// Append or optional vowel
		stt.pos++
		if nextCh == 'I' {
			return &SuffixTemplateToken{
				Type:     I_WOVEL,
				Value:    0,
				Optional: true,
			}
		} else if nextCh == 'A' {
			return &SuffixTemplateToken{
				Type:     A_WOVEL,
				Value:    0,
				Optional: true,
			}
		}
		return &SuffixTemplateToken{
			Type:  APPEND,
			Value: nextCh,
		}
	case '>':
		// Devoice first: if last letter is voiceless, devoice this letter
		stt.pos++
		return &SuffixTemplateToken{
			Type:  DEVOICE_FIRST,
			Value: nextCh,
		}
	case 'A':
		// A-type vowel (a/e)
		return &SuffixTemplateToken{
			Type:  A_WOVEL,
			Value: 0,
		}
	case 'I':
		// I-type vowel (ı/i/u/ü)
		return &SuffixTemplateToken{
			Type:  I_WOVEL,
			Value: 0,
		}
	case '~':
		// Last letter voiced
		stt.pos++
		return &SuffixTemplateToken{
			Type:  LAST_VOICED,
			Value: nextCh,
		}
	default:
		// Regular letter
		return &SuffixTemplateToken{
			Type:  LETTER,
			Value: ch,
		}
	}
}

// TokenizeSuffixTemplate tokenizes a template string
func TokenizeSuffixTemplate(template string) []*SuffixTemplateToken {
	tokenizer := NewSuffixTemplateTokenizer(template)
	tokens := make([]*SuffixTemplateToken, 0)
	for tokenizer.HasNext() {
		tokens = append(tokens, tokenizer.Next())
	}
	return tokens
}
