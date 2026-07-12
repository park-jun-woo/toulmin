//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderOperandVal — tests per-kind valCode dispatch and the unsupported error
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderOperandVal(t *testing.T) {
	vc, err := renderOperandVal(ast.Path{Steps: []string{"x"}}, nil, "l", true, false)
	if err != nil || !strings.Contains(vc.Setup, "tanglCtxVal(") {
		t.Errorf("path dispatch: %+v (%v)", vc, err)
	}
	vc, err = renderOperandVal(ast.SpecRef{Steps: []string{"d"}}, nil, "l", true, false)
	if err != nil || !strings.Contains(vc.Setup, "tanglSpecVal(") {
		t.Errorf("spec-ref dispatch: %+v (%v)", vc, err)
	}
	vc, err = renderOperandVal(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, nil, "l", true, true)
	if err != nil || !strings.Contains(vc.Setup, "nlp.score(ctx, specs)") {
		t.Errorf("call dispatch: %+v (%v)", vc, err)
	}
	if _, err := renderOperandVal(struct{ ast.Path }{}, nil, "l", true, false); err == nil {
		t.Error("want unsupported-operand error")
	}
}
