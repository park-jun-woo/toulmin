//ff:func feature=cli type=command control=sequence
//ff:what TestRunEncode — encode rejects non-.md input, unreadable files, undetectable locales, and reports "no codec registered" while no codec exists
package tanglcli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunEncode covers runEncode's branches that do not need the target locale's
// codec: a non-.md input is rejected, a missing file surfaces the read error, and
// an undetectable locale is propagated from detectLocale. In the default build both
// `en` and `ko` have real codecs (registerCodecs), so the "no codec registered"
// branch is no longer reachable here — it is exercised under the `fakecodec` build
// tag by TestRunEncodeNoCodec (run_encode_no_codec_test.go), where "ko" is left
// unregistered. The codec-found success/error branches likewise move to the
// `fakecodec` build's TestRunEncodeCodec (run_encode_with_codec_test.go).
func TestRunEncode(t *testing.T) {
	writeMD := func(t *testing.T, body string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "policy.md")
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		return p
	}

	t.Run("non-.md input is rejected", func(t *testing.T) {
		cmd := NewEncodeCmd()
		if err := runEncode(cmd, []string{"policy.tangeul"}); err == nil || !strings.Contains(err.Error(), "expects a .md file") {
			t.Fatalf("err = %v, want '.md file' rejection", err)
		}
	})

	t.Run("missing file surfaces read error", func(t *testing.T) {
		cmd := NewEncodeCmd()
		if err := runEncode(cmd, []string{filepath.Join(t.TempDir(), "nope.md")}); err == nil {
			t.Fatal("expected a read error, got nil")
		}
	})

	t.Run("undetectable locale is propagated", func(t *testing.T) {
		p := writeMD(t, "no tangl header here\n")
		cmd := NewEncodeCmd()
		err := runEncode(cmd, []string{p})
		if err == nil || !strings.Contains(err.Error(), "cannot determine locale") {
			t.Fatalf("err = %v, want locale-detection failure", err)
		}
	})
}
