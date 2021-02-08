package parser

import (
	"fmt"
)

type FuncCallParseNode struct {
	ParseNode
}

func (receiver FuncCallParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}
