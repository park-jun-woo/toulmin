//ff:type feature=tangl type=codegen
//ff:what genFlags — which optional imports/shared helpers the generated file needs
package gen

// genFlags gates which optional imports and shared helper functions the
// generated file needs, based on a single scan of the source Document. The
// Needs* rule-helper flags each gate one emitted tangl* runtime helper (and
// the standard imports it requires), so the generated file never contains an
// unused function or import.
type genFlags struct {
	NeedsTangl      bool
	NeedsTime       bool
	NeedsCaseHelper bool
	NeedsCompareVal bool
	NeedsCtxVal     bool
	NeedsSpecVal    bool
	NeedsWalk       bool
	NeedsAnd        bool
	NeedsOr         bool
	NeedsNot        bool
	NeedsAny        bool
	NeedsHolds      bool
	NeedsItoa       bool
}
