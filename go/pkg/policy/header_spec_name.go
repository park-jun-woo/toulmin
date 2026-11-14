//ff:func feature=policy type=model control=sequence
//ff:what SpecName — returns HeaderSpec type identifier
package policy

// SpecName returns the type identifier.
func (b *HeaderSpec) SpecName() string { return "HeaderSpec" }
