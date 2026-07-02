//ff:func feature=tangl type=analyzer control=sequence
//ff:what TestClosureBranches — covers Closure's revisit-skip, unknown-case, and nested error branches
package effects

import "testing"

// TestClosureBranches exercises Closure's edge-case branches as independent
// subtests: unknown case references (top-level and nested), self-cycle
// skipping, shared-case dedup across a diamond, and the UndoExec entry kind.
// Each subtest body lives in its own helper file (see closure_*_test.go) to
// keep this dispatcher small.
func TestClosureBranches(t *testing.T) {
	t.Run("UnknownCaseAtTop", closureUnknownCaseAtTop)
	t.Run("UnknownCaseNested", closureUnknownCaseNested)
	t.Run("SelfCycleSkipped", closureSelfCycleSkipped)
	t.Run("DiamondSharedCaseVisitedOnce", closureDiamondSharedCaseVisitedOnce)
	t.Run("UndoExecEntry", closureUndoExecEntry)
}
