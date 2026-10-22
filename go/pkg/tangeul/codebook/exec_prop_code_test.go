//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestExecPropCode — ExecKind↔PropCode golden (do/undo/run), roundtrip, unknown errors
package codebook

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestExecPropCode(t *testing.T) {
	cases := []struct {
		kind ast.ExecKind
		prop toulmin.PropCode
	}{
		{ast.DoExec, toulmin.PropDo},
		{ast.UndoExec, toulmin.PropUndo},
		{ast.RunExec, toulmin.PropRun},
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
	// A non-execution toulmin.PropCode (e.g. attacks) is not an ExecKind.
	if _, err := ExecKindFromProp(toulmin.PropAttacks); err == nil {
		t.Error("ExecKindFromProp(toulmin.PropAttacks): expected error")
	}
}
