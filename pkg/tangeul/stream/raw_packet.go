//ff:type feature=tangl type=model
//ff:what rawPacket — one framed packet lifted from the stream (its class and complete words)
package stream

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// rawPacket is a single framed packet sliced out of the stream: its Class (from
// the frame header) and the complete word run (header plus body). The reader
// dispatches on Class, and for the Payload class resolves Text vs EXPR by
// statement context, so framing is separated from typed decoding.
type rawPacket struct {
	Class word.Class
	Words []word.Word
}
