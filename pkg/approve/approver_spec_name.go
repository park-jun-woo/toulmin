//ff:func feature=approve type=model control=sequence
//ff:what ApproverSpec.SpecName: spec 타입 식별자 반환
package approve

// SpecName returns the type identifier for ApproverSpec.
func (b *ApproverSpec) SpecName() string {
	return "ApproverSpec"
}
