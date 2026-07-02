//ff:func feature=cli type=command control=sequence
//ff:what TestRunEffects — runEffects reads the dir flag, resolves the target document, computes the effect closure, and prints the formatted table
package tanglcli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestRunEffects covers every branch of runEffects: a missing "dir" flag
// definition surfaces GetString's error, a resolveEffectsTarget failure is
// propagated, an unknown endpoint surfaces effects.Closure's error, and the
// success path prints the formatted effect table.
func TestRunEffects(t *testing.T) {
	t.Run("missing dir flag surfaces GetString error", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runEffects(cmd, []string{"transfer.`transfer`"}); err == nil {
			t.Fatal("expected an error from GetString(\"dir\"), got nil")
		}
	})

	t.Run("resolveEffectsTarget error is propagated", func(t *testing.T) {
		cmd := NewEffectsCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runEffects(cmd, []string{"does-not-exist.md", "transfer"}); err == nil {
			t.Fatal("expected a resolve error, got nil")
		}
	})

	t.Run("unknown endpoint surfaces effects.Closure error", func(t *testing.T) {
		cmd := NewEffectsCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		err := runEffects(cmd, []string{"../../pkg/tangl/parser/testdata/transfer.md", "bogus"})
		if err == nil {
			t.Fatal("expected an unknown-endpoint error, got nil")
		}
		if !strings.Contains(err.Error(), "unknown endpoint") {
			t.Errorf("err = %v, want 'unknown endpoint' message", err)
		}
	})

	t.Run("success prints the formatted effect table", func(t *testing.T) {
		cmd := NewEffectsCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runEffects(cmd, []string{"../../pkg/tangl/parser/testdata/transfer.md", "transfer"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out.String() != wantTransferEffects {
			t.Errorf("effects output =\n%q\nwant\n%q", out.String(), wantTransferEffects)
		}
	})
}
