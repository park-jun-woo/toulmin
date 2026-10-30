//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what buildRefIndex — index every line-bearing AST element by its source line to a StmtRef
package en

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildRefIndex maps each line-bearing terminal AST element to its StmtRef by the
// element's 1-based source line (README contract §2 path table). Subject has no
// Line (Encode assigns it by section position) and Provides/Internal run/check
// bodies are `[]string` without a Line (matchBodyRefs fills them); everything else
// — See, Definitions and their fields, Rules heads, Cases with their require/node/
// attack/exec slots, endpoint and internal declarations, and endpoint requires —
// is indexed here.
func buildRefIndex(a *ast.Document) map[int]doc.StmtRef {
	idx := make(map[int]doc.StmtRef)
	for i, s := range a.Sees {
		idx[s.Line] = doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{i}}
	}
	for di, d := range a.Defs {
		idx[d.Line] = doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{di}}
		for fi, f := range d.Fields {
			idx[f.Line] = doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{di, fi}}
		}
	}
	for ri, r := range a.Rules {
		idx[r.Line] = doc.StmtRef{Sec: toulmin.SectionRules, Path: []int{ri}}
	}
	for ci, c := range a.Cases {
		indexCaseRefs(idx, ci, c)
	}
	for ei, ep := range a.Provides {
		idx[ep.Line] = doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{ei}}
		for si, req := range ep.Requires {
			idx[req.Line] = doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{ei, 0, si}}
		}
	}
	for ii, in := range a.Internals {
		idx[in.Line] = doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{ii}}
	}
	return idx
}
