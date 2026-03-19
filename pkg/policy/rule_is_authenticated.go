//ff:func feature=policy type=rule control=sequence
//ff:what IsAuthenticated: 사용자가 인증되었는지 판정
package policy

// IsAuthenticated returns true if the request has a non-nil user.
func IsAuthenticated(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RequestContext)
	return ctx.User != nil, nil
}
