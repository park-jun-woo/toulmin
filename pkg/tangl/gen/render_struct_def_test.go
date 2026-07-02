//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderStructDef — tests renderStructDef for no-fields and multi-field branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderStructDef(t *testing.T) {
	t.Run("no fields", func(t *testing.T) {
		var w strings.Builder
		d := ast.Definition{Name: "empty struct"}
		info := renderStructDef(&w, d)
		want := "type EmptyStruct struct {\n}\n\n"
		if w.String() != want {
			t.Errorf("got %q, want %q", w.String(), want)
		}
		if info.Kind != ast.StructDef {
			t.Errorf("expected StructDef kind, got %v", info.Kind)
		}
	})

	t.Run("multiple fields", func(t *testing.T) {
		var w strings.Builder
		d := ast.Definition{
			Name: "order",
			Fields: []ast.Field{
				{Name: "amount", Type: "Currency"},
				{Name: "label", Type: "Text"},
			},
		}
		renderStructDef(&w, d)
		out := w.String()
		if !strings.Contains(out, "type Order struct {") {
			t.Errorf("missing struct header: %q", out)
		}
		if !strings.Contains(out, "\tAmount float64\n") {
			t.Errorf("missing amount field: %q", out)
		}
		if !strings.Contains(out, "\tLabel string\n") {
			t.Errorf("missing label field: %q", out)
		}
	})
}
