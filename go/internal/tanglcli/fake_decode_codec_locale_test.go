//ff:func feature=cli type=command control=sequence
//ff:what fakeDecodeCodec.Locale — return the fake's configured locale tag
package tanglcli

// Locale returns the fake's configured locale tag.
func (f fakeDecodeCodec) Locale() string { return f.locale }
