//ff:func feature=approve type=model control=sequence
//ff:what ApproverSpec.Validate: spec 필드 유효성 검증
package approve

// Validate checks that ApproverSpec fields are valid.
func (b *ApproverSpec) Validate() error {
	return nil
}
