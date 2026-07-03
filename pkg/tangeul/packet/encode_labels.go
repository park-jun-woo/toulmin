//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what encodeLabels — serialize the EXPR label table ([count] + per label: byte length + padded bytes)
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// encodeLabels serializes an EXPR label table: a count word, then for each label
// a NUM32 byte length followed by the verbatim UTF-8 bytes packed big-endian
// into words with a trailing 0x00 pad when odd. This self-contained table lets
// an ExprPacket round-trip without external Text packets.
func encodeLabels(labels []string) []word.Word {
	out := []word.Word{word.Word(len(labels))}
	for _, l := range labels {
		b := []byte(l)
		n := word.EncodeNum32(uint32(len(b)))
		out = append(out, n[0], n[1])
		if len(b)%2 != 0 {
			b = append(append([]byte(nil), b...), 0x00)
		}
		data, _ := word.ReadWords(b)
		out = append(out, data...)
	}
	return out
}
