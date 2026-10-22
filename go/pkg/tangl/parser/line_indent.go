//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what lineIndent — count leading whitespace width of a line
package parser

// lineIndent counts the leading whitespace width of line, counting a
// space as 1 and a tab as 4, so nested list items can be compared by depth.
func lineIndent(line string) int {
	width := 0
	for _, ch := range line {
		switch ch {
		case ' ':
			width++
		case '\t':
			width += 4
		default:
			return width
		}
	}
	return width
}
