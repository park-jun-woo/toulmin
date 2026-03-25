//ff:type feature=price type=model
//ff:what MemberSpec: IsMemberLevel rule의 spec 타입
package price

// MemberSpec carries membership criteria.
type MemberSpec struct {
	Level    string           // membership level to match ("basic", "gold", "vip")
	Discount *DiscountSpec // discount to apply if matched
}
