//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestWriteRuleHelpers — tests flag-gated emission of every tangl* rule helper
package gen

import (
	"strings"
	"testing"
)

func TestWriteRuleHelpers(t *testing.T) {
	var w strings.Builder
	writeRuleHelpers(&w, genFlags{})
	if w.String() != "" {
		t.Errorf("no flags should emit nothing, got %q", w.String())
	}
	w.Reset()
	writeRuleHelpers(&w, genFlags{
		NeedsCompareVal: true, NeedsCtxVal: true, NeedsSpecVal: true, NeedsWalk: true,
		NeedsAnd: true, NeedsOr: true, NeedsNot: true, NeedsAny: true, NeedsHolds: true,
	})
	got := w.String()
	for _, fn := range []string{"tanglCompareVal", "tanglCtxVal", "tanglSpecVal", "tanglWalk", "tanglAnd", "tanglOr", "tanglNot", "tanglAny", "tanglHolds"} {
		if !strings.Contains(got, "func "+fn+"(") {
			t.Errorf("expected %s helper to be emitted", fn)
		}
	}
}
