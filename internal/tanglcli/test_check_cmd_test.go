//ff:func feature=cli type=command control=sequence
//ff:what TestCheckCmdTransferOK — check on the clean transfer .tangeul fixture prints "ok" and returns no error
package tanglcli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestCheckCmdTransferOK saves the clean transfer fixture to a .tangeul file and
// runs "tangl check" against it, expecting a single "ok" line on stdout with no
// error.
func TestCheckCmdTransferOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "transfer.tangeul")
	saveTangeul(t, path, buildTransferDoc())

	cmd := NewRootCmd()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs([]string{"check", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v (stderr: %s)", err, errOut.String())
	}
	if got := strings.TrimSpace(out.String()); got != "ok" {
		t.Errorf("stdout = %q, want %q", got, "ok")
	}
}
