//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestCollectAnchors — Subject is assigned by position and a ref-less tangl item is a drift error
package en

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// TestCollectAnchors covers the two special paths: a Subject item is anchored to
// {Subject,nil} by section position even with an empty ref index, while a non-struck
// item in any other section with no ref is a scanner/parser drift reported as error.
// It also covers an unknown `## tangl:` section name, which is rejected before
// sectionAnchors is ever consulted.
func TestCollectAnchors(t *testing.T) {
	subj := []byte("## tangl:Subject\n- this document is `x`\n")
	_, secs := scan.ScanSections(scan.ScanLines(subj))
	anchors, err := collectAnchors(secs, map[int]doc.StmtRef{})
	if err != nil {
		t.Fatalf("collectAnchors(Subject): %v", err)
	}
	var stmt *doc.Block
	for i := range anchors {
		if anchors[i].Kind == doc.StmtBlock {
			stmt = &anchors[i]
		}
	}
	if stmt == nil {
		t.Fatal("no StmtBlock produced for the Subject declaration")
	}
	if stmt.Ref.Sec != toulmin.SectionSubject || len(stmt.Ref.Path) != 0 {
		t.Errorf("Subject ref = %+v, want {Subject, nil}", stmt.Ref)
	}

	cases := []byte("## tangl:Cases\n- in case of `c`\n")
	_, secs = scan.ScanSections(scan.ScanLines(cases))
	if _, err := collectAnchors(secs, map[int]doc.StmtRef{}); err == nil {
		t.Fatal("collectAnchors accepted a non-struck Cases item with no AST ref")
	}

	unknown := []byte("## tangl:Bogus\n- this document is `x`\n")
	_, secs = scan.ScanSections(scan.ScanLines(unknown))
	if _, err := collectAnchors(secs, map[int]doc.StmtRef{}); err == nil {
		t.Fatal("collectAnchors accepted an unknown tangl section name")
	}
}
