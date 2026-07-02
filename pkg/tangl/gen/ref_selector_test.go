//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestRefSelector — tests refSelector for nil ref, aliased ref, and unaliased ref branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRefSelector(t *testing.T) {
	tests := []struct {
		name string
		ref  *ast.Ref
		want string
	}{
		{"nil ref", nil, ""},
		{"aliased ref", &ast.Ref{Alias: "pkg", Name: "Fn"}, "pkg.Fn"},
		{"unaliased ref", &ast.Ref{Name: "Fn"}, "Fn"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := refSelector(tt.ref)
			if got != tt.want {
				t.Errorf("refSelector(%+v) = %q, want %q", tt.ref, got, tt.want)
			}
		})
	}
}
