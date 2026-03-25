//ff:func feature=policy type=model control=sequence
//ff:what RoleSpec.SpecName: spec 타입 식별자 반환
package policy

// SpecName returns the type identifier for RoleSpec.
func (b *RoleSpec) SpecName() string { return "RoleSpec" }
