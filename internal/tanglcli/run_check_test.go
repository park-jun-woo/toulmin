//ff:func feature=cli type=command control=sequence
//ff:what TestRunCheck — runCheck loads, lints, and validates a .tangeul, printing "ok" or propagating a load/validation error
package tanglcli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunCheck covers every branch of runCheck: a load error is propagated
// immediately; a clean document with no lint warnings takes the zero-iteration
// lint loop and prints "ok"; and a document with a lint warning and a validation
// error takes the non-zero-iteration lint loop, writes each warning to stderr, and
// returns the validation error instead of printing "ok".
func TestRunCheck(t *testing.T) {
	t.Run("load error is propagated", func(t *testing.T) {
		cmd := NewCheckCmd()
		var out, errOut bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&errOut)
		if err := runCheck(cmd, []string{"does-not-exist.tangeul"}); err == nil {
			t.Fatal("expected a load error, got nil")
		}
	})

	t.Run("md input is rejected with encode guidance", func(t *testing.T) {
		cmd := NewCheckCmd()
		var out, errOut bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&errOut)
		err := runCheck(cmd, []string{"policy.md"})
		if err == nil {
			t.Fatal("expected a rejection error for a .md input, got nil")
		}
		if !strings.Contains(err.Error(), "text surface is not executable") {
			t.Errorf("err = %v, want the 'text surface is not executable' guidance", err)
		}
	})

	t.Run("clean document with no lint warnings prints ok", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "transfer.tangeul")
		saveTangeul(t, path, buildTransferDoc())

		cmd := NewCheckCmd()
		var out, errOut bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&errOut)
		if err := runCheck(cmd, []string{path}); err != nil {
			t.Fatalf("unexpected error: %v (stderr: %s)", err, errOut.String())
		}
		if got := strings.TrimSpace(out.String()); got != "ok" {
			t.Errorf("stdout = %q, want %q", got, "ok")
		}
		if errOut.String() != "" {
			t.Errorf("stderr = %q, want empty (no lint warnings expected)", errOut.String())
		}
	})

	t.Run("document with a lint warning and a validation error", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "bad.tangeul")
		saveTangeul(t, path, buildLintErrorDoc())

		cmd := NewCheckCmd()
		var out, errOut bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&errOut)
		err := runCheck(cmd, []string{path})
		if err == nil {
			t.Fatal("expected a validation error, got nil")
		}
		if !strings.Contains(err.Error(), `undo when "n2" has no preceding do`) {
			t.Errorf("err = %v, want the undo-without-do gap", err)
		}
		if errOut.Len() == 0 {
			t.Error("expected at least one lint warning on stderr")
		}
		if strings.TrimSpace(out.String()) == "ok" {
			t.Error("stdout must not print \"ok\" when validation fails")
		}
	})
}
