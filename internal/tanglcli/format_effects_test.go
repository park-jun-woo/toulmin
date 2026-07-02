//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestFormatEffects — formatEffects renders a tab-aligned do/undo effect table
package tanglcli

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/effects"
)

// TestFormatEffects covers every branch of formatEffects: the zero-iteration
// loop over an empty slice, and both outcomes of the e.Once conditional
// within a non-empty slice.
func TestFormatEffects(t *testing.T) {
	cases := []struct {
		name    string
		entries []effects.Entry
		want    string
	}{
		{
			name:    "empty entries produce empty output",
			entries: nil,
			want:    "",
		},
		{
			name: "once entry renders the once marker",
			entries: []effects.Entry{
				{Kind: "do", Func: ast.Ref{Alias: "bank", Name: "withdraw"}, Once: true, Case: "can withdraw", Node: "balance sufficient"},
			},
			want: "do bank.withdraw once (can withdraw / balance sufficient)\n",
		},
		{
			name: "non-once entry renders an empty once field",
			entries: []effects.Entry{
				{Kind: "undo", Func: ast.Ref{Alias: "bank", Name: "refund"}, Once: false, Case: "can withdraw", Node: "balance sufficient"},
			},
			want: "undo bank.refund  (can withdraw / balance sufficient)\n",
		},
		{
			name: "multiple entries with mixed once flags and bare refs",
			entries: []effects.Entry{
				{Kind: "do", Func: ast.Ref{Name: "log"}, Once: true, Case: "c1", Node: "n1"},
				{Kind: "undo", Func: ast.Ref{Alias: "bank", Name: "refund"}, Once: false, Case: "c2", Node: "n2"},
			},
			want: "do   log         once (c1 / n1)\n" +
				"undo bank.refund      (c2 / n2)\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := formatEffects(c.entries); got != c.want {
				t.Errorf("formatEffects() =\n%q\nwant\n%q", got, c.want)
			}
		})
	}
}
