package parser

import "fmt"

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
	first := receiver.children[0].(*IdentParseNode).Value

	if len(receiver.children) < 2 {
		return []string{first}
	}

	nestedVarName, ok := receiver.children[len(receiver.children)-1].(*VarNameParseNode)
	if ok {
		return append([]string{first}, nestedVarName.GetVarNameParts()...)
	}

	second := receiver.children[2].(*IdentParseNode).Value
	return []string{first, second}
	// flattened := receiver.GetFlattenedChildren()
	// nameParts := make([]string, 0, len(flattened))

	// for _, child := range flattened {
	// 	nameParts = append(nameParts, child.getIdentifierName())
	// }
	// return nameParts
}

// GetFlattenedChildren returns an array of nested VarNameParseNodes starting
// with the called node
func (receiver *VarNameParseNode) GetFlattenedChildren() []*VarNameParseNode {
	if len(receiver.flattenedChildren) == 0 {
		receiver.cacheFlattenedChildren()
	}

	return receiver.flattenedChildren
}

func (receiver *VarNameParseNode) cacheFlattenedChildren() {
	receiver.flattenedChildren = []*VarNameParseNode{receiver}
	current := receiver

	for current.hasNestedChildren() {
		next := current.children[len(current.children)-1].(*VarNameParseNode)
		receiver.flattenedChildren = append(receiver.flattenedChildren, next)
		current = next
	}
}

func (receiver *VarNameParseNode) hasNestedChildren() bool {
	if len(receiver.children) == 0 {
		return false
	}

	_, ok := receiver.children[len(receiver.children)-1].(*VarNameParseNode)
	return ok
}

func (receiver *VarNameParseNode) getIdentifierName() string {
	if len(receiver.children) > 0 {
		return receiver.children[0].(*IdentParseNode).Value
	}
	return ""
}
