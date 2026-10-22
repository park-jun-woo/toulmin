//ff:func feature=tangl type=util control=sequence
//ff:what validBoolean — reports whether v is a Go bool
package types

// validBoolean reports whether v is a Go bool.
func validBoolean(v any) bool {
	_, ok := v.(bool)
	return ok
}
