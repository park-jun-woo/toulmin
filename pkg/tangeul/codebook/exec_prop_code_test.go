//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestExecPropCode — ExecKind↔PropCode golden (do/undo/run), roundtrip, unknown errors
package codebook

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestExecPropCode(t *testing.T) {
	cases := []struct {
		kind ast.ExecKind
		prop PropCode
	}{
		{ast.DoExec, PropDo},
		{ast.UndoExec, PropUndo},
		{ast.RunExec, PropRun},
	}
	for _, c := range cases {
		prop, err := ExecPropCode(c.kind)
		if err != nil || prop != c.prop {
			t.Fatalf("ExecPropCode(%v): got 0x%04X err=%v, want 0x%04X", c.kind, uint16(prop), err, uint16(c.prop))
		}
		kind, err := ExecKindFromProp(c.prop)
		if err != nil || kind != c.kind {
			t.Fatalf("ExecKindFromProp(0x%04X): got %v err=%v, want %v", uint16(c.prop), kind, err, c.kind)
		}
	}
	if _, err := ExecPropCode(ast.ExecKind(99)); err == nil {
		t.Error("ExecPropCode(unknown): expected error")
	}
	// A non-execution PropCode (e.g. attacks) is not an ExecKind.
	if _, err := ExecKindFromProp(PropAttacks); err == nil {
		t.Error("ExecKindFromProp(PropAttacks): expected error")
	}
}
