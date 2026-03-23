//ff:func feature=state type=model control=sequence
//ff:what ExpiryBacking.BackingName: backing 타입 식별자 반환
package state

// BackingName returns the type identifier for ExpiryBacking.
func (b *ExpiryBacking) BackingName() string { return "ExpiryBacking" }
