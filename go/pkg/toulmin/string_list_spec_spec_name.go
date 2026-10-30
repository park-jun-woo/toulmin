//ff:func feature=engine type=model control=sequence
//ff:what SpecName — returns the list Spec's name
package toulmin

// SpecName returns the list Spec's name.
func (s StringListSpec) SpecName() string { return s.name }
