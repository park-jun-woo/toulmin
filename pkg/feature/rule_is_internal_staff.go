//ff:func feature=feature type=rule control=sequence
//ff:what IsInternalStaff: 내부 직원인지 판정
package feature

// IsInternalStaff checks if the user is internal staff.
// Checks Attributes["internal"].
func IsInternalStaff(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	internal, _ := ctx.Attributes["internal"].(bool)
	return internal, nil
}
