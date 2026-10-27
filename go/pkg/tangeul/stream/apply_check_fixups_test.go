//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestApplyCheckFixups — a zero-fixup pass is a no-op; a non-empty pass fills each Checking field by name
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestApplyCheckFixups exercises applyCheckFixups' loop with both its
// zero-iteration case (no deferred checkFixups, Nodes untouched) and its
// nonzero-iteration case (two deferred fixups across cases/slots, each filled
// from d.nodeName by tid).
func TestApplyCheckFixups(t *testing.T) {
	cases := []struct {
		name        string
		checkFixups []nameFixup
		wantChecks  []string // parallel to the two nodes below, "" means untouched
	}{
		{
			name:        "zero fixups is a no-op",
			checkFixups: nil,
			wantChecks:  []string{"", ""},
		},
		{
			name: "two fixups fill Checking by resolved name",
			checkFixups: []nameFixup{
				{caseIdx: 0, slotIdx: 0, tid: 4},
				{caseIdx: 0, slotIdx: 1, tid: 9},
			},
			wantChecks: []string{"can check", "can deposit"},
		},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		d := &dec{
			ast: &ast.Document{Cases: []ast.Case{{
				Name:  "c0",
				Nodes: []ast.Node{{Name: "n0"}, {Name: "n1"}},
			}}},
			nodeName:    map[toulmin.TID]string{4: "can check", 9: "can deposit"},
			checkFixups: c.checkFixups,
		}

		applyCheckFixups(d)

		for j, want := range c.wantChecks {
			got := d.ast.Cases[0].Nodes[j].Checking
			if want == "" {
				if got != nil {
					t.Errorf("Nodes[%d].Checking = %+v, want nil", j, got)
				}
				continue
			}
			if got == nil || got.Name != want {
				t.Errorf("Nodes[%d].Checking = %+v, want %q", j, got, want)
			}
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}
}
