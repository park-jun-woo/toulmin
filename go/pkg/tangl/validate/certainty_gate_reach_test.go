//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestCertaintyGateReach — never/always classification of each gate operator around the (50%, 100%] firing range
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCertaintyGateReach(t *testing.T) {
	cases := []struct {
		op            string
		pct           int
		never, always bool
	}{
		{"at least", 75, false, false},
		{"at least", 50, false, true},
		{"at least", 101, true, false},
		{"above", 100, true, false},
		{"above", 50, false, true},
		{"above", 60, false, false},
		{"less than", 50, true, false},
		{"less than", 101, false, true},
		{"less than", 90, false, false},
		{"at most", 50, true, false},
		{"at most", 100, false, true},
		{"at most", 80, false, false},
		{"sometimes", 10, false, false},
	}
	for _, tc := range cases {
		never, always := certaintyGateReach(ast.Certainty{Op: tc.op, Percent: tc.pct})
		if never != tc.never || always != tc.always {
			t.Errorf("%s %d%%: got (never %v, always %v), want (%v, %v)", tc.op, tc.pct, never, always, tc.never, tc.always)
		}
	}
}
