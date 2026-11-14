//ff:func feature=tangl type=codegen control=sequence
//ff:what TestPathValCode — tests binder dispatch and ctx-rooted path value code
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestPathValCode(t *testing.T) {
	env := map[string]binderInfo{"c": {Elem: "e0", Idx: "i0"}}

	vc := pathValCode(ast.Path{Steps: []string{"c"}}, env, "l", true)
	if vc.V != "e0" {
		t.Errorf("binder-headed path should dispatch to the binder code: %+v", vc)
	}
	vc = pathValCode(ast.Path{Steps: []string{"author", "trustScore"}}, env, "l", true)
	if !strings.Contains(vc.Setup, `lv, lok := tanglCtxVal(ctx, "author", "trustScore")`) ||
		vc.V != "lv" || vc.OK != "lok" || vc.Actual != "lv" {
		t.Errorf("ctx-rooted with ok: %+v", vc)
	}
	vc = pathValCode(ast.Path{Steps: []string{"author"}}, env, "r", false)
	if !strings.Contains(vc.Setup, "rv, _ := tanglCtxVal(") || vc.OK != "" {
		t.Errorf("ctx-rooted without ok: %+v", vc)
	}
	vc = pathValCode(ast.Path{}, env, "s", true)
	if !strings.Contains(vc.Setup, "tanglCtxVal(ctx, )") {
		t.Errorf("empty path should stay ctx-rooted: %+v", vc)
	}
}
