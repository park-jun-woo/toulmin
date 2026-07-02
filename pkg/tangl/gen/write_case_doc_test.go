//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteCaseDoc — tests writeCaseDoc for empty Requires and populated Requires branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWriteCaseDoc(t *testing.T) {
	t.Run("no requires", func(t *testing.T) {
		var w strings.Builder
		writeCaseDoc(&w, ast.Case{Name: "sample"})
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("with requires", func(t *testing.T) {
		var w strings.Builder
		c := ast.Case{
			Name: "sample",
			Requires: []ast.Require{
				{Field: "amount"},
				{Field: "status"},
			},
		}
		writeCaseDoc(&w, c)
		got := w.String()
		want := "// requires ctx fields: amount, status\n"
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}
