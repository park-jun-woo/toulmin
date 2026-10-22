//ff:func feature=tangl type=parser control=sequence
//ff:what leafItem — build a childless Item from one list-item line at the given indent
package scan

import "strings"

// leafItem builds a childless Item from one list-item line: it strips the list
// marker, unwraps any `~~...~~` strikethrough, and keeps the whole-line span.
// The caller (buildItemTree) attaches Children.
func leafItem(ln Line, indent int) Item {
	content, ordered, number, _ := parseMarker(strings.TrimSpace(ln.Text))
	text, struck := strike(content)
	return Item{
		Span:    ln.Span,
		Text:    text,
		No:      ln.No,
		Indent:  indent,
		Ordered: ordered,
		Index:   number,
		Struck:  struck,
	}
}
