//ff:type feature=tangl type=parser
//ff:what Span — a raw-source byte range [Off, Off+Len) preserved by the scanner
package scan

// Span is a byte range into the raw source bytes: the half-open interval
// [Off, Off+Len). It carries the same meaning as doc.Span but stays in the
// locale-agnostic scanner layer; the Doc method converts it (see span_doc.go).
type Span struct {
	Off uint32
	Len uint32
}
