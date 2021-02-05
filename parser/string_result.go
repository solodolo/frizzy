package parser

import (
	"fmt"
	"strconv"
)

type StringResult string

func (result StringResult) GetResult() interface{} {
	return result
}

func (left StringResult) Add(right Result) (Result, error) {
	switch typedRight := right.(type) {
	case IntResult:
		return StringResult(string(left) + strconv.Itoa(int(typedRight))), nil
	case StringResult:
		return StringResult(left + typedRight), nil
	default:
		return nil, fmt.Errorf("Cannot add a %T to a %T", left, right)
	}
}
