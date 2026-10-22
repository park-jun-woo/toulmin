//ff:type feature=tangl type=model
//ff:what headInfo — what the reader takes from a TANGEUL document's HEAD: Edge Id width and locale
package stream

// headInfo is the HEAD of a TANGEUL document as the reader needs it: the Edge Id
// width from Meta ID-WIDTH (1 when absent) and the authoring locale from Meta
// DEFAULT-LANGUAGE ("" when absent). SPEC-VERSION and PROFILE are checked, not
// kept.
type headInfo struct {
	idWidth int
	locale  string
	spec    bool
	profile bool
}
