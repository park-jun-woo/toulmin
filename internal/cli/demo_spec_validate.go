//ff:func feature=cli type=model control=sequence
//ff:what demoSpec.Validate — validates demo spec (always valid)
package cli

func (b *demoSpec) Validate() error { return nil }
