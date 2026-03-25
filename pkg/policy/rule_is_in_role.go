//ff:func feature=policy type=rule control=sequence
//ff:what IsInRole: spec(RoleSpec)으로 사용자 역할 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsInRole checks if the user has the role specified by spec (*RoleSpec).
// Reads the role from RequestContext.Role.
func IsInRole(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	user, _ := ctx.Get("user")
	if user == nil {
		return false, nil
	}
	role, _ := ctx.Get("role")
	rb := specs[0].(*RoleSpec)
	return role == rb.Role, nil
}
