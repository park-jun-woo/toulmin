//go:build fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what registerSurfaceCodecs — wire EN only, leaving "ko" free for fake-codec tests (fakecodec build)
package tanglcli

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/en"
)

// registerSurfaceCodecs registers only the EN codec under the `fakecodec` build
// tag, deliberately leaving locale "ko" unregistered so the fake-codec CLI tests
// (run_encode_with_codec_test.go, run_encode_no_codec_test.go) can occupy "ko"
// with a mutable stub to drive runEncode's Encode-error / VerifyTiling-error /
// no-codec branches — behaviour the real KO codec cannot force. The default build
// (register_surface_codecs.go) wires the real KO codec instead.
func registerSurfaceCodecs() {
	doc.RegisterCodec(en.ENCodec{})
	// KO left unregistered so fake-codec tests may occupy locale "ko".
}
