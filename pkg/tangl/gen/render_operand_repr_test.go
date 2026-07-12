//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderOperandRepr — tests Where-notation rendering for every operand form
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderOperandRepr(t *testing.T) {
	constEnv := map[string]binderInfo{"c": {Elem: "e0", Idx: "i0", Src: reprCode{Expr: `"coupons"`, Const: true, Val: "coupons"}}}
	exprEnv := map[string]binderInfo{"c": {Elem: "e1", Idx: "i1", Src: reprCode{Expr: `srcExpr`}}}

	if r := renderOperandRepr(ast.Path{}, nil); !r.Const || r.Val != "" || r.Expr != `""` {
		t.Errorf("empty path: %+v", r)
	}
	if r := renderOperandRepr(ast.Path{Steps: []string{"a", "b"}}, nil); !r.Const || r.Val != "a.b" {
		t.Errorf("plain path: %+v", r)
	}
	if r := renderOperandRepr(ast.Path{Steps: []string{"c"}}, constEnv); r.Const || r.Expr != `"coupons["+strconv.Itoa(i0)+"]"` {
		t.Errorf("binder-only path over const source: %+v", r)
	}
	if r := renderOperandRepr(ast.Path{Steps: []string{"c", "minPrice"}}, constEnv); r.Const || !strings.Contains(r.Expr, `"].minPrice"`) {
		t.Errorf("binder path with rest: %+v", r)
	}
	if r := renderOperandRepr(ast.Path{Steps: []string{"c"}}, exprEnv); r.Const || !strings.HasPrefix(r.Expr, `srcExpr+"["+strconv.Itoa(i1)+`) {
		t.Errorf("binder path over expr source: %+v", r)
	}
	if r := renderOperandRepr(ast.SpecRef{Steps: []string{"d", "max"}}, nil); !r.Const || r.Val != "d.max" {
		t.Errorf("spec ref: %+v", r)
	}
	if r := renderOperandRepr(ast.Call{Fn: ast.Ref{Name: "score"}}, nil); !r.Const || r.Val != "score" {
		t.Errorf("bare call: %+v", r)
	}
	if r := renderOperandRepr(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, nil); !r.Const || r.Val != "nlp.score" {
		t.Errorf("aliased call: %+v", r)
	}
	if r := renderOperandRepr(struct{ ast.Path }{}, nil); !r.Const || r.Val != "" || r.Expr != `""` {
		t.Errorf("unknown operand: %+v", r)
	}
}
