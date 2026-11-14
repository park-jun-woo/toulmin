//ff:func feature=tangl type=util control=sequence
//ff:what CheckUnitInterval — rejects a weight (qualifier/force) that is NaN, outside [0, 1], or finer than 4 decimals
package ast

import (
	"errors"
	"math"
)

// CheckUnitInterval reports whether v is a valid TANGEUL weight — a node
// qualifier (`qualified`/`가중치`) or an attack force (`with force`/`세기`). Both
// are carried in GEUL as a SIGNED-DECIMAL at four decimal places, so a value
// that is NaN, outside [0, 1], or not exact at four decimal places is an error
// rather than a silently rounded weight.
func CheckUnitInterval(v float64) error {
	if math.IsNaN(v) {
		return errors.New("weight is NaN")
	}
	if v < 0 || v > 1 {
		return errors.New("weight out of range [0, 1]")
	}
	scaled := v * 10000
	if math.Abs(scaled-math.Round(scaled)) > 1e-6 {
		return errors.New("weight has more than 4 decimal places")
	}
	return nil
}
