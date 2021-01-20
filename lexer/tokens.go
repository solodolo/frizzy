package lexer

type Token interface {
	GetLineNum() int
	GetLineCol() int
	GetValue() string
}

type TokenData struct {
	LineNum int
	LineCol int
}

func (token TokenData) GetLineNum() int {
	return token.LineNum
}

func (token TokenData) GetLineCol() int {
	return token.LineCol
}

func (token TokenData) GetValue() string {
	return ""
}

// Represents an operator like
// =
type AssignOpToken struct {
	Operator string
	TokenData
}

func (token AssignOpToken) GetValue() string {
	return token.Operator
}

// Represents an operator like
// !
type UnaryOpToken struct {
	Operator string
	TokenData
}

func (token UnaryOpToken) GetValue() string {
	return token.Operator
}

// Represents an operator like
// ||, &&
type LogicOpToken struct {
	Operator string
	TokenData
}

func (token LogicOpToken) GetValue() string {
	return token.Operator
}

// Represents an operator like
// <, >, <= etc
type RelOpToken struct {
	Operator string
	TokenData
}

func (token RelOpToken) GetValue() string {
	return token.Operator
}

// Represents operators like
// +, -
type AddOpToken struct {
	Operator string
	TokenData
}

func (token AddOpToken) GetValue() string {
	return token.Operator
}

// Represents operators like
// *, /, %
type MultOpToken struct {
	Operator string
	TokenData
}

func (token MultOpToken) GetValue() string {
	return token.Operator
}

// Represents an identifier like
// print, get_cur_date, if, else_if, end, etc...
type IdentToken struct {
	Identifier string
	TokenData
}

func (token IdentToken) GetValue() string {
	return token.Identifier
}

// Represents a variable like
// post.title
type VarToken struct {
	Variable string
	TokenData
}

func (token VarToken) GetValue() string {
	return token.Variable
}

// Represents a string in double quotes like
// "foo", "0129", "bl()#$)(!@"
type StrToken struct {
	Str string
	TokenData
}

func (token StrToken) GetValue() string {
	return token.Str
}

// Represents an integer number
type NumToken struct {
	Num string
	TokenData
}

func (token NumToken) GetValue() string {
	return token.Num
}

// Represents a true/false value
type BoolToken struct {
	Value string
	TokenData
}

func (token BoolToken) GetValue() string {
	return token.Value
}

// A single grammar symbol like
// ';', '(', ')'
type SymbolToken struct {
	Symbol string
	TokenData
}

func (token SymbolToken) GetValue() string {
	return token.Symbol
}

// An open or closing block marker like
// {{, {{:, }}
type BlockToken struct {
	Block string
	TokenData
}

func (token BlockToken) GetValue() string {
	return token.Block
}

// Catchall for things that arent explicitly defined
type PassthroughToken struct {
	Value string
	TokenData
}

func (token PassthroughToken) GetValue() string {
	return token.Value
}
