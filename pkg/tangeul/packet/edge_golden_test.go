//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEdgeGolden — the three §6.2 edges pin to their exact PropCode/src/dst/flags/SPAN words
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEdgeGolden pins the three §6.2 edges to their exact word sequences
// (PropCode, src, dst, flags, SPAN). Spans use TID=1 with zero offset/length so
// the structural words are the assertion. §6.2's frame Len annotations are not
// binding (they are inconsistent in the spec); the computed Len is used.
func TestEdgeGolden(t *testing.T) {
	span := word.Span{TID: 1}
	cases := []struct {
		name string
		edge EdgePacket
		want []word.Word
	}{
		{
			"attacks",
			EdgePacket{Prop: codebook.PropAttacks, Src: 0x0014, DstTID: 0x0013, Span: span},
			[]word.Word{0x1407, 0x0101, 0x0014, 0x0013, 0x0000, 0x0001, 0x0000, 0x0000},
		},
		{
			"do-once",
			EdgePacket{Prop: codebook.PropDo, Src: 0x0013, DstFunc: &FuncRef{Alias: 0x0002, Verb: 0xC200}, Flags: codebook.EdgeFlagOnce, Span: span},
			[]word.Word{0x1408, 0x0102, 0x0013, 0x0002, 0xC200, 0x0001, 0x0001, 0x0000, 0x0000},
		},
		{
			"run",
			EdgePacket{Prop: codebook.PropRun, Src: 0x0013, DstTID: 0x0020, Span: span},
			[]word.Word{0x1407, 0x0104, 0x0013, 0x0020, 0x0000, 0x0001, 0x0000, 0x0000},
		},
	}
	for _, c := range cases {
		ws, err := EncodeEdge(c.edge)
		if err != nil {
			t.Fatalf("%s: EncodeEdge: %v", c.name, err)
		}
		if !reflect.DeepEqual(ws, c.want) {
			t.Fatalf("%s: EncodeEdge =\n %04X\nwant\n %04X", c.name, ws, c.want)
		}
	}
}
