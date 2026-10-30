//ff:func feature=engine type=util control=sequence
//ff:what attackForce — resolves the optional force argument of Attacks (default 1, panics if invalid)
package toulmin

import "math"

// attackForce returns the attack force for Attacks' variadic argument.
// No value → 1. Panics if more than one value is given, or the value is NaN
// or outside [0.0, 1.0] (same convention as Qualifier).
func attackForce(force []float64) float64 {
	if len(force) == 0 {
		return 1.0
	}
	if len(force) > 1 {
		panic("toulmin: Attacks accepts at most one force value")
	}
	f := force[0]
	if math.IsNaN(f) || f < 0.0 || f > 1.0 {
		panic("toulmin: force must be between 0.0 and 1.0")
	}
	return f
}
