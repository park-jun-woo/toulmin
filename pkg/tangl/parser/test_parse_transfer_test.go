//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseTransfer — parses the spec's bank transfer example and checks key fields
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseTransfer parses the spec's "계좌 이체" (compensation) example
// verbatim and checks section counts plus the once/undo fields it exercises.
func TestParseTransfer(t *testing.T) {
	doc, err := Parse("testdata/transfer.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if doc.Subject != "transfer" {
		t.Errorf("Subject = %q, want transfer", doc.Subject)
	}
	if len(doc.Cases) != 2 {
		t.Fatalf("len(Cases) = %d, want 2", len(doc.Cases))
	}

	withdraw := doc.Cases[0]
	if withdraw.Name != "can withdraw" {
		t.Errorf("Cases[0].Name = %q, want 'can withdraw'", withdraw.Name)
	}
	if len(withdraw.Requires) != 2 {
		t.Fatalf("withdraw.Requires = %+v, want 2 entries", withdraw.Requires)
	}
	if len(withdraw.Execs) != 3 {
		t.Fatalf("withdraw.Execs = %+v, want 3 entries", withdraw.Execs)
	}
	doExec, undoExec, runExec := withdraw.Execs[0], withdraw.Execs[1], withdraw.Execs[2]
	if doExec.Kind != ast.DoExec || !doExec.Once || doExec.Func == nil || doExec.Func.Alias != "bank" || doExec.Func.Name != "withdraw" {
		t.Errorf("Execs[0] = %+v, want DoExec once bank.withdraw", doExec)
	}
	if undoExec.Kind != ast.UndoExec || undoExec.Func == nil || undoExec.Func.Name != "refund" || undoExec.Node != "balance sufficient" {
		t.Errorf("Execs[1] = %+v, want UndoExec bank.refund when 'balance sufficient'", undoExec)
	}
	if runExec.Kind != ast.RunExec || runExec.Case != "can deposit" {
		t.Errorf("Execs[2] = %+v, want RunExec to 'can deposit'", runExec)
	}

	deposit := doc.Cases[1]
	if deposit.Name != "can deposit" || len(deposit.Execs) != 2 {
		t.Fatalf("Cases[1] = %+v, want 'can deposit' with 2 execs", deposit)
	}
	if !deposit.Execs[0].Once || deposit.Execs[1].Once {
		t.Errorf("deposit execs once flags = [%v,%v], want [true,false]", deposit.Execs[0].Once, deposit.Execs[1].Once)
	}

	if len(doc.Provides) != 1 {
		t.Fatalf("len(Provides) = %d, want 1", len(doc.Provides))
	}
	ep := doc.Provides[0]
	if ep.Name != "transfer" || len(ep.Requires) != 3 || len(ep.Runs) != 1 || ep.Runs[0] != "can withdraw" {
		t.Errorf("Provides[0] = %+v", ep)
	}
}
