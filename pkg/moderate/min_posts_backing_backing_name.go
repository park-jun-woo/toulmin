//ff:func feature=moderate type=model control=sequence
//ff:what MinPostsBacking.BackingName: backing 타입 식별자 반환
package moderate

// BackingName returns the type identifier for MinPostsBacking.
func (b *MinPostsBacking) BackingName() string {
	return "MinPostsBacking"
}
