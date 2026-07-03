//ff:func feature=tangl type=parser control=sequence
//ff:what TestItemRef — Subject is assigned by section, everything else looks up refByLine
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestItemRef covers all three outcomes: SectionSubject returns {Subject,nil}
// without consulting refByLine (even when it is empty), a line present in
// refByLine returns that ref, and a line absent from refByLine is a drift error.
func TestItemRef(t *testing.T) {
	ref, err := itemRef(codebook.SectionSubject, 1, map[int]doc.StmtRef{})
	if err != nil || ref.Sec != codebook.SectionSubject || len(ref.Path) != 0 {
		t.Errorf("Subject: got (%+v, %v), want ({Subject,nil}, nil)", ref, err)
	}

	refByLine := map[int]doc.StmtRef{5: {Sec: codebook.SectionCases, Path: []int{0}}}
	ref, err = itemRef(codebook.SectionCases, 5, refByLine)
	if err != nil || ref.Sec != codebook.SectionCases {
		t.Errorf("hit: got (%+v, %v), want ({Cases,[0]}, nil)", ref, err)
	}

	if _, err := itemRef(codebook.SectionCases, 99, refByLine); err == nil {
		t.Error("itemRef accepted a line missing from refByLine")
	}
}
