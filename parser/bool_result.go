package parser

import (
	"fmt"
)

type BoolResult bool

func (receiver BoolResult) GetResult() interface{} {
	return receiver
}

func convertToBool(result Result) (bool, bool) {
	switch typedResult := result.(type) {
	case StringResult:
		if typedResult == "true" {
			return true, true
		} else if typedResult == "false" {
			return false, true
		}
		return false, false
	case IntResult:
		return typedResult != 0, true
	default:
		return false, false
	}
}

// LessThan checks if the provided result is logically less than
// the receiver
func (receiver BoolResult) LessThan(right Result) (Result, error) {
	return nil, fmt.Errorf("Cannot determine %T < %T", receiver, right)
}

// GreaterThan checks if the provided result is logically greater than
// the receiver
func (receiver BoolResult) GreaterThan(right Result) (Result, error) {
	return nil, fmt.Errorf("Cannot determine %T > %T", receiver, right)
}

// EqualTo checks if the provided result is logically equal to
// the receiver
func (receiver BoolResult) EqualTo(right Result) (Result, error) {
	if rightInt, ok := convertToBool(right); ok {
		return BoolResult(bool(receiver) == rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T == %T", receiver, right)
}

// LessThanEqual checks if the provided result is logically less than or equal to
// the receiver
func (receiver BoolResult) LessThanEqual(right Result) (Result, error) {
	return nil, fmt.Errorf("Cannot determine %T <= %T", receiver, right)
}

// GreaterThanEqual checks if the provided result is logically greater than or equal to
// the receiver
func (receiver BoolResult) GreaterThanEqual(right Result) (Result, error) {
	return nil, fmt.Errorf("Cannot determine %T >= %T", receiver, right)
}
