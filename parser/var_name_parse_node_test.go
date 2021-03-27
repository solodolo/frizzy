package parser

import (
	"testing"
)

func TestUnestedNodeReturnsCorrectNumberAndValueVarNameParts(t *testing.T) {
	ident := &IdentParseNode{Value: "foo"}
	node := VarNameParseNode{
		ParseNode: ParseNode{
			children: []TreeNode{ident},
		},
	}

	nameParts := node.GetVarNameParts()

	if len(nameParts) != 1 {
		t.Errorf("expected %d name part, got %d", 1, len(nameParts))
	} else if nameParts[0] != ident.Value {
		t.Errorf("expected name part to equal %q, got %q", ident.Value, nameParts[0])
	}
}

func TestNestedNodeReturnsCorrectNumberAndValueVarNameParts(t *testing.T) {
	ident1 := &IdentParseNode{Value: "foo"}
	ident2 := &IdentParseNode{Value: "title"}
	ident3 := &IdentParseNode{Value: "date"}

	node := VarNameParseNode{
		ParseNode: ParseNode{
			children: []TreeNode{
				ident1,
				&SymbolParseNode{Value: "."},
				&VarNameParseNode{
					ParseNode: ParseNode{
						children: []TreeNode{
							ident2,
							&SymbolParseNode{Value: "."},
							&VarNameParseNode{
								ParseNode: ParseNode{
									children: []TreeNode{
										ident3,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	nameParts := node.GetVarNameParts()

	if len(nameParts) != 3 {
		t.Errorf("expected %d name parts, got %d", 3, len(nameParts))
	} else if nameParts[0] != ident1.Value {
		t.Errorf("expected first name part to equal %q, got %q", ident1.Value, nameParts[0])
	} else if nameParts[1] != ident2.Value {
		t.Errorf("expected second name part to equal %q, got %q", ident2.Value, nameParts[1])
	} else if nameParts[2] != ident3.Value {
		t.Errorf("expected third name part to equal %q, got %q", ident3.Value, nameParts[2])
	}
}
