//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildRuleOrInline — buildRuleOrInline appends a case node with an open case, else a Rules entry
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBuildRuleOrInline checks a Rule node inside an open case is routed to
// appendCaseNode (a case Node is appended), and a Rule node with no open case
// becomes a Rules InlineRule carrying the recovered condition tree.
func TestBuildRuleOrInline(t *testing.T) {
	t.Run("inside an open case", func(t *testing.T) {
		d := newDec(testWire)
		d.curCase = 0
		d.ast.Cases = []ast.Case{{}}
		p := toulmin.NodePacket{Role: toulmin.RoleGeneral}
		labels := map[toulmin.TID][]string{}

		ref, err := buildRuleOrInline(d, p, "myrule", labels, nil)
		if err != nil {
			t.Fatalf("buildRuleOrInline: %v", err)
		}
		if len(d.ast.Cases[0].Nodes) != 1 || d.ast.Cases[0].Nodes[0].Name != "myrule" {
			t.Errorf("Cases[0].Nodes = %+v", d.ast.Cases[0].Nodes)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 1, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("outside any case", func(t *testing.T) {
		d := newDec(testWire) // curCase == -1
		p := toulmin.NodePacket{}
		labels := map[toulmin.TID][]string{}
		cond := ast.Compare{Field: "a", Op: "equals", Value: "1"}

		ref, err := buildRuleOrInline(d, p, "myinline", labels, cond)
		if err != nil {
			t.Fatalf("buildRuleOrInline: %v", err)
		}
		want := []ast.InlineRule{{Name: "myinline", Cond: cond}}
		if !reflect.DeepEqual(d.ast.Rules, want) {
			t.Errorf("Rules = %+v, want %+v", d.ast.Rules, want)
		}
		wantRef := doc.StmtRef{Sec: toulmin.SectionRules, Path: []int{0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})
	host := newDec(testWire)
	host.texts[5], host.texts[6] = "limit", "source-body"
	host.funcRef[9] = funcInfo{name: "IsBlocked"}
	hp := toulmin.NodePacket{Func: &toulmin.FuncRef{Alias: 9}, With: []toulmin.TID{5}, Grounds: []toulmin.TID{6}}
	if _, err := buildRuleOrInline(host, hp, "ip blocked", nil, nil); err != nil {
		t.Fatalf("host rule: %v", err)
	}
	if r := host.ast.Rules[0]; r.Using == nil || r.Using.Name != "IsBlocked" || r.With[0] != "limit" || r.Grounds[0] != "source-body" {
		t.Fatalf("host rule: %+v", r)
	}
	if _, err := buildRuleOrInline(newDec(testWire), toulmin.NodePacket{Func: &toulmin.FuncRef{Alias: 77}}, "x", nil, nil); err == nil {
		t.Fatal("unknown FUNC node: want error")
	}
}
