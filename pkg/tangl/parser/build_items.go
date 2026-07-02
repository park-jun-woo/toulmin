//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what buildItems — recursively build a nested item tree from line entries
package parser

// buildItems consumes entries at their minimal (first) indent level as
// siblings, recursively parsing each item's deeper-indented follower lines
// as its Children. Non-list lines are ignored (spec: prose is not code).
func buildItems(entries []lineEntry) []item {
	if len(entries) == 0 {
		return nil
	}
	baseIndent := entries[0].Indent
	var items []item
	i := 0
	for i < len(entries) {
		e := entries[i]
		if e.Indent != baseIndent {
			i++
			continue
		}
		content, ordered, number, ok := parseItemLine(e.Text)
		i++
		if !ok {
			continue
		}
		text, struck := stripStrike(content)
		it := item{Text: text, Ordered: ordered, Number: number, Line: e.Line, Struck: struck}
		childStart := i
		for i < len(entries) && entries[i].Indent > baseIndent {
			i++
		}
		it.Children = buildItems(entries[childStart:i])
		items = append(items, it)
	}
	return items
}
