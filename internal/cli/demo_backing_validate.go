//ff:func feature=cli type=model control=sequence
//ff:what demoBacking.Validate — validates demo backing (always valid)
package cli

func (b *demoBacking) Validate() error { return nil }
