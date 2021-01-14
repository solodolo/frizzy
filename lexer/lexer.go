package lexer

import (
	"bufio"
	"os"
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
var relOp = regexp.MustCompile(`^>=|<=|!=|==|<|>`)
var logicOp = regexp.MustCompile(`^\|\||&&`)
var assignOp = regexp.MustCompile(`^=`)
var unaryOp = regexp.MustCompile(`^!|-`)

var ident = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_]*`)
var varName = regexp.MustCompile(`^([a-zA-Z]+[a-zA-Z0-9_]*)(\.[a-zA-Z][a-zA-Z0-9_]*)*`)
var strExp = regexp.MustCompile(`^"[^"]*"`)
var numExp = regexp.MustCompile(`^[0-9]+`)

var ifExp = regexp.MustCompile(`^if`)
var elseIfExp = regexp.MustCompile(`^else_if`)
var elseExp = regexp.MustCompile(`^else`)

var forExp = regexp.MustCompile(`^for`)
var inExp = regexp.MustCompile(`^in`)
var endExp = regexp.MustCompile(`^end`)

var boolExp = regexp.MustCompile(`^true|false`)

var symbolExp = regexp.MustCompile(`^[;()]`)

var blockExp = regexp.MustCompile(`^{{:|{{|}}`)

// Runs through file and creates a stream of tokens
// from the input
func LexFile(file *os.File, tokChan chan []Token) {
	// Read text line by line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Send lines to channel
	lineChan := make(chan string)
	go func() {
		for scanner.Scan() {
			lineChan <- scanner.Text()
		}
	}()

	// Read as lines come into the channel and
	// convert into a slice of Tokens
	for line := range lineChan {
		tokChan <- getLineTokens(line)
	}
}

// Check the front of the line for each of the tokens
// When found, erase found token from line and repeat until
// the line is empty
func getLineTokens(line string) []Token {
	tokens := make([]Token, 0)

	for len(line) > 0 {
		if loc := multOp.FindStringIndex(line); loc != nil {
			operator, remaining := extractToken(loc, line)
			line = remaining
			token := OpToken{Operator: operator}
			tokens = append(tokens, token)
		} else if loc := addOp.FindStringIndex(line); loc != nil {

		}
	}
	return tokens
}

// Extract the token between [loc[0],loc[1]) from the line
// and return the remaining characters in the line
func extractToken(loc []int, line string) (string, string) {
	token := line[loc[0]:loc[1]]
	return token, line[loc[1]:]
}
