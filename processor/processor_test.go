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
		lexer.AddOpToken{},
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
		lexer.SubOpToken{},
		lexer.NumToken{Num: "10"},
		lexer.AddOpToken{},
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
		lexer.SubOpToken{},
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
		lexer.SubOpToken{},
		lexer.NumToken{Num: "10"},
		lexer.SubOpToken{},
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
		lexer.SubOpToken{},
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
		lexer.SubOpToken{},
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
		lexer.SubOpToken{},
		lexer.NumToken{Num: "13"},
		lexer.MultOpToken{Operator: "*"},
		lexer.SubOpToken{},
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
		lexer.NegationOpToken{Operator: "!"},
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
		lexer.NegationOpToken{Operator: "!"},
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

func TestAssignmentAddsValueIntoStore(t *testing.T) {
	partial := []lexer.Token{
		lexer.IdentToken{Identifier: "foo"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "title"},
		lexer.AssignOpToken{Operator: "="},
	}

	vals := []lexer.Token{
		lexer.StrToken{Str: "fizzbuzz"},
		lexer.NumToken{Num: "100"},
	}

	for _, val := range vals {
		head := generateTree(append(partial, val))
		resultChan, exportStore := runProcessWithExportStore(head, "foo")

		<-resultChan

		expected := val.GetValue()
		context := exportStore.GetFileContext("foo")
		if got, ok := context.At("foo.title"); !ok || got.result == nil || got.result.String() != expected {
			t.Errorf("expected export to contain key %q and value %q", "foo.title", expected)
		}
	}
}

func TestProcessedVarNodeReturnsContextValue(t *testing.T) {
	context := &Context{"foo": &ContextNode{result: StringResult("val")}}
	head := generateTree([]lexer.Token{
		lexer.IdentToken{Identifier: "foo"},
	})

	resultChan := runProcessWithContext(head, context)
	result := <-resultChan

	if result.String() != "val" {
		t.Errorf("expected result to be val, got %s", result.String())
	}
}

func TestProcessedVarNodeReturnsNestedContextValue(t *testing.T) {
	expectedResult := StringResult("fizzbuzz")
	context := &Context{
		"foo": &ContextNode{child: &Context{
			"bar": &ContextNode{child: &Context{
				"baz": &ContextNode{result: expectedResult},
			}},
		}},
	}
	head := generateTree([]lexer.Token{
		lexer.IdentToken{Identifier: "foo"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "bar"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "baz"},
	})

	resultChan := runProcessWithContext(head, context)
	result := <-resultChan

	if result.String() != expectedResult.String() {
		t.Errorf("expected result to be %s, got %s", expectedResult.String(), result.String())
	}
}

func TestProcessedVarNodeReturnsContainerValue(t *testing.T) {
	expectedResult := StringResult(fmt.Sprintf("%T", ContainerResult{}))

	context := &Context{
		"foo": &ContextNode{child: &Context{
			"bar": &ContextNode{child: &Context{
				"baz": &ContextNode{result: expectedResult},
			}},
		}},
	}

	head := generateTree([]lexer.Token{
		lexer.IdentToken{Identifier: "foo"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "bar"},
	})

	resultChan := runProcessWithContext(head, context)
	result := <-resultChan

	if result.String() != expectedResult.String() {
		t.Errorf("expected result to be %s, got %s", expectedResult.String(), result.String())
	}
}

func TestTrueIfReturnsIfBody(t *testing.T) {
	expected := "the if body"
	condition := []lexer.Token{lexer.BoolToken{Value: "true"}}
	body := []lexer.Token{lexer.PassthroughToken{Value: expected}}
	ifTokens := generateIfTokens(condition, body, true)
	head := generateTree(ifTokens)

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestTrueIfReturnsIfMultilineBody(t *testing.T) {
	expected := "the if body\nmore if body"
	condition := []lexer.Token{lexer.BoolToken{Value: "true"}}
	body := []lexer.Token{
		lexer.PassthroughToken{Value: "the if body\n"},
		lexer.PassthroughToken{Value: "more if body"},
	}

	head := generateTree(generateIfTokens(condition, body, true))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestFalseIfReturnsEmptyBody(t *testing.T) {
	expected := ""
	condition := []lexer.Token{lexer.BoolToken{Value: "false"}}
	body := []lexer.Token{lexer.StrToken{Str: "the if body"}}

	head := generateTree(generateIfTokens(condition, body, true))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestTrueIfDoesNotReturnElseBody(t *testing.T) {
	expected := "the if body"
	condition := []lexer.Token{lexer.BoolToken{Value: "true"}}
	ifBody := []lexer.Token{lexer.PassthroughToken{Value: expected}}
	ifToks := generateIfTokens(condition, ifBody, false)
	elseBody := []lexer.Token{lexer.PassthroughToken{Value: "the else body"}}
	elseToks := generateElseTokens(elseBody)

	head := generateTree(append(ifToks, elseToks...))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}

}

func TestFalseIfReturnsElseBody(t *testing.T) {
	expected := "the else body"
	condition := []lexer.Token{lexer.BoolToken{Value: "false"}}
	ifBody := []lexer.Token{lexer.PassthroughToken{Value: "the if body"}}
	ifToks := generateIfTokens(condition, ifBody, false)
	elseBody := []lexer.Token{lexer.PassthroughToken{Value: expected}}
	elseToks := generateElseTokens(elseBody)

	head := generateTree(append(ifToks, elseToks...))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestMultipleTrueElseIfReturnsFirstTrue(t *testing.T) {
	expected := "this is the one"
	ifCondition := []lexer.Token{lexer.BoolToken{Value: "false"}}
	ifBody := []lexer.Token{lexer.PassthroughToken{Value: "the if body"}}
	ifToks := generateIfTokens(ifCondition, ifBody, false)

	elseIfConditions := [][]lexer.Token{
		{lexer.BoolToken{Value: "false"}},
		{lexer.BoolToken{Value: "true"}},
		{lexer.BoolToken{Value: "true"}},
	}

	elseIfBodies := [][]lexer.Token{
		{lexer.PassthroughToken{Value: "a"}},
		{lexer.PassthroughToken{Value: expected}},
		{lexer.PassthroughToken{Value: "a"}},
	}

	elseIfToks := generateElseIfTokens(elseIfConditions, elseIfBodies, true)

	head := generateTree(append(ifToks, elseIfToks...))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestMultipleFalseElseIfReturnsEmptyString(t *testing.T) {
	expected := ""
	ifCondition := []lexer.Token{lexer.BoolToken{Value: "false"}}
	ifBody := []lexer.Token{lexer.PassthroughToken{Value: "the if body"}}
	ifToks := generateIfTokens(ifCondition, ifBody, false)

	elseIfConditions := [][]lexer.Token{
		{lexer.BoolToken{Value: "false"}},
		{lexer.BoolToken{Value: "false"}},
		{lexer.BoolToken{Value: "false"}},
	}

	elseIfBodies := [][]lexer.Token{
		{lexer.PassthroughToken{Value: "a"}},
		{lexer.PassthroughToken{Value: "b"}},
		{lexer.PassthroughToken{Value: "c"}},
	}

	elseIfToks := generateElseIfTokens(elseIfConditions, elseIfBodies, true)

	head := generateTree(append(ifToks, elseIfToks...))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestMultipleTrueElseIfReturnsTrueIfBody(t *testing.T) {
	expected := "the if body"
	ifCondition := []lexer.Token{lexer.BoolToken{Value: "true"}}
	ifBody := []lexer.Token{lexer.PassthroughToken{Value: expected}}
	ifToks := generateIfTokens(ifCondition, ifBody, false)

	elseIfConditions := [][]lexer.Token{
		{lexer.BoolToken{Value: "true"}},
		{lexer.BoolToken{Value: "true"}},
		{lexer.BoolToken{Value: "true"}},
	}

	elseIfBodies := [][]lexer.Token{
		{lexer.PassthroughToken{Value: "a"}},
		{lexer.PassthroughToken{Value: "b"}},
		{lexer.PassthroughToken{Value: "c"}},
	}

	elseIfToks := generateElseIfTokens(elseIfConditions, elseIfBodies, true)

	head := generateTree(append(ifToks, elseIfToks...))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestMultipleFalseElseIfReturnsElseBody(t *testing.T) {
	expected := "the else body"
	ifCondition := []lexer.Token{lexer.BoolToken{Value: "false"}}
	ifBody := []lexer.Token{lexer.PassthroughToken{Value: "the if body"}}
	ifToks := generateIfTokens(ifCondition, ifBody, false)

	elseIfConditions := [][]lexer.Token{
		{lexer.BoolToken{Value: "false"}},
		{lexer.BoolToken{Value: "false"}},
		{lexer.BoolToken{Value: "false"}},
	}

	elseIfBodies := [][]lexer.Token{
		{lexer.PassthroughToken{Value: "a"}},
		{lexer.PassthroughToken{Value: "b"}},
		{lexer.PassthroughToken{Value: "c"}},
	}

	elseIfToks := generateElseIfTokens(elseIfConditions, elseIfBodies, false)
	elseToks := generateElseTokens([]lexer.Token{lexer.PassthroughToken{Value: expected}})

	head := generateTree(append(ifToks, append(elseIfToks, elseToks...)...))

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected result to be %q, got %q", expected, result.String())
	}
}

func TestForLoopGeneratesCorrectNumberOfLines(t *testing.T) {
	bodyText := "this is the body\n"

	condition := []lexer.Token{
		lexer.IdentToken{Identifier: "foo"},
		lexer.InToken{},
		lexer.StrToken{Str: "bar"},
	}

	body := []lexer.Token{
		lexer.PassthroughToken{Value: bodyText},
	}

	expected := ""
	pathReader := getTestPathReader(3)
	for i := 0; i < len(pathReader("")); i++ {
		expected += bodyText
	}

	forToks := generateForLoopTokens(condition, body)
	head := generateTree(forToks)
	resultChan := runProcessWithGetPathFunc(head, nil, pathReader)

	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected for loop result to be %q, got %q", expected, result.String())
	}
}

func TestForLoopGeneratesCorrectFileContextualBody(t *testing.T) {
	condition := []lexer.Token{
		lexer.IdentToken{Identifier: "foo"},
		lexer.InToken{},
		lexer.StrToken{Str: "bar"},
	}

	body := []lexer.Token{
		lexer.BlockToken{Block: "{{:"},
		lexer.IdentToken{Identifier: "foo"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "date"},
		lexer.BlockToken{Block: "}}"},
	}

	expected := ""
	pathReader := getTestPathReader(3)
	exportStore := &ExportFileStore{FilePath: ""}
	exportStore.Insert([]string{"date"}, StringResult("some-date-someday"))

	for i := 0; i < len(pathReader("")); i++ {
		expected += "some-date-someday\n"
	}

	forToks := generateForLoopTokens(condition, body)
	head := generateTree(forToks)
	resultChan := runProcessWithGetPathFunc(head, exportStore, pathReader)

	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected for loop result to be %q, got %q", expected, result.String())
	}
}

func TestForLoopWithNormalBlockGeneratesEmptyBody(t *testing.T) {
	condition := []lexer.Token{
		lexer.IdentToken{Identifier: "foo"},
		lexer.InToken{},
		lexer.StrToken{Str: "bar"},
	}

	body := []lexer.Token{
		lexer.BlockToken{Block: "{{"},
		lexer.IdentToken{Identifier: "foo"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "date"},
		lexer.BlockToken{Block: "}}"},
	}

	expected := ""
	pathReader := getTestPathReader(3)
	exportStore := &ExportFileStore{FilePath: ""}
	exportStore.Insert([]string{"date"}, StringResult("some-date-someday"))

	for i := 0; i < len(pathReader("")); i++ {
		expected += ""
	}

	forToks := generateForLoopTokens(condition, body)
	head := generateTree(forToks)
	resultChan := runProcessWithGetPathFunc(head, exportStore, pathReader)

	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected for loop result to be %q, got %q", expected, result.String())
	}
}

func TestContentReturnsCorrectResult(t *testing.T) {
	head := &parser.ContentParseNode{}
	next := &parser.ContentParseNode{}

	head.SetChildren([]parser.TreeNode{
		next,
		&parser.StringParseNode{Value: "baz"},
	})

	cur := next

	for i := 1; i < 3; i++ {
		cur.SetChildren([]parser.TreeNode{next})
		cur = next
	}

	str := &parser.StringParseNode{Value: "foobar\n"}
	cur.SetChildren([]parser.TreeNode{str})

	resultChan := runProcess(head)

	result := <-resultChan

	if result.String() != "foobar\nbaz" {
		t.Errorf("expected content to be %q, got %q", "foobar\nbaz", result.String())
	}
}

func TestForLoopGeneratesCorrectContextBody(t *testing.T) {
	expected := "first\nsecond\nthird\n"
	context := &Context{
		"page": &ContextNode{child: &Context{
			"content": &ContextNode{child: &Context{
				"0": &ContextNode{child: &Context{
					"title": &ContextNode{result: StringResult("first")},
				}},
				"1": &ContextNode{child: &Context{
					"title": &ContextNode{result: StringResult("second")},
				}},
				"2": &ContextNode{child: &Context{
					"title": &ContextNode{result: StringResult("third")},
				}},
			}},
		}},
	}

	condition := []lexer.Token{
		lexer.IdentToken{Identifier: "page"},
		lexer.InToken{},
		lexer.IdentToken{Identifier: "page"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "content"},
	}

	body := []lexer.Token{
		lexer.BlockToken{Block: "{{:"},
		lexer.IdentToken{Identifier: "page"},
		lexer.SymbolToken{Symbol: "."},
		lexer.IdentToken{Identifier: "title"},
		lexer.BlockToken{Block: "}}"},
	}

	forToks := generateForLoopTokens(condition, body)
	head := generateTree(forToks)

	resultChan := runProcessWithContext(head, context)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected for loop result to be %q, got %q", expected, result.String())
	}
}

func TestPrintFuncCallReturnsValue(t *testing.T) {
	expected := "1"
	head := generateTree([]lexer.Token{
		lexer.IdentToken{Identifier: "print"},
		lexer.SymbolToken{Symbol: "("},
		lexer.NumToken{Num: "6423"},
		lexer.SubOpToken{},
		lexer.NumToken{Num: "6422"},
		lexer.SymbolToken{Symbol: ")"},
	})

	resultChan := runProcess(head)
	result := <-resultChan

	if result.String() != expected {
		t.Errorf("expected print result to be %q, got %q", expected, result)
	}
}

func runProcess(head parser.TreeNode) chan Result {
	return runProcessWithContext(head, nil)
}

func runProcessWithExportStore(head parser.TreeNode, filename string) (chan Result, ExportStorage) {
	nodeChan := getNodeChan([]parser.TreeNode{head})
	exportStorage := &ExportFileStore{filename}

	context := &Context{}
	processor := &NodeProcessor{Context: context, ExportStore: exportStorage}
	resultChan := make(chan Result)
	go processor.Process(nodeChan, resultChan)

	return resultChan, exportStorage
}

func runProcessWithContext(head parser.TreeNode, context *Context) chan Result {
	nodeChan := getNodeChan([]parser.TreeNode{head})
	resultChan := make(chan Result)

	processor := NodeProcessor{Context: context}
	go processor.Process(nodeChan, resultChan)

	return resultChan
}

func runProcessWithGetPathFunc(head parser.TreeNode, exportStore ExportStorage, getPathFunc func(string) []string) chan Result {
	nodeChan := getNodeChan([]parser.TreeNode{head})
	resultChan := make(chan Result)

	processor := NodeProcessor{Context: &Context{}, ExportStore: exportStore}
	processor.PathReader = getPathFunc
	go processor.Process(nodeChan, resultChan)

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
		// tokens := []lexer.Token{lexer.BlockToken{Block: "{{"}}
		// tokens = append(tokens, tok...)
		tok = append(tok, []lexer.Token{lexer.EOLToken{}}...)
		tokChan <- tok
	}()

	return <-nodeChan
}

func generateIfTokens(condition, body []lexer.Token, includeEnd bool) []lexer.Token {
	ifTokens := append(
		[]lexer.Token{
			lexer.IfToken{},
		},
		condition...,
	)

	ifTokens = append(ifTokens, lexer.BlockToken{Block: "}}"})

	for _, bodyTok := range body {
		ifTokens = append(ifTokens, bodyTok)
	}

	if includeEnd {
		return append(ifTokens, lexer.EndToken{})
	}

	return ifTokens
}

func generateElseTokens(body []lexer.Token) []lexer.Token {
	elseTokens := []lexer.Token{lexer.ElseToken{}}

	for _, bodyTok := range body {
		elseTokens = append(elseTokens, bodyTok)
	}

	return append(elseTokens, lexer.EndToken{})
}

func generateElseIfTokens(conditions, bodies [][]lexer.Token, includeEnd bool) []lexer.Token {
	elseIfTokens := []lexer.Token{}
	for i, condition := range conditions {
		elseIfTokens = append(elseIfTokens, lexer.ElseIfToken{})
		elseIfTokens = append(elseIfTokens, condition...)
		elseIfTokens = append(elseIfTokens, lexer.BlockToken{Block: "}}"})
		elseIfTokens = append(elseIfTokens, bodies[i]...)
	}

	if includeEnd {
		return append(elseIfTokens, lexer.EndToken{})
	}

	return elseIfTokens
}

func generateForLoopTokens(condition []lexer.Token, body []lexer.Token) []lexer.Token {
	forLoopTokens := []lexer.Token{
		lexer.ForToken{},
	}

	forLoopTokens = append(forLoopTokens, condition...)
	forLoopTokens = append(forLoopTokens, lexer.BlockToken{Block: "}}"})
	forLoopTokens = append(forLoopTokens, body...)

	return append(forLoopTokens, lexer.EndToken{})
}

func getTestPathReader(numPaths int) func(string) []string {
	paths := make([]string, numPaths)
	return func(string) []string { return paths }
}
