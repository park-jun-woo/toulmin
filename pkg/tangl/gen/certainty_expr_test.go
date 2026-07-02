//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestCertaintyExpr — tests certaintyExpr for all four Op branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCertaintyExpr(t *testing.T) {
	tests := []struct {
		name string
		op   string
		want string
	}{
		{"above", "above", "self.Verdict > 1"},
		{"less than", "less than", "self.Verdict < 1"},
		{"at most", "at most", "self.Verdict <= 1"},
		{"at least (default)", "at least", "self.Verdict >= 1"},
		{"unknown falls to default", "something else", "self.Verdict >= 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ast.Certainty{Op: tt.op, Percent: 100}
			got := certaintyExpr(c)
			if got != tt.want {
				t.Errorf("certaintyExpr(%q) = %q, want %q", tt.op, got, tt.want)
			}
		})
	}
}
