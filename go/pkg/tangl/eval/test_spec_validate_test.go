//ff:func feature=tangl type=engine control=sequence
//ff:what Validate — test spec validation (always nil)
package eval

// Validate reports the test spec as always valid.
func (s *testSpec) Validate() error { return nil }
