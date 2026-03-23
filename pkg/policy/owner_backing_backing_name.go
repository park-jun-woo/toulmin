//ff:func feature=policy type=model control=sequence
//ff:what OwnerBacking.BackingName: backing 타입 식별자 반환
package policy

// BackingName returns the type identifier for OwnerBacking.
func (b *OwnerBacking) BackingName() string { return "OwnerBacking" }
