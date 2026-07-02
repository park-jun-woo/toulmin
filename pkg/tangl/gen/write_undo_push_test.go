//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteUndoPush — tests writeUndoPush's generated statement
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWriteUndoPush(t *testing.T) {
	var w strings.Builder
	writeUndoPush(&w, ast.Exec{Kind: ast.UndoExec, Func: &ast.Ref{Name: "cleanup"}, Node: "node1"})
	got := w.String()
	if !strings.Contains(got, "tangl.PushCompensation(t.Ctx(), func() error {") {
		t.Errorf("expected push compensation call, got %q", got)
	}
	if !strings.Contains(got, "return cleanup(t.Ctx())") {
		t.Errorf("expected cleanup call, got %q", got)
	}
}
