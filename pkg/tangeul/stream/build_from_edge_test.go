//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildFromEdge — buildFromEdge routes a single relation edge by PropCode, or errors on shape
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestBuildFromEdge covers the relation-count guard, every PropCode switch arm
// (attacks in and out of a case, do/undo/run delegated to buildExecEdge, runs,
// evaluates delegated to buildCheckEdge, triggers) and the default arm for an
// unexpected PropCode.
func TestBuildFromEdge(t *testing.T) {
	t.Run("wrong relation count", func(t *testing.T) {
		rel := []packet.EdgePacket{{Prop: codebook.PropAttacks}, {Prop: codebook.PropRuns}}
		if _, err := buildFromEdge(newDec(), rel); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("attacks outside any case", func(t *testing.T) {
		if _, err := buildFromEdge(newDec(), []packet.EdgePacket{{Prop: codebook.PropAttacks}}); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("attacks inside a case", func(t *testing.T) {
		d := newDec()
		d.curCase = 0
		d.ast.Cases = []ast.Case{{}}
		d.nodeName[1] = "attacker"
		d.nodeName[2] = "target"
		ref, err := buildFromEdge(d, []packet.EdgePacket{{Prop: codebook.PropAttacks, Src: 1, DstTID: 2}})
		if err != nil {
			t.Fatalf("buildFromEdge: %v", err)
		}
		if len(d.ast.Cases[0].Attacks) != 1 {
			t.Fatalf("Attacks len = %d, want 1", len(d.ast.Cases[0].Attacks))
		}
		if ref.Sec != codebook.SectionCases {
			t.Errorf("StmtRef.Sec = %v", ref.Sec)
		}
	})

	t.Run("do/undo/run delegates to buildExecEdge", func(t *testing.T) {
		// newDec starts with curCase == -1, so the delegated call errors.
		if _, err := buildFromEdge(newDec(), []packet.EdgePacket{{Prop: codebook.PropDo}}); err == nil {
			t.Fatalf("want error propagated from buildExecEdge, got nil")
		}
	})

	t.Run("runs", func(t *testing.T) {
		d := newDec()
		d.epByTID[1] = 0
		d.ast.Provides = []ast.Endpoint{{Name: "ep"}}
		d.nodeName[2] = "caseA"
		ref, err := buildFromEdge(d, []packet.EdgePacket{{Prop: codebook.PropRuns, Src: 1, DstTID: 2}})
		if err != nil {
			t.Fatalf("buildFromEdge: %v", err)
		}
		if len(d.ast.Provides[0].Runs) != 1 || d.ast.Provides[0].Runs[0] != "caseA" {
			t.Errorf("Provides[0].Runs = %v, want [caseA]", d.ast.Provides[0].Runs)
		}
		if ref.Sec != codebook.SectionProvides {
			t.Errorf("StmtRef.Sec = %v", ref.Sec)
		}
	})

	t.Run("evaluates delegates to buildCheckEdge", func(t *testing.T) {
		d := newDec()
		d.nodeType[1] = codebook.NodeTrigger
		d.trByTID[1] = 0
		d.ast.Internals = []ast.Internal{{}}
		d.nodeName[2] = "caseB"
		ref, err := buildFromEdge(d, []packet.EdgePacket{{Prop: codebook.PropEvaluates, Src: 1, DstTID: 2}})
		if err != nil {
			t.Fatalf("buildFromEdge: %v", err)
		}
		if len(d.ast.Internals[0].Checks) != 1 || d.ast.Internals[0].Checks[0] != "caseB" {
			t.Errorf("Internals[0].Checks = %v, want [caseB]", d.ast.Internals[0].Checks)
		}
		if ref.Sec != codebook.SectionInternal {
			t.Errorf("StmtRef.Sec = %v", ref.Sec)
		}
	})

	t.Run("triggers", func(t *testing.T) {
		d := newDec()
		d.trByTID[1] = 0
		d.ast.Internals = []ast.Internal{{}}
		d.nodeName[2] = "caseC"
		ref, err := buildFromEdge(d, []packet.EdgePacket{{Prop: codebook.PropTriggers, Src: 1, DstTID: 2}})
		if err != nil {
			t.Fatalf("buildFromEdge: %v", err)
		}
		if len(d.ast.Internals[0].Runs) != 1 || d.ast.Internals[0].Runs[0] != "caseC" {
			t.Errorf("Internals[0].Runs = %v, want [caseC]", d.ast.Internals[0].Runs)
		}
		if ref.Sec != codebook.SectionInternal {
			t.Errorf("StmtRef.Sec = %v", ref.Sec)
		}
	})

	t.Run("unexpected relation prop", func(t *testing.T) {
		if _, err := buildFromEdge(newDec(), []packet.EdgePacket{{Prop: codebook.PropMeans}}); err == nil {
			t.Fatalf("want error, got nil")
		}
	})
}
