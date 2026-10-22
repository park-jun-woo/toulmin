//go:build !fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what registerSurfaceCodecs — wire the real EN and KO surface codecs (default build)
package tanglcli

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/en"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/ko"
)

// registerSurfaceCodecs registers the production surface codecs — EN
// (plans/tangl-en) and KO (plans/tangl-ko) — so `tangl encode` on an English or
// Korean .md and `tangl decode --locale en|ko` both find a codec. This is the
// default build; the `fakecodec` build overrides it to leave "ko" free for the
// fake-codec CLI tests.
func registerSurfaceCodecs() {
	doc.RegisterCodec(en.ENCodec{})
	doc.RegisterCodec(ko.KoCodec{})
}
