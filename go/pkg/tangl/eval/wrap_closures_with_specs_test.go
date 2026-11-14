//ff:func feature=tangl type=engine control=sequence
//ff:what TestWrapClosuresWithSpecs — owned and supplied Specs merge into the scope; conflicts return the error as evidence
package eval

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestWrapClosuresWithSpecs(t *testing.T) {
	var seen toulmin.Specs
	compiled := map[string]evalFunc{"r": func(sc *scope) (bool, any) { seen = sc.specs; return true, "ok" }}
	derived := map[string]toulmin.Specs{"r": {&derivedSpec{name: "owned"}}}
	fns := wrapClosuresWithSpecs([]ast.InlineRule{{Name: "r"}}, compiled, derived)

	got, ev := fns["r"](toulmin.NewContext(), toulmin.Specs{&testSpec{Name: "supplied"}})
	if !got || ev != "ok" || len(seen) != 2 || seen[0].SpecName() != "owned" || seen[1].SpecName() != "supplied" {
		t.Fatalf("merged call = %v,%v specs=%v", got, ev, seen)
	}

	got, ev = fns["r"](toulmin.NewContext(), toulmin.Specs{&testSpec{Name: "owned"}})
	err, isErr := ev.(error)
	if got || !isErr || !strings.Contains(err.Error(), `conflicting Spec "owned"`) {
		t.Fatalf("conflict = %v,%v, want false with the merge error", got, ev)
	}
}
