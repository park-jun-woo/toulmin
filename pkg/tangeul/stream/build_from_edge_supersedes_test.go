//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildFromEdge_Supersedes — buildFromEdge routes PropSupersedes into a case, or errors outside one
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildFromEdge_Supersedes(t *testing.T) {
	t.Run("supersedes outside any case", func(t *testing.T) {
		if _, err := buildFromEdge(newDec(), []packet.EdgePacket{{Prop: codebook.PropSupersedes}}); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("supersedes inside a case", func(t *testing.T) {
		d := newDec()
		d.curCase = 0
		d.ast.Cases = []ast.Case{{}}
		d.nodeName[1] = "superseder"
		d.nodeName[2] = "superseded"
		ref, err := buildFromEdge(d, []packet.EdgePacket{{Prop: codebook.PropSupersedes, Src: 1, DstTID: 2}})
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
		if ref.Sec != codebook.SectionCases || len(ref.Path) != 3 || ref.Path[1] != 4 {
			t.Errorf("StmtRef = %+v, want SectionCases path [0 4 0]", ref)
		}
	})
}
