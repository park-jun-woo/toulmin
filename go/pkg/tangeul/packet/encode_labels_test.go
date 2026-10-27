//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeLabels — encodeLabels writes a count word then one UTF8-TEXT Value per label
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestEncodeLabels checks the empty table and an even ("ab", codec 0x40) plus an
// odd ("c", codec 0x41 with a 0x00 pad byte) label.
func TestEncodeLabels(t *testing.T) {
	if got, err := encodeLabels(nil); err != nil || len(got) != 1 || got[0] != 0 {
		t.Fatalf("empty table: got %v %v", got, err)
	}
	got, err := encodeLabels([]string{"ab", "c"})
	want := []core.Word{2, 0x08C0, 0x6162, 0x08C1, 0x6300}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %04X %v, want %04X", got, err, want)
	}
}
