//ff:type feature=tangl type=model
//ff:what testCoupon — a Coupon-style struct for exists-iteration tests
package eval

// testCoupon mirrors the demo Coupon: a struct with a MinPrice field iterated
// by an exists quantifier.
type testCoupon struct {
	MinPrice float64
}
