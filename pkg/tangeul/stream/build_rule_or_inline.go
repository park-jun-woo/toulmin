//ff:func feature=tangl type=util control=sequence
//ff:what buildRuleOrInline — a rule node inside a case is a case node, outside it is a Rules rule
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildRuleOrInline routes a Rule node by context: with an open case it is a
// case rule node (appendCaseNode); with no open case it is a Rules inline rule
// carrying the recovered condition tree.
func buildRuleOrInline(d *dec, p packet.NodePacket, name string, labels map[uint16][]string, cond ast.Expr) (doc.StmtRef, error) {
	if d.curCase >= 0 {
		return appendCaseNode(d, p, name, labels)
	}
	idx := len(d.ast.Rules)
	d.ast.Rules = append(d.ast.Rules, ast.InlineRule{Name: name, Cond: cond})
	return doc.StmtRef{Sec: codebook.SectionRules, Path: []int{idx}}, nil
}
