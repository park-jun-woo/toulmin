//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderEndpoint — tests renderEndpoint for runs-present, checks-only, and no-runs-no-checks branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderEndpoint(t *testing.T) {
	t.Run("runs present takes precedence", func(t *testing.T) {
		var w strings.Builder
		ep := ast.Endpoint{Name: "myEndpoint", Runs: []string{"caseA"}, Checks: []string{"caseB"}}
		renderEndpoint(&w, ep)
		out := w.String()
		if out == "" {
			t.Fatal("expected output for runs-present endpoint")
		}
		if !strings.Contains(out, "MyEndpoint") {
			t.Errorf("expected function name derived from endpoint name: %q", out)
		}
	})

	t.Run("checks only, no runs", func(t *testing.T) {
		var w strings.Builder
		ep := ast.Endpoint{Name: "myCheckEndpoint", Checks: []string{"caseB"}}
		renderEndpoint(&w, ep)
		out := w.String()
		if out == "" {
			t.Fatal("expected output for checks-only endpoint")
		}
		if !strings.Contains(out, "MyCheckEndpoint") {
			t.Errorf("expected function name derived from endpoint name: %q", out)
		}
	})

	t.Run("no runs no checks writes nothing", func(t *testing.T) {
		var w strings.Builder
		ep := ast.Endpoint{Name: "emptyEndpoint"}
		renderEndpoint(&w, ep)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})
}
