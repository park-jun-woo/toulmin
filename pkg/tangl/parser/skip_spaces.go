//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what skipSpaces — advance index past consecutive space characters
package parser

// skipSpaces advances from start past consecutive spaces and returns the new index.
func skipSpaces(text string, start int) int {
	for start < len(text) && text[start] == ' ' {
		start++
	}
	return start
}
