package lexer

import (
	"bufio"
	"regexp"
)

// Theses are the tokens of our grammar
// MUL_OPERATORS : *, /, %
// ADD_OPERATORS : +, -
// REL_OPERATORS : <, >, !=, ==, <=, >=
// LOGIC_OPERATORS : ||, &&
// ASSIGNMENT: =
// UNARY_OPERATORS: !, -
// ID : [a-zA-Z]+[a-zA-Z0-9_]*
// VAR_NAME : ([a-zA-Z]+[a-zA-Z0-9_]*)(\.[a-zA-Z][a-zA-Z0-9_]*)*
// STRING : "[^”]*"
// NUM : [0-9]+
// IF : if
// ELSE_IF : else_if
// ELSE : else
// FOR : for
// IN : in
// END : end
// FALSE : false
// TRUE : true
// SEMI : ;
// L_PAREN : (
// R_PAREN : )
// OPEN_BLOCK : {{
// PRINT_OPEN : {{:
// CLOSE_BLOCK : }}
// OTHER

var multOp = regexp.MustCompile(`^[*\/%]`)
var addOp = regexp.MustCompile(`^[+-]`)
var relOp = regexp.MustCompile(`^(>=|<=|!=|==|<|>)`)
var logicOp = regexp.MustCompile(`^(\|\||&&)`)
var assignOp = regexp.MustCompile(`^=`)
var unaryOp = regexp.MustCompile(`^!`)

var identExp = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_]*`)

var strExp = regexp.MustCompile(`(?m)^"[^"]*"`)
var numExp = regexp.MustCompile(`^[0-9]+`)

var ifExp = regexp.MustCompile(`^if`)
var elseIfExp = regexp.MustCompile(`^else_if`)
var elseExp = regexp.MustCompile(`^else`)

var forExp = regexp.MustCompile(`^for`)
var inExp = regexp.MustCompile(`^in`)
var endExp = regexp.MustCompile(`^end`)

var boolExp = regexp.MustCompile(`^(true|false)`)

var symbolExp = regexp.MustCompile(`^[;(),\.]`)

var blockExp = regexp.MustCompile(`^({{:|{{|}})`)

var whitespaceExp = regexp.MustCompile(`^\s+`)

// Lexer states
type LexerState int

const (
	passthrough LexerState = iota
	inBlock
	inStr
)

type InputLine struct {
	line    string
	lineNum int
}

// Lexer turns a stream of text lines into a stream of tokens
type Lexer struct {
	lineChan <-chan InputLine
	state    LexerState
}

func (receiver *Lexer) Lex(inputReader *bufio.Reader) (<-chan []Token, <-chan error) {
	inputBuffer := bufio.NewScanner(inputReader)
	inputBuffer.Split(bufio.ScanLines)

	lineChan, lineErrChan := readLines(inputBuffer)

	receiver.lineChan = lineChan
	receiver.state = passthrough

	tokChan, tokErrChan := receiver.processLines()
	lexerErrChan := make(chan error, 1)

	go func() {
		defer close(lexerErrChan)

		select {
		case err := <-lineErrChan:
			lexerErrChan <- err
		case err := <-tokErrChan:
			lexerErrChan <- err
		}
	}()

	return tokChan, lexerErrChan
}

func readLines(inputBuffer *bufio.Scanner) (<-chan InputLine, <-chan error) {
	lineChan := make(chan InputLine)
	errChan := make(chan error, 1)

	go func() {
		defer close(lineChan)
		defer close(errChan)

		i := 1
		for inputBuffer.Scan() {
			line := inputBuffer.Text() + "\n"
			lineChan <- InputLine{line: line, lineNum: i}
			i++
		}

		if inputBuffer.Err() != nil {
			errChan <- inputBuffer.Err()
		}
	}()

	return lineChan, errChan
}

func (receiver *Lexer) getLineToProcess(remaining InputLine) (InputLine, bool) {
	if remaining.line == "" {
		newLine, ok := <-receiver.lineChan
		return newLine, ok
	}

	return remaining, true
}

func (receiver *Lexer) processLines() (<-chan []Token, <-chan error) {
	tokChan := make(chan []Token)
	errChan := make(chan error, 1)

	go func(receiver *Lexer) {
		defer close(tokChan)
		defer close(errChan)

		inputLine, ok := receiver.getLineToProcess(InputLine{})
		for ok {
			if receiver.state == passthrough {
				tok, remaining := receiver.processPassthroughTokens(inputLine)
				tokChan <- []Token{tok}
				inputLine = remaining
			} else if receiver.state == inBlock {
				toks, remaining := receiver.processTokensInBlock(inputLine)
				tokChan <- toks
				inputLine = remaining
			}

			inputLine, ok = receiver.getLineToProcess(inputLine)
		}
	}(receiver)

	return tokChan, errChan

}

func (receiver *Lexer) processPassthroughTokens(inputLine InputLine) (PassthroughToken, InputLine) {
	openBlockExp := regexp.MustCompile(`{{:|{{`)
	if loc := openBlockExp.FindStringIndex(inputLine.line); loc != nil {
		tok := PassthroughToken{Value: inputLine.line[:loc[0]]}
		tok.LineNum = inputLine.lineNum
		inputLine.line = inputLine.line[loc[0]:]
		receiver.state = inBlock
		return tok, inputLine
	}

	tok := PassthroughToken{Value: inputLine.line}
	tok.LineNum = inputLine.lineNum
	inputLine.line = ""
	return tok, inputLine
}

func (receiver *Lexer) processTokensInBlock(inputLine InputLine) ([]Token, InputLine) {
	toks := []Token{}
	currentLine, ok := receiver.getLineToProcess(inputLine)

	openRawStringExp := regexp.MustCompile("^`")

	for ok && receiver.state == inBlock {
		var (
			tok       Token
			remaining InputLine
		)

		// Ignore whitespace
		woWhitespace := whitespaceExp.ReplaceAllString(currentLine.line, "")
		currentLine.line = woWhitespace

		if openRawStringExp.MatchString(currentLine.line) {
			receiver.state = inStr
			// drop opening quote
			currentLine.line = currentLine.line[1:]
			tok, remaining = receiver.getRawStringToken(currentLine)
		} else {
			tok, remaining = receiver.getNextBlockToken(currentLine)
		}

		toks = append(toks, tok)
		currentLine, ok = receiver.getLineToProcess(remaining)
	}

	return toks, currentLine
}

func (receiver *Lexer) getRawStringToken(inputLine InputLine) (Token, InputLine) {
	rawStr := ""
	closeRawStringExp := regexp.MustCompile("`")

	currentLine, ok := receiver.getLineToProcess(inputLine)
	var tok Token

	for ok && receiver.state == inStr {
		if loc := closeRawStringExp.FindStringIndex(currentLine.line); loc != nil {
			rawStr += currentLine.line[:loc[0]]
			currentLine.line = currentLine.line[loc[1]:]
			tok = StrToken{Str: rawStr, TokenData: TokenData{LineNum: inputLine.lineNum}}
			receiver.state = inBlock
		} else {
			rawStr += currentLine.line
			currentLine, ok = receiver.getLineToProcess(InputLine{})
		}
	}

	return tok, currentLine
}

func (receiver *Lexer) getNextBlockToken(inputLine InputLine) (Token, InputLine) {
	tokData := TokenData{LineNum: inputLine.lineNum}

	if loc := multOp.FindStringIndex(inputLine.line); loc != nil {
		operator, remaining := extractToken(loc, inputLine)
		token := MultOpToken{Operator: operator, TokenData: tokData}
		return token, remaining
	} else if loc := addOp.FindStringIndex(inputLine.line); loc != nil {
		operator, remaining := extractToken(loc, inputLine)

		var token Token
		// if operator == "" && isActuallyUnary(tokens) {
		// 	token = UnaryOpToken{Operator: operator, TokenData: tokData}
		// } else {
		token = AddOpToken{Operator: operator, TokenData: tokData}
		// }

		return token, remaining
	} else if loc := relOp.FindStringIndex(inputLine.line); loc != nil {
		operator, remaining := extractToken(loc, inputLine)
		token := RelOpToken{Operator: operator, TokenData: tokData}
		return token, remaining
	} else if loc := logicOp.FindStringIndex(inputLine.line); loc != nil {
		operator, remaining := extractToken(loc, inputLine)
		token := LogicOpToken{Operator: operator, TokenData: tokData}
		return token, remaining
	} else if loc := assignOp.FindStringIndex(inputLine.line); loc != nil {
		operator, remaining := extractToken(loc, inputLine)
		token := AssignOpToken{Operator: operator, TokenData: tokData}
		return token, remaining
	} else if loc := unaryOp.FindStringIndex(inputLine.line); loc != nil {
		operator, remaining := extractToken(loc, inputLine)
		token := UnaryOpToken{Operator: operator, TokenData: tokData}
		return token, remaining
	} else if loc := strExp.FindStringIndex(inputLine.line); loc != nil {
		str, remaining := extractToken(loc, inputLine)
		token := StrToken{Str: str, TokenData: tokData}
		return token, remaining
	} else if loc := numExp.FindStringIndex(inputLine.line); loc != nil {
		num, remaining := extractToken(loc, inputLine)
		token := NumToken{Num: num, TokenData: tokData}
		return token, remaining
	} else if loc := boolExp.FindStringIndex(inputLine.line); loc != nil {
		boolVal, remaining := extractToken(loc, inputLine)
		token := BoolToken{Value: boolVal, TokenData: tokData}
		return token, remaining
	} else if loc := ifExp.FindStringIndex(inputLine.line); loc != nil {
		_, remaining := extractToken(loc, inputLine)
		token := IfToken{TokenData: tokData}
		return token, remaining
	} else if loc := elseIfExp.FindStringIndex(inputLine.line); loc != nil {
		_, remaining := extractToken(loc, inputLine)
		token := ElseIfToken{TokenData: tokData}
		return token, remaining
	} else if loc := elseExp.FindStringIndex(inputLine.line); loc != nil {
		_, remaining := extractToken(loc, inputLine)
		token := ElseToken{TokenData: tokData}
		return token, remaining
	} else if loc := forExp.FindStringIndex(inputLine.line); loc != nil {
		_, remaining := extractToken(loc, inputLine)
		token := ForToken{TokenData: tokData}
		return token, remaining
	} else if loc := inExp.FindStringIndex(inputLine.line); loc != nil {
		_, remaining := extractToken(loc, inputLine)
		token := InToken{TokenData: tokData}
		return token, remaining
	} else if loc := endExp.FindStringIndex(inputLine.line); loc != nil {
		_, remaining := extractToken(loc, inputLine)
		token := EndToken{TokenData: tokData}
		return token, remaining
	} else if loc := identExp.FindStringIndex(inputLine.line); loc != nil {
		// Ident should come after more specific tokens like bool and var
		ident, remaining := extractToken(loc, inputLine)
		token := IdentToken{Identifier: ident, TokenData: tokData}
		return token, remaining
	} else if loc := symbolExp.FindStringIndex(inputLine.line); loc != nil {
		symbol, remaining := extractToken(loc, inputLine)
		token := SymbolToken{Symbol: symbol, TokenData: tokData}
		return token, remaining
	} else if loc := blockExp.FindStringIndex(inputLine.line); loc != nil {
		block, remaining := extractToken(loc, inputLine)
		token := BlockToken{Block: block, TokenData: tokData}

		if block == "}}" {
			receiver.state = passthrough
		}
		return token, remaining
	}

	return nil, inputLine
}

// Extract the token between [loc[0],loc[1]) from the line
// and return the remaining characters in the line
func extractToken(loc []int, inputLine InputLine) (string, InputLine) {
	token := inputLine.line[loc[0]:loc[1]]
	inputLine.line = inputLine.line[loc[1]:]
	return token, inputLine
}
