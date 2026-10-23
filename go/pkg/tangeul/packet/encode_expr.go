//ff:func feature=tangl type=util control=sequence
//ff:what EncodeExpr — encode a Rules condition tree as a PREDICATE-DECLARE record
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// EncodeExpr encodes p as a PREDICATE-DECLARE record (GEUL t002 TOULMIN EXPR):
// the Edge Id, a node-count word, the self-contained label table, then the
// ast.Expr flattened into prefix-order records. DecodeExpr is its exact inverse.
func EncodeExpr(p ExprPacket, w toulmin.Wire) ([]core.Word, error) {
	var labels []string
	recs, count, err := emitExpr(p.Cond, &labels)
	if err != nil {
		return nil, err
	}
	body := append(toulmin.PutID(nil, p.ID), core.Word(count))
	lw, err := encodeLabels(labels)
	if err != nil {
		return nil, err
	}
	body = append(append(body, lw...), recs...)
	return toulmin.EdgeRecord(toulmin.CodecRef{Codec: toulmin.CodecPredicateDeclare}, body)
}
