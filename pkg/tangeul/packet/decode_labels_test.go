//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestDecodeLabels — decodeLabels is encodeLabels' exact inverse and rejects a missing count/truncated length/truncated bytes
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestDecodeLabels checks decodeLabels round-trips an empty table and a
// mixed even/odd-length label table through encodeLabels, and rejects an
// input with no count word, a count whose length word is truncated, and a
// count whose label bytes are truncated.
func TestDecodeLabels(t *testing.T) {
	for _, want := range [][]string{nil, {"ab", "c"}} {
		ws := encodeLabels(want)
		got, n, err := decodeLabels(ws)
		if err != nil || n != len(ws) {
			t.Fatalf("decodeLabels(%v): got %v, %d, %v", want, got, n, err)
		}
		if len(want) == 0 {
			want = []string{}
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("decodeLabels(%v) = %v, want %v", ws, got, want)
		}
	}

	if _, _, err := decodeLabels(nil); err == nil {
		t.Fatalf("missing count word: want error")
	}
	if _, _, err := decodeLabels([]word.Word{1}); err == nil {
		t.Fatalf("truncated label length: want error")
	}
	if _, _, err := decodeLabels([]word.Word{1, 0, 4}); err == nil {
		t.Fatalf("truncated label bytes: want error")
	}
}
