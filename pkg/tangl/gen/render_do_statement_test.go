//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderDoStatement — tests renderDoStatement for certainty/once branch combinations
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderDoStatement(t *testing.T) {
	fn := &ast.Ref{Name: "doWork"}

	t.Run("no certainty, no once", func(t *testing.T) {
		var w strings.Builder
		e := ast.Exec{Func: fn}
		renderDoStatement(&w, e, "key1")
		out := w.String()
		if !strings.Contains(out, "if err := doWork(t.Ctx()); err != nil {") {
			t.Errorf("missing plain call: %q", out)
		}
		if strings.Contains(out, "self.Verdict") {
			t.Errorf("unexpected certainty guard: %q", out)
		}
		if strings.Contains(out, "OnceDone") {
			t.Errorf("unexpected once guard: %q", out)
		}
	})

	t.Run("no certainty, with once", func(t *testing.T) {
		var w strings.Builder
		e := ast.Exec{Func: fn, Once: true}
		renderDoStatement(&w, e, "key2")
		out := w.String()
		if !strings.Contains(out, `if !tangl.OnceDone(t.Ctx(), "key2") {`) {
			t.Errorf("missing once guard: %q", out)
		}
		if !strings.Contains(out, `tangl.OnceMark(t.Ctx(), "key2")`) {
			t.Errorf("missing once mark: %q", out)
		}
		if strings.Contains(out, "self.Verdict") {
			t.Errorf("unexpected certainty guard: %q", out)
		}
	})

	t.Run("with certainty, no once", func(t *testing.T) {
		var w strings.Builder
		e := ast.Exec{Func: fn, Certainty: &ast.Certainty{Op: "above", Percent: 75}}
		renderDoStatement(&w, e, "key3")
		out := w.String()
		if !strings.Contains(out, "if self.Verdict > 0.5 {") {
			t.Errorf("missing certainty guard: %q", out)
		}
		if strings.Contains(out, "OnceDone") {
			t.Errorf("unexpected once guard: %q", out)
		}
		// closing brace for certainty guard
		if !strings.HasSuffix(strings.TrimRight(out, "\n"), "}") {
			t.Errorf("expected closing brace at end: %q", out)
		}
	})

	t.Run("with certainty and once", func(t *testing.T) {
		var w strings.Builder
		e := ast.Exec{Func: fn, Once: true, Certainty: &ast.Certainty{Op: "at least", Percent: 50}}
		renderDoStatement(&w, e, "key4")
		out := w.String()
		if !strings.Contains(out, "if self.Verdict >= 0 {") {
			t.Errorf("missing certainty guard: %q", out)
		}
		if !strings.Contains(out, `if !tangl.OnceDone(t.Ctx(), "key4") {`) {
			t.Errorf("missing once guard: %q", out)
		}
		if !strings.Contains(out, `tangl.OnceMark(t.Ctx(), "key4")`) {
			t.Errorf("missing once mark: %q", out)
		}
	})
}
