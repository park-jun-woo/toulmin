//ff:type feature=tangl type=model
//ff:what headInfo — what the reader takes from a TANGEUL document's HEAD: the locale
package stream

// headInfo is the HEAD of a TANGEUL document as the reader needs it: the
// authoring locale from Meta DEFAULT-LANGUAGE ("" when absent). SPEC-VERSION and PROFILE are checked, not
// kept.
type headInfo struct {
	locale  string
	spec    bool
	profile bool
}
