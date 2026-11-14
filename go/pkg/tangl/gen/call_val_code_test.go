//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCallValCode — tests value-call code in actual, ok-only, and bare positions
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCallValCode(t *testing.T) {
	c := ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}

	vc := callValCode(c, "l", true, true)
	if !strings.Contains(vc.Setup, "lv, lerr := nlp.score(ctx, specs)") ||
		!strings.Contains(vc.Setup, "lok := lerr == nil") ||
		!strings.Contains(vc.Setup, "lv, lactual = nil, lerr") ||
		vc.V != "lv" || vc.OK != "lok" || vc.Actual != "lactual" {
		t.Errorf("needActual form: %+v", vc)
	}
	vc = callValCode(c, "r", true, false)
	if !strings.Contains(vc.Setup, "rok := rerr == nil") ||
		strings.Contains(vc.Setup, "ractual") ||
		vc.V != "rv" || vc.OK != "rok" || vc.Actual != "" {
		t.Errorf("needOK form: %+v", vc)
	}
	vc = callValCode(c, "s", false, false)
	if !strings.Contains(vc.Setup, "if serr != nil {\nsv = nil\n}") ||
		vc.V != "sv" || vc.OK != "" || vc.Actual != "" {
		t.Errorf("bare form: %+v", vc)
	}
}
