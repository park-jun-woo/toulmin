//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestClassifyCaseTID — issues TIDs for a case decl, a require and a rule node, and skips attacks/execs
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestClassifyCaseTID walks every branch of classifyCaseTID's switch: a
// one-element path (the case declaration itself, recorded by name), a
// two-element path with slotKind 0 (a required field, TID only), slotKind 1
// (a rule node, recorded case-scoped by name), and slotKind 2 (an attack,
// which declares nothing and must not advance the TID counter or write any
// map entry).
func TestClassifyCaseTID(t *testing.T) {
	cases := []struct {
		name       string
		ref        doc.StmtRef
		elem       any
		wantAssign bool
	}{
		{"case declaration", doc.StmtRef{Sec: 5, Path: []int{0}}, ast.Case{Name: "can withdraw"}, true},
		{"require (slotKind 0)", doc.StmtRef{Sec: 5, Path: []int{0, 0, 1}}, ast.Require{Field: "amount"}, true},
		{"rule node (slotKind 1)", doc.StmtRef{Sec: 5, Path: []int{0, 1, 0}}, ast.Node{Name: "balance sufficient"}, true},
		{"attack (slotKind 2, declares nothing)", doc.StmtRef{Sec: 5, Path: []int{0, 2, 0}}, ast.Attack{Target: "x", Attacker: "y"}, false},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		tt := &tidTable{
			byRef:      map[string]toulmin.TID{},
			caseByName: map[string]toulmin.TID{},
			nodeByCase: map[string]toulmin.TID{},
		}
		next := toulmin.TID(1)

		classifyCaseTID(tt, &next, c.ref, c.elem)

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

	t.Run("case declaration populates caseByName", func(t *testing.T) {
		tt := &tidTable{byRef: map[string]toulmin.TID{}, caseByName: map[string]toulmin.TID{}, nodeByCase: map[string]toulmin.TID{}}
		next := toulmin.TID(1)
		ref := doc.StmtRef{Sec: 5, Path: []int{0}}
		classifyCaseTID(tt, &next, ref, ast.Case{Name: "can withdraw"})
		if tt.caseByName["can withdraw"] != 1 {
			t.Errorf("caseByName[%q] = %d, want 1", "can withdraw", tt.caseByName["can withdraw"])
		}
	})

	t.Run("rule node populates nodeByCase", func(t *testing.T) {
		tt := &tidTable{byRef: map[string]toulmin.TID{}, caseByName: map[string]toulmin.TID{}, nodeByCase: map[string]toulmin.TID{}}
		next := toulmin.TID(1)
		ref := doc.StmtRef{Sec: 5, Path: []int{0, 1, 0}}
		classifyCaseTID(tt, &next, ref, ast.Node{Name: "balance sufficient"})
		if tt.nodeByCase[nodeKey(0, "balance sufficient")] != 1 {
			t.Errorf("nodeByCase missing entry: %v", tt.nodeByCase)
		}
	})
}
