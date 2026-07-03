//ff:func feature=tangl type=util control=sequence
//ff:what LookupCodec — return the registered Codec for a locale, if any
package doc

// LookupCodec returns the Codec registered for locale and whether one exists.
func LookupCodec(locale string) (Codec, bool) {
	c, ok := codecRegistry[locale]
	return c, ok
}
