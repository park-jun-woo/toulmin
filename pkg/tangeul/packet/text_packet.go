//ff:type feature=tangl type=model
//ff:what TextPacket — a Text payload (Source Text or a label string), verbatim bytes
package packet

// TextPacket is a Text payload (Class 10, §7.4): a format word, a NUM32 byte
// length, then the verbatim bytes padded with 0x00 to a word boundary. Both the
// Source Text capsule and every label/path string ride in this packet class.
// Bytes are stored verbatim (no Unicode normalization, §7.4). TID is the
// segment-local id the stream layer assigns to this packet so SPAN anchors can
// target it; it is NOT part of the encoded body, so a bare packet round-trip
// preserves only Bytes.
type TextPacket struct {
	TID   TID
	Bytes []byte
}
