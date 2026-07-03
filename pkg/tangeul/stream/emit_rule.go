//ff:func feature=tangl type=util control=sequence
//ff:what emitRule — encode a Rules inline rule (condition EXPR payload + Rule node + name label)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitRule encodes a Rules inline rule: its condition tree as an EXPR payload
// (emitted first, so the reader recovers it by position) followed by a Rule node
// and its name label. The rule has no open case, which is how the reader tells a
// Rules rule from an in-case rule node.
func emitRule(e *enc, ref doc.StmtRef, elem any, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	rule := elem.(ast.InlineRule)
	exprWords, err := packet.EncodeExpr(packet.ExprPacket{Cond: rule.Cond})
	if err != nil {
		return nil, err
	}
	pre := exprWords
	var body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeRule, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, rule.Name, span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
