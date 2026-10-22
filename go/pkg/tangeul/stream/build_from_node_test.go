//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildFromNode — buildFromNode dispatches every NodeType to its builder and errors on an unknown type
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBuildFromNode drives every case of buildFromNode's NodeType switch —
// namespace, alias, term, field, rule (both inside and outside a case via
// buildRuleOrInline's own routing is not re-tested here, just the dispatch),
// case, endpoint, trigger — plus the default arm for an unrecognized type, and
// checks that nodeType/nodeName are always recorded first.
func TestBuildFromNode(t *testing.T) {
	labels := map[toulmin.TID][]string{1: {"x"}}

	// build runs buildFromNode with the shared labels, failing on error.
	build := func(t *testing.T, d *dec, p toulmin.NodePacket, cond ast.Expr) doc.StmtRef {
		t.Helper()
		ref, err := buildFromNode(d, p, labels, nil, cond)
		if err != nil {
			t.Fatalf("buildFromNode: %v", err)
		}
		return ref
	}

	t.Run("namespace", func(t *testing.T) {
		d := newDec(testWire)
		ref := build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeNamespace}, nil)
		if d.ast.Subject != "x" || ref.Sec != toulmin.SectionSubject {
			t.Errorf("Subject=%q ref=%+v", d.ast.Subject, ref)
		}
	})

	t.Run("alias", func(t *testing.T) {
		d := newDec(testWire)
		ref := build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeAlias}, nil)
		if len(d.ast.Sees) != 1 || ref.Sec != toulmin.SectionSee {
			t.Errorf("Sees=%+v ref=%+v", d.ast.Sees, ref)
		}
	})

	t.Run("term", func(t *testing.T) {
		d := newDec(testWire)
		ref := build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeTerm, Role: toulmin.RoleCode(ast.ConstDef)}, nil)
		if len(d.ast.Defs) != 1 || ref.Sec != toulmin.SectionDefinitions {
			t.Errorf("Defs=%+v ref=%+v", d.ast.Defs, ref)
		}
	})

	t.Run("field", func(t *testing.T) {
		// No has/requires edge: buildField errors, but the dispatch itself is
		// what this case is exercising.
		p := toulmin.NodePacket{TID: 1, Type: toulmin.NodeField}
		if _, err := buildFromNode(newDec(testWire), p, labels, nil, nil); err == nil {
			t.Fatalf("want error propagated from buildField, got nil")
		}
	})

	t.Run("rule", func(t *testing.T) {
		d := newDec(testWire) // curCase == -1 -> inline rule
		cond := ast.Compare{Field: "a", Op: "equals", Value: "1"}
		ref := build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeRule}, cond)
		if len(d.ast.Rules) != 1 || ref.Sec != toulmin.SectionRules {
			t.Errorf("Rules=%+v ref=%+v", d.ast.Rules, ref)
		}
	})

	t.Run("case", func(t *testing.T) {
		d := newDec(testWire)
		ref := build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeCase}, nil)
		if len(d.ast.Cases) != 1 || d.curCase != 0 || ref.Sec != toulmin.SectionCases {
			t.Errorf("Cases=%+v curCase=%d ref=%+v", d.ast.Cases, d.curCase, ref)
		}
		if d.caseByTID[1] != 0 {
			t.Errorf("caseByTID[1] = %d, want 0", d.caseByTID[1])
		}
	})

	t.Run("endpoint", func(t *testing.T) {
		d := newDec(testWire)
		ref := build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeEndpoint}, nil)
		if len(d.ast.Provides) != 1 || ref.Sec != toulmin.SectionProvides {
			t.Errorf("Provides=%+v ref=%+v", d.ast.Provides, ref)
		}
		if d.epByTID[1] != 0 {
			t.Errorf("epByTID[1] = %d, want 0", d.epByTID[1])
		}
	})

	t.Run("trigger", func(t *testing.T) {
		d := newDec(testWire)
		ref := build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeTrigger, Role: toulmin.RoleCode(ast.OnEvent)}, nil)
		if len(d.ast.Internals) != 1 || ref.Sec != toulmin.SectionInternal {
			t.Errorf("Internals=%+v ref=%+v", d.ast.Internals, ref)
		}
	})

	t.Run("unknown node type", func(t *testing.T) {
		if _, err := buildFromNode(newDec(testWire), toulmin.NodePacket{TID: 1, Type: toulmin.NodeFunc}, labels, nil, nil); err == nil {
			t.Fatalf("want error, got nil")
		}
	})

	t.Run("records nodeType and nodeName regardless of type", func(t *testing.T) {
		d := newDec(testWire)
		build(t, d, toulmin.NodePacket{TID: 1, Type: toulmin.NodeCase}, nil)
		if d.nodeType[1] != toulmin.NodeCase || d.nodeName[1] != "x" {
			t.Errorf("nodeType=%v nodeName=%q", d.nodeType[1], d.nodeName[1])
		}
	})
}
