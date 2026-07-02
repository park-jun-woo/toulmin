//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildInternals — tests buildInternals for EveryTick and OnEvent branches, plus empty list
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildInternals(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		buildInternals(&w, &ast.Document{})
		if w.String() != "" {
			t.Errorf("expected no output, got %q", w.String())
		}
	})

	t.Run("every and on", func(t *testing.T) {
		doc := &ast.Document{
			Internals: []ast.Internal{
				{Kind: ast.EveryTick, Interval: "30s"},
				{Kind: ast.OnEvent, Event: "start"},
			},
		}
		var w strings.Builder
		buildInternals(&w, doc)
		out := w.String()
		if !strings.Contains(out, "ticker := time.NewTicker") {
			t.Errorf("expected every-tick runner rendered, got:\n%s", out)
		}
		if !strings.Contains(out, `handles the "start" event`) {
			t.Errorf("expected on-event handler rendered, got:\n%s", out)
		}
	})
}
