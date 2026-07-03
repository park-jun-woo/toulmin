//ff:func feature=cli type=command control=sequence
//ff:what TestEffectsCmdDirForm — effects on "<subject>.<endpoint> --dir <d>" scans a directory by Subject
package tanglcli

import (
	"bytes"
	"path/filepath"
	"testing"
)

// TestEffectsCmdDirForm saves the transfer fixture into a directory, then runs
// "tangl effects transfer.`transfer` --dir <d>" and checks it locates the
// .tangeul by its Subject and reproduces the same effect summary as the direct
// file form.
func TestEffectsCmdDirForm(t *testing.T) {
	dir := t.TempDir()
	saveTangeul(t, filepath.Join(dir, "transfer.tangeul"), buildTransferDoc())

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"effects", "transfer.`transfer`", "--dir", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.String() != wantTransferEffects {
		t.Errorf("effects --dir output =\n%q\nwant\n%q", out.String(), wantTransferEffects)
	}
}
