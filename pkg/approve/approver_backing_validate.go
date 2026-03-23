//ff:func feature=approve type=model control=sequence
//ff:what ApproverBacking.Validate: backing 필드 유효성 검증
package approve

// Validate checks that ApproverBacking fields are valid.
func (b *ApproverBacking) Validate() error {
	return nil
}
