//ff:func feature=tangl type=util control=selection
//ff:what buildFromNode — append the AST element a primary node declares and return its StmtRef
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildFromNode records a primary node's name and type, then dispatches to the
// builder for its NodeType (namespace, alias, term, field, rule, case, endpoint,
// trigger), returning the StmtRef of the AST element it appends. Reference-heavy
// types delegate to helpers that resolve their linking edges via findRel.
func buildFromNode(d *dec, p packet.NodePacket, labels map[uint16][]string, rel []packet.EdgePacket, cond ast.Expr) (doc.StmtRef, error) {
	tid := uint16(p.TID)
	name := firstLabel(labels, tid)
	d.nodeType[tid] = p.Type
	d.nodeName[tid] = name

	switch p.Type {
	case codebook.NodeNamespace:
		d.ast.Subject = name
		return doc.StmtRef{Sec: codebook.SectionSubject}, nil
	case codebook.NodeAlias:
		return buildAlias(d, name, rel)
	case codebook.NodeTerm:
		return buildTerm(d, tid, name, p.Role, rel)
	case codebook.NodeField:
		return buildField(d, tid, name, labels, rel)
	case codebook.NodeRule:
		return buildRuleOrInline(d, p, name, labels, cond)
	case codebook.NodeCase:
		idx := len(d.ast.Cases)
		d.ast.Cases = append(d.ast.Cases, ast.Case{Name: name})
		d.caseByTID[tid] = idx
		d.curCase = idx
		return doc.StmtRef{Sec: codebook.SectionCases, Path: []int{idx}}, nil
	case codebook.NodeEndpoint:
		idx := len(d.ast.Provides)
		d.ast.Provides = append(d.ast.Provides, ast.Endpoint{Name: name})
		d.epByTID[tid] = idx
		return doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{idx}}, nil
	case codebook.NodeTrigger:
		return buildTrigger(d, tid, name, p.Role, rel)
	default:
		return doc.StmtRef{}, fmt.Errorf("tangeul/stream: unexpected primary node type 0x%04X", uint16(p.Type))
	}
}
