//ff:func feature=moderate type=rule control=sequence
//ff:what IsVerifiedUser: 인증된 사용자인지 판정
package moderate

// IsVerifiedUser returns true if the author is verified.
func IsVerifiedUser(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ContentContext)
	return ctx.Author.Verified, nil
}
