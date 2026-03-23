//ff:func feature=policy type=model control=sequence
//ff:what IPListBacking.BackingName: backing 타입 식별자 반환
package policy

// BackingName returns the type identifier for IPListBacking.
func (b *IPListBacking) BackingName() string { return "IPListBacking" }
