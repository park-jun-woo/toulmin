//ff:func feature=price type=model control=sequence
//ff:what MemberSpec.SpecName: spec 타입 식별자 반환
package price

// SpecName returns the type identifier for MemberSpec.
func (b *MemberSpec) SpecName() string {
	return "MemberSpec"
}
