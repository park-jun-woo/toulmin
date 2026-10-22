//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildCheckEdge — buildCheckEdge routes an evaluates edge by its src node's endpoint or internal type
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBuildCheckEdge checks that an evaluates edge whose src is an Endpoint node
// appends to Provides[i].Checks, while one whose src is any other node type
// (a trigger, in practice) appends to Internals[i].Checks, each returning the
// matching StmtRef.
func TestBuildCheckEdge(t *testing.T) {
	t.Run("endpoint src", func(t *testing.T) {
		d := newDec(testWire)
		d.nodeType[1] = toulmin.NodeEndpoint
		d.epByTID[1] = 0
		d.ast.Provides = []ast.Endpoint{{Name: "ep"}}
		d.nodeName[2] = "caseA"

		ref, err := buildCheckEdge(d, toulmin.EdgePacket{Src: 1, DstTID: 2})
		if err != nil {
			t.Fatalf("buildCheckEdge: %v", err)
		}
		if len(d.ast.Provides[0].Checks) != 1 || d.ast.Provides[0].Checks[0] != "caseA" {
			t.Errorf("Provides[0].Checks = %v, want [caseA]", d.ast.Provides[0].Checks)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{0, 2, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("internal (non-endpoint) src", func(t *testing.T) {
		d := newDec(testWire)
		d.nodeType[1] = toulmin.NodeTrigger
		d.trByTID[1] = 0
		d.ast.Internals = []ast.Internal{{}}
		d.nodeName[2] = "caseB"

		ref, err := buildCheckEdge(d, toulmin.EdgePacket{Src: 1, DstTID: 2})
		if err != nil {
			t.Fatalf("buildCheckEdge: %v", err)
		}
		if len(d.ast.Internals[0].Checks) != 1 || d.ast.Internals[0].Checks[0] != "caseB" {
			t.Errorf("Internals[0].Checks = %v, want [caseB]", d.ast.Internals[0].Checks)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{0, 2, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})
}
