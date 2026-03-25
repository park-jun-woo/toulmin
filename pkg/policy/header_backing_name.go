//ff:func feature=policy type=model control=sequence
//ff:what BackingName — returns HeaderBacking type identifier
package policy

// BackingName returns the type identifier.
func (b *HeaderBacking) BackingName() string { return "HeaderBacking" }
