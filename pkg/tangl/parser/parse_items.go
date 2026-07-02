//ff:func feature=tangl type=parser control=sequence
//ff:what parseItems — build, filter, and validate a section's item tree
package parser

// parseItems turns a section's raw content lines into a validated item
// tree: build the nested list structure, drop strikethrough comment items,
// then verify ordered-list numbering.
func parseItems(lines []string, lineOffset int, path string) ([]item, error) {
	entries := entriesFromLines(lines, lineOffset)
	items := buildItems(entries)
	items = filterStruck(items)
	if err := checkOrderedSiblings(items, path); err != nil {
		return nil, err
	}
	return items, nil
}
