package parser

import "fmt"

// ContentParseNode represents a tree structure containing
// other ContentParseNodes, passthrough values, or blocks
type ContentParseNode struct {
	ParseNode
	flattenedChildren []*ContentParseNode
}

func (receiver *ContentParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}
