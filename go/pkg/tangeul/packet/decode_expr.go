//ff:func feature=tangl type=util control=sequence
//ff:what DecodeExpr — decode a PREDICATE-DECLARE record back into its condition tree
package packet

import (
	"errors"
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// DecodeExpr decodes r as a PREDICATE-DECLARE record (Edge Id, node count, label
// table, prefix records) back into an ExprPacket. The rebuilt tree must consume
// every record word and match the declared node count.
func DecodeExpr(r core.Record, w toulmin.Wire) (ExprPacket, error) {
	ref, err := toulmin.RecordRef(r)
	if err != nil {
		return ExprPacket{}, err
	}
	if ref.Codec != toulmin.CodecPredicateDeclare {
		return ExprPacket{}, errors.New("tangeul/packet: record is not a PREDICATE-DECLARE")
	}
	ws, pos := r.Payload, 0
	id, err := toulmin.TakeID(ws, &pos)
	if err != nil {
		return ExprPacket{}, err
	}
	n, err := toulmin.TakeWord(ws, &pos, "EXPR node count")
	if err != nil {
		return ExprPacket{}, err
	}
	labels, consumed, err := decodeLabels(ws[pos:])
	if err != nil {
		return ExprPacket{}, err
	}
	recs, rpos := ws[pos+consumed:], 0
	cond, count, err := parseExpr(recs, &rpos, labels)
	if err != nil {
		return ExprPacket{}, err
	}
	if rpos != len(recs) {
		return ExprPacket{}, errors.New("tangeul/packet: EXPR records not fully consumed")
	}
	if count != int(n) {
		return ExprPacket{}, fmt.Errorf("tangeul/packet: EXPR node count mismatch (declared %d, decoded %d)", n, count)
	}
	return ExprPacket{ID: id, Cond: cond}, nil
}
