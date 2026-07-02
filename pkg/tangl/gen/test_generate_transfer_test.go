//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateTransfer — generates the spec's bank transfer example and checks the compensation wrapper
package gen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestGenerateTransfer compiles the spec's "계좌 이체" (compensation)
// example and checks the once-guard key, the PushCompensation arming for
// the undo edge, and the full endpoint compensation wrapper.
func TestGenerateTransfer(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/transfer.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	src, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := format.Source([]byte(src)); err != nil {
		t.Fatalf("re-format.Source: %v", err)
	}
	if !strings.Contains(src, "package transfer") {
		t.Errorf("missing package clause; got:\n%s", src)
	}
	if !strings.Contains(src, `"once:transfer.can withdraw.balance sufficient#0"`) {
		t.Errorf("missing the once-guard key for the withdraw do edge")
	}
	if !strings.Contains(src, "tangl.PushCompensation(t.Ctx(), func() error {") ||
		!strings.Contains(src, "return bank.refund(t.Ctx())") {
		t.Errorf("missing the refund compensation closure armed after withdraw succeeds")
	}
	if !strings.Contains(src, "func Transfer(ctx toulmin.Context) error {") {
		t.Errorf("missing the transfer endpoint function")
	}
	if !strings.Contains(src, "tangl.InitCompensation(ctx)") ||
		!strings.Contains(src, "tangl.Compensate(ctx)") ||
		!strings.Contains(src, "tangl.Review(ctx, err, cerr)") ||
		!strings.Contains(src, "tangl.CommitCompensation(ctx)") {
		t.Errorf("missing the full Init/Compensate/Review/Commit wrapper cycle")
	}
	if !strings.Contains(src, ".Run(canDepositGraph)") {
		t.Errorf("missing the can-withdraw to can-deposit run cascade")
	}
}
