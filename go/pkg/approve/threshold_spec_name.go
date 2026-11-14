//ff:func feature=approve type=model control=sequence
//ff:what SpecName — returns ThresholdSpec type identifier
package approve

// SpecName returns the type identifier.
func (b *ThresholdSpec) SpecName() string { return "ThresholdSpec" }
