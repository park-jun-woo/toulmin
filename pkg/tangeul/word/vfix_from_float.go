//ff:func feature=tangl type=util control=sequence
//ff:what VFixFromFloat — convert a float64 in [-1.0, 1.0] to a VFix, rejecting out-of-range input
package word

import (
	"errors"
	"math"
)

// VFixFromFloat converts f to a VFix (scale 10000), rounding to the nearest
// ten-thousandth. It errors if f falls outside the documented verdict range
// [-1.0, +1.0] so callers never silently truncate an out-of-domain value.
func VFixFromFloat(f float64) (VFix, error) {
	scaled := math.Round(f * vfixScale)
	if scaled < -vfixScale || scaled > vfixScale {
		return 0, errors.New("tangeul/word: VFix value out of range [-1.0, 1.0]")
	}
	return VFix(int16(scaled)), nil
}
