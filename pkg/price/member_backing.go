//ff:type feature=price type=model
//ff:what MemberBacking: IsMemberLevel rule의 backing 타입
package price

// MemberBacking carries membership criteria and extraction function.
type MemberBacking struct {
	Level          string              // membership level to match ("basic", "gold", "vip")
	MembershipFunc func(any) string   // extracts membership from domain User
	Discount       *DiscountBacking    // discount to apply if matched
}
