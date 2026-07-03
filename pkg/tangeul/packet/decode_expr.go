//ff:func feature=tangl type=util control=sequence
//ff:what DecodeExpr — decode an EXPR Payload back into an ExprPacket, verifying node/word counts
package packet

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// DecodeExpr decodes an EXPR Payload from the head of ws (node count, label
// table, prefix records) back into an ExprPacket and returns the total words
// consumed. It verifies the rebuilt subtree consumes every record word and that
// the node count matches the framed N — a mismatch is a hard error.
func DecodeExpr(ws []word.Word) (ExprPacket, int, error) {
	f, head, err := word.DecodeFrame(ws)
	if err != nil {
		return ExprPacket{}, 0, err
	}
	if f.Class != word.ClassPayload {
		return ExprPacket{}, 0, fmt.Errorf("tangeul/packet: expected Payload class, got %d", uint8(f.Class))
	}
	total := head + int(f.Len)
	if len(ws) < total {
		return ExprPacket{}, 0, errors.New("tangeul/packet: EXPR body shorter than framed length")
	}
	body := ws[head:total]
	if len(body) < 1 {
		return ExprPacket{}, 0, errors.New("tangeul/packet: EXPR body requires a node count")
	}
	n := int(body[0])
	labels, consumed, err := decodeLabels(body[1:])
	if err != nil {
		return ExprPacket{}, 0, err
	}
	recs := body[1+consumed:]
	pos := 0
	cond, count, err := parseExpr(recs, &pos, labels)
	if err != nil {
		return ExprPacket{}, 0, err
	}
	if pos != len(recs) {
		return ExprPacket{}, 0, errors.New("tangeul/packet: EXPR records not fully consumed")
	}
	if count != n {
		return ExprPacket{}, 0, fmt.Errorf("tangeul/packet: EXPR node count mismatch (framed %d, decoded %d)", n, count)
	}
	return ExprPacket{Cond: cond}, total, nil
}
