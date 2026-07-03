//ff:func feature=cli type=command control=sequence
//ff:what registerCodecs — assembly point registering the available surface codecs (once)
package tanglcli

import "sync"

// registerCodecsOnce guards registerCodecs so NewRootCmd may be built repeatedly
// (each test does) without doc.RegisterCodec panicking on a duplicate locale.
var registerCodecsOnce sync.Once

// registerCodecs is the CLI assembly point for surface codecs (README contract
// §3 — explicit registration, no init()). The concrete set of codecs is chosen by
// registerSurfaceCodecs, which has two build-tagged variants: the default build
// wires both the EN (plans/tangl-en) and KO (plans/tangl-ko) codecs, while the
// `fakecodec` test build wires EN only and leaves locale "ko" free for the
// fake-codec CLI tests to occupy (they exercise runEncode's post-lookup branches,
// which need a mutable stub the real KO codec cannot provide). The sync.Once bounds
// registration to one pass so repeated NewRootCmd builds do not panic on a
// duplicate locale.
func registerCodecs() {
	registerCodecsOnce.Do(registerSurfaceCodecs)
}
