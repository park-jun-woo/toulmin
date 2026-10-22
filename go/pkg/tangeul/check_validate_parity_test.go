//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what checkValidateParity — assert Validate and Lint agree (after normalizeLoc) for one fixture
package tangeul

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/validate"
)

// checkValidateParity compares validate.Validate and validate.Lint output between
// the text-parse AST and the reloaded golden .geul AST of the named fixture.
// Both routes bury doc.Path and Line into message prefixes that differ by
// construction (.md+real-line vs .geul+0), so the comparison is over
// normalizeLoc-rewritten strings — semantic content only. For the step-1 fixtures
// (transfer, access) both routes are clean (no errors, no lint), which this pins.
func checkValidateParity(t *testing.T, name string) {
	t.Helper()
	old, loaded := astPair(t, name)
	if oldMsg, newMsg := validateMsg(old), validateMsg(loaded); oldMsg != newMsg {
		t.Fatalf("%s: Validate differs\nold=%q\nnew=%q", name, oldMsg, newMsg)
	}
	oldLint, newLint := validate.Lint(old), validate.Lint(loaded)
	if len(oldLint) != len(newLint) {
		t.Fatalf("%s: Lint count differs: old=%d new=%d", name, len(oldLint), len(newLint))
	}
	for i := range oldLint {
		if normalizeLoc(oldLint[i]) != normalizeLoc(newLint[i]) {
			t.Fatalf("%s: Lint[%d] differs\nold=%q\nnew=%q", name, i, oldLint[i], newLint[i])
		}
	}
}
