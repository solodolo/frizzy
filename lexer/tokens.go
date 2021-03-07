package lexer

type Token interface {
	GetLineNum() int
	GetLineCol() int
	GetValue() string
	GetGrammarSymbol() string
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

// The actual value the token represents
// 5, "foo", post.title, *, etc
func (token TokenData) GetValue() string {
	return ""
}

// The symbol as it appears in the parse table
// =, +, MULT_OP, ID, NUM, etc
func (token TokenData) GetGrammarSymbol() string {
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

func (token AssignOpToken) GetGrammarSymbol() string {
	return "="
}

// NegationOpToken represents the ! operator
type NegationOpToken struct {
	Operator string
	TokenData
}

func (token NegationOpToken) GetValue() string {
	return "!"
}

func (token NegationOpToken) GetGrammarSymbol() string {
	return token.GetValue()
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

func (token LogicOpToken) GetGrammarSymbol() string {
	return "LOGIC_OP"
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

func (token RelOpToken) GetGrammarSymbol() string {
	return "REL_OP"
}

// AddOpToken represents the addition operator
type AddOpToken struct {
	TokenData
}

func (token AddOpToken) GetValue() string {
	return "+"
}

func (token AddOpToken) GetGrammarSymbol() string {
	return token.GetValue()
}

// SubOpToken represents the subtraction operator
type SubOpToken struct {
	TokenData
}

func (token SubOpToken) GetValue() string {
	return "-"
}

func (token SubOpToken) GetGrammarSymbol() string {
	return token.GetValue()
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

func (token MultOpToken) GetGrammarSymbol() string {
	return "MULT_OP"
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

func (token IdentToken) GetGrammarSymbol() string {
	return "ID"
}

type ForToken struct {
	TokenData
}

func (token ForToken) GetValue() string {
	return "for"
}

func (token ForToken) GetGrammarSymbol() string {
	return "FOR"
}

type InToken struct {
	TokenData
}

func (token InToken) GetValue() string {
	return "in"
}

func (token InToken) GetGrammarSymbol() string {
	return "IN"
}

type IfToken struct {
	TokenData
}

func (token IfToken) GetValue() string {
	return "if"
}

func (token IfToken) GetGrammarSymbol() string {
	return "IF"
}

type ElseIfToken struct {
	TokenData
}

func (token ElseIfToken) GetValue() string {
	return "else_if"
}

func (token ElseIfToken) GetGrammarSymbol() string {
	return "ELSE_IF"
}

type ElseToken struct {
	TokenData
}

func (token ElseToken) GetValue() string {
	return "else"
}

func (token ElseToken) GetGrammarSymbol() string {
	return "ELSE"
}

type EndToken struct {
	TokenData
}

func (token EndToken) GetValue() string {
	return "end"
}

func (token EndToken) GetGrammarSymbol() string {
	return "END"
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

func (token StrToken) GetGrammarSymbol() string {
	return "STRING"
}

// Represents an integer number
type NumToken struct {
	Num string
	TokenData
}

func (token NumToken) GetValue() string {
	return token.Num
}

func (token NumToken) GetGrammarSymbol() string {
	return "NUM"
}

// Represents a true/false value
type BoolToken struct {
	Value string
	TokenData
}

func (token BoolToken) GetValue() string {
	return token.Value
}

func (token BoolToken) GetGrammarSymbol() string {
	return "BOOL"
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

func (token SymbolToken) GetGrammarSymbol() string {
	return token.GetValue()
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

func (token BlockToken) GetGrammarSymbol() string {
	return token.GetValue()
}

// Represents the EOL
type EOLToken struct {
	TokenData
}

func (token EOLToken) GetValue() string {
	return ""
}

func (token EOLToken) GetGrammarSymbol() string {
	return "$"
}

// Catchall for things that arent explicitly defined
type PassthroughToken struct {
	Value string
	TokenData
}

func (token PassthroughToken) GetValue() string {
	return token.Value
}

func (token PassthroughToken) GetGrammarSymbol() string {
	return ""
}
