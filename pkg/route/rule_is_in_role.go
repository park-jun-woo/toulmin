//ff:func feature=route type=rule control=sequence
//ff:what IsInRole: backing(string)으로 지정된 역할을 가졌는지 판정
package route

// IsInRole checks if the user has the role specified by backing (string).
// Same function, different backing — "admin" vs "editor" via graph declaration.
func IsInRole(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	role := backing.(string)
	return ctx.User != nil && ctx.User.Role == role, nil
}
