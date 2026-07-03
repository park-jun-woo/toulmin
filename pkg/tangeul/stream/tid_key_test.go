//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestTidKey — tidKey renders the section code plus zero, one and multiple path indices into a colon-joined key
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestTidKey exercises the for-range over ref.Path directly: a nil path (zero
// iterations, key is just the section code), a single-index path (one
// iteration) and a multi-index path (multiple iterations), checking the
// colon-separated result each time.
func TestTidKey(t *testing.T) {
	cases := []struct {
		name string
		ref  doc.StmtRef
		want string
	}{
		{"nil path", doc.StmtRef{Sec: codebook.SectionSubject, Path: nil}, "1"},
		{"empty path", doc.StmtRef{Sec: codebook.SectionSee, Path: []int{}}, "2"},
		{"one index", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{3}}, "5:3"},
		{"multi index", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{3, 1, 2}}, "5:3:1:2"},
	}
	for _, c := range cases {
		if got := tidKey(c.ref); got != c.want {
			t.Errorf("%s: tidKey = %q, want %q", c.name, got, c.want)
		}
	}
}
