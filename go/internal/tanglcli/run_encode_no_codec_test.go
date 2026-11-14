//go:build fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what TestRunEncodeNoCodec — runEncode reports "no codec registered" when the detected locale has none
package tanglcli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunEncodeNoCodec exercises runEncode's "no codec registered for locale"
// branch, only reachable when a detectable locale has no registered codec. Under
// the `fakecodec` build registerSurfaceCodecs wires EN only, leaving "ko" free, so
// a Korean header (detectLocale → "ko") reaches an empty lookup. This file sorts
// before run_encode_with_codec_test.go, whose TestRunEncodeCodec registers a fake
// "ko" codec at its top — so this assertion runs first, while "ko" is still empty.
func TestRunEncodeNoCodec(t *testing.T) {
	p := filepath.Join(t.TempDir(), "policy.md")
	if err := os.WriteFile(p, []byte("## tangl:주제\n- this document is `x`\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	cmd := NewEncodeCmd()
	err := runEncode(cmd, []string{p})
	if err == nil || !strings.Contains(err.Error(), `no codec registered for locale "ko"`) {
		t.Fatalf("err = %v, want 'no codec registered for locale \"ko\"'", err)
	}
}
