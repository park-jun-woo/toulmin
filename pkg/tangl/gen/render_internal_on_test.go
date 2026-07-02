//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderInternalOn — tests renderInternalOn for runs/no-runs and checks/no-checks branch combinations
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderInternalOn(t *testing.T) {
	t.Run("no runs no checks returns error", func(t *testing.T) {
		var w strings.Builder
		in := ast.Internal{Kind: ast.OnEvent, Event: "sensor.press"}
		renderInternalOn(&w, in, 0)
		out := w.String()
		if !strings.Contains(out, ") error {") {
			t.Errorf("expected error-return signature: %q", out)
		}
		if strings.Contains(out, "InitCompensation") {
			t.Errorf("unexpected compensation block: %q", out)
		}
	})

	t.Run("with runs no checks returns error", func(t *testing.T) {
		var w strings.Builder
		in := ast.Internal{Kind: ast.OnEvent, Event: "sensor.press", Runs: []string{"caseA"}}
		renderInternalOn(&w, in, 1)
		out := w.String()
		if !strings.Contains(out, ") error {") {
			t.Errorf("expected error-return signature: %q", out)
		}
		if !strings.Contains(out, "tangl.InitCompensation(ctx)") {
			t.Errorf("expected compensation block: %q", out)
		}
		if !strings.Contains(out, "tangl.CommitCompensation(ctx)") {
			t.Errorf("expected commit compensation: %q", out)
		}
	})

	t.Run("no runs with checks returns results", func(t *testing.T) {
		var w strings.Builder
		in := ast.Internal{Kind: ast.OnEvent, Event: "sensor.press", Checks: []string{"caseB"}}
		renderInternalOn(&w, in, 2)
		out := w.String()
		if !strings.Contains(out, ") ([]toulmin.EvalResult, error) {") {
			t.Errorf("expected results-return signature: %q", out)
		}
		if strings.Contains(out, "InitCompensation") {
			t.Errorf("unexpected compensation block: %q", out)
		}
		if !strings.Contains(out, "caseBGraph.Evaluate(ctx)") {
			t.Errorf("expected check evaluate call: %q", out)
		}
	})

	t.Run("with runs and checks returns results", func(t *testing.T) {
		var w strings.Builder
		in := ast.Internal{Kind: ast.OnEvent, Event: "sensor.press", Runs: []string{"caseA"}, Checks: []string{"caseB"}}
		renderInternalOn(&w, in, 3)
		out := w.String()
		if !strings.Contains(out, ") ([]toulmin.EvalResult, error) {") {
			t.Errorf("expected results-return signature: %q", out)
		}
		if !strings.Contains(out, "tangl.InitCompensation(ctx)") {
			t.Errorf("expected compensation block: %q", out)
		}
		if !strings.Contains(out, "caseBGraph.Evaluate(ctx)") {
			t.Errorf("expected check evaluate call: %q", out)
		}
	})
}
