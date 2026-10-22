//ff:func feature=tangl type=codegen control=sequence
//ff:what writeRuleHelpers — writes every flag-gated tangl* rule evaluation helper
package gen

import "strings"

// writeRuleHelpers writes each shared tangl* rule evaluation helper the
// document actually needs, gated by the collectRuleFlags scan so the
// generated file never contains an unused function: value resolution
// (tanglCompareVal/tanglCtxVal/tanglSpecVal/tanglWalk), the short-circuit
// combinators (tanglAnd/tanglOr/tanglNot/tanglAny), and the predicate-call
// passthrough (tanglHolds).
func writeRuleHelpers(w *strings.Builder, f genFlags) {
	if f.NeedsCompareVal {
		writeCompareValHelper(w)
	}
	if f.NeedsCtxVal {
		writeCtxValHelper(w)
	}
	if f.NeedsSpecVal {
		writeSpecValHelper(w)
	}
	if f.NeedsWalk {
		writeWalkHelper(w)
	}
	if f.NeedsAnd {
		writeAndHelper(w)
	}
	if f.NeedsOr {
		writeOrHelper(w)
	}
	if f.NeedsNot {
		writeNotHelper(w)
	}
	if f.NeedsAny {
		writeAnyHelper(w)
	}
	if f.NeedsHolds {
		writeHoldsHelper(w)
	}
}
