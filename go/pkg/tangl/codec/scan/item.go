//ff:type feature=tangl type=parser
//ff:what Item — one markdown list item with its whole-line span and nested children
package scan

// Item is one markdown list entry (`- ` or `N. `) within a section. Span covers
// the whole line, including its marker, indentation, and trailing newline, so
// items tile without gaps. Text is the marker-stripped, strike-unwrapped
// content; No is the 1-based line number; Indent is the leading whitespace
// width (space 1, tab 4); Ordered/Index carry an `N. ` marker's number; Struck
// marks a `~~...~~` comment item; Children are deeper-indented items.
type Item struct {
	Span     Span
	Text     string
	No       int
	Indent   int
	Ordered  bool
	Index    int
	Struck   bool
	Children []Item
}
