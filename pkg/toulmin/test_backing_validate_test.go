//ff:func feature=engine type=model control=sequence
//ff:what testBacking.Validate — validates test backing (always valid)
package toulmin

func (b *testBacking) Validate() error { return nil }
