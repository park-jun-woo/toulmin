//ff:type feature=policy type=model
//ff:what RoleSpec: IsInRole rule의 spec 타입 (역할명)
package policy

// RoleSpec carries the role name for role checks.
type RoleSpec struct {
	Role string
}
