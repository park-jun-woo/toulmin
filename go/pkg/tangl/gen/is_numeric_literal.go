//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what isNumericLiteral — reports whether a Definitions literal is a plain number
package gen

// isNumericLiteral reports whether s is a plain decimal integer or float
// literal (optional leading sign), e.g. "650" or "3.14". Anything else —
// including a unit-suffixed literal like 65 degrees C — is not numeric
// and must be emitted as a Go string constant instead.
func isNumericLiteral(s string) bool {
	if s == "" {
		return false
	}
	i := 0
	if s[0] == '+' || s[0] == '-' {
		i++
	}
	if i == len(s) {
		return false
	}
	seenDigit := false
	seenDot := false
	for ; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			seenDigit = true
		case c == '.' && !seenDot:
			seenDot = true
		default:
			return false
		}
	}
	return seenDigit
}
