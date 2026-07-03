//ff:func feature=tangl type=util control=sequence
//ff:what EncodeText — encode a TextPacket (format + NUM32 byte length + verbatim padded bytes)
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// EncodeText encodes p as a Text payload (§7.4): a UTF-8 format word, a NUM32
// byte length, then the verbatim bytes packed big-endian into words with a
// trailing 0x00 pad when the length is odd. The byte length is a NUM32 (not the
// single word of §7.4) so it can address payloads that the 32-bit frame-length
// extension carries (README §design 2, large-payload consistency). TID is not
// emitted — it is stream-assigned.
func EncodeText(p TextPacket) []word.Word {
	padded := p.Bytes
	if len(padded)%2 != 0 {
		padded = append(append([]byte(nil), padded...), 0x00)
	}
	data, _ := word.ReadWords(padded)
	n := word.EncodeNum32(uint32(len(p.Bytes)))
	body := append([]word.Word{word.Word(TextFormatUTF8), n[0], n[1]}, data...)
	return append(word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: uint32(len(body))}), body...)
}
