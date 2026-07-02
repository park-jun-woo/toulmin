//ff:func feature=tangl type=codegen control=sequence dimension=1
//ff:what TestRenderCheckingWrapper — tests renderCheckingWrapper emits the checking wrapper function body
package gen

import (
	"strings"
	"testing"
)

func TestRenderCheckingWrapper(t *testing.T) {
	var w strings.Builder
	renderCheckingWrapper(&w, "checkOrderWrapper", "targetCase")
	out := w.String()

	want := "func checkOrderWrapper(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {\n" +
		"\tresults, err := targetCaseGraph.Evaluate(ctx)\n" +
		"\tif err != nil {\n" +
		"\t\treturn false, err\n" +
		"\t}\n" +
		"\treturn tanglCaseActive(results), results\n" +
		"}\n\n"

	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}
