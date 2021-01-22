package parser

// program -> blocks

// blocks -> block | print_block
// block -> {{ statement | statement_list }}
// print_block -> {{: statement }}

// func_call -> ID ( args )
// args -> arg_list | ε
// arg_list -> arg_list, expression | expression

// expression -> VAR_NAME = expression | unary_expression
// unary_expression -> UNARY_OPERATORS unary_expression | logic_expression
// logic_expression -> logic_expression LOGIC_OPERATORS rel_expression | rel_expression
// rel_expression -> rel_expression REL_OPERATORS add_expression | add_expression
// add_expression -> add_expression ADD_OPERATORS mult_expression | mult_expression
// mult_expression -> mult_expression MULT_OPERATORS term_expression | term_expression
// term_expression -> VAR_NAME | STRING | NUM | ( expression )

// if_statement -> if( expression ) statement_list else_if_statement else_statement end
// else_if_statement -> else_if ( expression ) statement_list | ε
// else_statement -> else statement_list | ε
// for_statement -> for( ID IN (STRING | VAR_NAME | func_call) ) statement_list end

// statement -> expression
// statement -> func_call

// statement_list -> statement_list statement; | ε
type ParentNode interface {
	GetChildren() []*ParseNode
}

type ParseNode struct {
	children []*ParseNode
}

func (node *ParseNode) GetChildren() []*ParseNode {
	return node.children
}

type StringParseNode struct {
	Value string
	ParseNode
}

type NumParseNode struct {
	Value int
	ParseNode
}

type VarParseNode struct {
	VarName string
	ParseNode
}
