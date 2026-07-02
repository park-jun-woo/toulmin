//ff:func feature=tangl type=util control=sequence
//ff:what validEmail — reports whether v is a string in a plausible email format
package types

import "regexp"

// emailPattern is a deliberately simple "does this look like an email" check,
// not a full RFC 5322 validator.
var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// validEmail reports whether v is a string matching a simple email format.
func validEmail(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	return emailPattern.MatchString(s)
}
