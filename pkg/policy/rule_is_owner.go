//ff:func feature=policy type=rule control=sequence
//ff:what IsOwner: backing(OwnerBacking)으로 사용자 소유권 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsOwner checks if the user owns the resource.
// Reads UserID and ResourceOwner from RequestContext.
func IsOwner(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	user, _ := ctx.Get("user")
	if user == nil {
		return false, nil
	}
	userID, _ := ctx.Get("userID")
	resourceOwner, _ := ctx.Get("resourceOwner")
	return userID == resourceOwner, nil
}
