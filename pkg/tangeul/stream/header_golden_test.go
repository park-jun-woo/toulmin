//ff:func feature=tangl type=util control=sequence
//ff:what TestHeaderGolden — the STREAM_HEADER/SEGMENT_BOUNDARY words pin to the v0.3 version stamp
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestHeaderGolden pins the leading control words to the codebook v0.3 version
// stamp: a Control frame (Len 6) carrying STREAM_HEADER, GEUL 0.3.0, TANGEUL
// profile 0x0001, codebook 0.3 and the authoring locale code (here
// LocaleUnknown = 0x0000), then a Control frame (Len 1) carrying
// SEGMENT_BOUNDARY. The codebook minor moved 2→3 with the Phase003 EXPR/QuantOp
// additions (§design 4); only this stamp word changes and the reader tolerates
// the minor drift (checkVersion gates on major only), so old streams still load.
func TestHeaderGolden(t *testing.T) {
	want := []word.Word{
		0x1C06, 0x0000, 0x0003, 0x0000, 0x0001, 0x0003, 0x0000,
		0x1C01, 0x0001,
	}
	if got := emitHeader(codebook.LocaleUnknown); !reflect.DeepEqual(got, want) {
		t.Errorf("emitHeader() =\n %04X\nwant\n %04X", got, want)
	}
}
