//ff:func feature=policy type=rule control=sequence
//ff:what IsOwner: backing(OwnerBacking)으로 사용자 소유권 판정
package policy

// IsOwner checks if the user owns the resource using backing (*OwnerBacking).
func IsOwner(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RequestContext)
	if ctx.User == nil {
		return false, nil
	}
	ob := backing.(*OwnerBacking)
	return ob.UserIDFunc(ctx.User) == ob.ResourceIDFunc(ctx), nil
}
