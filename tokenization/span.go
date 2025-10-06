package tokenization

import "fmt"

// Span represents specified chunks of a string
type Span struct {
	Start int
	End   int
}

// NewSpan creates a new Span
func NewSpan(start, end int) (*Span, error) {
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("span start and end values cannot be negative")
	}
	if end < start {
		return nil, fmt.Errorf("span end value cannot be smaller than start value")
	}
	return &Span{
		Start: start,
		End:   end,
	}, nil
}

// GetLength returns the length of the span
func (s *Span) GetLength() int {
	return s.End - s.Start
}

// MiddleValue returns the middle value of the span
func (s *Span) MiddleValue() int {
	return s.End + (s.End-s.Start)/2
}

// GetSubString returns the substring from the given string
func (s *Span) GetSubString(str string) string {
	return str[s.Start:s.End]
}

// InSpan returns true if the index is within the span
func (s *Span) InSpan(i int) bool {
	return s.Start <= i && i < s.End
}

// Copy returns a copy of the span with an offset
func (s *Span) Copy(offset int) *Span {
	return &Span{
		Start: offset + s.Start,
		End:   offset + s.End,
	}
}
