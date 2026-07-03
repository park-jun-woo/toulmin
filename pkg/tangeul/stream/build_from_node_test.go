//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildFromNode — buildFromNode dispatches every NodeType to its builder and errors on an unknown type
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestBuildFromNode drives every case of buildFromNode's NodeType switch —
// namespace, alias, term, field, rule (both inside and outside a case via
// buildRuleOrInline's own routing is not re-tested here, just the dispatch),
// case, endpoint, trigger — plus the default arm for an unrecognized type, and
// checks that nodeType/nodeName are always recorded first.
func TestBuildFromNode(t *testing.T) {
	labels := map[uint16][]string{1: {"x"}}

	// build runs buildFromNode with the shared labels, failing on error.
	build := func(t *testing.T, d *dec, p packet.NodePacket, cond ast.Expr) doc.StmtRef {
		t.Helper()
		ref, err := buildFromNode(d, p, labels, nil, cond)
		if err != nil {
			t.Fatalf("buildFromNode: %v", err)
		}
		return ref
	}

	t.Run("namespace", func(t *testing.T) {
		d := newDec()
		ref := build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeNamespace}, nil)
		if d.ast.Subject != "x" || ref.Sec != codebook.SectionSubject {
			t.Errorf("Subject=%q ref=%+v", d.ast.Subject, ref)
		}
	})

	t.Run("alias", func(t *testing.T) {
		d := newDec()
		ref := build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeAlias}, nil)
		if len(d.ast.Sees) != 1 || ref.Sec != codebook.SectionSee {
			t.Errorf("Sees=%+v ref=%+v", d.ast.Sees, ref)
		}
	})

	t.Run("term", func(t *testing.T) {
		d := newDec()
		ref := build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeTerm, Role: codebook.RoleCode(ast.ConstDef)}, nil)
		if len(d.ast.Defs) != 1 || ref.Sec != codebook.SectionDefinitions {
			t.Errorf("Defs=%+v ref=%+v", d.ast.Defs, ref)
		}
	})

	t.Run("field", func(t *testing.T) {
		// No has/requires edge: buildField errors, but the dispatch itself is
		// what this case is exercising.
		p := packet.NodePacket{TID: 1, Type: codebook.NodeField}
		if _, err := buildFromNode(newDec(), p, labels, nil, nil); err == nil {
			t.Fatalf("want error propagated from buildField, got nil")
		}
	})

	t.Run("rule", func(t *testing.T) {
		d := newDec() // curCase == -1 -> inline rule
		cond := ast.Compare{Field: "a", Op: "equals", Value: "1"}
		ref := build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeRule}, cond)
		if len(d.ast.Rules) != 1 || ref.Sec != codebook.SectionRules {
			t.Errorf("Rules=%+v ref=%+v", d.ast.Rules, ref)
		}
	})

	t.Run("case", func(t *testing.T) {
		d := newDec()
		ref := build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeCase}, nil)
		if len(d.ast.Cases) != 1 || d.curCase != 0 || ref.Sec != codebook.SectionCases {
			t.Errorf("Cases=%+v curCase=%d ref=%+v", d.ast.Cases, d.curCase, ref)
		}
		if d.caseByTID[1] != 0 {
			t.Errorf("caseByTID[1] = %d, want 0", d.caseByTID[1])
		}
	})

	t.Run("endpoint", func(t *testing.T) {
		d := newDec()
		ref := build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeEndpoint}, nil)
		if len(d.ast.Provides) != 1 || ref.Sec != codebook.SectionProvides {
			t.Errorf("Provides=%+v ref=%+v", d.ast.Provides, ref)
		}
		if d.epByTID[1] != 0 {
			t.Errorf("epByTID[1] = %d, want 0", d.epByTID[1])
		}
	})

	t.Run("trigger", func(t *testing.T) {
		d := newDec()
		ref := build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeTrigger, Role: codebook.RoleCode(ast.OnEvent)}, nil)
		if len(d.ast.Internals) != 1 || ref.Sec != codebook.SectionInternal {
			t.Errorf("Internals=%+v ref=%+v", d.ast.Internals, ref)
		}
	})

	t.Run("unknown node type", func(t *testing.T) {
		if _, err := buildFromNode(newDec(), packet.NodePacket{TID: 1, Type: codebook.NodeFunc}, labels, nil, nil); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("records nodeType and nodeName regardless of type", func(t *testing.T) {
		d := newDec()
		build(t, d, packet.NodePacket{TID: 1, Type: codebook.NodeCase}, nil)
		if d.nodeType[1] != codebook.NodeCase || d.nodeName[1] != "x" {
			t.Errorf("nodeType=%v nodeName=%q", d.nodeType[1], d.nodeName[1])
		}
	})
}
