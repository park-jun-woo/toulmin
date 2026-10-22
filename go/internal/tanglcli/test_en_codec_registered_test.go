//ff:func feature=cli type=command control=sequence
//ff:what TestENCodecRegistered — registerCodecs wires the en codec into the doc registry (Phase 004 §1)
package tanglcli

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestENCodecRegistered pins the Phase 004 §1 registration wiring: after
// registerCodecs (idempotent, also called by NewRootCmd), doc.LookupCodec finds a
// codec for locale "en", and that codec reports its locale as "en". This is the
// link `tangl encode` relies on once detectLocale returns "en".
func TestENCodecRegistered(t *testing.T) {
	registerCodecs()
	c, ok := doc.LookupCodec("en")
	if !ok {
		t.Fatal(`no codec registered for locale "en" after registerCodecs`)
	}
	if c.Locale() != "en" {
		t.Errorf("registered codec Locale() = %q, want en", c.Locale())
	}
}
