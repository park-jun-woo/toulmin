//ff:func feature=price type=rule control=sequence
//ff:what IsMemberLevel: backing(DiscountBacking).Name과 사용자 멤버십 등급 비교
package price

// IsMemberLevel checks if the user's membership matches backing.Name.
func IsMemberLevel(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*PriceContext)
	db := backing.(*DiscountBacking)
	return ctx.User.Membership == db.Name, db
}
