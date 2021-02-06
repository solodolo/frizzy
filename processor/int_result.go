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

// Subtract takes a result and subtracts it from this integer representation
// It returns a Result type or an error if right cannot be subtracted
// to receiver
func (receiver IntResult) Subtract(right Result) (Result, error) {
	switch typedRight := right.(type) {
	case IntResult:
		return IntResult(receiver - typedRight), nil
	default:
		return nil, fmt.Errorf("Cannt subtract %T and %T", receiver, right)
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

// Divide takes a result and divides it with this integer representation
// Returns a Result type or an error if right cannot be divided with receiver
// Floats are currently not supported so all division is integer divison
func (receiver IntResult) Divide(right Result) (Result, error) {
	switch typedRight := right.(type) {
	case IntResult:
		return IntResult(receiver / typedRight), nil
	default:
		return nil, fmt.Errorf("Cannot divide a %T and a %T", receiver, right)
	}
}

func convertToInt(right Result) (int, bool) {
	switch typedResult := right.(type) {
	case IntResult:
		return int(typedResult), true
	case StringResult:
		if num, err := strconv.Atoi(string(typedResult)); err == nil {
			return num, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// LessThan checks if the provided result is logically less than
// the receiver
func (receiver IntResult) LessThan(right Result) (Result, error) {
	if rightInt, ok := convertToInt(right); ok {
		return BoolResult(int(receiver) < rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T < %T", receiver, right)
}

// GreaterThan checks if the provided result is logically greater than
// the receiver
func (receiver IntResult) GreaterThan(right Result) (Result, error) {
	if rightInt, ok := convertToInt(right); ok {
		return BoolResult(int(receiver) > rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T > %T", receiver, right)
}

// EqualTo checks if the provided result is logically equal to
// the receiver
func (receiver IntResult) EqualTo(right Result) (Result, error) {
	if rightInt, ok := convertToInt(right); ok {
		return BoolResult(int(receiver) == rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T == %T", receiver, right)
}

// LessThanEqual checks if the provided result is logically less than or equal to
// the receiver
func (receiver IntResult) LessThanEqual(right Result) (Result, error) {
	if rightInt, ok := convertToInt(right); ok {
		return BoolResult(int(receiver) <= rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T <= %T", receiver, right)
}

// GreaterThanEqual checks if the provided result is logically greater than or equal to
// the receiver
func (receiver IntResult) GreaterThanEqual(right Result) (Result, error) {
	if rightInt, ok := convertToInt(right); ok {
		return BoolResult(int(receiver) >= rightInt), nil
	}
	return nil, fmt.Errorf("Cannot determine %T >= %T", receiver, right)
}
