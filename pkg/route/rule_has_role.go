//ff:func feature=route type=rule control=sequence
//ff:what HasRole: backing(string)으로 지정된 역할을 가졌는지 판정 — 클로저 불필요
package route

// HasRole is a backing-aware rule. backing is the role name (string).
// Replaces the closure factory pattern — same func, different backing.
func HasRole(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	role := backing.(string)
	return ctx.User != nil && ctx.User.Role == role, nil
}
