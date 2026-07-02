//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderConstDef — tests renderConstDef for numeric/string literal and specRef-present/absent branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderConstDef(t *testing.T) {
	t.Run("numeric literal without specRef", func(t *testing.T) {
		var w strings.Builder
		d := ast.Definition{Name: "threshold", Value: "650"}
		info := renderConstDef(&w, d)
		if !strings.Contains(w.String(), "const threshold = 650\n") {
			t.Errorf("unexpected output: %q", w.String())
		}
		if info.Const != "threshold" || info.Spec != "" || info.Kind != ast.ConstDef {
			t.Errorf("unexpected info: %+v", info)
		}
	})

	t.Run("string literal without specRef", func(t *testing.T) {
		var w strings.Builder
		d := ast.Definition{Name: "label", Value: "gold"}
		info := renderConstDef(&w, d)
		if !strings.Contains(w.String(), `const label = "gold"`) {
			t.Errorf("unexpected output: %q", w.String())
		}
		if info.Spec != "" {
			t.Errorf("expected empty spec, got %q", info.Spec)
		}
	})

	t.Run("numeric literal with specRef", func(t *testing.T) {
		var w strings.Builder
		d := ast.Definition{Name: "threshold", Value: "650", SpecRef: &ast.Ref{Name: "Threshold"}}
		info := renderConstDef(&w, d)
		out := w.String()
		if !strings.Contains(out, "const threshold = 650\n") {
			t.Errorf("missing const decl: %q", out)
		}
		if !strings.Contains(out, "var thresholdSpec = Threshold(threshold)\n") {
			t.Errorf("missing spec var decl: %q", out)
		}
		if info.Spec != "thresholdSpec" {
			t.Errorf("expected spec ident thresholdSpec, got %q", info.Spec)
		}
	})

	t.Run("string literal with specRef", func(t *testing.T) {
		var w strings.Builder
		d := ast.Definition{Name: "label", Value: "gold", SpecRef: &ast.Ref{Alias: "pkg", Name: "Threshold"}}
		info := renderConstDef(&w, d)
		out := w.String()
		if !strings.Contains(out, `const label = "gold"`) {
			t.Errorf("missing const decl: %q", out)
		}
		if !strings.Contains(out, "var labelSpec = pkg.Threshold(label)\n") {
			t.Errorf("missing spec var decl: %q", out)
		}
		if info.Spec != "labelSpec" {
			t.Errorf("expected spec ident labelSpec, got %q", info.Spec)
		}
	})
}
