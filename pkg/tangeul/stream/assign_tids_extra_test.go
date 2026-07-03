//ff:func feature=tangl type=util control=sequence
//ff:what TestAssignTIDsExtra — assignTIDs on an empty Doc (zero iterations) and on a ResolveStmt-rejected StmtRef
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestAssignTIDsExtra supplements TestAssignTIDs (tid_table_test.go), which
// only walks a well-formed fixture. This file covers the branches that leaves:
// a Doc with no Blocks at all (the range loop's zero-iteration case, primaryEnd
// stays 1), and a StmtBlock whose StmtRef doc.ResolveStmt rejects (a Subject
// ref carrying a non-nil Path), which must short-circuit assignTIDs with that
// error before classifyTID ever runs.
//
// classifyTID's own "unknown section" error branch is not reachable through
// assignTIDs: doc.ResolveStmt and classifyTID switch over the exact same set
// of seven codebook.SectionCode values, so any ref.Sec that would fall to
// classifyTID's default has already been rejected by ResolveStmt's default —
// assignTIDs never reaches classifyTID with an unrecognized section.
func TestAssignTIDsExtra(t *testing.T) {
	t.Run("empty Doc, zero iterations", func(t *testing.T) {
		d := &doc.Doc{AST: &ast.Document{}}
		tt, err := assignTIDs(d)
		if err != nil {
			t.Fatalf("assignTIDs: %v", err)
		}
		if tt.primaryEnd != 1 {
			t.Errorf("primaryEnd = %d, want 1", tt.primaryEnd)
		}
		if len(tt.byRef) != 0 {
			t.Errorf("byRef should be empty, got %v", tt.byRef)
		}
	})

	t.Run("ResolveStmt rejects a malformed StmtRef", func(t *testing.T) {
		d := &doc.Doc{
			AST: &ast.Document{Subject: "sys"},
			Blocks: []doc.Block{
				{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: codebook.SectionSubject, Path: []int{0}}},
			},
		}
		if _, err := assignTIDs(d); err == nil {
			t.Fatalf("want error from ResolveStmt on malformed Subject Path")
		}
	})
}
