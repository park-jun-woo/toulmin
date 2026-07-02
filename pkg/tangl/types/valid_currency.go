//ff:func feature=tangl type=util control=sequence
//ff:what validCurrency — reports whether v is a plain number
package types

// validCurrency reports whether v is a plain Go number (int or float kind).
func validCurrency(v any) bool {
	return isNumeric(v)
}
