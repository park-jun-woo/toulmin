//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what trimStart — trim leading spaces and tabs from a string
package parser

// trimStart trims leading spaces and tabs (but not other whitespace) so
// list-item content is normalized without touching non-ASCII content.
func trimStart(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return s[i:]
}
