//ff:func feature=tangl type=parser control=sequence
//ff:what Locale — report the codec's locale tag ("ko")
package ko

// Locale returns the codec's locale tag, "ko", stamped into Doc.Locale by Encode
// and used at codec registration (README contract §3).
func (KoCodec) Locale() string {
	return "ko"
}
