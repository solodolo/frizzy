package lexer

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
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
var varExp = regexp.MustCompile(`^([a-zA-Z]+[a-zA-Z0-9_]*)\.[a-zA-Z][a-zA-Z0-9_]*`)
var strExp = regexp.MustCompile(`^"[^"]*"`)
var numExp = regexp.MustCompile(`^[0-9]+`)

var ifExp = regexp.MustCompile(`^if`)
var elseIfExp = regexp.MustCompile(`^else_if`)
var elseExp = regexp.MustCompile(`^else`)

var forExp = regexp.MustCompile(`^for`)
var inExp = regexp.MustCompile(`^in`)
var endExp = regexp.MustCompile(`^end`)

var boolExp = regexp.MustCompile(`^(true|false)`)

var symbolExp = regexp.MustCompile(`^[;(),]`)

var blockExp = regexp.MustCompile(`^({{:|{{|}})`)

var whitespaceExp = regexp.MustCompile(`^\s+`)

type InputLine struct {
	line    string
	lineNum int
}

// Runs through file and creates a stream of tokens
// from the input
func Lex(scanner *bufio.Scanner, tokChan chan<- []Token, errChan chan<- error) {
	defer close(errChan)
	// Read line by line
	scanner.Split(bufio.ScanLines)
	lineChan := make(chan InputLine)
	lineErrChan := make(chan error)

	go getLines(scanner, lineChan, lineErrChan)
	go tokenizeLines(lineChan, tokChan)

	err := <-lineErrChan
	errChan <- err
}

// Reads lines from lineChan and converts each one into a slice of tokens.
// These slices are sent to tokChan
func tokenizeLines(lineChan <-chan InputLine, tokChan chan<- []Token) {
	defer close(tokChan)
	openBlock := false

	for line := range lineChan {
		tokens, stillOpen := processLine(line, openBlock)
		openBlock = stillOpen
		tokChan <- tokens
	}
}

// Reads lines from scanner and sends them to the string channel.
// Any errors will be sent to error channel or nil once everything
// has been read.
func getLines(scanner *bufio.Scanner, lineChan chan<- InputLine, errChan chan<- error) {
	defer close(lineChan)
	defer close(errChan)

	lineNum := 1
	for scanner.Scan() {
		lineChan <- InputLine{
			line:    scanner.Text(),
			lineNum: lineNum,
		}

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		errChan <- fmt.Errorf("error reading lines for lexing: %s", err.Error())
	} else {
		errChan <- nil
	}
}

// All text that comes before a {{: or {{ is passed through
// and all text that comes after the corresponding }} is
// also passed through
func processLine(line InputLine, openBlock bool) ([]Token, bool) {
	// Check for indices of open and close blocks
	start, end := getBlockIndices(line.line)
	// Get text that should be passed through vs fully tokenized
	preBlock, inBlock, postBlock := partitionLine(line.line, start, end, openBlock)

	lineTokens := []Token{}

	// If there is text before block, pass it through
	if len(preBlock) > 0 {
		tokData := TokenData{LineNum: line.lineNum, LineCol: 0}
		token := PassthroughToken{Value: preBlock, TokenData: tokData}
		lineTokens = append(lineTokens, token)
	}

	// Tokenize the text in the blocks including any open/close blocks
	lineTokens = append(lineTokens, getLineTokens(inBlock, line.lineNum)...)

	// If there is text after block, pass it through
	if len(postBlock) > 0 {
		tokData := TokenData{LineNum: line.lineNum, LineCol: len(preBlock) + len(inBlock)}
		token := PassthroughToken{Value: postBlock, TokenData: tokData}
		lineTokens = append(lineTokens, token)
	}

	return lineTokens, isBlockStillOpen(start, end, openBlock)
}

/* Returns the indices for the open and close blocks in line
* Examples
* {{abcde}} => 0,9
* abc{{de}} => 3,9
* ab{{cd}}e => 2,8
* {{abcde => 0,7
* abcde}} => 0,7
* a{{bcde => 1,7
* abcd}}e => 0,6 */
func getBlockIndices(line string) (int, int) {
	openingBlocks := []string{"{{:", "{{"}
	openingIndex, closingIndex := -1, -1

	for _, block := range openingBlocks {
		if i := strings.Index(line, block); i != -1 {
			openingIndex = i
			break
		}
	}

	if i := strings.Index(line, "}}"); i != -1 {
		closingIndex = i + 2 // + 2 to include '}}'
	}

	return openingIndex, closingIndex
}

// Partitions the line into text before opening block text in block and
// text after closing block.
// Considers if there was an open block in some previous line.
func partitionLine(line string, start, end int, openBlock bool) (string, string, string) {
	var preBlock, inBlock, postBlock string
	if openBlock {
		if end == -1 {
			inBlock = line
		} else {
			inBlock = line[:end]
			postBlock = line[end:]
		}
	} else {
		if start == -1 {
			preBlock = line
		} else if start != -1 && end != -1 {
			preBlock = line[:start]
			inBlock = line[start:end]
			postBlock = line[end:]
		} else if start != -1 {
			preBlock = line[:start]
			inBlock = line[start:]
		}
	}

	return preBlock, inBlock, postBlock
}

// Determines the new open block status based on current
// status and the presence of open and/or close blocks in
// the current line
func isBlockStillOpen(start, end int, openBlock bool) bool {
	if openBlock && end != -1 {
		openBlock = false
	} else if !openBlock && start != -1 && end == -1 {
		openBlock = true
	}

	return openBlock
}

// Check the front of the line for each of the tokens
// When found, erase found token from line and repeat until
// the line is empty.
// Assumes the text in line is within a block.
func getLineTokens(line string, lineNum int) []Token {
	tokens := make([]Token, 0)

	col := 1
	// Check each regex against the line
	for len(line) > 0 {
		// Ignore whitespace
		woWhitespace := whitespaceExp.ReplaceAllString(line, "")
		col += len(line) - len(woWhitespace)
		line = woWhitespace

		tokData := TokenData{LineNum: lineNum, LineCol: col}
		newLine := ""
		if loc := multOp.FindStringIndex(line); loc != nil {
			operator, remaining := extractToken(loc, line)
			newLine = remaining
			token := MultOpToken{Operator: operator, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := addOp.FindStringIndex(line); loc != nil {
			operator, remaining := extractToken(loc, line)
			newLine = remaining

			lastTok := tokens[len(tokens)-1]
			var token Token
			if _, ok := lastTok.(NumToken); !ok && operator == "-" {
				token = UnaryOpToken{Operator: operator, TokenData: tokData}
			} else {
				token = AddOpToken{Operator: operator, TokenData: tokData}
			}
			tokens = append(tokens, token)
		} else if loc := relOp.FindStringIndex(line); loc != nil {
			operator, remaining := extractToken(loc, line)
			newLine = remaining
			token := RelOpToken{Operator: operator, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := logicOp.FindStringIndex(line); loc != nil {
			operator, remaining := extractToken(loc, line)
			newLine = remaining
			token := LogicOpToken{Operator: operator, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := assignOp.FindStringIndex(line); loc != nil {
			operator, remaining := extractToken(loc, line)
			newLine = remaining
			token := AssignOpToken{Operator: operator, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := unaryOp.FindStringIndex(line); loc != nil {
			operator, remaining := extractToken(loc, line)
			newLine = remaining
			token := UnaryOpToken{Operator: operator, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := strExp.FindStringIndex(line); loc != nil {
			str, remaining := extractToken(loc, line)
			newLine = remaining
			token := StrToken{Str: str, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := numExp.FindStringIndex(line); loc != nil {
			num, remaining := extractToken(loc, line)
			newLine = remaining
			token := NumToken{Num: num, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := boolExp.FindStringIndex(line); loc != nil {
			boolVal, remaining := extractToken(loc, line)
			newLine = remaining
			token := BoolToken{Value: boolVal, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := ifExp.FindStringIndex(line); loc != nil {
			_, remaining := extractToken(loc, line)
			newLine = remaining
			token := IfToken{TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := elseIfExp.FindStringIndex(line); loc != nil {
			_, remaining := extractToken(loc, line)
			newLine = remaining
			token := ElseIfToken{TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := elseExp.FindStringIndex(line); loc != nil {
			_, remaining := extractToken(loc, line)
			newLine = remaining
			token := ElseToken{TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := forExp.FindStringIndex(line); loc != nil {
			_, remaining := extractToken(loc, line)
			newLine = remaining
			token := ForToken{TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := inExp.FindStringIndex(line); loc != nil {
			_, remaining := extractToken(loc, line)
			newLine = remaining
			token := InToken{TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := endExp.FindStringIndex(line); loc != nil {
			_, remaining := extractToken(loc, line)
			newLine = remaining
			token := EndToken{TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := varExp.FindStringIndex(line); loc != nil {
			variable, remaining := extractToken(loc, line)
			newLine = remaining
			token := VarToken{Variable: variable, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := identExp.FindStringIndex(line); loc != nil {
			// Ident should come after more specific tokens like bool and var
			ident, remaining := extractToken(loc, line)
			newLine = remaining
			token := IdentToken{Identifier: ident, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := symbolExp.FindStringIndex(line); loc != nil {
			symbol, remaining := extractToken(loc, line)
			newLine = remaining
			token := SymbolToken{Symbol: symbol, TokenData: tokData}
			tokens = append(tokens, token)
		} else if loc := blockExp.FindStringIndex(line); loc != nil {
			block, remaining := extractToken(loc, line)
			newLine = remaining
			token := BlockToken{Block: block, TokenData: tokData}
			tokens = append(tokens, token)
		} else if len(line) > 0 {
			// No match so just pass through the char at the front
			// of the line
			val, remaining := extractToken([]int{0, 1}, line)
			newLine = remaining
			token := PassthroughToken{Value: val, TokenData: tokData}
			tokens = append(tokens, token)
		}

		col += len(line) - len(newLine)
		line = newLine
	}
	return tokens
}

// Extract the token between [loc[0],loc[1]) from the line
// and return the remaining characters in the line
func extractToken(loc []int, line string) (string, string) {
	token := line[loc[0]:loc[1]]
	return token, line[loc[1]:]
}
