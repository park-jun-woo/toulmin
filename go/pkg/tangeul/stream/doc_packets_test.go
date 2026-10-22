//ff:func feature=tangl type=util control=sequence
//ff:what docPackets — test helper: the classified BODY records and reader state of an encoded document
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// docPackets parses b (WriteDoc output) down to its classified BODY records and
// returns them with a fresh reader state at the document's widths, as ReadDoc
// would have them before walking statements.
func docPackets(t *testing.T, b []byte) (*dec, []rawPacket) {
	t.Helper()
	ws, err := core.ReadWords(b)
	if err != nil {
		t.Fatal(err)
	}
	gd, err := core.SplitDocument(ws)
	if err != nil {
		t.Fatal(err)
	}
	h, err := readHead(gd.Head)
	if err != nil {
		t.Fatal(err)
	}
	packets, err := readPackets(gd.Body)
	if err != nil {
		t.Fatal(err)
	}
	_, sw, err := readPreamble(packets, h.idWidth)
	if err != nil {
		t.Fatal(err)
	}
	return newDec(toulmin.Wire{ID: h.idWidth, Span: sw}), packets
}
