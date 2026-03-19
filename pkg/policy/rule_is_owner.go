//ff:func feature=policy type=rule control=sequence
//ff:what IsOwner: 사용자가 리소스 소유자인지 판정
package policy

// IsOwner checks if the user owns the resource.
func IsOwner(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RequestContext)
	if ctx.User == nil {
		return false, nil
	}
	return ctx.User.ID == ctx.ResourceOwnerID, nil
}
