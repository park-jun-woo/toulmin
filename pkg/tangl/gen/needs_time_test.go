//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestNeedsTime — tests needsTime for empty internals, non-tick internal, and every-tick internal branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestNeedsTime(t *testing.T) {
	tests := []struct {
		name string
		doc  *ast.Document
		want bool
	}{
		{"no internals", &ast.Document{}, false},
		{"on-event internal only", &ast.Document{Internals: []ast.Internal{{Kind: ast.OnEvent}}}, false},
		{"every-tick internal triggers true", &ast.Document{Internals: []ast.Internal{{Kind: ast.OnEvent}, {Kind: ast.EveryTick}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsTime(tt.doc)
			if got != tt.want {
				t.Errorf("needsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
