//ff:func feature=engine type=model control=sequence
//ff:what testBacking.BackingName — returns type identifier for test backing
package toulmin

func (b *testBacking) BackingName() string { return "testBacking" }
