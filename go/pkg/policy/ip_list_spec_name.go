//ff:func feature=policy type=model control=sequence
//ff:what IPListSpec.SpecName: spec 타입 식별자 반환
package policy

// SpecName returns the type identifier for IPListSpec.
func (b *IPListSpec) SpecName() string { return "IPListSpec" }
