package parser

// ValueType describes a type of a value.
type ValueType string

// The valid value types.
const (
	ValueTypeNone   ValueType = "none"
	ValueTypeVector ValueType = "vector"
	ValueTypeScalar ValueType = "scalar"
	ValueTypeMatrix ValueType = "matrix"
	ValueTypeString ValueType = "string"
)

type Value interface {
	Type() ValueType
	String() string
}
