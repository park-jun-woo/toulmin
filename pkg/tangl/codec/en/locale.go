//ff:func feature=tangl type=parser control=sequence
//ff:what Locale — report the codec's locale tag ("en")
package en

// Locale returns the codec's locale tag, "en", stamped into Doc.Locale by Encode
// and used at codec registration (README contract §3).
func (ENCodec) Locale() string {
	return "en"
}
