//ff:func feature=tangl type=util control=sequence
//ff:what TestHandleControl — decode errors, the pending-heading guard, every switch case, and the unexpected-control default
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestHandleControl walks every branch of handleControl: a DecodeControl
// failure, the pending-SECTION-heading guard tripping on a non-PROSE control,
// each of the four handled Control codes, and the default unexpected-control
// error.
func TestHandleControl(t *testing.T) {
	t.Run("decodeError", func(t *testing.T) {
		d := newDec()
		pending := -1
		good := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlCaseBegin, Arg: 1})
		bad := rawPacket{Class: word.ClassControl, Words: good.Words[:len(good.Words)-1]}
		if err := handleControl(d, &pending, bad); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("pendingSectionGuardTrips", func(t *testing.T) {
		d := newDec()
		pending := int(codebook.SectionSubject)
		rp := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlCaseBegin, Arg: 1})
		if err := handleControl(d, &pending, rp); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("section", func(t *testing.T) {
		d := newDec()
		pending := -1
		rp := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlSection, Arg: uint16(codebook.SectionSee)})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if pending != int(codebook.SectionSee) {
			t.Errorf("pending = %d, want %d", pending, codebook.SectionSee)
		}
	})

	t.Run("proseStandalone", func(t *testing.T) {
		d := newDec()
		pending := -1
		span := word.Span{Off: 0, Len: 4}
		rp := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlProse, Span: &span})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if pending != -1 {
			t.Errorf("pending = %d, want -1", pending)
		}
		if len(d.blocks) != 1 {
			t.Fatalf("blocks = %d, want 1", len(d.blocks))
		}
	})

	t.Run("proseSectionHeading", func(t *testing.T) {
		d := newDec()
		pending := int(codebook.SectionSubject)
		span := word.Span{Off: 0, Len: 4}
		rp := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlProse, Span: &span})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if pending != -1 {
			t.Errorf("pending = %d, want -1 (consumed)", pending)
		}
		if len(d.blocks) != 1 {
			t.Fatalf("blocks = %d, want 1", len(d.blocks))
		}
	})

	t.Run("caseBeginNoOp", func(t *testing.T) {
		d := newDec()
		d.curCase = 5
		pending := -1
		rp := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlCaseBegin, Arg: 9})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if d.curCase != 5 {
			t.Errorf("curCase = %d, want unchanged 5", d.curCase)
		}
	})

	t.Run("caseEndClearsCurCase", func(t *testing.T) {
		d := newDec()
		d.curCase = 5
		pending := -1
		rp := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlCaseEnd, Arg: 9})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if d.curCase != -1 {
			t.Errorf("curCase = %d, want -1", d.curCase)
		}
	})

	t.Run("unexpectedControlDefault", func(t *testing.T) {
		d := newDec()
		pending := -1
		rp := hcControlRaw(t, packet.ControlPacket{Code: codebook.ControlSegmentBoundary, Arg: 1})
		if err := handleControl(d, &pending, rp); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
