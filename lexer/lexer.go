package lexer

import (
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
