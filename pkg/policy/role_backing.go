//ff:type feature=policy type=model
//ff:what RoleBacking: IsInRole rule의 backing 타입 (역할명 + 추출 함수)
package policy

// RoleBacking carries the role name and extraction function for role checks.
type RoleBacking struct {
	Role     string
	RoleFunc func(any) string
}
