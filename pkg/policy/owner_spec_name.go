//ff:func feature=policy type=model control=sequence
//ff:what OwnerSpec.SpecName: spec 타입 식별자 반환
package policy

// SpecName returns the type identifier for OwnerSpec.
func (b *OwnerSpec) SpecName() string { return "OwnerSpec" }
