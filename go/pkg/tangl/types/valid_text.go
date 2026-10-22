//ff:func feature=tangl type=util control=sequence
//ff:what validText — reports whether v is a Go string
package types

// validText reports whether v is a Go string.
func validText(v any) bool {
	_, ok := v.(string)
	return ok
}
