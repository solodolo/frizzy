package processor

import (
	"fmt"
	"strconv"
)

// BoolResult represents a boolean processor result
type BoolResult bool

// GetResult returns this result
func (receiver BoolResult) GetResult() interface{} {
	return receiver
}

func (receiver BoolResult) String() string {
	return strconv.FormatBool(bool(receiver))
}

func convertToBool(result Result) (bool, bool) {
	switch typedResult := result.(type) {
	case BoolResult:
		return bool(typedResult), true
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

// EqualTo checks if the provided result is logically equal to
// the receiver
func (receiver BoolResult) EqualTo(right Result) (Result, error) {
	if rightInt, ok := convertToBool(right); ok {
		return BoolResult(bool(receiver) == rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T == %T", receiver, right)
}

// NotEqualTo checks if the provided result is logically equal to
// the receiver
func (receiver BoolResult) NotEqualTo(right Result) (Result, error) {
	if rightInt, ok := convertToBool(right); ok {
		return BoolResult(bool(receiver) != rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T != %T", receiver, right)
}

// LogicalAnd determines if left and right are both logically true
func (receiver BoolResult) LogicalAnd(right Result) (Result, error) {
	if rightInt, ok := convertToBool(right); ok {
		return BoolResult(bool(receiver) && rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T && %T", receiver, right)
}

// LogicalOr determines if left and right are both logically true
func (receiver BoolResult) LogicalOr(right Result) (Result, error) {
	if rightInt, ok := convertToBool(right); ok {
		return BoolResult(bool(receiver) || rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T || %T", receiver, right)
}

// Not returns the inverse of the receiver
func (receiver BoolResult) Not() (Result, error) {
	return !receiver, nil
}
