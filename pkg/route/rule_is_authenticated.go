//ff:func feature=route type=rule control=sequence
//ff:what IsAuthenticated: 사용자가 인증되었는지 판정
package route

// IsAuthenticated returns true if the request has an authenticated user.
func IsAuthenticated(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	return ctx.User != nil, nil
}
