//ff:func feature=feature type=model control=sequence
//ff:what PercentageSpec.SpecName: spec 타입 식별자 반환
package feature

// SpecName returns the type identifier for PercentageSpec.
func (b *PercentageSpec) SpecName() string {
	return "PercentageSpec"
}
