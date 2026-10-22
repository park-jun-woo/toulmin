//ff:func feature=tangl type=util control=sequence
//ff:what validTime — reports whether v is a string in HH:MM:SS or HH:MM format
package types

import "time"

// validTime reports whether v is a string parseable as "15:04:05" or "15:04".
func validTime(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	if _, err := time.Parse("15:04:05", s); err == nil {
		return true
	}
	_, err := time.Parse("15:04", s)
	return err == nil
}
