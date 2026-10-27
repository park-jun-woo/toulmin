//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildFromEdge_Supersedes — buildFromEdge routes PropSupersedes into a case, or errors outside one
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestBuildFromEdge_Supersedes(t *testing.T) {
	t.Run("supersedes outside any case", func(t *testing.T) {
		if _, err := buildFromEdge(newDec(testWire), []toulmin.EdgePacket{{Prop: toulmin.PropSupersedes}}); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("supersedes inside a case", func(t *testing.T) {
		d := newDec(testWire)
		d.curCase = 0
		d.ast.Cases = []ast.Case{{}}
		d.nodeName[1] = "superseder"
		d.nodeName[2] = "superseded"
		ref, err := buildFromEdge(d, []toulmin.EdgePacket{{Prop: toulmin.PropSupersedes, Src: 1, DstTID: 2}})
		if err != nil {
			t.Fatalf("buildFromEdge: %v", err)
		}
		if len(d.ast.Cases[0].Supersedes) != 1 {
			t.Fatalf("Supersedes len = %d, want 1", len(d.ast.Cases[0].Supersedes))
		}
		s := d.ast.Cases[0].Supersedes[0]
		if s.Superseder != "superseder" || s.Superseded != "superseded" {
			t.Errorf("Supersede = %+v, want {superseder superseded}", s)
		}
		if ref.Sec != toulmin.SectionCases || len(ref.Path) != 3 || ref.Path[1] != 4 {
			t.Errorf("StmtRef = %+v, want SectionCases path [0 4 0]", ref)
		}
	})
}
