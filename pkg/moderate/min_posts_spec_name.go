//ff:func feature=moderate type=model control=sequence
//ff:what MinPostsSpec.SpecName: spec 타입 식별자 반환
package moderate

// SpecName returns the type identifier for MinPostsSpec.
func (b *MinPostsSpec) SpecName() string {
	return "MinPostsSpec"
}
