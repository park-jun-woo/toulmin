//ff:func feature=feature type=model control=sequence
//ff:what PercentageBacking.BackingName: backing 타입 식별자 반환
package feature

// BackingName returns the type identifier for PercentageBacking.
func (b *PercentageBacking) BackingName() string {
	return "PercentageBacking"
}
