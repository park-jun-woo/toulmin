//ff:func feature=cli type=command control=sequence
//ff:what TestRunCheck — runCheck parses, lints, and validates a TANGL file, printing "ok" or propagating a validation/parse error
package tanglcli

import (
	"bytes"
	"strings"
	"testing"
)

// TestRunCheck covers every branch of runCheck: a parse error is propagated
// immediately; a clean document with no tangl:Internal section takes the
// zero-iteration lint loop and prints "ok"; and a document with lint
// warnings and a validation error takes the non-zero-iteration lint loop,
// writes each warning to stderr, and returns the validation error instead
// of printing "ok".
func TestRunCheck(t *testing.T) {
	t.Run("parse error is propagated", func(t *testing.T) {
		cmd := NewCheckCmd()
		var out, errOut bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&errOut)
		if err := runCheck(cmd, []string{"does-not-exist.md"}); err == nil {
			t.Fatal("expected a parse error, got nil")
		}
	})

	t.Run("clean document with no lint warnings prints ok", func(t *testing.T) {
		cmd := NewCheckCmd()
		var out, errOut bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&errOut)
		if err := runCheck(cmd, []string{"../../pkg/tangl/parser/testdata/transfer.md"}); err != nil {
			t.Fatalf("unexpected error: %v (stderr: %s)", err, errOut.String())
		}
		if got := strings.TrimSpace(out.String()); got != "ok" {
			t.Errorf("stdout = %q, want %q", got, "ok")
		}
		if errOut.String() != "" {
			t.Errorf("stderr = %q, want empty (no lint warnings expected)", errOut.String())
		}
	})

	t.Run("document with lint warnings and a validation error", func(t *testing.T) {
		cmd := NewCheckCmd()
		var out, errOut bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&errOut)
		err := runCheck(cmd, []string{"../../pkg/tangl/parser/testdata/americano.md"})
		if err == nil {
			t.Fatal("expected a validation error, got nil")
		}
		if !strings.Contains(err.Error(), `undeclared package alias "log"`) {
			t.Errorf("err = %v, want the known undeclared 'log' alias gap", err)
		}
		if errOut.Len() == 0 {
			t.Error("expected at least one lint warning on stderr")
		}
		if strings.TrimSpace(out.String()) == "ok" {
			t.Error("stdout must not print \"ok\" when validation fails")
		}
	})
}
