//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestRefString — refString renders a Ref as "alias.name" or bare "name"
package tanglcli

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRefString covers both branches of refString: an empty Alias renders
// the bare Name, and a non-empty Alias renders "alias.name".
func TestRefString(t *testing.T) {
	cases := []struct {
		name string
		ref  ast.Ref
		want string
	}{
		{name: "no alias", ref: ast.Ref{Name: "withdraw"}, want: "withdraw"},
		{name: "with alias", ref: ast.Ref{Alias: "bank", Name: "withdraw"}, want: "bank.withdraw"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := refString(c.ref); got != c.want {
				t.Errorf("refString(%+v) = %q, want %q", c.ref, got, c.want)
			}
		})
	}
}
