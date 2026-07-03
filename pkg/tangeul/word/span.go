//ff:type feature=tangl type=model
//ff:what Span — a 3-word SPAN anchor (target TID + offset + length)
package word

// Span is a GEUL SPAN anchor (3 words): the target document TID, an offset and
// a length (GEUL §5.5). Every TANGEUL packet anchors into the Source Text
// packet via a Span. The Off/Len widths are provisionally 16 bits each while
// the L0 layout is unconfirmed (§11-1); this isolation lets the width change
// without touching packet/stream/doc.
type Span struct {
	TID uint16
	Off uint16
	Len uint16
}

// spanWords is the fixed word count of a SPAN field.
const spanWords = 3
