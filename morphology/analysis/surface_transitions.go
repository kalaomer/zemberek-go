package analysis

import (
	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

// SurfaceTransition represents a surface form transition
type SurfaceTransition struct {
	Surface           string
	LexicalTransition morphotactics.MorphemeTransition
}

// NewSurfaceTransition creates a new SurfaceTransition
func NewSurfaceTransition(surface string, lexicalTransition morphotactics.MorphemeTransition) *SurfaceTransition {
	return &SurfaceTransition{
		Surface:           surface,
		LexicalTransition: lexicalTransition,
	}
}

// GetState returns the state
func (st *SurfaceTransition) GetState() *morphotactics.MorphemeState {
	if st.LexicalTransition == nil {
		return nil
	}
	return st.LexicalTransition.GetState()
}

// GetMorpheme returns the morpheme
func (st *SurfaceTransition) GetMorpheme() *morphotactics.Morpheme {
	state := st.GetState()
	if state == nil {
		return nil
	}
	return state.Morpheme
}

// GetSurface returns the surface string (for TransitionInterface)
func (st *SurfaceTransition) GetSurface() string {
	return st.Surface
}

// IsDerivative checks if this is a derivative transition
func (st *SurfaceTransition) IsDerivative() bool {
	state := st.GetState()
	return state != nil && state.Derivative
}

// String returns string representation
func (st *SurfaceTransition) String() string {
	if len(st.Surface) > 0 {
		return st.Surface + ":" + st.GetMorpheme().ID
	}
	return st.GetMorpheme().ID
}

// SuffixTemplateToken represents a token in suffix template
type SuffixTemplateToken struct {
	Type      TemplateTokenType
	Value     rune
	Optional  bool
}

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

// GenerateSurface generates surface form from suffix transition and phonetic attributes
func GenerateSurface(transition *morphotactics.SuffixTransition, phoneticAttrs map[turkish.PhoneticAttribute]bool) string {
	// Check cache first
	cached := transition.GetFromSurfaceCache(phoneticAttrs)
	if cached != "" {
		return cached
	}

	var result []rune

	for index, token := range transition.TokenList {
		// Get current morphemic attributes
		attrs := GetMorphemicAttributes(string(result), phoneticAttrs)

		switch token.Type {
		case morphotactics.LETTER:
			result = append(result, token.Value)

		case morphotactics.A_WOVEL:
			// A-type vowel: a or e
			if index != 0 || !phoneticAttrs[turkish.LastLetterVowel] {
				if attrs[turkish.LastVowelBack] {
					result = append(result, 'a')
				} else if attrs[turkish.LastVowelFrontal] {
					result = append(result, 'e')
				}
			}

		case morphotactics.I_WOVEL:
			// I-type vowel: ı, i, u, or ü
			if index != 0 || !phoneticAttrs[turkish.LastLetterVowel] {
				if attrs[turkish.LastVowelFrontal] && attrs[turkish.LastVowelUnrounded] {
					result = append(result, 'i')
				} else if attrs[turkish.LastVowelBack] && attrs[turkish.LastVowelUnrounded] {
					result = append(result, 'ı')
				} else if attrs[turkish.LastVowelBack] && attrs[turkish.LastVowelRounded] {
					result = append(result, 'u')
				} else if attrs[turkish.LastVowelFrontal] && attrs[turkish.LastVowelRounded] {
					result = append(result, 'ü')
				}
			}

		case morphotactics.APPEND:
			// Append letter only if last letter is vowel
			if attrs[turkish.LastLetterVowel] {
				result = append(result, token.Value)
			}

		case morphotactics.DEVOICE_FIRST:
			// Devoice letter if last letter is voiceless
			letter := token.Value
			if attrs[turkish.LastLetterVoiceless] {
				letter = turkish.Instance.Devoice(letter)
			}
			result = append(result, letter)

		case morphotactics.LAST_VOICED, morphotactics.LAST_NOT_VOICED:
			result = append(result, token.Value)
		}
	}

	surface := string(result)
	transition.AddToSurfaceCache(phoneticAttrs, surface)
	return surface
}
