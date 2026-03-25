//ff:func feature=state type=rule control=sequence
//ff:what IsOwner: backing(OwnerBacking)의 소유자 판정
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsOwner checks if the user owns the resource.
// Reads UserID and ResourceOwnerID from TransitionContext.
func IsOwner(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*TransitionContext)
	return ctx.UserID == ctx.ResourceOwnerID, nil
}
