package lexer

type Token struct {
	LineNum int
	LineCol int
}

// Represents an operator like
// +,-,*,%,==,!=,<, etc
type OpToken struct {
	Operator string
	Token
}

// Represents an identifier like
// print, get_cur_date
type IdentToken struct {
	Identifier string
	Token
}

// Represents a variable like
// post.title
type VarToken struct {
	Variable string
	Token
}

// Represents a string in double quotes like
// "foo", "0129", "bl()#$)(!@"
type StrToken struct {
	Str string
	Token
}

// Represents an integer number
type NumToken struct {
	Num int
	Token
}

type IfToken Token
type ElseIfToken Token
type ElseToken Token

type ForToken Token
type InToken Token
type EndToken Token

// Represents a true/false value
type BoolToken struct {
	Value bool
	Token
}

// A single grammar symbol like
// ';', '(', ')'
type SymbolExp struct {
	Symbol rune
	Token
}

// An open or closing block marker like
// {{, {{:, }}
type BlockExp struct {
	Block string
	Token
}
