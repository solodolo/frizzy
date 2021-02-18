package processor

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"mettlach.codes/frizzy/lexer"
	"mettlach.codes/frizzy/parser"
)

func getNodeChan(nodes []parser.TreeNode) chan parser.TreeNode {
	nodeChan := make(chan parser.TreeNode)
	go func(nodes []parser.TreeNode) {
		defer close(nodeChan)

		for _, head := range nodes {
			nodeChan <- head
		}
	}(nodes)

	return nodeChan
}

func TestProcessStringNodeReturnsString(t *testing.T) {
	head := generateStringTree("foo")
	resultChan := runProcess(head)
	result := <-resultChan
	if result.String() != "foo" {
		t.Errorf("expected result to be \"foo\", got %s", result.String())
	}
}

func TestProcessNumNodeReturnsNumAsString(t *testing.T) {
	head := generateNumTree(123)
	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "123" {
		t.Errorf("expected result to be 123, got %s", result.String())
	}
}

func TestProcessBoolNodeReturnsBoolAsString(t *testing.T) {
	trueHead := generateTree([]lexer.Token{lexer.BoolToken{Value: "true"}})
	falseHead := generateTree([]lexer.Token{lexer.BoolToken{Value: "false"}})
	trueResultChan := runProcess(trueHead)
	falseResultChan := runProcess(falseHead)

	trueResult := <-trueResultChan
	falseResult := <-falseResultChan

	if trueResult.String() != "true" {
		t.Errorf("expected true result to be \"true\", got %q", trueResult.String())
	}

	if falseResult.String() != "false" {
		t.Errorf("expected false result to be \"false\", got %q", falseResult.String())
	}
}

func TestAddTwoNumbersReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.NumToken{Num: "10"},
		lexer.AddOpToken{Operator: "+"},
		lexer.NumToken{Num: "3"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "13" {
		t.Errorf("expected result to be 13, got %s", result.String())
	}
}

func TestAddOneNegativeNumbersReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.UnaryOpToken{Operator: "-"},
		lexer.NumToken{Num: "10"},
		lexer.AddOpToken{Operator: "+"},
		lexer.NumToken{Num: "3"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "-7" {
		t.Errorf("expected result to be -7, got %s", result.String())
	}
}

func TestSubtractTwoNumbersReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.NumToken{Num: "10"},
		lexer.AddOpToken{Operator: "-"},
		lexer.NumToken{Num: "3"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "7" {
		t.Errorf("expected result to be 7, got %s", result.String())
	}
}

func TestSubtractFromNegativeNumberReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.UnaryOpToken{Operator: "-"},
		lexer.NumToken{Num: "10"},
		lexer.AddOpToken{Operator: "-"},
		lexer.NumToken{Num: "99"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "-109" {
		t.Errorf("expected result to be -109, got %s", result.String())
	}
}

func TestMultiplyTwoNumbersReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.NumToken{Num: "99"},
		lexer.MultOpToken{Operator: "*"},
		lexer.NumToken{Num: "71"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "7029" {
		t.Errorf("expected result to be 7029, got %s", result.String())
	}
}

func TestMultiplyPositiveNegativeNumbersReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.NumToken{Num: "6"},
		lexer.MultOpToken{Operator: "*"},
		lexer.UnaryOpToken{Operator: "-"},
		lexer.NumToken{Num: "4"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "-24" {
		t.Errorf("expected result to be -24, got %s", result.String())
	}
}

func TestMultiplyNegativePositiveNumbersReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.UnaryOpToken{Operator: "-"},
		lexer.NumToken{Num: "6001"},
		lexer.MultOpToken{Operator: "*"},
		lexer.NumToken{Num: "30"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "-180030" {
		t.Errorf("expected result to be -180030 got %s", result.String())
	}
}

func TestMultiplyTwoNegativeNumbersReturnsCorrectResult(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.UnaryOpToken{Operator: "-"},
		lexer.NumToken{Num: "13"},
		lexer.MultOpToken{Operator: "*"},
		lexer.UnaryOpToken{Operator: "-"},
		lexer.NumToken{Num: "44"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "572" {
		t.Errorf("expected result to be 572, got %s", result.String())
	}
}

func TestNegationOfTrueBoolReturnsFalse(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.UnaryOpToken{Operator: "!"},
		lexer.BoolToken{Value: "true"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "false" {
		t.Errorf("expected result to be false, got %s", result.String())
	}
}
func TestNegationOfFalseBoolReturnsTrue(t *testing.T) {
	head := generateTree([]lexer.Token{
		lexer.UnaryOpToken{Operator: "!"},
		lexer.BoolToken{Value: "false"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != "true" {
		t.Errorf("expected result to be true, got %s", result.String())
	}
}

func TestLTOfTwoNumbersReturnsCorrectResult(t *testing.T) {
	nums := [][]int{
		{6, 7},
		{100, 2},
		{42, 42},
	}

	for _, test := range nums {
		expected := test[0] < test[1]

		head := generateTree([]lexer.Token{
			lexer.NumToken{Num: strconv.Itoa(test[0])},
			lexer.RelOpToken{Operator: "<"},
			lexer.NumToken{Num: strconv.Itoa(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}
func TestLTEOfTwoNumbersReturnsCorrectResult(t *testing.T) {
	nums := [][]int{
		{983, 456},
		{3520, 9874},
		{42, 42},
	}

	for _, test := range nums {
		expected := test[0] <= test[1]

		head := generateTree([]lexer.Token{
			lexer.NumToken{Num: strconv.Itoa(test[0])},
			lexer.RelOpToken{Operator: "<="},
			lexer.NumToken{Num: strconv.Itoa(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestGTOfTwoNumbersReturnsCorrectResult(t *testing.T) {
	nums := [][]int{
		{1021, 6789},
		{38, 19},
		{42, 42},
	}

	for _, test := range nums {
		expected := test[0] > test[1]

		head := generateTree([]lexer.Token{
			lexer.NumToken{Num: strconv.Itoa(test[0])},
			lexer.RelOpToken{Operator: ">"},
			lexer.NumToken{Num: strconv.Itoa(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestGTEOfTwoNumbersReturnsCorrectResult(t *testing.T) {
	nums := [][]int{
		{math.MinInt64, math.MaxInt64},
		{1209934, 19},
		{42, 42},
	}

	for _, test := range nums {
		expected := test[0] >= test[1]

		head := generateTree([]lexer.Token{
			lexer.NumToken{Num: strconv.Itoa(test[0])},
			lexer.RelOpToken{Operator: ">="},
			lexer.NumToken{Num: strconv.Itoa(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestEqOfTwoNumbersReturnsCorrectResult(t *testing.T) {
	nums := [][]int{
		{123, 321},
		{897012, 5243},
		{42, 42},
	}

	for _, test := range nums {
		expected := test[0] == test[1]

		head := generateTree([]lexer.Token{
			lexer.NumToken{Num: strconv.Itoa(test[0])},
			lexer.RelOpToken{Operator: "=="},
			lexer.NumToken{Num: strconv.Itoa(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestNotEqOfTwoNumbersReturnsCorrectResult(t *testing.T) {
	nums := [][]int{
		{783, 65847},
		{28796, 543},
		{42, 42},
	}

	for _, test := range nums {
		expected := test[0] != test[1]

		head := generateTree([]lexer.Token{
			lexer.NumToken{Num: strconv.Itoa(test[0])},
			lexer.RelOpToken{Operator: "!="},
			lexer.NumToken{Num: strconv.Itoa(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestLTOfTwoStringsReturnsCorrectResult(t *testing.T) {
	strs := [][]string{
		{"foo", "bar"},
		{"baz", "Baz"},
		{"fizzbuzz", "fizzbuzz"},
	}

	for _, test := range strs {
		expected := test[0] < test[1]

		head := generateTree([]lexer.Token{
			lexer.StrToken{Str: test[0]},
			lexer.RelOpToken{Operator: "<"},
			lexer.StrToken{Str: test[1]},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestLTEqOfTwoStringsReturnsCorrectResult(t *testing.T) {
	strs := [][]string{
		{"foo", "bar"},
		{"baz", "Baz"},
		{"fizzbuzz", "fizzbuzz"},
	}

	for _, test := range strs {
		expected := test[0] <= test[1]

		head := generateTree([]lexer.Token{
			lexer.StrToken{Str: test[0]},
			lexer.RelOpToken{Operator: "<="},
			lexer.StrToken{Str: test[1]},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestGTOfTwoStringsReturnsCorrectResult(t *testing.T) {
	strs := [][]string{
		{"foo", "bar"},
		{"baz", "Baz"},
		{"fizzbuzz", "fizzbuzz"},
	}

	for _, test := range strs {
		expected := test[0] > test[1]

		head := generateTree([]lexer.Token{
			lexer.StrToken{Str: test[0]},
			lexer.RelOpToken{Operator: ">"},
			lexer.StrToken{Str: test[1]},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestGTEqOfTwoStringsReturnsCorrectResult(t *testing.T) {
	strs := [][]string{
		{"foo", "\u21E7bar"},
		{"baz", "Baz"},
		{"fizzbuzz", "fizzbuzz"},
	}

	for _, test := range strs {
		expected := test[0] >= test[1]

		head := generateTree([]lexer.Token{
			lexer.StrToken{Str: test[0]},
			lexer.RelOpToken{Operator: ">="},
			lexer.StrToken{Str: test[1]},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestEqOfTwoStringsReturnsCorrectResult(t *testing.T) {
	strs := [][]string{
		{"foo", "bar"},
		{"baz", "Baz"},
		{"fizzbuzz", "fizzbuzz"},
		{"\u2366", "\u2366"},
	}

	for _, test := range strs {
		expected := test[0] == test[1]

		head := generateTree([]lexer.Token{
			lexer.StrToken{Str: test[0]},
			lexer.RelOpToken{Operator: "=="},
			lexer.StrToken{Str: test[1]},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestNotEqOfTwoStringsReturnsCorrectResult(t *testing.T) {
	strs := [][]string{
		{"foo", "bar"},
		{"baz", "Baz"},
		{"fizzbuzz", "fizzbuzz"},
	}

	for _, test := range strs {
		expected := test[0] != test[1]

		head := generateTree([]lexer.Token{
			lexer.StrToken{Str: test[0]},
			lexer.RelOpToken{Operator: "!="},
			lexer.StrToken{Str: test[1]},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestEqOfTwoBoolsReturnsCorrectResult(t *testing.T) {
	bools := [][]bool{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}

	for _, test := range bools {
		expected := test[0] == test[1]

		head := generateTree([]lexer.Token{
			lexer.BoolToken{Value: strconv.FormatBool(test[0])},
			lexer.RelOpToken{Operator: "=="},
			lexer.BoolToken{Value: strconv.FormatBool(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}
func TestNotEqOfTwoBoolsReturnsCorrectResult(t *testing.T) {
	bools := [][]bool{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}

	for _, test := range bools {
		expected := test[0] != test[1]

		head := generateTree([]lexer.Token{
			lexer.BoolToken{Value: strconv.FormatBool(test[0])},
			lexer.RelOpToken{Operator: "!="},
			lexer.BoolToken{Value: strconv.FormatBool(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}
func TestLogicalAndOfTwoBoolsReturnsCorrectResult(t *testing.T) {
	bools := [][]bool{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}

	for _, test := range bools {
		expected := test[0] && test[1]

		head := generateTree([]lexer.Token{
			lexer.BoolToken{Value: strconv.FormatBool(test[0])},
			lexer.LogicOpToken{Operator: "&&"},
			lexer.BoolToken{Value: strconv.FormatBool(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func TestLogicalOrOfTwoBoolsReturnsCorrectResult(t *testing.T) {
	bools := [][]bool{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}

	for _, test := range bools {
		expected := test[0] || test[1]

		head := generateTree([]lexer.Token{
			lexer.BoolToken{Value: strconv.FormatBool(test[0])},
			lexer.LogicOpToken{Operator: "||"},
			lexer.BoolToken{Value: strconv.FormatBool(test[1])},
		})

		resultChan := runProcess(head)
		result := <-resultChan

		if result.String() != fmt.Sprint(expected) {
			t.Errorf("expected result to be %s, got %s", fmt.Sprint(expected), result)
		}
	}
}

func runProcess(head parser.TreeNode) chan Result {
	nodeChan := getNodeChan([]parser.TreeNode{head})
	resultChan := make(chan Result)

	context := &Context{}

	go Process(nodeChan, resultChan, context)

	return resultChan
}

func generateStringTree(str string) parser.TreeNode {
	tok := lexer.StrToken{Str: str}
	return generateTree([]lexer.Token{tok})
}

func generateNumTree(num int) parser.TreeNode {
	tok := lexer.NumToken{Num: strconv.Itoa(num)}
	return generateTree([]lexer.Token{tok})
}

func generateTree(tok []lexer.Token) parser.TreeNode {
	tokChan := make(chan []lexer.Token)
	nodeChan := make(chan parser.TreeNode)
	errChan := make(chan error)

	go parser.Parse(tokChan, nodeChan, errChan)
	go func() {
		defer close(tokChan)
		tok = append([]lexer.Token{lexer.BlockToken{Block: "{{"}}, tok...)
		tok = append(tok, []lexer.Token{lexer.BlockToken{Block: "}}"}, lexer.EOLToken{}}...)
		tokChan <- tok
	}()

	return <-nodeChan
}
