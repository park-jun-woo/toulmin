//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what indentWidth — count the leading whitespace width of a line's text
package scan

// indentWidth counts the leading whitespace width of text, counting a space as
// 1 and a tab as 4, so nested list items can be compared by depth. It mirrors
// the parser's lineIndent.
func indentWidth(text string) int {
	width := 0
	for _, ch := range text {
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
