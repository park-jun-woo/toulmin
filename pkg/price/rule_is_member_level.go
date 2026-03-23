//ff:func feature=price type=rule control=sequence
//ff:what IsMemberLevel: backing(MemberBacking)으로 사용자 멤버십 등급 비교
package price

// IsMemberLevel checks if the user's membership matches backing.Level.
// Returns the Discount as evidence if matched.
func IsMemberLevel(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*PriceContext)
	mb := backing.(*MemberBacking)
	if ctx.Membership == "" {
		return false, nil
	}
	if ctx.Membership == mb.Level {
		return true, mb.Discount
	}
	return false, nil
}
