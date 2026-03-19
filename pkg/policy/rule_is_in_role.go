//ff:func feature=policy type=rule control=sequence
//ff:what IsInRole: backing(string)으로 지정된 역할을 가졌는지 판정
package policy

// IsInRole checks if the user has the role specified by backing (string).
func IsInRole(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RequestContext)
	role := backing.(string)
	return ctx.User != nil && ctx.User.Role == role, nil
}
