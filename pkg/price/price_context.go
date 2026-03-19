//ff:type feature=price type=model
//ff:what PriceContext: 가격 판정에 필요한 런타임 컨텍스트
package price

// PriceContext holds per-request facts for price evaluation.
type PriceContext struct {
	User       *User
	Coupons    []Coupon
	Promotions []Promotion
	Metadata   map[string]any
}
