//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderInternalEvery — tests renderInternalEvery for valid/invalid interval and with/without until branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderInternalEvery(t *testing.T) {
	t.Run("valid interval without until", func(t *testing.T) {
		var w strings.Builder
		in := ast.Internal{Kind: ast.EveryTick, Interval: "30s", Runs: []string{"caseA"}}
		renderInternalEvery(&w, in, 0)
		out := w.String()
		if !strings.Contains(out, `interval: "30s"`) {
			t.Errorf("expected valid interval comment: %q", out)
		}
		if strings.Contains(out, "defaulting to 24h") {
			t.Errorf("unexpected default fallback comment: %q", out)
		}
		if strings.Contains(out, "Graph.Evaluate(ctx)") {
			t.Errorf("unexpected until guard when Until is empty: %q", out)
		}
	})

	t.Run("invalid interval falls back to 24h", func(t *testing.T) {
		var w strings.Builder
		in := ast.Internal{Kind: ast.EveryTick, Interval: "not-a-duration", Checks: []string{"caseB"}}
		renderInternalEvery(&w, in, 1)
		out := w.String()
		if !strings.Contains(out, "is not a supported duration/clock schedule, defaulting to 24h") {
			t.Errorf("expected fallback comment: %q", out)
		}
	})

	t.Run("with until emits termination guard", func(t *testing.T) {
		var w strings.Builder
		in := ast.Internal{Kind: ast.EveryTick, Interval: "1m", Runs: []string{"caseA"}, Until: "doneCase"}
		renderInternalEvery(&w, in, 2)
		out := w.String()
		if !strings.Contains(out, "doneCaseGraph.Evaluate(ctx)") {
			t.Errorf("expected until guard evaluate call: %q", out)
		}
		if !strings.Contains(out, "if err == nil && tanglCaseActive(results) {") {
			t.Errorf("expected until guard condition: %q", out)
		}
	})
}
