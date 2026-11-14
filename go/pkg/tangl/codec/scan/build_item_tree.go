//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what buildItemTree — recursively build a nested Item tree from list-item lines
package scan

// buildItemTree consumes item lines at their first (minimal) indent as siblings,
// recursively parsing each item's deeper-indented follower lines as Children. It
// mirrors the parser's buildItems, but every input line is already a list item
// (blank and prose lines were removed by ScanItems) and each Item keeps its
// whole-line span for gapless tiling.
func buildItemTree(lines []Line) []Item {
	if len(lines) == 0 {
		return nil
	}
	baseIndent := indentWidth(lines[0].Text)
	var items []Item
	i := 0
	for i < len(lines) {
		if indentWidth(lines[i].Text) != baseIndent {
			i++
			continue
		}
		it := leafItem(lines[i], baseIndent)
		i++
		childStart := i
		for i < len(lines) && indentWidth(lines[i].Text) > baseIndent {
			i++
		}
		it.Children = buildItemTree(lines[childStart:i])
		items = append(items, it)
	}
	return items
}
