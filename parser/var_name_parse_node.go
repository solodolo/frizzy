package parser

import (
	"fmt"
	"math"
)

// VarNameParseNode represents a variable name
// Can be dot separated like "post.title" or a single
// identifier like "title"
type VarNameParseNode struct {
	ParseNode
	flattenedChildren []*VarNameParseNode
}

func (receiver *VarNameParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

func (node *VarNameParseNode) IsTerminal() bool {
	return false
}

// GetVarNameParts returns an array of string ident names represented
// by this VarNameParseNode tree
// e.g. "foo" will return ["foo"] and "foo.bar" will return ["foo", "bar"]
func (receiver *VarNameParseNode) GetVarNameParts() []string {
	numParts := float64(len(receiver.children)) / float64(2)
	numNames := int(math.Ceil(numParts))
	nameParts := make([]string, 0, numNames)

	for i := 0; i < len(receiver.children); i += 2 {
		nameParts = append(nameParts, receiver.children[i].(*IdentParseNode).Value)
	}

	return nameParts
}
