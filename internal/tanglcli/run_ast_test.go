//ff:func feature=cli type=command control=sequence
//ff:what TestRunAst — runAst parses a TANGL file and prints its AST as JSON, or propagates a parse error
package tanglcli

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestRunAst covers both branches of runAst: a parse error from
// parser.Parse is propagated unchanged, and a successful parse is
// JSON-encoded to the command's output.
func TestRunAst(t *testing.T) {
	t.Run("parse error is propagated", func(t *testing.T) {
		cmd := NewAstCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runAst(cmd, []string{"does-not-exist.md"}); err == nil {
			t.Fatal("expected a parse error, got nil")
		}
	})

	t.Run("successful parse is encoded as JSON", func(t *testing.T) {
		cmd := NewAstCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		if err := runAst(cmd, []string{"../../pkg/tangl/parser/testdata/transfer.md"}); err != nil {
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
