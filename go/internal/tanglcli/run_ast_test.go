//ff:func feature=cli type=command control=sequence
//ff:what TestRunAst — runAst loads a .geul and prints its AST as JSON, or propagates a load error
package tanglcli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestRunAst covers both branches of runAst: a load error from loadExecutable is
// propagated unchanged, and a successful load is JSON-encoded to the command's
// output.
func TestRunAst(t *testing.T) {
	t.Run("load error is propagated", func(t *testing.T) {
		cmd := NewAstCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runAst(cmd, []string{"does-not-exist.geul"}); err == nil {
			t.Fatal("expected a load error, got nil")
		}
	})

	t.Run("successful load is encoded as JSON", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "transfer.geul")
		saveTangeul(t, path, buildTransferDoc())

		cmd := NewAstCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runAst(cmd, []string{path}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var decoded struct {
			Subject string `json:"subject"`
		}
		if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
			t.Fatalf("json.Unmarshal: %v (output: %s)", err, out.String())
		}
		if decoded.Subject != "transfer" {
			t.Errorf("subject = %q, want %q", decoded.Subject, "transfer")
		}
	})
}
