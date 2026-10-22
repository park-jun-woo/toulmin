//ff:func feature=tangl type=parser control=sequence
//ff:what TestMatchSectionBodies — resolve a section's live heads to their AST endpoint and assign body slots
package en

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// TestMatchSectionBodies covers: a struck head is skipped before it is ever
// matched against the AST, a live head that resolves cleanly has its body slots
// assigned, a live head with no matching AST line is a drift error, and an
// assignBodySlots failure (a run name disagreeing with the AST) propagates.
func TestMatchSectionBodies(t *testing.T) {
	src := []byte("## tangl:Provides\n" +
		"- ~~provides `dead`~~\n" +
		"- provides `ep`\n" +
		"  - run `a`\n")
	_, secs := scan.ScanSections(scan.ScanLines(src))
	if len(secs) != 1 {
		t.Fatalf("ScanSections = %d sections, want 1", len(secs))
	}

	a := &ast.Document{Provides: []ast.Endpoint{{Name: "ep", Line: 3, Runs: []string{"a"}}}}
	refByLine := map[int]doc.StmtRef{}
	if err := matchSectionBodies(a, toulmin.SectionProvides, secs[0], refByLine); err != nil {
		t.Fatalf("matchSectionBodies: %v", err)
	}
	want := map[int]doc.StmtRef{4: {Sec: toulmin.SectionProvides, Path: []int{0, 1, 0}}}
	if !reflect.DeepEqual(refByLine, want) {
		t.Fatalf("refByLine = %v, want %v (the struck head at line 2 must not appear)", refByLine, want)
	}

	// No AST endpoint declared at line 3: bodySlices misses, ok=false.
	if err := matchSectionBodies(&ast.Document{}, toulmin.SectionProvides, secs[0], map[int]doc.StmtRef{}); err == nil {
		t.Fatal("matchSectionBodies accepted a head with no matching AST entry")
	}

	// assignBodySlots itself fails: the AST Runs disagree with the scanned name.
	bad := &ast.Document{Provides: []ast.Endpoint{{Name: "ep", Line: 3, Runs: []string{"wrong"}}}}
	if err := matchSectionBodies(bad, toulmin.SectionProvides, secs[0], map[int]doc.StmtRef{}); err == nil {
		t.Fatal("matchSectionBodies accepted a run name that disagrees with the AST")
	}

	// An empty section has no top-level items, so the loop body never runs.
	if err := matchSectionBodies(&ast.Document{}, toulmin.SectionProvides, scan.Section{}, map[int]doc.StmtRef{}); err != nil {
		t.Fatalf("matchSectionBodies(empty section): %v", err)
	}
}
