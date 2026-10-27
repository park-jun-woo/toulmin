//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestExecKindFromProp — do/undo/run map to exec kinds; any other relation fails
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestExecKindFromProp(t *testing.T) {
	want := map[toulmin.PropCode]ast.ExecKind{toulmin.PropDo: ast.DoExec, toulmin.PropUndo: ast.UndoExec, toulmin.PropRun: ast.RunExec}
	for p, k := range want {
		if got, err := ExecKindFromProp(p); err != nil || got != k {
			t.Fatalf("0x%04X: %v %v", uint16(p), got, err)
		}
	}
	if _, err := ExecKindFromProp(toulmin.PropAttacks); err == nil {
		t.Fatal("attacks: want error")
	}
}
