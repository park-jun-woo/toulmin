//ff:func feature=tangl type=util control=selection
//ff:what emitCaseSlot — dispatch a case-body slot (require, rule node, attack, exec)
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitCaseSlot routes a case-body statement by slotKind (0 Require, 1 Node,
// 2 Attack, 3 Exec, 4 Supersede) to its emitter, restoring the document order the
// split AST slices lose.
func emitCaseSlot(e *enc, ref doc.StmtRef, elem any, span word.Span) ([]word.Word, error) {
	switch ref.Path[1] {
	case 0:
		self := e.tt.byRef[tidKey(ref)]
		owner := e.tt.byRef[tidKey(doc.StmtRef{Sec: codebook.SectionCases, Path: ref.Path[:1]})]
		return emitRequire(e, elem.(ast.Require), self, owner, span)
	case 1:
		return emitRuleNode(e, ref, elem.(ast.Node), span)
	case 2:
		return emitAttack(e, ref, elem.(ast.Attack), span)
	case 3:
		return emitExecStmt(e, ref, elem.(ast.Exec), span)
	case 4:
		return emitSupersede(e, ref, elem.(ast.Supersede), span)
	default:
		return nil, fmt.Errorf("tangeul/stream: unknown case slotKind %d", ref.Path[1])
	}
}
