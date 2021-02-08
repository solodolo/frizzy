package processor

import (
	"testing"
)

func TestStringResultLessThan(t *testing.T) {
	var tests = []struct {
		left     StringResult
		right    Result
		expected bool
	}{
		{StringResult("A"), StringResult("B"), true},
		{StringResult("B"), StringResult("A"), false},
		{StringResult("A"), StringResult("A"), false},
		{StringResult("5"), IntResult(6), true},
		{StringResult("5"), IntResult(4), false},
	}

	for _, test := range tests {
		result, _ := test.left.LessThan(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v < %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestStringResultGreaterThan(t *testing.T) {
	var tests = []struct {
		left     StringResult
		right    Result
		expected bool
	}{
		{StringResult("A"), StringResult("B"), false},
		{StringResult("B"), StringResult("A"), true},
		{StringResult("A"), StringResult("A"), false},
		{StringResult("5"), IntResult(6), false},
		{StringResult("5"), IntResult(4), true},
	}

	for _, test := range tests {
		result, _ := test.left.GreaterThan(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v > %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestStringResultEqualTo(t *testing.T) {
	var tests = []struct {
		left     StringResult
		right    Result
		expected bool
	}{
		{StringResult("A"), StringResult("B"), false},
		{StringResult("B"), StringResult("A"), false},
		{StringResult("A"), StringResult("A"), true},
		{StringResult("5"), IntResult(5), true},
		{StringResult("5"), IntResult(4), false},
	}

	for _, test := range tests {
		result, _ := test.left.EqualTo(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v == %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestStringResultLessThanOrEqualTo(t *testing.T) {
	var tests = []struct {
		left     StringResult
		right    Result
		expected bool
	}{
		{StringResult("A"), StringResult("B"), true},
		{StringResult("B"), StringResult("A"), false},
		{StringResult("A"), StringResult("A"), true},
		{StringResult("5"), IntResult(6), true},
		{StringResult("5"), IntResult(4), false},
	}

	for _, test := range tests {
		result, _ := test.left.LessThanEqual(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v <= %v to be %v", test.left, test.right, test.expected)
		}
	}
}

func TestStringResultGreaterThanOrEqual(t *testing.T) {
	var tests = []struct {
		left     StringResult
		right    Result
		expected bool
	}{
		{StringResult("A"), StringResult("B"), false},
		{StringResult("B"), StringResult("A"), true},
		{StringResult("A"), StringResult("A"), true},
		{StringResult("5"), IntResult(6), false},
		{StringResult("5"), IntResult(4), true},
	}

	for _, test := range tests {
		result, _ := test.left.GreaterThanEqual(test.right)
		boolResult := bool(result.(BoolResult))

		if boolResult != test.expected {
			t.Errorf("expected %v >= %v to be %v", test.left, test.right, test.expected)
		}
	}
}
