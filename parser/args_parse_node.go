package parser

import "fmt"

type ArgsParseNode struct {
	ParseNode
}

func (receiver *ArgsParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}
