//ff:func feature=feature type=model control=sequence
//ff:what AttributeBacking.BackingName: backing 타입 식별자 반환
package feature

// BackingName returns the type identifier for AttributeBacking.
func (b *AttributeBacking) BackingName() string {
	return "AttributeBacking"
}
