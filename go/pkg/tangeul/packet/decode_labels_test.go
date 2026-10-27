//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestDecodeLabels — decodeLabels inverts encodeLabels and rejects a missing count, a bad record and a non-text Value
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestDecodeLabels round-trips an empty and a mixed even/odd table through
// encodeLabels, and rejects no count word, a count past the end, and a Value
// that is not UTF8-TEXT.
func TestDecodeLabels(t *testing.T) {
	for _, want := range [][]string{{}, {"ab", "c"}} {
		ws, _ := encodeLabels(want)
		got, n, err := decodeLabels(ws)
		if err != nil || n != len(ws) || !reflect.DeepEqual(got, want) {
			t.Fatalf("%v: got %v %d %v", want, got, n, err)
		}
	}
	if _, _, err := decodeLabels(nil); err == nil {
		t.Fatal("missing count word: want error")
	}
	if _, _, err := decodeLabels([]core.Word{1}); err == nil {
		t.Fatal("count past the end: want error")
	}
	if _, _, err := decodeLabels([]core.Word{1, 0x0882, 7}); err == nil {
		t.Fatal("non-text Value: want error")
	}
}
