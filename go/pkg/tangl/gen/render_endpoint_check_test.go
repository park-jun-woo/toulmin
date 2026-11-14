//ff:func feature=tangl type=codegen control=sequence dimension=1
//ff:what TestRenderEndpointCheck — tests renderEndpointCheck emits the check-endpoint function with guard and checks
package gen

import (
	"strings"
	"testing"
)

func TestRenderEndpointCheck(t *testing.T) {
	t.Run("with fields and checks", func(t *testing.T) {
		var w strings.Builder
		renderEndpointCheck(&w, "CheckOrder", []string{"amount"}, []string{"caseA"})
		out := w.String()
		if !strings.Contains(out, "func CheckOrder(ctx toulmin.Context) ([]toulmin.EvalResult, error) {") {
			t.Errorf("missing function signature: %q", out)
		}
		if !strings.Contains(out, "tangl.Required(ctx,") {
			t.Errorf("missing required guard: %q", out)
		}
		if !strings.Contains(out, "caseAGraph.Evaluate(ctx)") {
			t.Errorf("missing check evaluate call: %q", out)
		}
	})

	t.Run("without fields and checks", func(t *testing.T) {
		var w strings.Builder
		renderEndpointCheck(&w, "CheckPlain", nil, nil)
		out := w.String()
		if !strings.Contains(out, "func CheckPlain(ctx toulmin.Context) ([]toulmin.EvalResult, error) {") {
			t.Errorf("missing function signature: %q", out)
		}
		if strings.Contains(out, "tangl.Required") {
			t.Errorf("unexpected required guard: %q", out)
		}
	})
}
