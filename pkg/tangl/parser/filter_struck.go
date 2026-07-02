//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what filterStruck — drop strikethrough (commented-out) items and their subtrees
package parser

// filterStruck removes struck-through items (and their entire subtree, since
// a commented-out item's children are comments too) at every nesting level.
func filterStruck(items []item) []item {
	var kept []item
	for _, it := range items {
		if it.Struck {
			continue
		}
		it.Children = filterStruck(it.Children)
		kept = append(kept, it)
	}
	return kept
}
