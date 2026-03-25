//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what skipCommaSpaces — advance index past consecutive comma and space characters
package parser

// skipCommaSpaces advances from start past consecutive commas and spaces and returns the new index.
func skipCommaSpaces(text string, start int) int {
	for start < len(text) && (text[start] == ',' || text[start] == ' ') {
		start++
	}
	return start
}
