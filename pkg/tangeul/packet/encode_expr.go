//ff:func feature=tangl type=util control=sequence
//ff:what EncodeExpr — encode an ExprPacket ([node count] + label table + prefix records) as a Payload
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// EncodeExpr encodes p as one Payload packet (§design 4): a node-count word, the
// self-contained label table, then the ast.Expr flattened into prefix-order
// records. DecodeExpr is its exact inverse over the ast.Expr tree.
func EncodeExpr(p ExprPacket) ([]word.Word, error) {
	var labels []string
	recs, count, err := emitExpr(p.Cond, &labels)
	if err != nil {
		return nil, err
	}
	body := []word.Word{word.Word(count)}
	body = append(body, encodeLabels(labels)...)
	body = append(body, recs...)
	return append(word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: uint32(len(body))}), body...), nil
}
