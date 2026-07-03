//ff:func feature=tangl type=parser control=sequence
//ff:what TestAssignBodySlots — assign per-kind run/check slot indices, skipping struck and non-body children
package en

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// TestAssignBodySlots covers: a `require` child is not a run/check line and is
// skipped (isBody false), a struck run child is skipped without even being
// classified, and interleaved run/check children get independently advancing
// per-kind slot indices. It then covers the two error paths: a malformed
// backtick propagates parseRunCheck's error, and a name disagreeing with the
// AST slice is a verifyName drift error.
func TestAssignBodySlots(t *testing.T) {
	items, _ := scan.ScanItems(scan.ScanLines([]byte("- provides `ep`\n" +
		"  - `q` is required\n" +
		"  - ~~run `dead`~~\n" +
		"  - run `a`\n" +
		"  - check `b`\n" +
		"  - run `c`\n")))
	if len(items) != 1 {
		t.Fatalf("ScanItems top-level = %d, want 1", len(items))
	}
	children := items[0].Children

	refByLine := map[int]doc.StmtRef{}
	runs := []string{"a", "c"}
	checks := []string{"b"}
	if err := assignBodySlots(children, codebook.SectionProvides, 0, runs, checks, refByLine); err != nil {
		t.Fatalf("assignBodySlots: %v", err)
	}
	want := map[int]doc.StmtRef{
		4: {Sec: codebook.SectionProvides, Path: []int{0, 1, 0}},
		5: {Sec: codebook.SectionProvides, Path: []int{0, 2, 0}},
		6: {Sec: codebook.SectionProvides, Path: []int{0, 1, 1}},
	}
	if !reflect.DeepEqual(refByLine, want) {
		t.Fatalf("refByLine = %v, want %v", refByLine, want)
	}

	malformedItems, _ := scan.ScanItems(scan.ScanLines([]byte("- provides `ep`\n  - run oops\n")))
	if err := assignBodySlots(malformedItems[0].Children, codebook.SectionProvides, 0, nil, nil, map[int]doc.StmtRef{}); err == nil {
		t.Fatal("assignBodySlots accepted a malformed run line")
	}

	mismatchedItems, _ := scan.ScanItems(scan.ScanLines([]byte("- provides `ep`\n  - run `wrong`\n")))
	if err := assignBodySlots(mismatchedItems[0].Children, codebook.SectionProvides, 0, []string{"right"}, nil, map[int]doc.StmtRef{}); err == nil {
		t.Fatal("assignBodySlots accepted a run name that disagrees with the AST")
	}
}
