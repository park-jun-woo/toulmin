//ff:func feature=state type=rule control=sequence
//ff:what IsOwner: backing(OwnerBacking)의 소유자 판정
package state

// IsOwner checks if the user owns the resource.
// Reads UserID and ResourceOwnerID from TransitionContext.
func IsOwner(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*TransitionContext)
	return ctx.UserID == ctx.ResourceOwnerID, nil
}
