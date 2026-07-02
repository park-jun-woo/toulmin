//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestHasDoOrUndo — tests hasDoOrUndo for empty input, do-exec match, undo-exec match, and no-match cases
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestHasDoOrUndo(t *testing.T) {
	tests := []struct {
		name  string
		execs []ast.Exec
		want  bool
	}{
		{"empty slice", nil, false},
		{"only run exec", []ast.Exec{{Kind: ast.RunExec}}, false},
		{"do exec present", []ast.Exec{{Kind: ast.RunExec}, {Kind: ast.DoExec}}, true},
		{"undo exec present", []ast.Exec{{Kind: ast.UndoExec}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasDoOrUndo(tt.execs)
			if got != tt.want {
				t.Errorf("hasDoOrUndo(%+v) = %v, want %v", tt.execs, got, tt.want)
			}
		})
	}
}
