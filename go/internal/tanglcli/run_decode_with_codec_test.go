//ff:func feature=cli type=command control=sequence
//ff:what TestRunDecodeCodec — mode-2 rendering via a registered codec (success and Render error), and the locale/out flag-error branches
package tanglcli

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/spf13/cobra"
)

// TestRunDecodeCodec covers runDecode's remaining branches, which need a
// registered codec: mode 2 with a found codec renders via doc.Render; a codec
// whose RenderStmt fails propagates that error; and the locale/out flag lookups
// error when the flag is undefined on cmd. It is split out from TestRunDecode to
// stay under the control=sequence line limit.
func TestRunDecodeCodec(t *testing.T) {
	d := buildTransferDoc()
	path := filepath.Join(t.TempDir(), "transfer.geul")
	saveTangeul(t, path, d)

	t.Run("mode 2 with a registered codec renders via Render", func(t *testing.T) {
		// The real en codec (installed by registerCodecs) drives the found-codec
		// success path: it has section headers and a RenderStmt for the transfer
		// fixture, so Render produces a non-empty mode-2 rendering. Using "en"
		// avoids registering a second "ko" codec, which would collide with the
		// fake "ko" encode codec of TestRunEncodeCodec.
		registerCodecs()
		cmd := NewDecodeCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := cmd.Flags().Set("locale", "en"); err != nil {
			t.Fatalf("Set(locale): %v", err)
		}
		if err := runDecode(cmd, []string{path}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out.Len() == 0 {
			t.Error("expected non-empty stdout output from Render")
		}
	})

	t.Run("mode 2 propagates a Render error", func(t *testing.T) {
		doc.RegisterCodec(fakeDecodeCodec{locale: "zz-decode-fail", fail: true})
		cmd := NewDecodeCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := cmd.Flags().Set("locale", "zz-decode-fail"); err != nil {
			t.Fatalf("Set(locale): %v", err)
		}
		if err := runDecode(cmd, []string{path}); err == nil {
			t.Fatal("expected a Render error, got nil")
		}
	})

	t.Run("locale flag GetString error is propagated", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runDecode(cmd, []string{path}); err == nil {
			t.Fatal("expected a flag error for the undefined 'locale' flag, got nil")
		}
	})

	t.Run("out flag GetString error is propagated", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("locale", "", "")
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runDecode(cmd, []string{path}); err == nil {
			t.Fatal("expected a flag error for the undefined 'out' flag, got nil")
		}
	})
}
