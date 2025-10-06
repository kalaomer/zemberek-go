package morphotactics

import (
	"strings"
	"unicode"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

// SuffixTransition represents a morpheme transition with surface template
type SuffixTransition struct {
	From             *MorphemeState
	To               *MorphemeState
	SurfaceTemplate  string
	Condition        Condition
	TokenList        []*SuffixTemplateToken
	ConditionCount   int
	SurfaceCache     *AttributeToSurfaceCache
}

// NewSuffixTransition creates a new suffix transition
func NewSuffixTransition(from, to *MorphemeState, template string, condition Condition) *SuffixTransition {
	st := &SuffixTransition{
		From:            from,
		To:              to,
		SurfaceTemplate: template,
		Condition:       condition,
		SurfaceCache:    NewAttributeToSurfaceCache(),
	}

	st.conditionsFromTemplate(template)
	st.TokenList = TokenizeSuffixTemplate(template)
	st.ConditionCount = st.countConditions()

	return st
}

// NewSuffixTransitionBuilder creates a builder
func NewSuffixTransitionBuilder(from, to *MorphemeState) *SuffixTransitionBuilder {
	return &SuffixTransitionBuilder{
		From: from,
		To:   to,
	}
}

// SuffixTransitionBuilder builds suffix transitions
type SuffixTransitionBuilder struct {
	From     *MorphemeState
	To       *MorphemeState
	Template string
	Cond     Condition
}

func (b *SuffixTransitionBuilder) SetTemplate(template string) *SuffixTransitionBuilder {
	b.Template = template
	return b
}

func (b *SuffixTransitionBuilder) SetCondition(condition Condition) *SuffixTransitionBuilder {
	b.Cond = condition
	return b
}

func (b *SuffixTransitionBuilder) Empty() *SuffixTransitionBuilder {
	b.Template = ""
	return b
}

func (b *SuffixTransitionBuilder) Build() *SuffixTransition {
	st := NewSuffixTransition(b.From, b.To, b.Template, b.Cond)
	st.Connect()
	return st
}

// GetState returns the target state
func (st *SuffixTransition) GetState() *MorphemeState {
	return st.To
}

// GetMorpheme returns the morpheme
func (st *SuffixTransition) GetMorpheme() *Morpheme {
	return st.To.Morpheme
}

// Connect connects this transition to states
func (st *SuffixTransition) Connect() {
	st.From.AddOutgoing(st)
	st.To.AddIncoming(st)
}

// GetCopy creates a copy of this transition
func (st *SuffixTransition) GetCopy() MorphemeTransition {
	copy := &SuffixTransition{
		From:            st.From,
		To:              st.To,
		SurfaceTemplate: st.SurfaceTemplate,
		Condition:       st.Condition,
		TokenList:       st.TokenList,
		ConditionCount:  st.ConditionCount,
		SurfaceCache:    st.SurfaceCache,
	}
	return copy
}

// countConditions counts the number of conditions
func (st *SuffixTransition) countConditions() int {
	if st.Condition == nil {
		return 0
	}
	if cc, ok := st.Condition.(*CombinedCondition); ok {
		return cc.Count()
	}
	return 1
}

// conditionsFromTemplate extracts conditions from template
func (st *SuffixTransition) conditionsFromTemplate(template string) {
	if template == "" {
		return
	}

	// Turkish lowercase mapping - convert to []rune to handle UTF-8 properly
	lowerStr := strings.Map(func(r rune) rune {
		switch r {
		case 'I':
			return 'ı'
		case 'İ':
			return 'i'
		default:
			return unicode.ToLower(r)
		}
	}, template)
	lower := []rune(lowerStr)  // Convert to rune slice for proper indexing

	var c Condition
	firstCharVowel := false
	if len(lower) > 0 {
		firstChar := lower[0]  // Already a rune
		firstCharVowel = turkish.Instance.IsVowel(firstChar)
	}

	// Check for ExpectsVowel constraint
	if (len(lower) > 0 && lower[0] == '>') || (len(lower) > 0 && !firstCharVowel) {
		expectsVowel := turkish.ExpectsVowel
		c = NotHave(&expectsVowel, nil)
	}

	// Check for ExpectsConsonant constraint - Java uses separate if, not else-if!
	if (len(lower) >= 3 && lower[0] == '+' && turkish.Instance.IsVowel(lower[2])) || firstCharVowel {
		expectsConsonant := turkish.ExpectsConsonant
		c = NotHave(&expectsConsonant, nil)
	}

	if c != nil {
		if st.Condition == nil {
			st.Condition = c
		} else {
			st.Condition = c.And(st.Condition)
		}
	}
}

// AddToSurfaceCache adds a surface form to cache
func (st *SuffixTransition) AddToSurfaceCache(attributes map[turkish.PhoneticAttribute]bool, value string) {
	st.SurfaceCache.AddSurface(attributes, value)
}

// GetFromSurfaceCache retrieves from cache
func (st *SuffixTransition) GetFromSurfaceCache(attributes map[turkish.PhoneticAttribute]bool) string {
	return st.SurfaceCache.GetSurface(attributes)
}

// GetLastTemplateToken returns the last template token
func (st *SuffixTransition) GetLastTemplateToken() *SuffixTemplateToken {
	if len(st.TokenList) == 0 {
		return nil
	}
	return st.TokenList[len(st.TokenList)-1]
}

// HasSurfaceForm checks if transition has surface form
func (st *SuffixTransition) HasSurfaceForm() bool {
	return len(st.TokenList) > 0
}

// CanPass checks if this transition can pass with given search path
func (st *SuffixTransition) CanPass(path SearchPathInterface) bool {
	return st.Condition == nil || st.Condition.Accept(path)
}

// String returns string representation
func (st *SuffixTransition) String() string {
	result := "[" + st.From.ID + "→" + st.To.ID
	if st.SurfaceTemplate != "" {
		result += ":" + st.SurfaceTemplate
	}
	result += "]"
	return result
}
