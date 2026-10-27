//ff:func feature=tangl type=util control=sequence
//ff:what buildRuleOrInline — a rule node inside a case is a case node, outside it is a Rules rule
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildRuleOrInline routes a Rule node by context: with an open case it is a
// case rule node (appendCaseNode); with no open case it is a Rules rule carrying
// the recovered condition tree or, for a host rule, its function, `with` specs
// and `grounded on` facts.
func buildRuleOrInline(d *dec, p toulmin.NodePacket, name string, labels map[toulmin.TID][]string, cond ast.Expr) (doc.StmtRef, error) {
	if d.curCase >= 0 {
		return appendCaseNode(d, p, name, labels)
	}
	r := ast.InlineRule{Name: name, Cond: cond, With: textsOf(d, p.With), Grounds: textsOf(d, p.Grounds)}
	if p.Func != nil {
		ref, err := resolveFuncRef(d, p.Func.Alias)
		if err != nil {
			return doc.StmtRef{}, err
		}
		r.Using = &ref
	}
	idx := len(d.ast.Rules)
	d.ast.Rules = append(d.ast.Rules, r)
	return doc.StmtRef{Sec: toulmin.SectionRules, Path: []int{idx}}, nil
}
