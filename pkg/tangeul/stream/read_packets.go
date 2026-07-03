//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what readPackets — split a word stream into framed packets, erroring on truncation
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// readPackets walks ws front to back, using the self-describing frame header
// (Len, 32-bit-extension aware) to slice each packet into a rawPacket without
// interpreting its body. A frame whose declared length runs past the end of the
// stream is a hard error (truncation), as is any prefix violation surfaced by
// DecodeFrame.
func readPackets(ws []word.Word) ([]rawPacket, error) {
	var out []rawPacket
	i := 0
	for i < len(ws) {
		f, head, err := word.DecodeFrame(ws[i:])
		if err != nil {
			return nil, err
		}
		total := head + int(f.Len)
		if i+total > len(ws) {
			return nil, errors.New("tangeul/stream: packet framed length runs past end of stream")
		}
		out = append(out, rawPacket{Class: f.Class, Words: ws[i : i+total]})
		i += total
	}
	return out, nil
}
