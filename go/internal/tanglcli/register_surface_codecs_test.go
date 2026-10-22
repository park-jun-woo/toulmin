//go:build !fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what TestRegisterSurfaceCodecs — registerSurfaceCodecs wires both EN and KO codecs into the doc registry (default build)
package tanglcli

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestRegisterSurfaceCodecs pins the default-build registerSurfaceCodecs wiring
// directly: both statements in its body (README contract §3) each register one
// locale, and after they run doc.LookupCodec must find both "en" and "ko". The
// entry point exercised is registerCodecs (its sync.Once wraps registerSurfaceCodecs)
// rather than registerSurfaceCodecs itself, because doc.RegisterCodec panics on a
// duplicate locale and registerSurfaceCodecs has already run by the time any other
// test in this package calls registerCodecs/NewRootCmd — registerCodecs is the safe,
// idempotent entry point used everywhere else in this package for the same reason.
func TestRegisterSurfaceCodecs(t *testing.T) {
	registerCodecs()

	enCodec, ok := doc.LookupCodec("en")
	if !ok {
		t.Fatal(`no codec registered for locale "en" after registerCodecs`)
	}
	if enCodec.Locale() != "en" {
		t.Errorf("en codec Locale() = %q, want en", enCodec.Locale())
	}

	koCodec, ok := doc.LookupCodec("ko")
	if !ok {
		t.Fatal(`no codec registered for locale "ko" after registerCodecs`)
	}
	if koCodec.Locale() != "ko" {
		t.Errorf("ko codec Locale() = %q, want ko", koCodec.Locale())
	}
}
