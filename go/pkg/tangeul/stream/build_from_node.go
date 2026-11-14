//ff:func feature=tangl type=util control=selection
//ff:what buildFromNode — append the AST element a primary node declares and return its StmtRef
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildFromNode records a primary node's name and type, then dispatches to the
// builder for its NodeType (namespace, alias, term, field, rule, case, endpoint,
// trigger), returning the StmtRef of the AST element it appends. Reference-heavy
// types delegate to helpers that resolve their linking edges via findRel.
func buildFromNode(d *dec, p toulmin.NodePacket, labels map[toulmin.TID][]string, rel []toulmin.EdgePacket, cond ast.Expr) (doc.StmtRef, error) {
	tid := toulmin.TID(p.TID)
	name := firstLabel(labels, tid)
	d.nodeType[tid] = p.Type
	d.nodeName[tid] = name

	switch p.Type {
	case toulmin.NodeNamespace:
		d.ast.Subject = name
		return doc.StmtRef{Sec: toulmin.SectionSubject}, nil
	case toulmin.NodeAlias:
		return buildAlias(d, name, rel)
	case toulmin.NodeTerm:
		return buildTerm(d, tid, name, p.Role, rel)
	case toulmin.NodeField:
		return buildField(d, tid, name, labels, rel)
	case toulmin.NodeRule:
		return buildRuleOrInline(d, p, name, labels, cond)
	case toulmin.NodeCase:
		idx := len(d.ast.Cases)
		d.ast.Cases = append(d.ast.Cases, ast.Case{Name: name})
		d.caseByTID[tid] = idx
		d.curCase = idx
		return doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{idx}}, nil
	case toulmin.NodeEndpoint:
		idx := len(d.ast.Provides)
		d.ast.Provides = append(d.ast.Provides, ast.Endpoint{Name: name})
		d.epByTID[tid] = idx
		return doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{idx}}, nil
	case toulmin.NodeTrigger:
		return buildTrigger(d, tid, name, p.Role, rel)
	default:
		return doc.StmtRef{}, fmt.Errorf("tangeul/stream: unexpected primary node type 0x%04X", uint16(p.Type))
	}
}
