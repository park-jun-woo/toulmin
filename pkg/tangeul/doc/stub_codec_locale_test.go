//ff:func feature=tangl type=util control=sequence
//ff:what stubCodec.Locale — return the stub's configured locale tag
package doc

// Locale returns the stub's configured locale tag.
func (s stubCodec) Locale() string {
	return s.locale
}
