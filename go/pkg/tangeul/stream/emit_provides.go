//ff:func feature=tangl type=util control=selection
//ff:what emitProvides — dispatch a Provides statement (endpoint decl, require, run, check)
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitProvides routes a Provides coordinate: the endpoint declaration, a required
// field (owner = endpoint), a run edge (PropRuns) or a check edge (PropEvaluates).
// Runs/Checks are case names resolved through the pass-1 case map.
func emitProvides(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	if len(ref.Path) == 1 {
		return emitEndpointDecl(e, ref, elem.(ast.Endpoint), span)
	}
	owner := e.tt.byRef[tidKey(doc.StmtRef{Sec: toulmin.SectionProvides, Path: ref.Path[:1]})]
	switch ref.Path[1] {
	case 0:
		return emitRequire(e, elem.(ast.Require), e.tt.byRef[tidKey(ref)], owner, span)
	case 1:
		return emitCaseEdge(e, toulmin.PropRuns, owner, elem.(string), span)
	case 2:
		return emitCaseEdge(e, toulmin.PropEvaluates, owner, elem.(string), span)
	default:
		return nil, fmt.Errorf("tangeul/stream: unknown provides slotKind %d", ref.Path[1])
	}
}
