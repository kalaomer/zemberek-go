package morphotactics

// Operator represents logical operators
type Operator int

const (
	AND Operator = iota
	OR
)

// String returns string representation
func (o Operator) String() string {
	switch o {
	case AND:
		return "AND"
	case OR:
		return "OR"
	default:
		return "UNKNOWN"
	}
}
