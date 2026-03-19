//ff:func feature=route type=rule control=sequence
//ff:what IsAdminOverride: 관리자 오버라이드 여부 판정
package route

// IsAdminOverride returns true if the user is an admin.
func IsAdminOverride(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	return ctx.User != nil && ctx.User.Role == "admin", nil
}
