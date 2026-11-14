//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestVerifyTiling — VerifyTiling accepts a clean tiling and distinguishes gap/overlap/reversal/trailing
package doc

import (
	"strings"
	"testing"
)

// TestVerifyTiling checks that a contiguous tiling passes and that each violation
// family (leading/interior gap, overlap, order reversal, trailing gap) is caught
// with a distinct message, and that an empty Source with no Blocks is valid.
func TestVerifyTiling(t *testing.T) {
	blk := func(off, ln uint32) Block { return Block{Span: Span{Off: off, Len: ln}} }

	cases := []struct {
		name   string
		srcLen int
		blocks []Block
		want   string // "" means no error
	}{
		{"normal", 10, []Block{blk(0, 4), blk(4, 6)}, ""},
		{"empty", 0, nil, ""},
		{"gap", 10, []Block{blk(0, 4), blk(5, 5)}, "gap"},
		{"overlap", 10, []Block{blk(0, 5), blk(3, 7)}, "overlap"},
		{"reversal", 8, []Block{blk(0, 4), blk(4, 4), blk(2, 2)}, "order reversal"},
		{"trailing", 10, []Block{blk(0, 4)}, "trailing gap"},
	}

	for _, tc := range cases {
		d := &Doc{Source: make([]byte, tc.srcLen), Blocks: tc.blocks}
		err := VerifyTiling(d)
		if tc.want == "" && err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
		}
		if tc.want != "" && err == nil {
			t.Errorf("%s: expected error containing %q, got nil", tc.name, tc.want)
		}
		if tc.want != "" && err != nil && !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: error %q does not contain %q", tc.name, err.Error(), tc.want)
		}
	}

	if err := VerifyTiling(nil); err == nil {
		t.Errorf("nil Doc: expected error")
	}
}
