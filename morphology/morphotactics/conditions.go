package morphotactics

import (
	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

// Condition interface for morphotactic conditions
type Condition interface {
	Accept(path SearchPathInterface) bool
	Not() Condition
	And(other Condition) Condition
	Or(other Condition) Condition
	AndNot(other Condition) Condition
}

// Singleton condition instances
var (
	HAS_TAIL                 Condition
	HAS_SURFACE              Condition
	HAS_NO_SURFACE           Condition
	CURRENT_GROUP_EMPTY      Condition
)

func init() {
	HAS_TAIL = &HasTail{}
	HAS_SURFACE = &HasAnySuffixSurface{}
	HAS_NO_SURFACE = CondNot(&HasAnySuffixSurface{})
	CURRENT_GROUP_EMPTY = &NoSurfaceAfterDerivation{}
}

// Helper functions for creating conditions
func CondNot(c Condition) Condition {
	return &NotCondition{Condition: c}
}

func CondAnd(left, right Condition) Condition {
	return NewCombinedCondition(AND, left, right)
}

func CondOr(left, right Condition) Condition {
	return NewCombinedCondition(OR, left, right)
}

// BaseCondition provides default implementations
type BaseCondition struct{}

func (c BaseCondition) Not() Condition {
	// This should never be called directly - subclasses override
	panic("BaseCondition.Not() should not be called directly")
}

func (c BaseCondition) And(other Condition) Condition {
	// This should never be called directly  
	panic("BaseCondition.And() should not be called directly")
}

func (c BaseCondition) Or(other Condition) Condition {
	// This should never be called directly
	panic("BaseCondition.Or() should not be called directly")
}

func (c BaseCondition) AndNot(other Condition) Condition {
	// This should never be called directly
	panic("BaseCondition.AndNot() should not be called directly")
}

// HasPhoneticAttribute checks for phonetic attribute
type HasPhoneticAttribute struct {
	Attribute turkish.PhoneticAttribute
}

func (c *HasPhoneticAttribute) Accept(path SearchPathInterface) bool {
	return path.GetPhoneticAttributes()[c.Attribute]
}

func (c *HasPhoneticAttribute) Not() Condition           { return CondNot(c) }
func (c *HasPhoneticAttribute) And(o Condition) Condition { return CondAnd(c, o) }
func (c *HasPhoneticAttribute) Or(o Condition) Condition { return CondOr(c, o) }
func (c *HasPhoneticAttribute) AndNot(o Condition) Condition { return c.And(o.Not()) }

// NotCondition negates a condition
type NotCondition struct {
	Condition Condition
}

func (c *NotCondition) Accept(path SearchPathInterface) bool {
	return !c.Condition.Accept(path)
}

func (c *NotCondition) Not() Condition { return CondNot(c) }
func (c *NotCondition) And(o Condition) Condition { return CondAnd(c, o) }
func (c *NotCondition) Or(o Condition) Condition { return CondOr(c, o) }
func (c *NotCondition) AndNot(o Condition) Condition { return c.And(o.Not()) }

// CombinedCondition combines multiple conditions with AND/OR
type CombinedCondition struct {
	Operator   Operator
	Conditions []Condition
}

func NewCombinedCondition(op Operator, left, right Condition) *CombinedCondition {
	cc := &CombinedCondition{
		Operator:   op,
		Conditions: make([]Condition, 0),
	}
	cc.add(op, left)
	cc.add(op, right)
	return cc
}

func (c *CombinedCondition) add(op Operator, condition Condition) {
	if cc, ok := condition.(*CombinedCondition); ok && cc.Operator == op {
		c.Conditions = append(c.Conditions, cc.Conditions...)
	} else {
		c.Conditions = append(c.Conditions, condition)
	}
}

func (c *CombinedCondition) Count() int {
	if len(c.Conditions) == 0 {
		return 0
	}
	if len(c.Conditions) == 1 {
		if cc, ok := c.Conditions[0].(*CombinedCondition); ok {
			return cc.Count()
		}
		return 1
	}
	count := 0
	for _, cond := range c.Conditions {
		if cc, ok := cond.(*CombinedCondition); ok {
			count += cc.Count()
		} else {
			count++
		}
	}
	return count
}

func (c *CombinedCondition) Accept(path SearchPathInterface) bool {
	if len(c.Conditions) == 0 {
		return true
	}
	if len(c.Conditions) == 1 {
		return c.Conditions[0].Accept(path)
	}

	if c.Operator == AND {
		for _, cond := range c.Conditions {
			if !cond.Accept(path) {
				return false
			}
		}
		return true
	} else {
		for _, cond := range c.Conditions {
			if cond.Accept(path) {
				return true
			}
		}
		return false
	}
}

func (c *CombinedCondition) Not() Condition { return CondNot(c) }
func (c *CombinedCondition) And(o Condition) Condition { return CondAnd(c, o) }
func (c *CombinedCondition) Or(o Condition) Condition { return CondOr(c, o) }
func (c *CombinedCondition) AndNot(o Condition) Condition { return c.And(o.Not()) }

// HasRootAttribute checks for root attribute
type HasRootAttribute struct {
	Attribute turkish.RootAttribute
}

func (c *HasRootAttribute) Accept(path SearchPathInterface) bool {
	return path.GetDictionaryItem().HasAttribute(c.Attribute)
}

func (c *HasRootAttribute) Not() Condition { return CondNot(c) }
func (c *HasRootAttribute) And(o Condition) Condition { return CondAnd(c, o) }
func (c *HasRootAttribute) Or(o Condition) Condition { return CondOr(c, o) }
func (c *HasRootAttribute) AndNot(o Condition) Condition { return c.And(o.Not()) }

// Implement the same pattern for all other condition types...
// I'll create a macro-like approach by defining all conditions similarly

// DictionaryItemIs checks if root is specific item
type DictionaryItemIs struct {
	Item *lexicon.DictionaryItem
}

func (c *DictionaryItemIs) Accept(path SearchPathInterface) bool {
	return c.Item != nil && path.HasDictionaryItem(c.Item)
}

func (c *DictionaryItemIs) Not() Condition { return CondNot(c) }
func (c *DictionaryItemIs) And(o Condition) Condition { return CondAnd(c, o) }
func (c *DictionaryItemIs) Or(o Condition) Condition { return CondOr(c, o) }
func (c *DictionaryItemIs) AndNot(o Condition) Condition { return c.And(o.Not()) }

// DictionaryItemIsAny checks if root is any of items
type DictionaryItemIsAny struct {
	Items map[*lexicon.DictionaryItem]bool
}

func (c *DictionaryItemIsAny) Accept(path SearchPathInterface) bool {
	return c.Items[path.GetDictionaryItem()]
}

func (c *DictionaryItemIsAny) Not() Condition { return CondNot(c) }
func (c *DictionaryItemIsAny) And(o Condition) Condition { return CondAnd(c, o) }
func (c *DictionaryItemIsAny) Or(o Condition) Condition { return CondOr(c, o) }
func (c *DictionaryItemIsAny) AndNot(o Condition) Condition { return c.And(o.Not()) }

// DictionaryItemIsNone checks if root is none of items
type DictionaryItemIsNone struct {
	Items map[*lexicon.DictionaryItem]bool
}

func (c *DictionaryItemIsNone) Accept(path SearchPathInterface) bool {
	return !c.Items[path.GetDictionaryItem()]
}

func (c *DictionaryItemIsNone) Not() Condition { return CondNot(c) }
func (c *DictionaryItemIsNone) And(o Condition) Condition { return CondAnd(c, o) }
func (c *DictionaryItemIsNone) Or(o Condition) Condition { return CondOr(c, o) }
func (c *DictionaryItemIsNone) AndNot(o Condition) Condition { return c.And(o.Not()) }

// PreviousMorphemeIs checks if previous morpheme matches
type PreviousMorphemeIs struct {
	Morpheme *Morpheme
}

func (c *PreviousMorphemeIs) Accept(path SearchPathInterface) bool {
	prev := path.GetPreviousState()
	return prev != nil && prev.Morpheme == c.Morpheme
}

func (c *PreviousMorphemeIs) Not() Condition { return CondNot(c) }
func (c *PreviousMorphemeIs) And(o Condition) Condition { return CondAnd(c, o) }
func (c *PreviousMorphemeIs) Or(o Condition) Condition { return CondOr(c, o) }
func (c *PreviousMorphemeIs) AndNot(o Condition) Condition { return c.And(o.Not()) }

// HasTail checks if path has remaining tail
type HasTail struct{}

func (c *HasTail) Accept(path SearchPathInterface) bool {
	return len(path.GetTail()) != 0
}

func (c *HasTail) Not() Condition { return CondNot(c) }
func (c *HasTail) And(o Condition) Condition { return CondAnd(c, o) }
func (c *HasTail) Or(o Condition) Condition { return CondOr(c, o) }
func (c *HasTail) AndNot(o Condition) Condition { return c.And(o.Not()) }

// HasAnySuffixSurface checks if any suffix has surface form
type HasAnySuffixSurface struct{}

func (c *HasAnySuffixSurface) Accept(path SearchPathInterface) bool {
	return path.GetContainsSuffixWithSurface()
}

func (c *HasAnySuffixSurface) Not() Condition { return CondNot(c) }
func (c *HasAnySuffixSurface) And(o Condition) Condition { return CondAnd(c, o) }
func (c *HasAnySuffixSurface) Or(o Condition) Condition { return CondOr(c, o) }
func (c *HasAnySuffixSurface) AndNot(o Condition) Condition { return c.And(o.Not()) }

// NoSurfaceAfterDerivation checks no surface after derivation
type NoSurfaceAfterDerivation struct{}

func (c *NoSurfaceAfterDerivation) Accept(path SearchPathInterface) bool {
	transitions := path.GetTransitions()
	for i := len(transitions) - 1; i >= 1; i-- {
		sf := transitions[i]
		if sf.GetState().Derivative {
			return true
		}
		if len(sf.GetSurface()) != 0 {
			return false
		}
	}
	return true
}

func (c *NoSurfaceAfterDerivation) Not() Condition { return CondNot(c) }
func (c *NoSurfaceAfterDerivation) And(o Condition) Condition { return CondAnd(c, o) }
func (c *NoSurfaceAfterDerivation) Or(o Condition) Condition { return CondOr(c, o) }
func (c *NoSurfaceAfterDerivation) AndNot(o Condition) Condition { return c.And(o.Not()) }

// Helper functions for creating conditions

func NotHave(pAttr *turkish.PhoneticAttribute, rAttr *turkish.RootAttribute) Condition {
	if rAttr != nil {
		return (&HasRootAttribute{Attribute: *rAttr}).Not()
	}
	return (&HasPhoneticAttribute{Attribute: *pAttr}).Not()
}

func Has(pAttr *turkish.PhoneticAttribute, rAttr *turkish.RootAttribute) Condition {
	if pAttr != nil {
		return &HasPhoneticAttribute{Attribute: *pAttr}
	}
	return &HasRootAttribute{Attribute: *rAttr}
}
