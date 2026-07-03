//ff:func feature=tangl type=util control=sequence
//ff:what TestStepRead — control application, its error, pending-section-before-statement error, and statement collection all advance/report correctly
package stream

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestStepRead drives stepRead over the real packet stream of the §10
// account-transfer fixture (packets[3] SECTION control, packets[4] its PROSE
// heading, packets[5..7] the Subject statement) to cover all four branches: a
// control that applies cleanly and advances by one, a control whose
// handleControl call fails (a second SECTION while one is already pending), a
// non-control packet arriving while a SECTION is still pending (stepRead's own
// hard error), and a non-control packet with no section pending, which is
// collected and interpreted as a statement.
func TestStepRead(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	ws, err := word.ReadWords(b)
	if err != nil {
		t.Fatalf("ReadWords: %v", err)
	}
	packets, err := readPackets(ws)
	if err != nil {
		t.Fatalf("readPackets: %v", err)
	}
	// Sanity on the fixture shape this test depends on.
	if packets[3].Class != word.ClassControl || packets[4].Class != word.ClassControl || packets[5].Class != word.ClassPayload {
		t.Fatalf("unexpected fixture packet classes: %d %d %d", packets[3].Class, packets[4].Class, packets[5].Class)
	}

	t.Run("control applies and advances", func(t *testing.T) {
		d := newDec()
		pendingSection := -1
		next, err := stepRead(d, packets, 3, &pendingSection)
		if err != nil {
			t.Fatalf("stepRead: %v", err)
		}
		if next != 4 {
			t.Errorf("next = %d, want 4", next)
		}
		if pendingSection != 1 {
			t.Errorf("pendingSection = %d, want 1 (SECTION armed)", pendingSection)
		}
	})

	t.Run("control error propagates", func(t *testing.T) {
		d := newDec()
		pendingSection := 0 // a section is already pending
		next, err := stepRead(d, packets, 3, &pendingSection)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if next != 0 {
			t.Errorf("next = %d, want 0 on error", next)
		}
	})

	t.Run("statement while section pending is a hard error", func(t *testing.T) {
		d := newDec()
		pendingSection := 0 // pretend a SECTION heading is still owed
		next, err := stepRead(d, packets, 5, &pendingSection)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "SECTION control must be followed") {
			t.Errorf("err = %v, want SECTION-pending message", err)
		}
		if next != 0 {
			t.Errorf("next = %d, want 0 on error", next)
		}
	})

	t.Run("statement collected when nothing pending", func(t *testing.T) {
		d := newDec()
		pendingSection := -1
		// Feed the SECTION then its PROSE heading first, matching the real
		// reader loop, so the Subject statement at index 5 sees pendingSection
		// cleared exactly as ReadDoc would leave it.
		next, err := stepRead(d, packets, 3, &pendingSection)
		if err != nil {
			t.Fatalf("stepRead(SECTION): %v", err)
		}
		next, err = stepRead(d, packets, next, &pendingSection)
		if err != nil {
			t.Fatalf("stepRead(PROSE): %v", err)
		}
		if pendingSection != -1 {
			t.Fatalf("pendingSection = %d, want -1 after heading", pendingSection)
		}
		next, err = stepRead(d, packets, next, &pendingSection)
		if err != nil {
			t.Fatalf("stepRead(statement): %v", err)
		}
		if next != 8 {
			t.Errorf("next = %d, want 8 past the Subject statement", next)
		}
		if len(d.blocks) == 0 {
			t.Errorf("expected a Block to be appended for the statement")
		}
	})
}
