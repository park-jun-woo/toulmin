//ff:func feature=cli type=command control=sequence
//ff:what TestRunDecode — decode does mode-1 verbatim by default, writes to -o, propagates load errors, and reports "no codec registered" for mode 2
package tanglcli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestRunDecode covers runDecode's branches: a load error is propagated; the
// default (no --locale) is mode 1, dumping the Source verbatim to stdout; -o
// writes that dump to a file instead; and --locale for an unregistered locale
// fails with "no codec registered for locale". The mode-2 success/error and
// flag-error branches are covered separately by TestRunDecodeCodec
// (run_decode_with_codec_test.go) to keep this func under the control=sequence
// line limit.
func TestRunDecode(t *testing.T) {
	d := buildTransferDoc()
	path := filepath.Join(t.TempDir(), "transfer.tangeul")
	saveTangeul(t, path, d)
	wantVerbatim := doc.RenderVerbatim(d)

	t.Run("load error is propagated", func(t *testing.T) {
		cmd := NewDecodeCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runDecode(cmd, []string{"does-not-exist.tangeul"}); err == nil {
			t.Fatal("expected a load error, got nil")
		}
	})

	t.Run("mode 1 verbatim to stdout", func(t *testing.T) {
		cmd := NewDecodeCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runDecode(cmd, []string{path}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(out.Bytes(), wantVerbatim) {
			t.Errorf("stdout = %q, want verbatim Source %q", out.Bytes(), wantVerbatim)
		}
	})

	t.Run("mode 1 verbatim to -o file", func(t *testing.T) {
		outPath := filepath.Join(t.TempDir(), "out.md")
		cmd := NewDecodeCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := cmd.Flags().Set("out", outPath); err != nil {
			t.Fatalf("Set(out): %v", err)
		}
		if err := runDecode(cmd, []string{path}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out.Len() != 0 {
			t.Errorf("expected no stdout output with -o, got %q", out.String())
		}
		data, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if !bytes.Equal(data, wantVerbatim) {
			t.Errorf("file = %q, want verbatim Source %q", data, wantVerbatim)
		}
	})

	t.Run("mode 2 with no codec registered", func(t *testing.T) {
		// "en" now has a registered codec, so an unregistered tag drives this branch.
		cmd := NewDecodeCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := cmd.Flags().Set("locale", "zz-none"); err != nil {
			t.Fatalf("Set(locale): %v", err)
		}
		err := runDecode(cmd, []string{path})
		if err == nil || err.Error() != `no codec registered for locale "zz-none"` {
			t.Fatalf("err = %v, want 'no codec registered for locale \"zz-none\"'", err)
		}
	})
}
