//ff:func feature=policy type=model control=sequence
//ff:what RoleBacking.BackingName: backing 타입 식별자 반환
package policy

// BackingName returns the type identifier for RoleBacking.
func (b *RoleBacking) BackingName() string { return "RoleBacking" }
