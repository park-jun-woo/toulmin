//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestStepRead — stepRead applies controls, rejects a second SOURCE and a statement under a pending SECTION
package stream

import (
	"strings"
	"testing"
)

// TestStepRead walks an encoded document's BODY: index 1 is SECTION, 2 its
// heading PROSE, 3 the Subject statement. It checks control application, a
// control error, a statement while a SECTION heading is owed, a second SOURCE,
// and a statement read once nothing is pending.
func TestStepRead(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	d, packets := docPackets(t, b)
	if packets[1].Class != classControl || packets[2].Class != classControl || packets[3].Class != classPayload {
		t.Fatalf("unexpected fixture classes: %d %d %d", packets[1].Class, packets[2].Class, packets[3].Class)
	}
	pending := 0
	if _, err := stepRead(d, packets, 1, &pending); err == nil {
		t.Fatal("control while a SECTION is pending: want error")
	}
	pending = 0
	if _, err := stepRead(d, packets, 3, &pending); err == nil || !strings.Contains(err.Error(), "SECTION control must be followed") {
		t.Fatalf("statement while SECTION pending: %v", err)
	}
	if _, err := stepRead(d, packets, 0, &pending); err == nil {
		t.Fatal("second SOURCE: want error")
	}
	d, packets = docPackets(t, b)
	pending = -1
	next := 1
	for _, want := range []int{2, 3, 6} {
		if next, err = stepRead(d, packets, next, &pending); err != nil || next != want {
			t.Fatalf("step to %d: got %d %v", want, next, err)
		}
	}
	if len(d.blocks) != 2 || pending != -1 {
		t.Fatalf("blocks=%d pending=%d, want SECTION heading + Subject", len(d.blocks), pending)
	}
}
