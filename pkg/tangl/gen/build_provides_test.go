//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildProvides — tests buildProvides for empty and non-empty Provides lists
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildProvides(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gc := &genContext{Doc: &ast.Document{}}
		var w strings.Builder
		buildProvides(&w, gc)
		if w.String() != "" {
			t.Errorf("expected no output, got %q", w.String())
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		gc := &genContext{
			Doc: &ast.Document{
				Provides: []ast.Endpoint{
					{Name: "endpoint one", Runs: []string{"case one"}},
				},
			},
		}
		var w strings.Builder
		buildProvides(&w, gc)
		out := w.String()
		if !strings.Contains(out, "func EndpointOne") {
			t.Errorf("expected endpoint rendered, got:\n%s", out)
		}
	})
}
