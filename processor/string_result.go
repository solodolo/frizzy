package processor

import (
	"fmt"
	"strconv"
)

type StringResult string

func (receiver StringResult) GetResult() interface{} {
	return receiver
}

func (receiver StringResult) String() string {
	return string(receiver)
}

func (receiver StringResult) Add(right Result) (Result, error) {
	switch typedRight := right.(type) {
	case IntResult:
		return StringResult(string(receiver) + strconv.Itoa(int(typedRight))), nil
	case StringResult:
		return StringResult(receiver + typedRight), nil
	default:
		return nil, fmt.Errorf("Cannot add a %T to a %T", receiver, right)
	}
}

func convertToString(result Result) (string, bool) {
	switch typedResult := result.(type) {
	case IntResult:
		return strconv.Itoa(int(typedResult)), true
	case StringResult:
		return string(typedResult), true
	default:
		return "", false
	}
}

// EqualTo checks if the provided result is logically equal to
// the receiver
func (receiver StringResult) EqualTo(right Result) (Result, error) {
	if rightStr, ok := convertToString(right); ok {
		return BoolResult(string(receiver) == rightStr), nil
	}

	return nil, fmt.Errorf("Cannot determine %T == %T", receiver, right)
}

// NotEqualTo checks if the provided result is logically not equal
// to the receiver
func (receiver StringResult) NotEqualTo(right Result) (Result, error) {
	if rightStr, ok := convertToString(right); ok {
		return BoolResult(string(receiver) != rightStr), nil
	}

	return nil, fmt.Errorf("Cannot determine %T != %T", receiver, right)
}

// LessThan checks if the provided result is logically less than
// the receiver
func (receiver StringResult) LessThan(right Result) (Result, error) {
	if rightStr, ok := convertToString(right); ok {
		return BoolResult(string(receiver) < rightStr), nil
	}

	return nil, fmt.Errorf("Cannot determine %T < %T", receiver, right)
}

// GreaterThan checks if the provided result is logically greater than
// the receiver
func (receiver StringResult) GreaterThan(right Result) (Result, error) {
	if rightStr, ok := convertToString(right); ok {
		return BoolResult(string(receiver) > rightStr), nil
	}

	return nil, fmt.Errorf("Cannot determine %T > %T", receiver, right)
}

// LessThanEqual checks if the provided result is logically less than or equal to
// the receiver
func (receiver StringResult) LessThanEqual(right Result) (Result, error) {
	if rightStr, ok := convertToString(right); ok {
		return BoolResult(string(receiver) <= rightStr), nil
	}

	return nil, fmt.Errorf("Cannot determine %T <= %T", receiver, right)
}

// GreaterThanEqual checks if the provided result is logically greater than or equal to
// the receiver
func (receiver StringResult) GreaterThanEqual(right Result) (Result, error) {
	if rightStr, ok := convertToString(right); ok {
		return BoolResult(string(receiver) >= rightStr), nil
	}

	return nil, fmt.Errorf("Cannot determine %T >= %T", receiver, right)
}
