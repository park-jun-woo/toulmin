//ff:type feature=tangl type=model
//ff:what Span — a logical-layer byte range [Off, Off+Len) into Doc.Source
package doc

// Span is a byte range into Doc.Source: the half-open interval [Off, Off+Len).
// It is the logical (Block) layer's own coordinate and is distinct from
// word.Span, which is the wire SPAN anchor of the packet layer.
type Span struct {
	Off uint32
	Len uint32
}
