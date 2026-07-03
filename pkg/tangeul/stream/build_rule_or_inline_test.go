//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildRuleOrInline — buildRuleOrInline appends a case node with an open case, else a Rules entry
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestBuildRuleOrInline checks a Rule node inside an open case is routed to
// appendCaseNode (a case Node is appended), and a Rule node with no open case
// becomes a Rules InlineRule carrying the recovered condition tree.
func TestBuildRuleOrInline(t *testing.T) {
	t.Run("inside an open case", func(t *testing.T) {
		d := newDec()
		d.curCase = 0
		d.ast.Cases = []ast.Case{{}}
		p := packet.NodePacket{Role: codebook.RoleGeneral}
		labels := map[uint16][]string{}

		ref, err := buildRuleOrInline(d, p, "myrule", labels, nil)
		if err != nil {
			t.Fatalf("buildRuleOrInline: %v", err)
		}
		if len(d.ast.Cases[0].Nodes) != 1 || d.ast.Cases[0].Nodes[0].Name != "myrule" {
			t.Errorf("Cases[0].Nodes = %+v", d.ast.Cases[0].Nodes)
		}
		wantRef := doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 1, 0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})

	t.Run("outside any case", func(t *testing.T) {
		d := newDec() // curCase == -1
		p := packet.NodePacket{}
		labels := map[uint16][]string{}
		cond := ast.Compare{Field: "a", Op: "equals", Value: "1"}

		ref, err := buildRuleOrInline(d, p, "myinline", labels, cond)
		if err != nil {
			t.Fatalf("buildRuleOrInline: %v", err)
		}
		want := []ast.InlineRule{{Name: "myinline", Cond: cond}}
		if !reflect.DeepEqual(d.ast.Rules, want) {
			t.Errorf("Rules = %+v, want %+v", d.ast.Rules, want)
		}
		wantRef := doc.StmtRef{Sec: codebook.SectionRules, Path: []int{0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("StmtRef = %+v, want %+v", ref, wantRef)
		}
	})
}
