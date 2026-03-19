//ff:func feature=state type=rule control=sequence
//ff:what IsOwner: backing(OwnerBacking)의 ID 추출 함수로 소유자 판정
package state

// IsOwner checks if the user owns the resource using backing (*OwnerBacking).
func IsOwner(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*TransitionContext)
	b := backing.(*OwnerBacking)
	return b.UserIDFunc(ctx.User) == b.OwnerIDFunc(ctx.Resource), nil
}
