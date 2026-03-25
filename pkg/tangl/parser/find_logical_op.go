//ff:func feature=tangl type=parser control=sequence
//ff:what findLogicalOp — find position of a logical operator in expression text
package parser

import "strings"

// findLogicalOp finds the position of a logical operator (and/or) in text.
func findLogicalOp(text string, op string) int {
	idx := strings.Index(text, op)
	if idx < 0 {
		return -1
	}
	before := text[:idx]
	if !strings.ContainsAny(before, " ") {
		return -1
	}
	return idx
}
