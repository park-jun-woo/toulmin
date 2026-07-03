//ff:func feature=tangl type=util control=selection
//ff:what LocaleTagFromCode — map a STREAM_HEADER locale code back to its tag
package codebook

// LocaleTagFromCode returns the locale tag for a STREAM_HEADER authoring-locale
// code, the inverse of LocaleCodeFromTag. LocaleUnknown and any unrecognized
// code map to the empty tag.
func LocaleTagFromCode(code uint16) string {
	switch code {
	case LocaleEN:
		return "en"
	case LocaleKO:
		return "ko"
	default:
		return ""
	}
}
