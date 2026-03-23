//ff:func feature=approve type=model control=sequence
//ff:what ApproverBacking.BackingName: backing 타입 식별자 반환
package approve

// BackingName returns the type identifier for ApproverBacking.
func (b *ApproverBacking) BackingName() string {
	return "ApproverBacking"
}
