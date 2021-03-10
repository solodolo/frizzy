package lexer

import (
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
		{"-", "SubOpToken"},
		{"!", "NegationOpToken"},
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
		{"post", "IdentToken"},
		{"true", "BoolToken"},
		{"false", "BoolToken"},
		{";", "SymbolToken"},
		{"(", "SymbolToken"},
		{")", "SymbolToken"},
		{"{{", "BlockToken"},
	}

	for _, test := range tests {
		lexer := Lexer{}
		tok, _ := lexer.getNextBlockToken(InputLine{line: test.symbol})
		// Get type of each token and trim off package name
		tokType := strings.TrimPrefix(fmt.Sprintf("%T", tok), "lexer.")

		if tokType != test.tokType {
			t.Errorf("Expected %q to return type %s. Got type %s.", test.symbol, test.tokType, tokType)
		}
	}
}

func TestProcessLineReturnsCorrectTokens(t *testing.T) {
	var tests = []struct {
		line       string
		tokenTypes []string
	}{
		{`<html></html>`, []string{"PassthroughToken"}},
		{`<html>{{"blah"}}</html>`, []string{"PassthroughToken", "BlockToken", "StrToken", "BlockToken", "EOLToken", "PassthroughToken"}},
		{`{{: "foo" }}`, []string{"BlockToken", "StrToken", "BlockToken", "EOLToken"}},
		{`{{ !a.b }}</html>`, []string{"BlockToken", "NegationOpToken", "IdentToken", "SymbolToken", "IdentToken", "BlockToken", "EOLToken", "PassthroughToken"}},
		{`<html>"blah"}}</html>`, []string{"PassthroughToken"}},
		{`<html>"blah"{{ print()`, []string{"PassthroughToken", "BlockToken", "IdentToken", "SymbolToken", "SymbolToken"}},
		{`{{ foo(a,b)`, []string{"BlockToken", "IdentToken", "SymbolToken", "IdentToken", "SymbolToken", "IdentToken", "SymbolToken"}},
		{"{{a\nb}}", []string{"BlockToken", "IdentToken", "IdentToken", "BlockToken", "EOLToken"}},
		{"a{{a\nb}}", []string{"PassthroughToken", "BlockToken", "IdentToken", "IdentToken", "BlockToken", "EOLToken"}},
		{"{{a\nb}}c", []string{"BlockToken", "IdentToken", "IdentToken", "BlockToken", "EOLToken", "PassthroughToken"}},
		{"{{if(true)}}", []string{"BlockToken", "IfToken", "SymbolToken", "BoolToken", "SymbolToken", "BlockToken", "EOLToken"}},
		{"{{if else_if else end}}", []string{"BlockToken", "IfToken", "ElseIfToken", "ElseToken", "EndToken", "BlockToken", "EOLToken"}},
	}

	for _, test := range tests {
		lineChan := make(chan InputLine)
		lexer := Lexer{lineChan: lineChan}

		go func(lineChan chan InputLine) {
			defer close(lineChan)
			lineChan <- InputLine{line: test.line}
		}(lineChan)

		tokChan, errChan := lexer.processLines()

		tokens := []Token{}

		for toks := range tokChan {
			tokens = append(tokens, toks...)
		}

		err := <-errChan

		if err != nil {
			t.Errorf("expected no errors, got %q", err)
		} else if equal, gotTypes := tokenTypesAreEqual(test.tokenTypes, tokens); !equal {
			t.Errorf("expected %v types, got %v", test.tokenTypes, gotTypes)
		}
	}
}

func TestLexHandlesReadFailure(t *testing.T) {
	pipeReader, _ := io.Pipe()
	pipeReader.Close()
	lexer := Lexer{}

	tokChan, errChan := lexer.Lex(pipeReader)
	expected := "lexer read error line 1: io: read/write on closed pipe"

	for tokens := range tokChan {
		t.Errorf("expecting error, got tokens %v.", tokens)
	}

	err := <-errChan
	if err == nil {
		t.Error("expected error, got nil.")
	} else if err.Error() != expected {
		t.Errorf("expecting error %q, got %q.", expected, err.Error())
	}
}

func TestLexReturnsCorrectTokenTypes(t *testing.T) {
	var tests = []struct {
		lines    string
		expected []string
	}{
		{
			"first {{a}}\nsecond", []string{
				"PassthroughToken", "BlockToken", "IdentToken", "BlockToken", "EOLToken",
				"PassthroughToken", "PassthroughToken",
			},
		},
		{
			"{{a-}\nsecond", []string{
				"BlockToken", "IdentToken", "BlockToken", "EOLToken", "PassthroughToken",
			},
		},
		{
			"{{a-}        \nsecond", []string{
				"BlockToken", "IdentToken", "BlockToken", "EOLToken", "PassthroughToken",
			},
		},
		{
			"{{a-}          second", []string{
				"BlockToken", "IdentToken", "BlockToken", "EOLToken", "PassthroughToken",
			},
		},
		{
			"foo bar", []string{"PassthroughToken"},
		},
		{
			"{{\nfor a in b", []string{"BlockToken", "ForToken", "IdentToken", "InToken", "IdentToken"},
		},
		{
			`{{: "Foo"}}`, []string{"BlockToken", "StrToken", "BlockToken", "EOLToken"},
		},
		{
			"{{: \"Foo\"}}\n", []string{"BlockToken", "StrToken", "BlockToken", "EOLToken", "PassthroughToken"},
		},
		{
			"{{: title}}\n{{: title}}\n{{: title}}", []string{
				"BlockToken", "IdentToken", "BlockToken", "EOLToken", "PassthroughToken",
				"BlockToken", "IdentToken", "BlockToken", "EOLToken", "PassthroughToken",
				"BlockToken", "IdentToken", "BlockToken", "EOLToken",
			},
		},
		{
			"{{print(a)\n}}blah", []string{
				"BlockToken", "IdentToken", "SymbolToken", "IdentToken", "SymbolToken",
				"BlockToken", "EOLToken", "PassthroughToken",
			},
		},
		{
			"{{`multi\nline\nstr`}}", []string{
				"BlockToken", "StrToken", "BlockToken", "EOLToken",
			},
		},
		{
			"{{post.", []string{"BlockToken", "IdentToken", "SymbolToken"},
		},
		{
			"{{(a < b)", []string{"BlockToken", "SymbolToken", "IdentToken", "RelOpToken", "IdentToken", "SymbolToken"},
		},
		{
			"", []string{},
		},
		{
			"{{a.b && b.a", []string{
				"BlockToken", "IdentToken", "SymbolToken", "IdentToken", "LogicOpToken", "IdentToken", "SymbolToken", "IdentToken",
			},
		},
		{
			"{{foo || false", []string{"BlockToken", "IdentToken", "LogicOpToken", "BoolToken"},
		},
		{
			"{{ title = `this is the title`;\n  date= `2021-03-08`; }}\n# This is a test", []string{
				"BlockToken", "IdentToken", "AssignOpToken", "StrToken", "SymbolToken", "IdentToken",
				"AssignOpToken", "StrToken", "SymbolToken", "BlockToken", "EOLToken", "PassthroughToken",
				"PassthroughToken",
			},
		},
	}

	for _, test := range tests {
		reader := strings.NewReader(test.lines)
		got := []Token{}
		lexer := Lexer{}

		tokChan, errChan := lexer.Lex(reader)

		for tokens := range tokChan {
			got = append(got, tokens...)
		}

		err := <-errChan

		if err != nil {
			t.Errorf("expected no errors, got %q", err.Error())
		}

		if len(got) != len(test.expected) {
			t.Errorf("expected %d tokens, got %d.", len(test.expected), len(got))
		}

		equal, tokTypes := tokenTypesAreEqual(test.expected, got)
		if !equal {
			t.Errorf("expected %q to return %v, got %v.", test.lines, test.expected, tokTypes)
		}
	}
}

func TestTokensAreAssignedCorrectLineNum(t *testing.T) {
	tests := []struct {
		lines    string
		lineNums []int
	}{
		{"first {{a}}\nsecond", []int{1, 1, 1, 1, 1, 1, 2}},
		{"foo bar", []int{1}},
		{"{{\nfor a in b", []int{1, 2, 2, 2, 2}},
		{`{{: "Foo"}}`, []int{1, 1, 1}},
		{"{{: post}}\n{{: post}}\n{{: post}}", []int{1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 3, 3, 3, 3}},
		{"{{print(a)\n\n}}blah", []int{1, 1, 1, 1, 1, 3, 3, 3}},
		{"a\nb\nc\n\n\nf", []int{1, 2, 3, 4, 5, 6}},
	}

	for testNum, test := range tests {
		reader := strings.NewReader(test.lines)
		lexer := Lexer{}

		tokChan, _ := lexer.Lex(reader)

		got := []Token{}
		for tokens := range tokChan {
			got = append(got, tokens...)
		}

		for i := range test.lineNums {
			if i >= len(got) {
				t.Errorf("%d: expected token with line number %d, but none found", testNum, test.lineNums[i])
			} else if test.lineNums[i] != got[i].GetLineNum() {
				t.Errorf("%d: expected token %q to have line number %d, got %d",
					testNum, got[i].GetValue(), test.lineNums[i], got[i].GetLineNum())
			}
		}
	}
}

// Helper function to compare the token types of a to the
// token types in b.
// Returns the comparison result and the found token types
func tokenTypesAreEqual(expected []string, got []Token) (bool, []string) {
	equal := len(expected) == len(got)
	tokenTypes := make([]string, len(got))

	for i := range got {
		tokenTypes[i] = strings.TrimPrefix(fmt.Sprintf("%T", got[i]), "lexer.")
		if equal {
			equal = equal && (expected[i] == tokenTypes[i])
		}
	}

	return equal, tokenTypes
}

func TestRawStringReturnsCorrectLexResultFromParam(t *testing.T) {
	var tests = []struct {
		tokText           string
		expectedRemaining string
	}{
		{"a single line string", "this is the remaining text"},
		{"", ""},
		{"", "foobar"},
		{"somestr", ""},
		{`a
		b
		c`, "this is the remaining text"},
		{`
		
		`, "more remaining"},
		{`
		
		`, ""},
	}

	for i, test := range tests {
		lineChan := make(chan InputLine)
		lexer := Lexer{lineChan: lineChan, state: inStr}
		inputLine := InputLine{line: fmt.Sprintf("`%s`%s", test.tokText, test.expectedRemaining)}

		tok, remaining := lexer.getRawStringToken(inputLine)

		if tok.GetValue() != test.tokText {
			t.Errorf("failed test %d: expected token text %q, got %q", i, test.tokText, tok)
		} else if remaining.line != test.expectedRemaining {
			t.Errorf("failed test %d: expected remaining text %q, got %q", i, test.expectedRemaining, remaining.line)
		}
	}
}

func TestRawStringReturnsCorrectLexResultFromChan(t *testing.T) {
	var tests = []struct {
		tokText           string
		expectedRemaining string
	}{
		{"a single line string", "this is the remaining text"},
		{"", ""},
		{"", "foobar"},
		{"somestr", ""},
		{`a
		b
		c`, "this is the remaining text"},
		{`
	
		`, "more remaining"},
		{`
	
		`, ""},
	}

	for i, test := range tests {
		lineChan := make(chan InputLine)
		lexer := Lexer{lineChan: lineChan, state: inStr}
		inputLine := InputLine{line: fmt.Sprintf("`%s`%s", test.tokText, test.expectedRemaining)}

		go func(lineChan chan InputLine, inputLine InputLine) {
			lineChan <- inputLine
		}(lineChan, inputLine)

		tok, remaining := lexer.getRawStringToken(InputLine{})

		if tok.GetValue() != test.tokText {
			t.Errorf("failed test %d: expected token text %q, got %q", i, test.tokText, tok)
		} else if remaining.line != test.expectedRemaining {
			t.Errorf("failed test %d: expected remaining text %q, got %q", i, test.expectedRemaining, remaining.line)
		}
	}
}
