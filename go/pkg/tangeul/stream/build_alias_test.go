//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestBuildAlias — buildAlias resolves the package path from importedFrom, or leaves it empty without one
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// TestBuildAlias checks buildAlias appends an ast.See named after the label,
// with its Package resolved from the importedFrom edge's Text when present, and
// left empty when the edge is absent.
func TestBuildAlias(t *testing.T) {
	cases := []struct {
		name    string
		rel     []toulmin.EdgePacket
		texts   map[toulmin.TID]string
		wantPkg string
	}{
		{
			name:    "with importedFrom edge",
			rel:     []toulmin.EdgePacket{{Prop: toulmin.PropImportedFrom, DstTID: 5}},
			texts:   map[toulmin.TID]string{5: "bank/pkg"},
			wantPkg: "bank/pkg",
		},
		{
			name:    "without importedFrom edge",
			rel:     nil,
			texts:   map[toulmin.TID]string{},
			wantPkg: "",
		},
	}
	for _, c := range cases {
		d := newDec(testWire)
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
		if ref.Sec != toulmin.SectionSee || len(ref.Path) != 1 || ref.Path[0] != 0 {
			t.Errorf("%s: StmtRef = %+v", c.name, ref)
		}
	}
}
