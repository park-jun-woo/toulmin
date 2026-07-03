//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what buildRefIndex — index every line-bearing AST element by its source line to a StmtRef
package en

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
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
		idx[s.Line] = doc.StmtRef{Sec: codebook.SectionSee, Path: []int{i}}
	}
	for di, d := range a.Defs {
		idx[d.Line] = doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{di}}
		for fi, f := range d.Fields {
			idx[f.Line] = doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{di, fi}}
		}
	}
	for ri, r := range a.Rules {
		idx[r.Line] = doc.StmtRef{Sec: codebook.SectionRules, Path: []int{ri}}
	}
	for ci, c := range a.Cases {
		idx[c.Line] = doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci}}
		for si, req := range c.Requires {
			idx[req.Line] = doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci, 0, si}}
		}
		for si, n := range c.Nodes {
			idx[n.Line] = doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci, 1, si}}
		}
		for si, at := range c.Attacks {
			idx[at.Line] = doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci, 2, si}}
		}
		for si, ex := range c.Execs {
			idx[ex.Line] = doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci, 3, si}}
		}
		for si, sp := range c.Supersedes {
			idx[sp.Line] = doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci, 4, si}}
		}
	}
	for ei, ep := range a.Provides {
		idx[ep.Line] = doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{ei}}
		for si, req := range ep.Requires {
			idx[req.Line] = doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{ei, 0, si}}
		}
	}
	for ii, in := range a.Internals {
		idx[in.Line] = doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{ii}}
	}
	return idx
}
