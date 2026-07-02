//ff:type feature=tangl type=parser
//ff:what lineEntry — a non-blank source line with its indent width and number
package parser

// lineEntry is a non-blank source line reduced to its indent width, trimmed
// text, and 1-based line number, ready for recursive list-tree parsing.
type lineEntry struct {
	Indent int
	Text   string
	Line   int
}
