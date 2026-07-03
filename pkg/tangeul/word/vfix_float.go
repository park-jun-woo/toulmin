//ff:func feature=tangl type=util control=sequence
//ff:what Float — convert a VFix back to its float64 value (scale 10000)
package word

// Float returns the float64 value of v by dividing the stored int16 by the
// fixed-point scale (10000). It is the exact inverse of VFixFromFloat for any
// value representable at four decimal digits.
func (v VFix) Float() float64 {
	return float64(v) / vfixScale
}
