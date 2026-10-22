//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildField — buildField routes a field node by its has/requires edge, or errors with neither
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBuildField checks a has edge appends a struct Definition field (picking up
// the optional second label as the type), a requires edge routes to its case
// owner via appendRequire (leaving Type empty when there is only one label), and
// a field node with neither edge is a hard error.
func TestBuildField(t *testing.T) {
	t.Run("has edge with typed label", func(t *testing.T) {
		d := newDec(testWire)
		d.defByTID[3] = 0
		d.ast.Defs = []ast.Definition{{Name: "config"}}
		labels := map[toulmin.TID][]string{9: {"host", "Text"}}
		rel := []toulmin.EdgePacket{{Prop: toulmin.PropHas, Src: 3}}

		ref, err := buildField(d, 9, "host", labels, rel)
		if err != nil {
			t.Fatalf("buildField: %v", err)
		}
		wantFields := []ast.Field{{Name: "host", Type: "Text"}}
		if !reflect.DeepEqual(d.ast.Defs[0].Fields, wantFields) {
			t.Errorf("Defs[0].Fields = %+v, want %+v", d.ast.Defs[0].Fields, wantFields)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{0, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("requires edge with untyped label", func(t *testing.T) {
		d := newDec(testWire)
		d.nodeType[4] = toulmin.NodeCase
		d.caseByTID[4] = 0
		d.ast.Cases = []ast.Case{{}}
		labels := map[toulmin.TID][]string{9: {"user"}}
		rel := []toulmin.EdgePacket{{Prop: toulmin.PropRequires, Src: 4}}

		ref, err := buildField(d, 9, "user", labels, rel)
		if err != nil {
			t.Fatalf("buildField: %v", err)
		}
		wantReq := []ast.Require{{Field: "user", Type: ""}}
		if !reflect.DeepEqual(d.ast.Cases[0].Requires, wantReq) {
			t.Errorf("Cases[0].Requires = %+v, want %+v", d.ast.Cases[0].Requires, wantReq)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 0, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("neither has nor requires edge", func(t *testing.T) {
		d := newDec(testWire)
		if _, err := buildField(d, 9, "orphan", nil, nil); err == nil {
			t.Fatalf("want error, got nil")
		}
	})
}
