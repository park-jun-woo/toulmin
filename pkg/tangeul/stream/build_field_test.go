//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildField — buildField routes a field node by its has/requires edge, or errors with neither
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestBuildField checks a has edge appends a struct Definition field (picking up
// the optional second label as the type), a requires edge routes to its case
// owner via appendRequire (leaving Type empty when there is only one label), and
// a field node with neither edge is a hard error.
func TestBuildField(t *testing.T) {
	t.Run("has edge with typed label", func(t *testing.T) {
		d := newDec()
		d.defByTID[3] = 0
		d.ast.Defs = []ast.Definition{{Name: "config"}}
		labels := map[uint16][]string{9: {"host", "Text"}}
		rel := []packet.EdgePacket{{Prop: codebook.PropHas, Src: 3}}

		ref, err := buildField(d, 9, "host", labels, rel)
		if err != nil {
			t.Fatalf("buildField: %v", err)
		}
		wantFields := []ast.Field{{Name: "host", Type: "Text"}}
		if !reflect.DeepEqual(d.ast.Defs[0].Fields, wantFields) {
			t.Errorf("Defs[0].Fields = %+v, want %+v", d.ast.Defs[0].Fields, wantFields)
		}
		wantRef := doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{0, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("requires edge with untyped label", func(t *testing.T) {
		d := newDec()
		d.nodeType[4] = codebook.NodeCase
		d.caseByTID[4] = 0
		d.ast.Cases = []ast.Case{{}}
		labels := map[uint16][]string{9: {"user"}}
		rel := []packet.EdgePacket{{Prop: codebook.PropRequires, Src: 4}}

		ref, err := buildField(d, 9, "user", labels, rel)
		if err != nil {
			t.Fatalf("buildField: %v", err)
		}
		wantReq := []ast.Require{{Field: "user", Type: ""}}
		if !reflect.DeepEqual(d.ast.Cases[0].Requires, wantReq) {
			t.Errorf("Cases[0].Requires = %+v, want %+v", d.ast.Cases[0].Requires, wantReq)
		}
		wantRef := doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 0, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("neither has nor requires edge", func(t *testing.T) {
		d := newDec()
		if _, err := buildField(d, 9, "orphan", nil, nil); err == nil {
			t.Fatalf("want error, got nil")
		}
	})
}
