//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteCaseHelper — tests writeCaseHelper emits the tanglCaseActive helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteCaseHelper(t *testing.T) {
	var w strings.Builder
	writeCaseHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglCaseActive(results []toulmin.EvalResult) bool {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, "r.Verdict > 0") {
		t.Errorf("expected verdict check, got %q", got)
	}
}
