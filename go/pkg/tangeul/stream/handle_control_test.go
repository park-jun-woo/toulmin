//ff:func feature=tangl type=util control=sequence
//ff:what TestHandleControl — decode errors, the pending-heading guard, every switch case, and the unexpected-control default
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestHandleControl walks every branch of handleControl: a DecodeControl
// failure, the pending-SECTION-heading guard tripping on a non-PROSE control,
// each of the four handled controls, and a repeated Edge Id.
func TestHandleControl(t *testing.T) {
	t.Run("decodeError", func(t *testing.T) {
		d := newDec(testWire)
		pending := -1
		good := hcControlRaw(t, toulmin.ControlPacket{Code: toulmin.ControlCaseBegin, Arg: 1})
		bad := good
		bad.Rec.Payload = bad.Rec.Payload[:len(bad.Rec.Payload)-1]
		if err := handleControl(d, &pending, bad); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("pendingSectionGuardTrips", func(t *testing.T) {
		d := newDec(testWire)
		pending := int(toulmin.SectionSubject)
		rp := hcControlRaw(t, toulmin.ControlPacket{Code: toulmin.ControlCaseBegin, Arg: 1})
		if err := handleControl(d, &pending, rp); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("section", func(t *testing.T) {
		d := newDec(testWire)
		pending := -1
		rp := hcControlRaw(t, toulmin.ControlPacket{Code: toulmin.ControlSection, Arg: uint64(toulmin.SectionSee)})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if pending != int(toulmin.SectionSee) {
			t.Errorf("pending = %d, want %d", pending, toulmin.SectionSee)
		}
	})
	t.Run("proseStandalone", func(t *testing.T) {
		d := newDec(testWire)
		pending := -1
		span := toulmin.Span{Off: 0, Len: 4}
		rp := hcControlRaw(t, toulmin.ControlPacket{Code: toulmin.ControlProse, Span: &span})
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
		d := newDec(testWire)
		pending := int(toulmin.SectionSubject)
		span := toulmin.Span{Off: 0, Len: 4}
		rp := hcControlRaw(t, toulmin.ControlPacket{Code: toulmin.ControlProse, Span: &span})
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
		d := newDec(testWire)
		d.curCase = 5
		pending := -1
		rp := hcControlRaw(t, toulmin.ControlPacket{Code: toulmin.ControlCaseBegin, Arg: 9})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if d.curCase != 5 {
			t.Errorf("curCase = %d, want unchanged 5", d.curCase)
		}
	})
	t.Run("caseEndClearsCurCase", func(t *testing.T) {
		d := newDec(testWire)
		d.curCase = 5
		pending := -1
		rp := hcControlRaw(t, toulmin.ControlPacket{Code: toulmin.ControlCaseEnd, Arg: 9})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("handleControl: %v", err)
		}
		if d.curCase != -1 {
			t.Errorf("curCase = %d, want -1", d.curCase)
		}
	})
	t.Run("repeatedEdgeId", func(t *testing.T) {
		d := newDec(testWire)
		pending := -1
		rp := hcControlRaw(t, toulmin.ControlPacket{ID: 4, Code: toulmin.ControlCaseEnd, Arg: 1})
		if err := handleControl(d, &pending, rp); err != nil {
			t.Fatalf("first: %v", err)
		}
		if err := handleControl(d, &pending, rp); err == nil {
			t.Fatalf("expected error for a repeated Edge Id, got nil")
		}
	})
}
