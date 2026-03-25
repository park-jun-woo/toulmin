//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what findDotAfterDigits — find index of first dot preceded by digits in a string
package parser

// findDotAfterDigits returns the index of the first '.' that follows one or more digits, or -1.
func findDotAfterDigits(s string) int {
	for i, ch := range s {
		if ch == '.' && i > 0 {
			return i
		}
		if ch < '0' || ch > '9' {
			return -1
		}
	}
	return -1
}
