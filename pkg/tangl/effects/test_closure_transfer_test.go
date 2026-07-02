//ff:func feature=tangl type=analyzer control=iteration dimension=1
//ff:what TestClosureTransfer — the transfer fixture's effect closure matches the spec's worked example
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestClosureTransfer parses the spec's bank transfer example and checks
// that Closure(doc, "transfer") reproduces the spec's `tangl effects
// transfer.transfer` worked example exactly:
//
//	do   bank.withdraw   once   (can withdraw / balance sufficient)
//	undo bank.refund            (can withdraw / balance sufficient)
//	do   bank.deposit    once   (can deposit / recipient valid)
//	do   log.TransferComplete   (can deposit / recipient valid)
func TestClosureTransfer(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/transfer.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	entries, err := Closure(doc, "transfer")
	if err != nil {
		t.Fatalf("Closure: %v", err)
	}
	want := []Entry{
		{Kind: "do", Func: ast.Ref{Alias: "bank", Name: "withdraw"}, Once: true, Case: "can withdraw", Node: "balance sufficient"},
		{Kind: "undo", Func: ast.Ref{Alias: "bank", Name: "refund"}, Case: "can withdraw", Node: "balance sufficient"},
		{Kind: "do", Func: ast.Ref{Alias: "bank", Name: "deposit"}, Once: true, Case: "can deposit", Node: "recipient valid"},
		{Kind: "do", Func: ast.Ref{Alias: "log", Name: "TransferComplete"}, Case: "can deposit", Node: "recipient valid"},
	}
	if len(entries) != len(want) {
		t.Fatalf("Closure returned %d entries, want %d: %+v", len(entries), len(want), entries)
	}
	for i, e := range entries {
		if e != want[i] {
			t.Errorf("entries[%d] = %+v, want %+v", i, e, want[i])
		}
	}
}
