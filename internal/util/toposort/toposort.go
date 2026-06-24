package toposort

// Order returns ids arranged so every parent precedes its children, given the
// child adjacency (parent id -> child ids). Input order is preserved among roots
// and siblings. Edges to ids outside the input set are ignored, and ids reachable
// only through a cycle are appended in input order so none are dropped.
func Order(ids []string, children map[string][]string) []string {
	inSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		inSet[id] = true
	}

	isChild := map[string]bool{}
	for parent, kids := range children {
		if !inSet[parent] {
			continue
		}
		for _, c := range kids {
			if inSet[c] {
				isChild[c] = true
			}
		}
	}

	ordered := make([]string, 0, len(ids))
	visited := make(map[string]bool, len(ids))

	var visit func(id string)
	visit = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true
		ordered = append(ordered, id)
		for _, child := range children[id] {
			if inSet[child] {
				visit(child)
			}
		}
	}

	// Roots (no parent within the set) first, so every parent precedes its children.
	for _, id := range ids {
		if !isChild[id] {
			visit(id)
		}
	}

	// Safety net for cycles: emit anything still unreached, in input order.
	for _, id := range ids {
		visit(id)
	}

	return ordered
}
