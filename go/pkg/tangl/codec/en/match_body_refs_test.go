//ff:func feature=tangl type=parser control=sequence
//ff:what TestMatchBodyRefs — assign interleaved run/check slots and reject a name mismatch
package en

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// TestMatchBodyRefs feeds a Provides body whose run/check lines interleave and
// checks each body line gets the right per-kind slot index, then checks that a
// backtick name disagreeing with the AST slice is reported as a drift error. It
// also feeds a Subject section ahead of the Provides body to cover the skip branch
// (a recognized but non-body section, and an unrecognized section, are both passed
// over without being dispatched to matchSectionBodies).
func TestMatchBodyRefs(t *testing.T) {
	src := []byte("## tangl:Subject\n" +
		"- this document is `x`\n" +
		"## tangl:Bogus\n" +
		"- unrecognized section\n" +
		"## tangl:Provides\n" +
		"- provides `ep`\n" +
		"  - run `a`\n" +
		"  - check `b`\n" +
		"  - run `c`\n")
	_, secs := scan.ScanSections(scan.ScanLines(src))

	a := &ast.Document{Provides: []ast.Endpoint{{
		Name:   "ep",
		Line:   6,
		Runs:   []string{"a", "c"},
		Checks: []string{"b"},
	}}}
	got := map[int]doc.StmtRef{}
	if err := matchBodyRefs(a, secs, got); err != nil {
		t.Fatalf("matchBodyRefs: %v", err)
	}
	want := map[int]doc.StmtRef{
		7: {Sec: toulmin.SectionProvides, Path: []int{0, 1, 0}},
		8: {Sec: toulmin.SectionProvides, Path: []int{0, 2, 0}},
		9: {Sec: toulmin.SectionProvides, Path: []int{0, 1, 1}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("matchBodyRefs slots\n got %v\nwant %v", got, want)
	}

	bad := &ast.Document{Provides: []ast.Endpoint{{
		Name: "ep",
		Line: 6,
		Runs: []string{"wrong", "c"},
	}}}
	if err := matchBodyRefs(bad, secs, map[int]doc.StmtRef{}); err == nil {
		t.Fatal("matchBodyRefs accepted a run name that disagrees with the AST")
	}
}
