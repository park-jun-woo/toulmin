//ff:func feature=tangl type=util control=selection
//ff:what validInteger — reports whether v is a Go signed or unsigned integer
package types

// validInteger reports whether v is any Go signed or unsigned integer kind.
func validInteger(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}
