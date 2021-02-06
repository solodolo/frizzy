package parser

import (
	"testing"
)

func TestLessThan(t *testing.T) {
	var tests = []struct {
		left     IntResult
		right    Result
		expected bool
	}{
		{IntResult(5), IntResult(4), false},
		{IntResult(4), IntResult(5), true},
		{IntResult(5), IntResult(5), false},
		{IntResult(-4), IntResult(-3), true},
		{IntResult(-1), IntResult(-3), false},
		{IntResult(5), StringResult("6"), true},
		{IntResult(6), StringResult("5"), false},
	}

	for _, test := range tests {
		result, _ := test.left.LessThan(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v < %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	var tests = []struct {
		left     IntResult
		right    Result
		expected bool
	}{
		{IntResult(5), IntResult(4), true},
		{IntResult(4), IntResult(5), false},
		{IntResult(5), IntResult(5), false},
		{IntResult(-4), IntResult(-3), false},
		{IntResult(-1), IntResult(-3), true},
		{IntResult(5), StringResult("6"), false},
		{IntResult(6), StringResult("5"), true},
	}

	for _, test := range tests {
		result, _ := test.left.GreaterThan(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v > %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestEqualTo(t *testing.T) {
	var tests = []struct {
		left     IntResult
		right    Result
		expected bool
	}{
		{IntResult(5), IntResult(4), false},
		{IntResult(4), IntResult(5), false},
		{IntResult(5), IntResult(5), true},
		{IntResult(-4), IntResult(-4), true},
		{IntResult(-1), IntResult(-3), false},
		{IntResult(0), IntResult(0), true},
		{IntResult(5), StringResult("5"), true},
		{IntResult(6), StringResult("5"), false},
	}

	for _, test := range tests {
		result, _ := test.left.EqualTo(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v == %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestLessThanOrEqualTo(t *testing.T) {
	var tests = []struct {
		left     IntResult
		right    Result
		expected bool
	}{
		{IntResult(5), IntResult(4), false},
		{IntResult(4), IntResult(5), true},
		{IntResult(5), IntResult(5), true},
		{IntResult(-4), IntResult(-3), true},
		{IntResult(-1), IntResult(-3), false},
		{IntResult(5), StringResult("6"), true},
		{IntResult(5), StringResult("5"), true},
		{IntResult(6), StringResult("5"), false},
	}

	for _, test := range tests {
		result, _ := test.left.LessThanEqual(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v <= %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestGreaterThanOrEqualTo(t *testing.T) {
	var tests = []struct {
		left     IntResult
		right    Result
		expected bool
	}{
		{IntResult(5), IntResult(4), true},
		{IntResult(4), IntResult(5), false},
		{IntResult(5), IntResult(5), true},
		{IntResult(-4), IntResult(-3), false},
		{IntResult(-1), IntResult(-3), true},
		{IntResult(5), StringResult("6"), false},
		{IntResult(6), StringResult("5"), true},
	}

	for _, test := range tests {
		result, _ := test.left.GreaterThanEqual(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v >= %v to be %v", test.left, test.right, test.expected)
		}
	}
}
