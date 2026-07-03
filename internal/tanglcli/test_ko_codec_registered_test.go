//go:build !fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what TestKOCodecRegistered — registerCodecs wires the ko codec into the doc registry (Phase 007 §1)
package tanglcli

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestKOCodecRegistered pins the Phase 007 §1 registration wiring: after
// registerCodecs (idempotent, also called by NewRootCmd), doc.LookupCodec finds a
// codec for locale "ko", and that codec reports its locale as "ko". This is the
// link `tangl encode` relies on once detectLocale returns "ko" for a Korean
// surface. It is a default-build test — the `fakecodec` build deliberately leaves
// "ko" unregistered.
func TestKOCodecRegistered(t *testing.T) {
	registerCodecs()
	c, ok := doc.LookupCodec("ko")
	if !ok {
		t.Fatal(`no codec registered for locale "ko" after registerCodecs`)
	}
	if c.Locale() != "ko" {
		t.Errorf("registered codec Locale() = %q, want ko", c.Locale())
	}
}
