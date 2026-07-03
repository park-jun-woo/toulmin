//ff:func feature=tangl type=util control=sequence
//ff:what TestHeaderGolden — the STREAM_HEADER/SEGMENT_BOUNDARY words pin to the v0.2 version stamp
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestHeaderGolden pins the leading control words to the codebook v0.2 version
// stamp: a Control frame (Len 6) carrying STREAM_HEADER, GEUL 0.3.0, TANGL
// profile 0x0001, codebook 0.2 and the authoring locale code (here
// LocaleUnknown = 0x0000), then a Control frame (Len 1) carrying
// SEGMENT_BOUNDARY. A change here is a wire-format change the reader's version
// check must move with.
func TestHeaderGolden(t *testing.T) {
	want := []word.Word{
		0x1C06, 0x0000, 0x0003, 0x0000, 0x0001, 0x0002, 0x0000,
		0x1C01, 0x0001,
	}
	if got := emitHeader(codebook.LocaleUnknown); !reflect.DeepEqual(got, want) {
		t.Errorf("emitHeader() =\n %04X\nwant\n %04X", got, want)
	}
}
