//ff:func feature=price type=model control=sequence
//ff:what MemberBacking.BackingName: backing 타입 식별자 반환
package price

// BackingName returns the type identifier for MemberBacking.
func (b *MemberBacking) BackingName() string {
	return "MemberBacking"
}
