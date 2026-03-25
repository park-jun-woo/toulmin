//ff:func feature=state type=model control=sequence
//ff:what ExpirySpec.SpecName: spec 타입 식별자 반환
package state

// SpecName returns the type identifier for ExpirySpec.
func (b *ExpirySpec) SpecName() string { return "ExpirySpec" }
