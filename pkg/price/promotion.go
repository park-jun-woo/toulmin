//ff:type feature=price type=model
//ff:what Promotion: 프로모션 정보
package price

// Promotion represents an active or inactive promotion.
type Promotion struct {
	Name   string
	Active bool
}
