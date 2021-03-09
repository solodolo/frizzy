package lexer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var (
	multOp               = regexp.MustCompile(`^[*\/%]`)
	addOp                = regexp.MustCompile(`^\+`)
	subOp                = regexp.MustCompile(`^-`)
	relOp                = regexp.MustCompile(`^(>=|<=|!=|==|<|>)`)
	logicOp              = regexp.MustCompile(`^(\|\||&&)`)
	assignOp             = regexp.MustCompile(`^=`)
	unaryOp              = regexp.MustCompile(`^!`)
	identExp             = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_]*`)
	strExp               = regexp.MustCompile(`(?m)^"[^"]*"`)
	numExp               = regexp.MustCompile(`^[0-9]+`)
	ifExp                = regexp.MustCompile(`^if`)
	elseIfExp            = regexp.MustCompile(`^else_if`)
	elseExp              = regexp.MustCompile(`^else`)
	forExp               = regexp.MustCompile(`^for`)
	inExp                = regexp.MustCompile(`^in`)
	endExp               = regexp.MustCompile(`^end`)
	boolExp              = regexp.MustCompile(`^(true|false)`)
	symbolExp            = regexp.MustCompile(`^[;(),\.]`)
	noWhitespaceBlockExp = regexp.MustCompile(`^-}`)
	blockExp             = regexp.MustCompile(`^({{:|{{|}})`)
	openBlockExp         = regexp.MustCompile(`{{:|{{`)
	openRawStringExp     = regexp.MustCompile("^`")
	closeRawStringExp    = regexp.MustCompile("`")
	whitespaceExp        = regexp.MustCompile(`^\s+`)
)

// Lexer states
type LexerState int

const (
	passthrough LexerState = iota
	passthroughNoWhitespace
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
	inputBuffer.Split(splitLinesKeepNL)

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

func splitLinesKeepNL(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	if err == nil && token != nil && !atEOF {
		// append a new line byte
		token = append(token, 10)
	}
	return
}

func readLines(inputBuffer *bufio.Scanner) (<-chan InputLine, <-chan error) {
	lineChan := make(chan InputLine)
	errChan := make(chan error, 1)

	go func() {
		defer close(lineChan)
		defer close(errChan)

		i := 1
		for inputBuffer.Scan() {
			line := inputBuffer.Text()
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
		if ok {
			return newLine, ok
		}
		return remaining, ok
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
			if receiver.state == passthrough || receiver.state == passthroughNoWhitespace {
				tok, remaining := receiver.processPassthroughTokens(inputLine)
				if tok != nil {
					tokChan <- []Token{tok}
				}
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

func (receiver *Lexer) processPassthroughTokens(inputLine InputLine) (Token, InputLine) {
	passthroughText := inputLine.line
	remainder := ""

	if receiver.state == passthroughNoWhitespace {
		passthroughText = whitespaceExp.ReplaceAllString(passthroughText, "")
	}

	if loc := openBlockExp.FindStringIndex(passthroughText); loc != nil {
		receiver.state = inBlock

		remainder = passthroughText[loc[0]:]
		passthroughText = passthroughText[:loc[0]]
	} else {
		receiver.state = passthrough
	}

	inputLine.line = remainder
	if passthroughText == "" {
		return nil, inputLine
	}

	tok := PassthroughToken{Value: passthroughText}
	tok.LineNum = inputLine.lineNum
	return tok, inputLine
}

func (receiver *Lexer) processTokensInBlock(inputLine InputLine) ([]Token, InputLine) {
	toks := []Token{}
	currentLine, ok := receiver.getLineToProcess(inputLine)

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

	// if the state changed from inBlock to some passthrough state
	if receiver.state == passthrough || receiver.state == passthroughNoWhitespace {
		toks = append(toks, EOLToken{TokenData: TokenData{LineNum: currentLine.lineNum}})
	}

	return toks, currentLine
}

func (receiver *Lexer) getRawStringToken(inputLine InputLine) (Token, InputLine) {
	rawStr := ""
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

		if loc := noWhitespaceBlockExp.FindStringIndex(currentLine.line); loc != nil { // should come before subOp
			receiver.state = passthroughNoWhitespace
			block, remaining := extractToken(loc, currentLine)
			token := BlockToken{Block: block}

			return token, remaining
		} else if loc := multOp.FindStringIndex(currentLine.line); loc != nil {
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

			// drop open and close quotes
			str = str[1 : len(str)-1]
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
