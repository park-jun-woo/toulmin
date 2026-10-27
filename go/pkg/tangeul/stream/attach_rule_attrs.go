//ff:func feature=tangl type=util control=sequence
//ff:what attachRuleAttrs — put a rule's host function, `with` specs and `grounded on` facts into its NODE
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// attachRuleAttrs fills p's judgment attributes shared by Rules-section rules and
// in-case rule nodes: the host function (`쓰는`/`using`) materialized as a FUNC
// node behind F_FUNC, the `with` spec names behind F_WITH and the `grounded on`
// fact names behind F_GROUNDS, each name carried as a TEXT record.
func attachRuleAttrs(e *enc, pre, body *[]core.Word, using *ast.Ref, with, grounds []string, span toulmin.Span, p *toulmin.NodePacket) error {
	if using != nil {
		fTID, err := emitFuncRef(e, pre, body, *using, span)
		if err != nil {
			return err
		}
		p.Flags |= toulmin.NodeFlagFunc
		p.Func = &toulmin.FuncRef{Alias: fTID}
	}
	if len(with) > 0 {
		p.Flags |= toulmin.NodeFlagWith
		p.With = withTexts(e, pre, with)
	}
	if len(grounds) > 0 {
		p.Flags |= toulmin.NodeFlagGrounds
		p.Grounds = withTexts(e, pre, grounds)
	}
	return nil
}
