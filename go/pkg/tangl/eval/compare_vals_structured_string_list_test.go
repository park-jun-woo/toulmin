//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompareValsStructuredStringList — "is in" against a structured string list is exact membership
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestCompareValsStructuredStringList(t *testing.T) {
	list := toulmin.NewStringListSpec("domains", []string{"a,b", " c "})
	if !compareVals("a,b", true, "is in", list) || !compareVals(" c ", true, "is in", list) || compareVals("c", true, "is in", list) {
		t.Fatal("structured string-list membership is not exact")
	}
}
