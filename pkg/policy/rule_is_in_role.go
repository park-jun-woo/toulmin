//ff:func feature=policy type=rule control=sequence
//ff:what IsInRole: backing(RoleBacking)으로 사용자 역할 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsInRole checks if the user has the role specified by backing (*RoleBacking).
// Reads the role from RequestContext.Role.
func IsInRole(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	user, _ := ctx.Get("user")
	if user == nil {
		return false, nil
	}
	role, _ := ctx.Get("role")
	rb := backing.(*RoleBacking)
	return role == rb.Role, nil
}
