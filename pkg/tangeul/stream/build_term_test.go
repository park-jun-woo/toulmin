//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildTerm — buildTerm resolves Value from means and SpecRef from sameAs, or propagates a Func error
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestBuildTerm checks buildTerm with neither a means nor a sameAs edge (bare
// definition), with a means edge (Value set), with a sameAs edge resolving
// through a registered Func node (SpecRef set), and with a sameAs edge pointing
// at an unregistered Func node (the termSpecRef/resolveFuncRef error propagates).
func TestBuildTerm(t *testing.T) {
	t.Run("no means, no sameAs", func(t *testing.T) {
		d := newDec()
		ref, err := buildTerm(d, 3, "flag", codebook.RoleGeneral, nil)
		if err != nil {
			t.Fatalf("buildTerm: %v", err)
		}
		want := ast.Definition{Name: "flag", Kind: ast.DefKind(codebook.RoleGeneral)}
		if len(d.ast.Defs) != 1 || !reflect.DeepEqual(d.ast.Defs[0], want) {
			t.Errorf("Defs = %+v, want [%+v]", d.ast.Defs, want)
		}
		if d.defByTID[3] != 0 {
			t.Errorf("defByTID[3] = %d, want 0", d.defByTID[3])
		}
		wantRef := doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("means edge sets Value", func(t *testing.T) {
		d := newDec()
		d.texts[8] = "650"
		rel := []packet.EdgePacket{{Prop: codebook.PropMeans, DstTID: 8}}
		_, err := buildTerm(d, 3, "threshold", codebook.RoleGeneral, rel)
		if err != nil {
			t.Fatalf("buildTerm: %v", err)
		}
		if d.ast.Defs[0].Value != "650" {
			t.Errorf("Value = %q, want 650", d.ast.Defs[0].Value)
		}
	})

	t.Run("sameAs edge resolves SpecRef", func(t *testing.T) {
		d := newDec()
		d.funcRef[9] = funcInfo{alias: 0, name: "Threshold"}
		rel := []packet.EdgePacket{{Prop: codebook.PropSameAs, DstTID: 9}}
		_, err := buildTerm(d, 3, "threshold", codebook.RoleGeneral, rel)
		if err != nil {
			t.Fatalf("buildTerm: %v", err)
		}
		want := &ast.Ref{Name: "Threshold"}
		if !reflect.DeepEqual(d.ast.Defs[0].SpecRef, want) {
			t.Errorf("SpecRef = %+v, want %+v", d.ast.Defs[0].SpecRef, want)
		}
	})

	t.Run("structured definition payloads", func(t *testing.T) {
		for _, tc := range []struct {
			kind    ast.DefKind
			payload string
			want    ast.Definition
		}{
			{ast.StringListDef, `["a,b","c"]`, ast.Definition{Name: "x", Kind: ast.StringListDef, Strings: []string{"a,b", "c"}}},
			{ast.FileSpecDef, `{"loadType":"JSON list","filePath":"./x.json"}`, ast.Definition{Name: "x", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "./x.json"}},
		} {
			d := newDec()
			d.texts[8] = tc.payload
			_, err := buildTerm(d, 3, "x", codebook.RoleCode(tc.kind), []packet.EdgePacket{{Prop: codebook.PropMeans, DstTID: 8}})
			if err != nil || !reflect.DeepEqual(d.ast.Defs[0], tc.want) {
				t.Fatalf("kind %v: got %#v err=%v", tc.kind, d.ast.Defs[0], err)
			}
		}
	})

	t.Run("sameAs edge to unregistered Func node errors", func(t *testing.T) {
		d := newDec()
		rel := []packet.EdgePacket{{Prop: codebook.PropSameAs, DstTID: 9}}
		if _, err := buildTerm(d, 3, "threshold", codebook.RoleGeneral, rel); err == nil {
			t.Fatalf("want error, got nil")
		}
	})
}
