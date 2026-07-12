//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestWriteCompareValHelper — tests writeCompareValHelper emits the tanglCompareVal source
package gen

import (
	"strings"
	"testing"
)

func TestWriteCompareValHelper(t *testing.T) {
	var w strings.Builder
	writeCompareValHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglCompareVal(lhs any, ok bool, op string, rhs any) bool {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	for _, op := range []string{`"is empty"`, `"is not empty"`, `"equals"`, `"is in"`, `"is greater than", "is less than", "is at most", "is at least"`} {
		if !strings.Contains(got, "case "+op+":") {
			t.Errorf("expected operator case %s, got %q", op, got)
		}
	}
}
