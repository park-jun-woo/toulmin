//ff:func feature=state type=rule control=sequence
//ff:what IsOwner: backing(OwnerBacking)의 소유자 판정
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsOwner checks if the user owns the resource.
// Reads UserID and ResourceOwnerID from TransitionContext.
func IsOwner(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	userID, _ := ctx.Get("userID")
	resourceOwnerID, _ := ctx.Get("resourceOwnerID")
	return userID == resourceOwnerID, nil
}
