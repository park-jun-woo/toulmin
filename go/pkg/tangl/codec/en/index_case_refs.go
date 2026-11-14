//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what indexCaseRefs — index one Cases entry and its require/node/attack/exec/supersede slots by source line
package en

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// indexCaseRefs records the StmtRef of case ci and of each of its child slots
// (path {ci, slot, si}) into idx. A supersede line is indexed unless the line
// already holds another element and the edge is not a chain head (chain heads,
// with len(Chain) >= 3, always win the line).
func indexCaseRefs(idx map[int]doc.StmtRef, ci int, c ast.Case) {
	idx[c.Line] = doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{ci}}
	for si, req := range c.Requires {
		idx[req.Line] = doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{ci, 0, si}}
	}
	for si, n := range c.Nodes {
		idx[n.Line] = doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{ci, 1, si}}
	}
	for si, at := range c.Attacks {
		idx[at.Line] = doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{ci, 2, si}}
	}
	for si, ex := range c.Execs {
		idx[ex.Line] = doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{ci, 3, si}}
	}
	for si, sp := range c.Supersedes {
		if _, exists := idx[sp.Line]; !exists || len(sp.Chain) >= 3 {
			idx[sp.Line] = doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{ci, 4, si}}
		}
	}
}
