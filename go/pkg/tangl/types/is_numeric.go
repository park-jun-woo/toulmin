//ff:func feature=tangl type=util control=selection
//ff:what isNumeric — reports whether v is any Go integer or floating-point kind
package types

// isNumeric reports whether v is any Go signed/unsigned integer or float kind.
// Shared by validCurrency and validQuantity, which both accept plain numbers.
func isNumeric(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}
