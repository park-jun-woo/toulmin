//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAddSpanBlock — a pending SECTION yields a SectionBlock and clears; otherwise a ProseBlock
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestAddSpanBlock covers both branches of addSpanBlock: a non-negative
// pendingSection consumes the pending SECTION into a SectionBlock and resets
// pendingSection to -1, while a negative pendingSection appends a standalone
// ProseBlock and leaves pendingSection untouched.
func TestAddSpanBlock(t *testing.T) {
	cases := []struct {
		name           string
		pendingSection int
		wantKind       doc.BlockKind
		wantPending    int
	}{
		{"pending section consumed", int(toulmin.SectionSee), doc.SectionBlock, -1},
		{"no pending section", -1, doc.ProseBlock, -1},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		d := &dec{}
		ps := c.pendingSection
		span := doc.Span{Off: 3, Len: 5}

		if err := addSpanBlock(d, &ps, toulmin.Span{Off: 3, Len: 5}); err != nil {
			t.Fatal(err)
		}

		if len(d.blocks) != 1 {
			t.Fatalf("blocks len = %d, want 1", len(d.blocks))
		}
		got := d.blocks[0]
		if got.Kind != c.wantKind {
			t.Errorf("Kind = %v, want %v", got.Kind, c.wantKind)
		}
		if got.Span != span {
			t.Errorf("Span = %+v, want %+v", got.Span, span)
		}
		if ps != c.wantPending {
			t.Errorf("pendingSection = %d, want %d", ps, c.wantPending)
		}
		if c.wantKind == doc.SectionBlock {
			want := toulmin.SectionCode(c.pendingSection)
			if got.Ref.Sec != want {
				t.Errorf("Ref.Sec = %v, want %v", got.Ref.Sec, want)
			}
		} else if !reflect.DeepEqual(got.Ref, doc.StmtRef{}) {
			t.Errorf("Ref = %+v, want zero value", got.Ref)
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}
	ps := -1
	if err := addSpanBlock(&dec{}, &ps, toulmin.Span{Off: 1 << 32}); err == nil {
		t.Error("span past 32 bits: want error")
	}
}
