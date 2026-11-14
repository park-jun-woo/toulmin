//ff:func feature=feature type=model control=sequence
//ff:what RegionSpec.SpecName: spec 타입 식별자 반환
package feature

// SpecName returns the type identifier for RegionSpec.
func (b *RegionSpec) SpecName() string {
	return "RegionSpec"
}
