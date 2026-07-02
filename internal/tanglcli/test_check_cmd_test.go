//ff:func feature=cli type=command control=sequence
//ff:what TestCheckCmdTransferOK — check on the clean transfer fixture prints "ok" and returns no error
package tanglcli

import (
	"bytes"
	"strings"
	"testing"
)

// TestCheckCmdTransferOK runs "tangl check" against the clean transfer.md
// fixture and expects a single "ok" line on stdout with no error.
func TestCheckCmdTransferOK(t *testing.T) {
	cmd := NewRootCmd()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs([]string{"check", "../../pkg/tangl/parser/testdata/transfer.md"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v (stderr: %s)", err, errOut.String())
	}
	if got := strings.TrimSpace(out.String()); got != "ok" {
		t.Errorf("stdout = %q, want %q", got, "ok")
	}
}
