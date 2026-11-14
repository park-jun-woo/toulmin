//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderCertainty — renderCertainty emits `if <op> <N>% certain` for each operator
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderCertainty checks each of the four threshold operators.
func TestRenderCertainty(t *testing.T) {
	cases := map[string]ast.Certainty{
		"if at least 75% certain":  {Op: "at least", Percent: 75},
		"if above 50% certain":     {Op: "above", Percent: 50},
		"if less than 25% certain": {Op: "less than", Percent: 25},
		"if at most 90% certain":   {Op: "at most", Percent: 90},
	}
	for want, c := range cases {
		if got := renderCertainty(c); got != want {
			t.Errorf("renderCertainty(%+v) = %q, want %q", c, got, want)
		}
	}
}
