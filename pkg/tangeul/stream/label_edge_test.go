//ff:func feature=tangl type=util control=sequence
//ff:what TestLabelEdge — labelEdge encodes a PropLabel edge from src to a Text TID that decodes back intact
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestLabelEdge encodes a label edge and checks it decodes back to the same
// PropLabel/src/dst/span.
func TestLabelEdge(t *testing.T) {
	span := word.Span{Off: 3, Len: 5}
	got, err := labelEdge(7, 12, span)
	if err != nil {
		t.Fatalf("labelEdge: %v", err)
	}
	ep, n, err := packet.DecodeEdge(got)
	if err != nil {
		t.Fatalf("DecodeEdge: %v", err)
	}
	if n != len(got) {
		t.Errorf("consumed %d words, want %d", n, len(got))
	}
	if ep.Prop != codebook.PropLabel || ep.Src != 7 || ep.DstTID != 12 || ep.Span != span {
		t.Errorf("got %+v, want Prop=PropLabel Src=7 DstTID=12 Span=%v", ep, span)
	}
}
