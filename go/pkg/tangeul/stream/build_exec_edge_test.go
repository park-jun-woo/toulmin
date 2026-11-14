//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildExecEdge — buildExecEdge covers its case guard, gate, do/undo FuncRef and run fixup paths
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBuildExecEdge exercises every branch of buildExecEdge: the outside-case
// guard, an invalid ExecKindFromProp prop, an invalid gate operator, a do/undo
// edge missing its FuncRef dst, a do/undo edge whose Func node is unregistered,
// a successful do edge (with Once and a valid gate) resolving its Func, and a
// successful run edge deferring a case-name fixup.
func TestBuildExecEdge(t *testing.T) {
	// inCase builds a decoder state positioned inside an (empty) case 0.
	inCase := func() *dec {
		d := newDec(testWire)
		d.curCase = 0
		d.ast.Cases = []ast.Case{{}}
		return d
	}

	t.Run("outside any case", func(t *testing.T) {
		// newDec starts with curCase == -1.
		if _, err := buildExecEdge(newDec(testWire), toulmin.EdgePacket{Prop: toulmin.PropDo}); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("not an execution prop", func(t *testing.T) {
		if _, err := buildExecEdge(inCase(), toulmin.EdgePacket{Prop: toulmin.PropAttacks}); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("invalid gate op", func(t *testing.T) {
		e := toulmin.EdgePacket{Prop: toulmin.PropDo, Gate: &toulmin.Gate{Op: toulmin.GateOp(99), Percent: 50}}
		if _, err := buildExecEdge(inCase(), e); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("do missing FuncRef dst", func(t *testing.T) {
		if _, err := buildExecEdge(inCase(), toulmin.EdgePacket{Prop: toulmin.PropDo}); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("do unresolvable FuncRef", func(t *testing.T) {
		e := toulmin.EdgePacket{Prop: toulmin.PropDo, DstFunc: &toulmin.FuncRef{Alias: 999}}
		if _, err := buildExecEdge(inCase(), e); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("do success with once and gate", func(t *testing.T) {
		d := inCase()
		d.nodeName[7] = "brewing complete"
		d.funcRef[5] = funcInfo{alias: 0, name: "serve"}
		e := toulmin.EdgePacket{
			Prop:    toulmin.PropDo,
			Src:     7,
			Flags:   toulmin.EdgeFlagOnce,
			Gate:    &toulmin.Gate{Op: toulmin.GateAtLeast, Percent: 75},
			DstFunc: &toulmin.FuncRef{Alias: 5},
		}
		ref, err := buildExecEdge(d, e)
		if err != nil {
			t.Fatalf("buildExecEdge: %v", err)
		}
		if len(d.ast.Cases[0].Execs) != 1 {
			t.Fatalf("Execs len = %d, want 1", len(d.ast.Cases[0].Execs))
		}
		x := d.ast.Cases[0].Execs[0]
		want := ast.Exec{
			Kind:      ast.DoExec,
			Node:      "brewing complete",
			Once:      true,
			Certainty: &ast.Certainty{Op: "at least", Percent: 75},
			Func:      &ast.Ref{Name: "serve"},
		}
		if !reflect.DeepEqual(x, want) {
			t.Errorf("Exec = %+v, want %+v", x, want)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 3, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("run success defers fixup", func(t *testing.T) {
		d := inCase()
		d.nodeName[7] = "brewing complete"
		ref, err := buildExecEdge(d, toulmin.EdgePacket{Prop: toulmin.PropRun, Src: 7, DstTID: 42})
		if err != nil {
			t.Fatalf("buildExecEdge: %v", err)
		}
		if len(d.ast.Cases[0].Execs) != 1 || d.ast.Cases[0].Execs[0].Kind != ast.RunExec {
			t.Fatalf("Execs = %+v, want one RunExec", d.ast.Cases[0].Execs)
		}
		wantFixup := []nameFixup{{caseIdx: 0, slotIdx: 0, tid: 42}}
		if !reflect.DeepEqual(d.execCaseFixups, wantFixup) {
			t.Errorf("execCaseFixups = %+v, want %+v", d.execCaseFixups, wantFixup)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 3, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})
}
