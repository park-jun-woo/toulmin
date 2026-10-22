//ff:type feature=price type=model
//ff:what Coupon: 쿠폰 정보
package price

// Coupon represents a coupon that may apply to a purchase.
type Coupon struct {
	Code     string
	MinPrice float64
}
