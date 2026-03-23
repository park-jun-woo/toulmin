//ff:type feature=policy type=model
//ff:what RoleBacking: IsInRole rule의 backing 타입 (역할명)
package policy

// RoleBacking carries the role name for role checks.
type RoleBacking struct {
	Role string
}
