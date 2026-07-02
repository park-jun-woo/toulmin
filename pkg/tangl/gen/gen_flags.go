//ff:type feature=tangl type=codegen
//ff:what genFlags — which optional imports/shared helpers the generated file needs
package gen

// genFlags gates which optional imports and shared helper functions the
// generated file needs, based on a single scan of the source Document.
type genFlags struct {
	NeedsTangl      bool
	NeedsTime       bool
	NeedsCompare    bool
	NeedsCaseHelper bool
}
