package analyser

// A criteria to match among nodes that are being traversed and searched.
type MatchCriteria interface {
	Match(*DocumentNode) bool
}

// Match this node against a set of criteria.
func (node *DocumentNode) Match(criteria ...MatchCriteria) bool {
	for _, c := range criteria {
		if !c.Match(node) {
			return false
		}
	}
	return true
}

// Match each leaf against set of criteria, return all matched leaves.
func (node *DocumentNode) SearchLeaves(criteria ...MatchCriteria) (matches []*DocumentNode) {
	matches = make([]*DocumentNode, 0, 0)
	if node.Leaves == nil {
		return
	}
	for _, leaf := range node.Leaves {
		if leaf.Match(criteria...) {
			matches = append(matches, leaf)
		}
	}
	return
}

// Match each leaf against set of criteria, recursively to leaves of the leaf, return all matched leaves.
func (node *DocumentNode) SearchLeavesRecursively(criteria ...MatchCriteria) (matches []*DocumentNode) {
	matches = make([]*DocumentNode, 0, 0)
	if node.Leaves == nil {
		return
	}
	for _, leaf := range node.Leaves {
		if leaf.Match(criteria...) {
			matches = append(matches, leaf)
		} else {
			matches = append(matches, leaf.SearchLeavesRecursively(criteria...)...)
		}
	}
	return
}
