//ff:func feature=tangl type=util control=selection
//ff:what validNumber — reports whether v is a Go floating-point number
package types

// validNumber reports whether v is a Go float32 or float64.
func validNumber(v any) bool {
	switch v.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}
