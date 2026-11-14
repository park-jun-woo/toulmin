//ff:type feature=tangl type=parser
//ff:what item — one parsed markdown list item within a tangl: section
package parser

// item is one markdown list entry (`- ` or `N. `) within a tangl: section,
// with its nested children (by indentation) already attached.
type item struct {
	Text     string
	Ordered  bool
	Number   int
	Line     int
	Struck   bool
	Children []item
}
