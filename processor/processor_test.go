package processor

import (
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

	if result.String() != "7" {
		t.Errorf("expected result to be 7, got %s", result.String())
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
