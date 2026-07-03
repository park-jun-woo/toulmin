//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendCaseNode — rebuilds an in-case node across role/func/with/qual/check branch combinations
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestAppendCaseNode iterates scenarios that together exercise every branch of
// appendCaseNode: an unknown RoleCode (hard error), a bare node with no
// optional slots, a node whose Func resolves successfully with a populated With
// list and a Qualified qualifier, a node whose Func TID is unregistered (a
// resolveFuncRef error), and a node carrying an F_CHECK deferred fixup.
func TestAppendCaseNode(t *testing.T) {
	newDec := func() *dec {
		return &dec{
			ast:   &ast.Document{Cases: []ast.Case{{Name: "c0"}}},
			texts: map[uint16]string{10: "term-a", 11: "term-b"},
		}
	}

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "unknown role errors",
			run: func(t *testing.T) {
				d := newDec()
				_, err := appendCaseNode(d, packet.NodePacket{Role: codebook.RoleCode(9)}, "n", nil)
				if err == nil {
					t.Fatalf("want error for unknown role")
				}
				if len(d.ast.Cases[0].Nodes) != 0 {
					t.Errorf("Nodes should be untouched on role error, got %v", d.ast.Cases[0].Nodes)
				}
			},
		},
		{
			name: "bare node, no optional slots",
			run: func(t *testing.T) {
				d := newDec()
				ref, err := appendCaseNode(d, packet.NodePacket{Role: codebook.RoleGeneral}, "bare", nil)
				if err != nil {
					t.Fatalf("appendCaseNode: %v", err)
				}
				if len(d.ast.Cases[0].Nodes) != 1 {
					t.Fatalf("Nodes len = %d, want 1", len(d.ast.Cases[0].Nodes))
				}
				n := d.ast.Cases[0].Nodes[0]
				if n.Name != "bare" || n.Role != ast.GeneralRule {
					t.Errorf("node = %+v", n)
				}
				if n.Using != nil || n.With != nil || n.Qualified != nil || n.Checking != "" {
					t.Errorf("expected all optional fields empty, got %+v", n)
				}
				wantRef := doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 1, 0}}
				if ref.Sec != wantRef.Sec || len(ref.Path) != 3 || ref.Path[2] != 0 {
					t.Errorf("ref = %+v, want %+v", ref, wantRef)
				}
				if len(d.checkFixups) != 0 {
					t.Errorf("checkFixups should be empty, got %v", d.checkFixups)
				}
			},
		},
		{
			name: "func resolves, with, qual and check all present",
			run: func(t *testing.T) {
				d := newDec()
				d.funcRef = map[uint16]funcInfo{5: {alias: 0, name: "hasBalance"}}
				qual := word.VFix(500)
				checkTID := packet.TID(3)
				p := packet.NodePacket{
					Role:  codebook.RoleCounter,
					Func:  &packet.FuncRef{Alias: 5},
					With:  []packet.TID{10, 11},
					Qual:  &qual,
					Check: &checkTID,
				}
				ref, err := appendCaseNode(d, p, "full", nil)
				if err != nil {
					t.Fatalf("appendCaseNode: %v", err)
				}
				n := d.ast.Cases[0].Nodes[0]
				if n.Using == nil || n.Using.Name != "hasBalance" {
					t.Errorf("Using = %+v", n.Using)
				}
				if len(n.With) != 2 || n.With[0] != "term-a" || n.With[1] != "term-b" {
					t.Errorf("With = %v", n.With)
				}
				if n.Qualified == nil {
					t.Errorf("Qualified = nil, want set")
				}
				if len(d.checkFixups) != 1 {
					t.Fatalf("checkFixups len = %d, want 1", len(d.checkFixups))
				}
				fx := d.checkFixups[0]
				if fx.caseIdx != 0 || fx.slotIdx != 0 || fx.tid != 3 {
					t.Errorf("fixup = %+v", fx)
				}
				if len(ref.Path) != 3 || ref.Path[2] != 0 {
					t.Errorf("ref = %+v", ref)
				}
			},
		},
		{
			name: "func unresolved errors",
			run: func(t *testing.T) {
				d := newDec()
				p := packet.NodePacket{Role: codebook.RoleGeneral, Func: &packet.FuncRef{Alias: 99}}
				_, err := appendCaseNode(d, p, "n", nil)
				if err == nil {
					t.Fatalf("want error for unresolved func ref")
				}
				if len(d.ast.Cases[0].Nodes) != 0 {
					t.Errorf("Nodes should be untouched on func resolution error, got %v", d.ast.Cases[0].Nodes)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
