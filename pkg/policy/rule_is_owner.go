//ff:func feature=policy type=rule control=sequence
//ff:what IsOwner: backing(OwnerBacking)으로 사용자 소유권 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsOwner checks if the user owns the resource.
// Reads UserID and ResourceOwner from RequestContext.
func IsOwner(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*RequestContext)
	if ctx.User == nil {
		return false, nil
	}
	return ctx.UserID == ctx.ResourceOwner, nil
}
