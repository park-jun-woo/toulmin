//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderClosureNot — tests not-node emission via tanglNot and an inner error
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderClosureNot(t *testing.T) {
	got, err := renderClosureNot(ast.Not{Term: ast.Compare{Field: "x", Op: "is empty"}}, nil, 0)
	if err != nil || !strings.Contains(got, "return tanglNot(") || !strings.Contains(got, "tanglCompareVal(") {
		t.Errorf("not closure: %q (%v)", got, err)
	}
	if _, err := renderClosureNot(ast.Not{Term: nil}, nil, 0); err == nil {
		t.Error("want inner render error")
	}
}
