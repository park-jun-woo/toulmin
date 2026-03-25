//ff:func feature=policy type=rule control=sequence
//ff:what IsInRole: backing(RoleBacking)으로 사용자 역할 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsInRole checks if the user has the role specified by backing (*RoleBacking).
// Reads the role from RequestContext.Role.
func IsInRole(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*RequestContext)
	if ctx.User == nil {
		return false, nil
	}
	rb := backing.(*RoleBacking)
	return ctx.Role == rb.Role, nil
}
