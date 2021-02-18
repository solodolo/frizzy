package processor

import (
	"fmt"
)

type Result interface {
	GetResult() interface{}
	fmt.Stringer
}

type AddableResult interface {
	Add(right Result) (Result, error)
}

type SubtractableResult interface {
	Subtract(right Result) (Result, error)
}

type MultipliableResult interface {
	Multiply(right Result) (Result, error)
	Divide(right Result) (Result, error)
}

type LogicalResult interface {
	LessThan(right Result) (Result, error)
	GreaterThan(right Result) (Result, error)
	EqualTo(right Result) (Result, error)
	LessThanEqual(right Result) (Result, error)
	GreaterThanEqual(right Result) (Result, error)
}

type NotResult interface {
	Not() (Result, error)
}

type NegativeResult interface {
	Negative() (Result, error)
}
