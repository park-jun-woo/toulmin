//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestTransferGolden — the §10 balance-sufficient node pins to its words and every §10 packet round-trips
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestTransferGolden fixes the §10 account-transfer example. The general
// RuleNode `balance sufficient` is pinned to its exact word sequence
// (TID/type/header/FuncRef/Claim/verdict, matching §10's documented word
// values; the spec's frame Len annotations are not binding). Every node, edge and
// control of the §10 case is then round-tripped. The provenance edge is excluded
// (v0.1 not implemented).
func TestTransferGolden(t *testing.T) {
	span := word.Span{TID: 1}
	verdict := word.VFix(0)
	frozenVerdict := word.VFix(0)

	balance := NodePacket{
		TID: 0x0013, Type: codebook.NodeRule, Role: codebook.RoleGeneral,
		Flags:   codebook.NodeFlagFunc | codebook.NodeFlagClaim | codebook.NodeFlagVerdict,
		Func:    &FuncRef{Alias: 0x0002, Verb: 0x7400},
		Claim:   &Claim{Subject: 0x0011, Predicate: 0x7400, Kind: ObjTID, ObjTID: 0x0012},
		Verdict: &verdict,
		Span:    span,
	}
	wantBalance := []word.Word{
		0x100D, 0x0013, 0x0003, 0x0019, 0x0002, 0x7400,
		0x0011, 0x7400, 0x0012, 0x0000, 0x0000, 0x0001, 0x0000, 0x0000,
	}
	ws, err := EncodeNode(balance)
	if err != nil {
		t.Fatalf("EncodeNode(balance): %v", err)
	}
	if !reflect.DeepEqual(ws, wantBalance) {
		t.Fatalf("balance node =\n %04X\nwant\n %04X", ws, wantBalance)
	}

	nodes := []NodePacket{
		{TID: 0x0011, Type: codebook.NodeField, Span: span},
		{TID: 0x0012, Type: codebook.NodeField, Span: span},
		balance,
		{
			TID: 0x0014, Type: codebook.NodeRule, Role: codebook.RoleCounter,
			Flags:   codebook.NodeFlagFunc | codebook.NodeFlagClaim | codebook.NodeFlagVerdict,
			Func:    &FuncRef{Alias: 0x0002, Verb: 0x7000},
			Claim:   &Claim{Subject: 0x0011, Predicate: 0x7000, Kind: ObjTID, ObjTID: 0x0012},
			Verdict: &frozenVerdict,
			Span:    span,
		},
	}
	for i, in := range nodes {
		w, err := EncodeNode(in)
		if err != nil {
			t.Fatalf("node %d: EncodeNode: %v", i, err)
		}
		got, _, err := DecodeNode(w)
		if err != nil || !reflect.DeepEqual(got, in) {
			t.Fatalf("node %d: DecodeNode = %+v, %v; want %+v", i, got, err, in)
		}
	}

	edges := []EdgePacket{
		{Prop: codebook.PropRequires, Src: 0x0010, DstTID: 0x0011, Span: span},
		{Prop: codebook.PropRequires, Src: 0x0010, DstTID: 0x0012, Span: span},
		{Prop: codebook.PropAttacks, Src: 0x0014, DstTID: 0x0013, Span: span},
		{Prop: codebook.PropDo, Src: 0x0013, DstFunc: &FuncRef{Alias: 0x0002, Verb: 0xC200}, Flags: codebook.EdgeFlagOnce, Span: span},
		{Prop: codebook.PropUndo, Src: 0x0013, DstFunc: &FuncRef{Alias: 0x0002, Verb: 0xA000}, Span: span},
		{Prop: codebook.PropRun, Src: 0x0013, DstTID: 0x0020, Span: span},
	}
	for i, in := range edges {
		w, err := EncodeEdge(in)
		if err != nil {
			t.Fatalf("edge %d: EncodeEdge: %v", i, err)
		}
		got, _, err := DecodeEdge(w)
		if err != nil || !reflect.DeepEqual(got, in) {
			t.Fatalf("edge %d: DecodeEdge = %+v, %v; want %+v", i, got, err, in)
		}
	}

	controls := []ControlPacket{
		{Code: codebook.ControlCaseBegin, Arg: 0x0010},
		{Code: codebook.ControlCaseEnd, Arg: 0x0010},
	}
	for i, in := range controls {
		w, err := EncodeControl(in)
		if err != nil {
			t.Fatalf("control %d: EncodeControl: %v", i, err)
		}
		got, _, err := DecodeControl(w)
		if err != nil || !reflect.DeepEqual(got, in) {
			t.Fatalf("control %d: DecodeControl = %+v, %v; want %+v", i, got, err, in)
		}
	}
}
