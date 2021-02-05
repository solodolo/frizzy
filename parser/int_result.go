package parser

import (
	"fmt"
	"strconv"
)

// IntResult represents a result containing an integer value
type IntResult int

// GetResult returns this result value
func (receiver IntResult) GetResult() interface{} {
	return receiver
}

func (receiver IntResult) String() string {
	return strconv.Itoa(int(receiver))
}

// Add takes a result and adds it to this integer representation
// It returns a Result type or an error if right cannot be added
// to receiver
func (receiver IntResult) Add(right Result) (Result, error) {
	switch typedRight := right.(type) {
	case IntResult:
		return IntResult(receiver + typedRight), nil
	case StringResult:
		return StringResult(strconv.Itoa(int(receiver)) + string(typedRight)), nil
	default:
		return nil, fmt.Errorf("Cannot add a %T to a %T", receiver, right)
	}
}

// Multiply takes a result and multiplies it with this integer representation
// Returns a Result type or an error if right cannot be multiplied with receiver
func (receiver IntResult) Multiply(right Result) (Result, error) {
	switch typedRight := right.(type) {
	case IntResult:
		return IntResult(receiver * typedRight), nil
	default:
		return nil, fmt.Errorf(("Cannot multiply a %T and a %T"), receiver, right)
	}
}
