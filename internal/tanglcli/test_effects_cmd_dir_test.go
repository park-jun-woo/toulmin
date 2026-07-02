//ff:func feature=cli type=command control=sequence
//ff:what TestEffectsCmdDirForm — effects on "<subject>.<endpoint> --dir <d>" scans a directory by Subject
package tanglcli

import (
	"bytes"
	"testing"
)

// TestEffectsCmdDirForm runs "tangl effects transfer.`transfer` --dir
// <testdata>" and checks it locates transfer.md by its Subject and
// reproduces the same effect summary as the direct file form.
func TestEffectsCmdDirForm(t *testing.T) {
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"effects", "transfer.`transfer`", "--dir", "../../pkg/tangl/parser/testdata"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.String() != wantTransferEffects {
		t.Errorf("effects --dir output =\n%q\nwant\n%q", out.String(), wantTransferEffects)
	}
}
