//ff:func feature=cli type=model control=sequence
//ff:what demoBacking.BackingName — returns type identifier
package cli

func (b *demoBacking) BackingName() string { return "demoBacking" }
