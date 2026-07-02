//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildDefinitions — tests buildDefinitions for empty, StructDef, and ConstDef branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildDefinitions(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		index := buildDefinitions(&w, nil)
		if len(index) != 0 {
			t.Errorf("expected empty index, got %+v", index)
		}
	})

	t.Run("struct and const", func(t *testing.T) {
		defs := []ast.Definition{
			{
				Name: "widget",
				Kind: ast.StructDef,
				Fields: []ast.Field{
					{Name: "size", Type: "number"},
				},
			},
			{
				Name:  "threshold",
				Kind:  ast.ConstDef,
				Value: "650",
			},
		}
		var w strings.Builder
		index := buildDefinitions(&w, defs)
		if len(index) != 2 {
			t.Fatalf("expected 2 entries, got %+v", index)
		}
		if info, ok := index["widget"]; !ok || info.Kind != ast.StructDef {
			t.Errorf("expected widget StructDef entry, got %+v", index["widget"])
		}
		if info, ok := index["threshold"]; !ok || info.Kind != ast.ConstDef {
			t.Errorf("expected threshold ConstDef entry, got %+v", index["threshold"])
		}
		out := w.String()
		if !strings.Contains(out, "type Widget struct") {
			t.Errorf("expected struct rendered, got:\n%s", out)
		}
		if !strings.Contains(out, "const threshold = 650") {
			t.Errorf("expected const rendered, got:\n%s", out)
		}
	})
}
