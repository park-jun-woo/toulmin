//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildRules — tests buildRules for empty, success, and error-propagation branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildRules(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		if err := buildRules(&w, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.String() != "" {
			t.Errorf("expected no output, got %q", w.String())
		}
	})

	t.Run("success", func(t *testing.T) {
		rules := []ast.InlineRule{
			{Name: "my rule", Cond: ast.Compare{Field: "x", Op: "==", Value: "1"}},
		}
		var w strings.Builder
		if err := buildRules(&w, rules); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(w.String(), "func myRule") {
			t.Errorf("expected rule function rendered, got:\n%s", w.String())
		}
	})

	t.Run("error propagates", func(t *testing.T) {
		rules := []ast.InlineRule{
			{Name: "bad rule", Cond: nil},
		}
		var w strings.Builder
		if err := buildRules(&w, rules); err == nil {
			t.Fatal("expected error for unsupported expression node")
		}
	})
}
