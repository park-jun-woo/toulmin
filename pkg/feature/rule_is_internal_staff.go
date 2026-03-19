//ff:func feature=feature type=rule control=sequence
//ff:what IsInternalStaff: backing(func)으로 추출한 역할이 "internal"인지 판정
package feature

// IsInternalStaff checks if the user is internal staff.
// backing is func(any) string that extracts role from the domain User.
// If backing is nil, checks Attributes["internal"].
func IsInternalStaff(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	if backing != nil {
		roleFunc := backing.(func(any) string)
		return roleFunc(ctx.User) == "internal", nil
	}
	internal, _ := ctx.Attributes["internal"].(bool)
	return internal, nil
}
