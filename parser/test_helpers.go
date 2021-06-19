package parser

func extractToken(head TreeNode, decendentIndices []int) TreeNode {
	children := head.GetChildren()

	if len(children) > 0 {
		blockNode := children[0].(*BlockParseNode)

		var nextNode TreeNode = blockNode
		for _, decendentIndex := range decendentIndices {
			grandChildren := nextNode.GetChildren()

			if len(grandChildren) > decendentIndex {
				nextNode = grandChildren[decendentIndex]
			}
		}

		return nextNode
	}

	return nil
}
