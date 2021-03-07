package lexer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var multOp = regexp.MustCompile(`^[*\/%]`)
var addOp = regexp.MustCompile(`^\+`)
var subOp = regexp.MustCompile(`^-`)
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

func (receiver *Lexer) Lex(inputReader io.Reader) (<-chan []Token, <-chan error) {
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
			errChan <- fmt.Errorf("lexer read error line %d: %s", i, inputBuffer.Err())
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
				if tok, remaining := receiver.processPassthroughTokens(inputLine); tok != nil {
					tokChan <- []Token{tok}
					inputLine = remaining
				}
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

func (receiver *Lexer) processPassthroughTokens(inputLine InputLine) (Token, InputLine) {
	openBlockExp := regexp.MustCompile(`{{:|{{`)
	if loc := openBlockExp.FindStringIndex(inputLine.line); loc != nil {
		receiver.state = inBlock

		if loc[0] == 0 {
			return nil, inputLine
		}

		tok := PassthroughToken{Value: inputLine.line[:loc[0]]}
		tok.LineNum = inputLine.lineNum
		inputLine.line = inputLine.line[loc[0]:]
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
			tok, remaining = receiver.getRawStringToken(currentLine)
		} else {
			tok, remaining = receiver.getNextBlockToken(currentLine)
		}

		if tok != nil {
			toks = append(toks, tok)
		}
		currentLine, ok = receiver.getLineToProcess(remaining)
	}

	if receiver.state == passthrough {
		toks = append(toks, EOLToken{TokenData: TokenData{LineNum: currentLine.lineNum}})
	}

	return toks, currentLine
}

func (receiver *Lexer) getRawStringToken(inputLine InputLine) (Token, InputLine) {
	rawStr := ""
	closeRawStringExp := regexp.MustCompile("`")
	currentLine, ok := receiver.getLineToProcess(inputLine)
	var tok Token

	for ok && receiver.state == inStr {
		if currentLine.line[0] == '`' {
			// drop opening quote
			currentLine.line = currentLine.line[1:]
		}

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
	if currentLine, ok := receiver.getLineToProcess(inputLine); ok {
		tokData := TokenData{LineNum: currentLine.lineNum}

		if loc := multOp.FindStringIndex(currentLine.line); loc != nil {
			operator, remaining := extractToken(loc, currentLine)
			token := MultOpToken{Operator: operator, TokenData: tokData}
			return token, remaining
		} else if loc := addOp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := AddOpToken{TokenData: tokData}
			return token, remaining
		} else if loc := subOp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := SubOpToken{TokenData: tokData}
			return token, remaining
		} else if loc := relOp.FindStringIndex(currentLine.line); loc != nil {
			operator, remaining := extractToken(loc, currentLine)
			token := RelOpToken{Operator: operator, TokenData: tokData}
			return token, remaining
		} else if loc := logicOp.FindStringIndex(currentLine.line); loc != nil {
			operator, remaining := extractToken(loc, currentLine)
			token := LogicOpToken{Operator: operator, TokenData: tokData}
			return token, remaining
		} else if loc := assignOp.FindStringIndex(currentLine.line); loc != nil {
			operator, remaining := extractToken(loc, currentLine)
			token := AssignOpToken{Operator: operator, TokenData: tokData}
			return token, remaining
		} else if loc := unaryOp.FindStringIndex(currentLine.line); loc != nil {
			operator, remaining := extractToken(loc, currentLine)
			token := NegationOpToken{Operator: operator, TokenData: tokData}
			return token, remaining
		} else if loc := strExp.FindStringIndex(currentLine.line); loc != nil {
			str, remaining := extractToken(loc, currentLine)
			token := StrToken{Str: str, TokenData: tokData}
			return token, remaining
		} else if loc := numExp.FindStringIndex(currentLine.line); loc != nil {
			num, remaining := extractToken(loc, currentLine)
			token := NumToken{Num: num, TokenData: tokData}
			return token, remaining
		} else if loc := boolExp.FindStringIndex(currentLine.line); loc != nil {
			boolVal, remaining := extractToken(loc, currentLine)
			token := BoolToken{Value: boolVal, TokenData: tokData}
			return token, remaining
		} else if loc := ifExp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := IfToken{TokenData: tokData}
			return token, remaining
		} else if loc := elseIfExp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := ElseIfToken{TokenData: tokData}
			return token, remaining
		} else if loc := elseExp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := ElseToken{TokenData: tokData}
			return token, remaining
		} else if loc := forExp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := ForToken{TokenData: tokData}
			return token, remaining
		} else if loc := inExp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := InToken{TokenData: tokData}
			return token, remaining
		} else if loc := endExp.FindStringIndex(currentLine.line); loc != nil {
			_, remaining := extractToken(loc, currentLine)
			token := EndToken{TokenData: tokData}
			return token, remaining
		} else if loc := identExp.FindStringIndex(currentLine.line); loc != nil {
			// Ident should come after more specific tokens like bool and var
			ident, remaining := extractToken(loc, currentLine)
			token := IdentToken{Identifier: ident, TokenData: tokData}
			return token, remaining
		} else if loc := symbolExp.FindStringIndex(currentLine.line); loc != nil {
			symbol, remaining := extractToken(loc, currentLine)
			token := SymbolToken{Symbol: symbol, TokenData: tokData}
			return token, remaining
		} else if loc := blockExp.FindStringIndex(currentLine.line); loc != nil {
			block, remaining := extractToken(loc, currentLine)
			token := BlockToken{Block: block, TokenData: tokData}

			if block == "}}" {
				receiver.state = passthrough
			}
			return token, remaining
		}
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
