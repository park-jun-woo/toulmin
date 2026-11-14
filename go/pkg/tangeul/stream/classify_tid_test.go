//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestClassifyTID — issues TIDs per section, resolves See/Cases name maps, and rejects an unknown section
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestClassifyTID walks every ref.Sec case classifyTID routes: Subject,
// Definitions and Rules (bare assignTID), See (aliasByName), Cases (delegates
// to classifyCaseTID), Provides both sub-conditions of `len(Path)==1 ||
// Path[1]==0` (true via each disjunct, and false), Internal's `len(Path)==1`
// (true and false), and an unrecognized section code (a hard error).
func TestClassifyTID(t *testing.T) {
	newTT := func() *tidTable {
		return &tidTable{
			byRef:       map[string]toulmin.TID{},
			aliasByName: map[string]toulmin.TID{},
			caseByName:  map[string]toulmin.TID{},
			nodeByCase:  map[string]toulmin.TID{},
		}
	}

	cases := []struct {
		name       string
		ref        doc.StmtRef
		elem       any
		wantErr    bool
		wantAssign bool
	}{
		{"subject", doc.StmtRef{Sec: toulmin.SectionSubject}, "sys", false, true},
		{"definitions", doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{0}}, ast.Definition{Name: "n"}, false, true},
		{"rules", doc.StmtRef{Sec: toulmin.SectionRules, Path: []int{0}}, ast.InlineRule{Name: "r"}, false, true},
		{"see", doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{0}}, ast.See{Alias: "bank"}, false, true},
		{"cases decl", doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}}, ast.Case{Name: "c"}, false, true},
		{"provides, one-element path", doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{0}}, ast.Endpoint{Name: "e"}, false, true},
		{"provides, Path[1]==0", doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{0, 0, 0}}, ast.Require{Field: "f"}, false, true},
		{"provides, Path[1]!=0", doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{0, 1, 0}}, "run target", false, false},
		{"internal, one-element path", doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{0}}, ast.Internal{Kind: ast.OnEvent}, false, true},
		{"internal, multi-element path", doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{0, 1, 0}}, "run target", false, false},
		{"unknown section", doc.StmtRef{Sec: toulmin.SectionCode(99)}, nil, true, false},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		tt := newTT()
		next := toulmin.TID(1)

		err := classifyTID(tt, &next, c.ref, c.elem)
		if c.wantErr {
			if err == nil {
				t.Fatalf("classifyTID(%+v): expected error, got nil", c.ref)
			}
			return
		}
		if err != nil {
			t.Fatalf("classifyTID(%+v): %v", c.ref, err)
		}

		_, assigned := tt.byRef[tidKey(c.ref)]
		if assigned != c.wantAssign {
			t.Errorf("%s: byRef assigned = %v, want %v", c.name, assigned, c.wantAssign)
		}
		if c.wantAssign && next != 2 {
			t.Errorf("%s: next = %d, want 2", c.name, next)
		}
		if !c.wantAssign && next != 1 {
			t.Errorf("%s: next = %d, want 1 (unchanged)", c.name, next)
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}

	t.Run("see populates aliasByName", func(t *testing.T) {
		tt := newTT()
		next := toulmin.TID(1)
		ref := doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{0}}
		if err := classifyTID(tt, &next, ref, ast.See{Alias: "bank"}); err != nil {
			t.Fatalf("classifyTID: %v", err)
		}
		if tt.aliasByName["bank"] != 1 {
			t.Errorf("aliasByName[bank] = %d, want 1", tt.aliasByName["bank"])
		}
	})
}
