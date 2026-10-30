//ff:func feature=tangl type=util control=sequence
//ff:what emitRuleNode — encode an in-case rule node (role, using, checking, with, qualifier, label)
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitRuleNode encodes a case rule node: its role, its crisp reporting level, an
// F_FUNC warrant reference materialized as a Func node, an F_CHECK checked case,
// F_WITH/F_GROUNDS name TEXTs, an F_QUAL qualifier, and its name label. An undeclared
// `checking` case is a hard error, and so is an alias-qualified (cross-document)
// one: the GEUL TOULMIN NODE F_CHECK field is a REF to a case record of the same
// document, so there is no slot for another document's case.
func emitRuleNode(e *enc, ref doc.StmtRef, n ast.Node, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	role, err := codebook.RoleFromAST(n.Role)
	if err != nil {
		return nil, err
	}
	level, err := codebook.LevelFromAST(n.Level)
	if err != nil {
		return nil, err
	}
	p := toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeRule, Role: role, Level: level, Span: span}
	var pre, body []core.Word
	if err := attachRuleAttrs(e, &pre, &body, n.Using, n.With, n.Grounds, span, &p); err != nil {
		return nil, err
	}
	if n.Checking != nil {
		if n.Checking.Alias != "" {
			return nil, fmt.Errorf("tangeul/stream: node %q: alias-qualified checking `%s`.`%s` cannot be encoded: GEUL TOULMIN F_CHECK carries a same-document case REF only", n.Name, n.Checking.Alias, n.Checking.Name)
		}
		dst, ok := e.tt.caseByName[n.Checking.Name]
		if !ok {
			return nil, fmt.Errorf("tangeul/stream: `checking` references undeclared case %q", n.Checking.Name)
		}
		p.Flags |= toulmin.NodeFlagCheck
		ct := toulmin.TID(dst)
		p.Check = &ct
	}
	if n.Qualified != nil {
		v, err := toulmin.VFixFromFloat(*n.Qualified)
		if err != nil {
			return nil, err
		}
		p.Flags |= toulmin.NodeFlagQual
		p.Qual = &v
	}
	appendNode(e, &body, p)
	appendLabel(e, &pre, &body, self, n.Name, span)
	return append(pre, body...), nil
}
