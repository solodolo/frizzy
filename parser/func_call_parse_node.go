package parser

import (
	"fmt"
)

// FuncCallParseNode represents a function call in our grammar
type FuncCallParseNode struct {
	ParseNode
}

func (receiver FuncCallParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}
