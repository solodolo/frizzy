package parser

import (
	"fmt"
)

type ForLoopParseNode struct {
	ParseNode
}

func (node ForLoopParseNode) String() string {
	return fmt.Sprintf("%T", node)
}

func (receiver ForLoopParseNode) GetNumLoops() int {
	return 0
}
