package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestGetLineTokensReturnsCorrectTokenType(t *testing.T) {
	var tests = []struct {
		symbol  string
		tokType string
	}{
		{"*", "MultOpToken"},
		{"/", "MultOpToken"},
		{"%", "MultOpToken"},
		{"+", "AddOpToken"},
		{"-", "AddOpToken"},
		{"!", "UnaryOpToken"},
		{"==", "RelOpToken"},
		{"!=", "RelOpToken"},
		{"<=", "RelOpToken"},
		{">=", "RelOpToken"},
		{"<", "RelOpToken"},
		{">", "RelOpToken"},
		{"=", "AssignOpToken"},
		{"||", "LogicOpToken"},
		{"&&", "LogicOpToken"},
		{"123", "NumToken"},
		{`"foobar"`, "StrToken"},
		{"for", "ForToken"},
		{"if", "IfToken"},
		{"else_if", "ElseIfToken"},
		{"else", "ElseToken"},
		{"end", "EndToken"},
		{"post.title", "VarToken"},
		{"true", "BoolToken"},
		{"false", "BoolToken"},
		{";", "SymbolToken"},
		{"(", "SymbolToken"},
		{")", "SymbolToken"},
		{"{{", "BlockToken"},
	}

	for _, test := range tests {
		tokens := getLineTokens(test.symbol, 1)
		// Get type of each token and trim off package name
		tokType := strings.TrimPrefix(fmt.Sprintf("%T", tokens[0]), "lexer.")

		if tokType != test.tokType {
			t.Errorf("Expected %q to return type %s. Got type %s.", test.symbol, test.tokType, tokType)
		}
	}
}

func TestGetLineTokensReturnsCorrectTokensForLine(t *testing.T) {
	var tests = []struct {
		line       string
		tokenTypes []string
	}{
		{"post.", []string{"IdentToken", "PassthroughToken"}},
		{"(a < b)", []string{"SymbolToken", "IdentToken", "RelOpToken", "IdentToken", "SymbolToken"}},
		{"", []string{}},
		{"a.b && b.a", []string{"VarToken", "LogicOpToken", "VarToken"}},
		{"foo || false", []string{"IdentToken", "LogicOpToken", "BoolToken"}},
	}

	for _, test := range tests {
		tokens := getLineTokens(test.line, 1)
		equal, tokenTypes := tokenTypesAreEqual(tokens, test.tokenTypes)
		if !equal {
			t.Errorf("Expected %q to return %v. Got %v.", test.line, test.tokenTypes, tokenTypes)
		}
	}
}

func TestBlockIndicesReturnsCorrectIndices(t *testing.T) {
	var tests = []struct {
		line  string
		start int
		end   int
	}{
		{"{{abcde}}", 0, 9},
		{"a{{bcd}}", 1, 8},
		{"a{{bcd}}e", 1, 8},
		{"{{bcd}}e", 0, 7},
		{"{{abcde", 0, -1},
		{"abcde}}", -1, 7},
		{"ab{{cd", 2, -1},
		{"ab}}cd", -1, 4},
	}

	for _, test := range tests {
		start, end := getBlockIndices(test.line)

		if start != test.start || end != test.end {
			t.Errorf("Expected %q to return %d, %d. Got %d, %d.", test.line, test.start, test.end, start, end)
		}
	}
}

func TestProcessLineReturnsCorrectTokens(t *testing.T) {
	var tests = []struct {
		line       string
		tokenTypes []string
	}{
		{`<html></html>`, []string{"PassthroughToken"}},
		{`<html>{{"blah"}}</html>`, []string{"PassthroughToken", "BlockToken", "StrToken", "BlockToken", "PassthroughToken"}},
		{`{{: "foo" }}`, []string{"BlockToken", "StrToken", "BlockToken"}},
		{`{{ !a.b }}</html>`, []string{"BlockToken", "UnaryOpToken", "VarToken", "BlockToken", "PassthroughToken"}},
		{`<html>"blah"}}</html>`, []string{"PassthroughToken"}},
		{`<html>"blah"{{ print()`, []string{"PassthroughToken", "BlockToken", "IdentToken", "SymbolToken", "SymbolToken"}},
		{`{{ foo(a,b)`, []string{"BlockToken", "IdentToken", "SymbolToken", "IdentToken", "SymbolToken", "IdentToken", "SymbolToken"}},
		{"{{a\nb}}", []string{"BlockToken", "IdentToken", "IdentToken", "BlockToken"}},
		{"a{{a\nb}}", []string{"PassthroughToken", "BlockToken", "IdentToken", "IdentToken", "BlockToken"}},
		{"{{a\nb}}c", []string{"BlockToken", "IdentToken", "IdentToken", "BlockToken", "PassthroughToken"}},
		{"{{if(true)}}", []string{"BlockToken", "IfToken", "SymbolToken", "BoolToken", "SymbolToken", "BlockToken"}},
		{"{{if else_if else end}}", []string{"BlockToken", "IfToken", "ElseIfToken", "ElseToken", "EndToken", "BlockToken"}},
	}

	for _, test := range tests {
		tokens, _ := processLine(InputLine{line: test.line}, false) // Not in open block
		equal, tokenTypes := tokenTypesAreEqual(tokens, test.tokenTypes)

		if !equal {
			t.Errorf("Expected %q to return %v. Got %v.", test.line, test.tokenTypes, tokenTypes)
		}
	}
}

func TestLexHandlesReadFailure(t *testing.T) {
	pipeReader, _ := io.Pipe()
	scanner := bufio.NewScanner(pipeReader)
	pipeReader.Close()

	tokChan := make(chan []Token)
	errChan := make(chan error)

	go Lex(scanner, tokChan, errChan)
	expected := "error reading lines for lexing: io: read/write on closed pipe"

	for tokens := range tokChan {
		t.Errorf("Expecting error. Got tokens %v.", tokens)
	}

	err := <-errChan
	if err == nil {
		t.Error("Expected error. Got nil.")
	} else if err.Error() != expected {
		t.Errorf("Expecting error %q. Got %q.", expected, err.Error())
	}
}

func TestLexReturnsCorrectTokens(t *testing.T) {
	var tests = []struct {
		lines    string
		expected [][]string
	}{
		{
			"first {{a}}\nsecond", [][]string{
				{"PassthroughToken", "BlockToken", "IdentToken", "BlockToken"},
				{"PassthroughToken"},
			},
		},
		{
			"foo bar", [][]string{{"PassthroughToken"}},
		},
		{
			"{{\nfor a in b", [][]string{
				{"BlockToken"},
				{"ForToken", "IdentToken", "InToken", "IdentToken"},
			},
		},
		{
			`{{: "Foo"}}`, [][]string{
				{"BlockToken", "StrToken", "BlockToken"},
			},
		},
		{
			"{{: post.title}}\n{{: post.title}}\n{{: post.title}}", [][]string{
				{"BlockToken", "VarToken", "BlockToken"},
				{"BlockToken", "VarToken", "BlockToken"},
				{"BlockToken", "VarToken", "BlockToken"},
			},
		},
		{
			"{{print(a)\n}}blah", [][]string{
				{"BlockToken", "IdentToken", "SymbolToken", "IdentToken", "SymbolToken"},
				{"BlockToken", "PassthroughToken"},
			},
		},
	}

	for _, test := range tests {
		scanner := bufio.NewScanner(strings.NewReader(test.lines))
		got := [][]Token{}
		tokChan := make(chan []Token)
		errChan := make(chan error)
		go Lex(scanner, tokChan, errChan)

		for tokens := range tokChan {
			got = append(got, tokens)
		}

		err := <-errChan

		if err != nil {
			t.Errorf("Expected no errors. Got %q", err.Error())
		}

		if len(got) != len(test.expected) {
			t.Errorf("Expected %d lines of tokens. Got %d.", len(test.expected), len(got))
		}

		for i, toks := range got {
			equal, tokTypes := tokenTypesAreEqual(toks, test.expected[i])
			if !equal {
				t.Errorf("Expected %q to return %v. Got %v.", test.lines, test.expected[i], tokTypes)
			}
		}
	}
}

func TestTokensAreAssignedCorrectLineNum(t *testing.T) {
	tests := []struct {
		lines    string
		lineNums int
	}{
		{"first {{a}}\nsecond", 2},
		{"foo bar", 1},
		{"{{\nfor a in b", 2},
		{`{{: "Foo"}}`, 1},
		{"{{: post.title}}\n{{: post.title}}\n{{: post.title}}", 3},
		{"{{print(a)\n}}blah", 2},
		{"a\nb\nc\n\n\nf", 6},
	}

	for _, test := range tests {
		scanner := bufio.NewScanner(strings.NewReader(test.lines))
		tokChan := make(chan []Token)
		errChan := make(chan error)
		go Lex(scanner, tokChan, errChan)

		got := [][]Token{}
		for tokens := range tokChan {
			got = append(got, tokens)
		}

		for i := 1; i <= test.lineNums; i++ {
			for _, token := range got[i-1] {
				if token.GetLineNum() != i {
					t.Errorf("Expected token to have line number %d. Got %d.", i, token.GetLineNum())
				}
			}
		}
	}
}

// Helper function to compare the token types of a to the
// token types in b.
// Returns the comparison result and the found token types
func tokenTypesAreEqual(a []Token, b []string) (bool, []string) {
	equal := len(a) == len(b)
	tokenTypes := make([]string, len(a))

	if equal {
		for i := range a {
			tokenTypes[i] = strings.TrimPrefix(fmt.Sprintf("%T", a[i]), "lexer.")
			equal = equal && (b[i] == tokenTypes[i])
		}
	}

	return equal, tokenTypes
}
