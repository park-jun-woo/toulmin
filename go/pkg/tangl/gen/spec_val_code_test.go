//ff:func feature=tangl type=codegen control=sequence
//ff:what TestSpecValCode — tests spec-ref value code with and without the ok flag
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestSpecValCode(t *testing.T) {
	sr := ast.SpecRef{Steps: []string{"discount", "minPrice"}}

	vc := specValCode(sr, "l", true)
	if !strings.Contains(vc.Setup, `lv, lok := tanglSpecVal(specs, "discount", "minPrice")`) ||
		vc.V != "lv" || vc.OK != "lok" || vc.Actual != "lv" {
		t.Errorf("with ok: %+v", vc)
	}
	vc = specValCode(sr, "r", false)
	if !strings.Contains(vc.Setup, "rv, _ := tanglSpecVal(") || vc.OK != "" {
		t.Errorf("without ok: %+v", vc)
	}
}
