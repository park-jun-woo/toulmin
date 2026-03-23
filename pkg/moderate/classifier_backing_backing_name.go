//ff:func feature=moderate type=model control=sequence
//ff:what ClassifierBacking.BackingName: backing 타입 식별자 반환
package moderate

// BackingName returns the type identifier for ClassifierBacking.
func (b *ClassifierBacking) BackingName() string {
	return "ClassifierBacking"
}
