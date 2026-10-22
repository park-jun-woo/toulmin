//ff:func feature=tangl type=util control=sequence
//ff:what emitRule — encode a Rules rule (optional condition PREDICATE-DECLARE + RULE-DECLARE node with host attrs + label)
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitRule encodes a Rules-section rule: a condition rule's tree as a
// PREDICATE-DECLARE record first (so the reader finds it in the prefix), then the
// RULE-DECLARE node — carrying a host rule's function, `with` specs and `grounded
// on` facts — and its name label. The rule has no open case, which is how the
// reader tells a Rules rule from an in-case rule node.
func emitRule(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	rule := elem.(ast.InlineRule)
	var pre, body []core.Word
	if rule.Cond != nil {
		exprWords, err := packet.EncodeExpr(packet.ExprPacket{ID: issueID(e), Cond: rule.Cond}, e.wire)
		if err != nil {
			return nil, err
		}
		pre = exprWords
	}
	p := toulmin.NodePacket{TID: self, Type: toulmin.NodeRule, Span: span}
	if err := attachRuleAttrs(e, &pre, &body, rule.Using, rule.With, rule.Grounds, span, &p); err != nil {
		return nil, err
	}
	if err := appendNode(e, &body, p); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, rule.Name, span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
