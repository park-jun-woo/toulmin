//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderEndpointRun — tests renderEndpointRun for run-only (error return) and run-with-checks (results return) branches
package gen

import (
	"strings"
	"testing"
)

func TestRenderEndpointRun(t *testing.T) {
	t.Run("run only returns error", func(t *testing.T) {
		var w strings.Builder
		renderEndpointRun(&w, "RunOrder", []string{"amount"}, []string{"caseA"}, nil)
		out := w.String()
		if !strings.Contains(out, "func RunOrder(ctx toulmin.Context) error {") {
			t.Errorf("expected error-return signature: %q", out)
		}
		if !strings.Contains(out, "tangl.InitCompensation(ctx)") {
			t.Errorf("missing init compensation: %q", out)
		}
		if !strings.Contains(out, "tangl.CommitCompensation(ctx)") {
			t.Errorf("missing commit compensation: %q", out)
		}
	})

	t.Run("run with extra checks returns results", func(t *testing.T) {
		var w strings.Builder
		renderEndpointRun(&w, "RunOrder", nil, []string{"caseA"}, []string{"caseB"})
		out := w.String()
		if !strings.Contains(out, "func RunOrder(ctx toulmin.Context) ([]toulmin.EvalResult, error) {") {
			t.Errorf("expected results-return signature: %q", out)
		}
		if !strings.Contains(out, "caseBGraph.Evaluate(ctx)") {
			t.Errorf("missing extra check evaluate: %q", out)
		}
	})
}
