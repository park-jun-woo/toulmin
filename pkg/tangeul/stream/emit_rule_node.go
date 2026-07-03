//ff:func feature=tangl type=util control=sequence
//ff:what emitRuleNode — encode an in-case rule node (role, using, checking, with, qualifier, label)
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitRuleNode encodes a case rule node: its role, its crisp reporting level, an
// F_FUNC warrant reference materialized as a Func node, an F_CHECK checked case,
// F_WITH term Texts, an F_QUAL qualifier, and its name label. An undeclared
// `checking` case is a hard error.
func emitRuleNode(e *enc, ref doc.StmtRef, n ast.Node, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	role, err := codebook.RoleFromAST(n.Role)
	if err != nil {
		return nil, err
	}
	level, err := codebook.LevelFromAST(n.Level)
	if err != nil {
		return nil, err
	}
	p := packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeRule, Role: role, Level: level, Span: span}
	var pre, body []word.Word
	if n.Using != nil {
		fTID, err := emitFuncRef(e, &pre, &body, *n.Using, span)
		if err != nil {
			return nil, err
		}
		p.Flags |= codebook.NodeFlagFunc
		p.Func = &packet.FuncRef{Alias: packet.TID(fTID)}
	}
	if n.Checking != "" {
		dst, ok := e.tt.caseByName[n.Checking]
		if !ok {
			return nil, fmt.Errorf("tangeul/stream: `checking` references undeclared case %q", n.Checking)
		}
		p.Flags |= codebook.NodeFlagCheck
		ct := packet.TID(dst)
		p.Check = &ct
	}
	if len(n.With) > 0 {
		p.Flags |= codebook.NodeFlagWith
		p.With = withTexts(e, &pre, n.With)
	}
	if n.Qualified != nil {
		v, err := word.VFixFromFloat(*n.Qualified)
		if err != nil {
			return nil, err
		}
		p.Flags |= codebook.NodeFlagQual
		p.Qual = &v
	}
	if err := appendNode(&body, p); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, n.Name, span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
