//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestBuildTerm — buildTerm resolves Value from means and SpecRef from sameAs, or propagates a Func error
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBuildTerm checks buildTerm with neither a means nor a sameAs edge (bare
// definition), with a means edge (Value set), with a sameAs edge resolving
// through a registered Func node (SpecRef set), and with a sameAs edge pointing
// at an unregistered Func node (the termSpecRef/resolveFuncRef error propagates).
func TestBuildTerm(t *testing.T) {
	t.Run("no means, no sameAs", func(t *testing.T) {
		d := newDec(testWire)
		ref, err := buildTerm(d, 3, "flag", toulmin.RoleGeneral, nil)
		if err != nil {
			t.Fatalf("buildTerm: %v", err)
		}
		want := ast.Definition{Name: "flag", Kind: ast.DefKind(toulmin.RoleGeneral)}
		if len(d.ast.Defs) != 1 || !reflect.DeepEqual(d.ast.Defs[0], want) {
			t.Errorf("Defs = %+v, want [%+v]", d.ast.Defs, want)
		}
		if d.defByTID[3] != 0 {
			t.Errorf("defByTID[3] = %d, want 0", d.defByTID[3])
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("means edge sets Value", func(t *testing.T) {
		d := newDec(testWire)
		d.texts[8] = "650"
		rel := []toulmin.EdgePacket{{Prop: toulmin.PropMeans, DstTID: 8}}
		_, err := buildTerm(d, 3, "threshold", toulmin.RoleGeneral, rel)
		if err != nil {
			t.Fatalf("buildTerm: %v", err)
		}
		if d.ast.Defs[0].Value != "650" {
			t.Errorf("Value = %q, want 650", d.ast.Defs[0].Value)
		}
	})

	t.Run("sameAs edge resolves SpecRef", func(t *testing.T) {
		d := newDec(testWire)
		d.funcRef[9] = funcInfo{alias: 0, name: "Threshold"}
		rel := []toulmin.EdgePacket{{Prop: toulmin.PropSameAs, DstTID: 9}}
		_, err := buildTerm(d, 3, "threshold", toulmin.RoleGeneral, rel)
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
			d := newDec(testWire)
			d.texts[8] = tc.payload
			_, err := buildTerm(d, 3, "x", toulmin.RoleCode(tc.kind), []toulmin.EdgePacket{{Prop: toulmin.PropMeans, DstTID: 8}})
			if err != nil || !reflect.DeepEqual(d.ast.Defs[0], tc.want) {
				t.Fatalf("kind %v: got %#v err=%v", tc.kind, d.ast.Defs[0], err)
			}
		}
	})

	t.Run("sameAs edge to unregistered Func node errors", func(t *testing.T) {
		d := newDec(testWire)
		rel := []toulmin.EdgePacket{{Prop: toulmin.PropSameAs, DstTID: 9}}
		if _, err := buildTerm(d, 3, "threshold", toulmin.RoleGeneral, rel); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	for _, kind := range []ast.DefKind{ast.StringListDef, ast.FileSpecDef} {
		bd := newDec(testWire)
		bd.texts[9] = "{not json"
		if _, err := buildTerm(bd, 1, "t", toulmin.RoleCode(kind), []toulmin.EdgePacket{{Prop: toulmin.PropMeans, DstTID: 9}}); err == nil {
			t.Fatalf("kind %d with bad JSON: want error", kind)
		}
	}
}
