//ff:func feature=tangl type=util control=selection
//ff:what LocaleCodeFromTag — map a locale tag ("en"/"ko") to its STREAM_HEADER code
package codebook

// LocaleCodeFromTag returns the STREAM_HEADER authoring-locale code for a locale
// tag. An empty or unrecognized tag maps to LocaleUnknown, so a Doc with no
// authoring locale serializes a 0 code.
func LocaleCodeFromTag(tag string) uint16 {
	switch tag {
	case "en":
		return LocaleEN
	case "ko":
		return LocaleKO
	default:
		return LocaleUnknown
	}
}
