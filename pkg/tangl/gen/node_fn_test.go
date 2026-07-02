//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestNodeFn — tests nodeFn for using-ref, checking-wrapper-hit, checking-wrapper-miss, and plain-name branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestNodeFn(t *testing.T) {
	tests := []struct {
		name string
		gc   *genContext
		n    ast.Node
		want string
	}{
		{
			name: "using ref takes precedence",
			gc:   &genContext{},
			n:    ast.Node{Name: "myNode", Using: &ast.Ref{Name: "otherFn"}},
			want: "otherFn",
		},
		{
			name: "using ref with alias",
			gc:   &genContext{},
			n:    ast.Node{Name: "myNode", Using: &ast.Ref{Alias: "pkg", Name: "otherFn"}},
			want: "pkg.otherFn",
		},
		{
			name: "checking wrapper found",
			gc:   &genContext{CheckWrappers: map[string]string{"caseA": "check1"}},
			n:    ast.Node{Name: "myNode", Checking: "caseA"},
			want: "check1",
		},
		{
			name: "checking set but wrapper missing falls back to goIdent",
			gc:   &genContext{CheckWrappers: map[string]string{}},
			n:    ast.Node{Name: "my node", Checking: "caseB"},
			want: "myNode",
		},
		{
			name: "no using no checking falls back to goIdent",
			gc:   &genContext{},
			n:    ast.Node{Name: "plain node"},
			want: "plainNode",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nodeFn(tt.gc, tt.n)
			if got != tt.want {
				t.Errorf("nodeFn() = %q, want %q", got, tt.want)
			}
		})
	}
}
