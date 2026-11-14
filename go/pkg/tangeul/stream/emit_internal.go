//ff:func feature=tangl type=util control=sequence
//ff:what emitInternal — dispatch a tangl:Internal statement (trigger decl, run, check)
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitInternal routes a tangl:Internal coordinate: the trigger head declaration,
// or a run (PropTriggers) or check (PropEvaluates) edge to a case. Runs and
// Checks are separate statements routed by the trigger's TID as the edge src.
func emitInternal(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	if len(ref.Path) == 1 {
		return emitInternalDecl(e, ref, elem.(ast.Internal), span)
	}
	owner := e.tt.byRef[tidKey(doc.StmtRef{Sec: toulmin.SectionInternal, Path: ref.Path[:1]})]
	prop := toulmin.PropTriggers
	if ref.Path[1] == 2 {
		prop = toulmin.PropEvaluates
	}
	return emitCaseEdge(e, prop, owner, elem.(string), span)
}
