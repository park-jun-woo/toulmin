//ff:func feature=cli type=command control=sequence
//ff:what TestRegisterCodecs — registerCodecs is idempotent (safe across repeated NewRootCmd builds)
package tanglcli

// TestRegisterCodecs asserts registerCodecs can be called repeatedly without
// panicking on a duplicate doc.RegisterCodec — the guard that lets each test
// build NewRootCmd afresh. It now registers the en codec, so the sync.Once
// bounding registration to one pass is what keeps repeated calls safe.

import "testing"

func TestRegisterCodecs(t *testing.T) {
	registerCodecs()
	registerCodecs()
	registerCodecs()
}
