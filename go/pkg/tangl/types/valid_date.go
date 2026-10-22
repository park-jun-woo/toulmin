//ff:func feature=tangl type=util control=sequence
//ff:what validDate — reports whether v is a string in YYYY-MM-DD format
package types

import "time"

// validDate reports whether v is a string parseable as "2006-01-02".
func validDate(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}
