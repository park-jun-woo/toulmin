//ff:func feature=tangl type=engine control=sequence
//ff:what SpecName — returns the test spec's name
package eval

// SpecName returns the test spec's name.
func (s *testSpec) SpecName() string { return s.Name }
