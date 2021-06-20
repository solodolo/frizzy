package parser

func extractToken(head TreeNode, decendentIndices []int) TreeNode {
	current := head
	for _, decendentIndex := range decendentIndices {
		children := current.GetChildren()
		if len(children) > decendentIndex {
			current = children[decendentIndex]
		}
	}
	return current
}

func nodeSlicesEqual(arr1, arr2 []TreeNode) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i].String() != arr2[i].String() {
			return false
		}
	}

	return true
}
