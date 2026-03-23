//ff:func feature=feature type=model control=sequence
//ff:what RegionBacking.BackingName: backing 타입 식별자 반환
package feature

// BackingName returns the type identifier for RegionBacking.
func (b *RegionBacking) BackingName() string {
	return "RegionBacking"
}
