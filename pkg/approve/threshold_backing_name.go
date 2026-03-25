//ff:func feature=approve type=model control=sequence
//ff:what BackingName — returns ThresholdBacking type identifier
package approve

// BackingName returns the type identifier.
func (b *ThresholdBacking) BackingName() string { return "ThresholdBacking" }
