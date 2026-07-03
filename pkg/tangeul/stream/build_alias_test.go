//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestBuildAlias — buildAlias resolves the package path from importedFrom, or leaves it empty without one
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestBuildAlias checks buildAlias appends an ast.See named after the label,
// with its Package resolved from the importedFrom edge's Text when present, and
// left empty when the edge is absent.
func TestBuildAlias(t *testing.T) {
	cases := []struct {
		name    string
		rel     []packet.EdgePacket
		texts   map[uint16]string
		wantPkg string
	}{
		{
			name:    "with importedFrom edge",
			rel:     []packet.EdgePacket{{Prop: codebook.PropImportedFrom, DstTID: 5}},
			texts:   map[uint16]string{5: "bank/pkg"},
			wantPkg: "bank/pkg",
		},
		{
			name:    "without importedFrom edge",
			rel:     nil,
			texts:   map[uint16]string{},
			wantPkg: "",
		},
	}
	for _, c := range cases {
		d := newDec()
		d.texts = c.texts
		ref, err := buildAlias(d, "bank", c.rel)
		if err != nil {
			t.Fatalf("%s: buildAlias: %v", c.name, err)
		}
		if len(d.ast.Sees) != 1 {
			t.Fatalf("%s: Sees len = %d, want 1", c.name, len(d.ast.Sees))
		}
		got := d.ast.Sees[0]
		if got.Alias != "bank" || got.Package != c.wantPkg {
			t.Errorf("%s: See = %+v, want Alias=bank Package=%q", c.name, got, c.wantPkg)
		}
		if ref.Sec != codebook.SectionSee || len(ref.Path) != 1 || ref.Path[0] != 0 {
			t.Errorf("%s: StmtRef = %+v", c.name, ref)
		}
	}
}
