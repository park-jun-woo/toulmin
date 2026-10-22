//ff:func feature=tangl type=util control=sequence
//ff:what validQuantity — reports whether v is a plain number
package types

// validQuantity reports whether v is a plain Go number (int or float kind).
func validQuantity(v any) bool {
	return isNumeric(v)
}
