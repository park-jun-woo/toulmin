//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBinderPathValCode — tests binder-only and walked binder-path value code
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestBinderPathValCode(t *testing.T) {
	b := binderInfo{Elem: "e0", Idx: "i0"}

	vc := binderPathValCode(ast.Path{Steps: []string{"c"}}, b, "l", true)
	if vc.Setup != "" || vc.V != "e0" || vc.OK != "true" || vc.Actual != "e0" {
		t.Errorf("binder-only with ok: %+v", vc)
	}
	vc = binderPathValCode(ast.Path{Steps: []string{"c"}}, b, "l", false)
	if vc.OK != "" || vc.V != "e0" {
		t.Errorf("binder-only without ok: %+v", vc)
	}
	vc = binderPathValCode(ast.Path{Steps: []string{"c", "minPrice"}}, b, "l", true)
	if !strings.Contains(vc.Setup, `lv, lok := tanglWalk(e0, "minPrice")`) || vc.V != "lv" || vc.OK != "lok" || vc.Actual != "lv" {
		t.Errorf("walked with ok: %+v", vc)
	}
	vc = binderPathValCode(ast.Path{Steps: []string{"c", "minPrice"}}, b, "l", false)
	if !strings.Contains(vc.Setup, "lv, _ := tanglWalk(") || vc.OK != "" {
		t.Errorf("walked without ok: %+v", vc)
	}
}
