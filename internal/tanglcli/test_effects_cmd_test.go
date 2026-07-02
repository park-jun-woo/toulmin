//ff:func feature=cli type=command control=sequence
//ff:what TestEffectsCmdTransfer — effects on both the "<file.md> <endpoint>" and backticked-endpoint forms reproduces the spec's transfer effect summary
package tanglcli

import (
	"bytes"
	"testing"
)

// wantTransferEffects is the transfer fixture's effect summary rendered
// by formatEffects, matching the spec's worked example content and order
// (kind, func, once, case/node) for the "transfer" endpoint.
const wantTransferEffects = "" +
	"do   bank.withdraw        once (can withdraw / balance sufficient)\n" +
	"undo bank.refund               (can withdraw / balance sufficient)\n" +
	"do   bank.deposit         once (can deposit / recipient valid)\n" +
	"do   log.TransferComplete      (can deposit / recipient valid)\n"

// TestEffectsCmdTransfer runs "tangl effects <file.md> <endpoint>" against
// transfer.md with both a bare endpoint and a backticked endpoint argument,
// checking the rendered table matches the spec's four-entry, in-order
// effect summary exactly in both cases.
func TestEffectsCmdTransfer(t *testing.T) {
	t.Run("file form with bare endpoint", func(t *testing.T) {
		cmd := NewRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"effects", "../../pkg/tangl/parser/testdata/transfer.md", "transfer"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute: %v", err)
		}
		if out.String() != wantTransferEffects {
			t.Errorf("effects output =\n%q\nwant\n%q", out.String(), wantTransferEffects)
		}
	})

	t.Run("file form with backticked endpoint", func(t *testing.T) {
		cmd := NewRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"effects", "../../pkg/tangl/parser/testdata/transfer.md", "`transfer`"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute: %v", err)
		}
		if out.String() != wantTransferEffects {
			t.Errorf("effects output =\n%q\nwant\n%q", out.String(), wantTransferEffects)
		}
	})
}
