//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestApplyExecFixups — a zero-fixup pass is a no-op; a non-empty pass fills each Exec.Case field by name
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestApplyExecFixups exercises applyExecFixups' loop with both its
// zero-iteration case (no deferred execCaseFixups, Execs untouched) and its
// nonzero-iteration case (two deferred fixups filled from d.nodeName by tid).
func TestApplyExecFixups(t *testing.T) {
	cases := []struct {
		name    string
		fixups  []nameFixup
		wantExe []string // parallel to the two execs below, "" means untouched
	}{
		{
			name:    "zero fixups is a no-op",
			fixups:  nil,
			wantExe: []string{"", ""},
		},
		{
			name: "two fixups fill Case by resolved name",
			fixups: []nameFixup{
				{caseIdx: 0, slotIdx: 0, tid: 4},
				{caseIdx: 0, slotIdx: 1, tid: 9},
			},
			wantExe: []string{"can withdraw", "can deposit"},
		},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		d := &dec{
			ast: &ast.Document{Cases: []ast.Case{{
				Name:  "c0",
				Execs: []ast.Exec{{Node: "n0"}, {Node: "n1"}},
			}}},
			nodeName:       map[toulmin.TID]string{4: "can withdraw", 9: "can deposit"},
			execCaseFixups: c.fixups,
		}

		applyExecFixups(d)

		for j, want := range c.wantExe {
			if got := d.ast.Cases[0].Execs[j].Case; got != want {
				t.Errorf("Execs[%d].Case = %q, want %q", j, got, want)
			}
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}
}
