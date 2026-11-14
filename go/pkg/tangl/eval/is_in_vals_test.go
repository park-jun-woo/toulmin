//ff:func feature=tangl type=engine control=sequence
//ff:what TestIsInVals — missing lhs never matches; structured lists match exactly; other rhs use comma-list membership
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestIsInVals(t *testing.T) {
	list := toulmin.NewStringListSpec("d", []string{"a,b", "7"})
	if isInVals("a,b", false, list) {
		t.Fatal("missing lhs must not match")
	}
	if !isInVals("a,b", true, list) || !isInVals(7, true, list) || isInVals("a", true, list) {
		t.Fatal("structured list membership must be exact over the formatted lhs")
	}
	if !isInVals("b", true, "a, b, c") || isInVals("z", true, "a, b, c") {
		t.Fatal("plain rhs must use comma-list membership")
	}
}
