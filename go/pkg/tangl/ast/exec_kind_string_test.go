//ff:func feature=tangl type=model control=sequence
//ff:what TestExecKind_String — tests ExecKind.String for all known kinds and unknown default
package ast

import "testing"

func TestExecKind_String(t *testing.T) {
	if got := DoExec.String(); got != "DoExec" {
		t.Errorf("DoExec: got %q, want %q", got, "DoExec")
	}
	if got := UndoExec.String(); got != "UndoExec" {
		t.Errorf("UndoExec: got %q, want %q", got, "UndoExec")
	}
	if got := RunExec.String(); got != "RunExec" {
		t.Errorf("RunExec: got %q, want %q", got, "RunExec")
	}
	if got := ExecKind(99).String(); got != "ExecKind(?)" {
		t.Errorf("unknown: got %q, want %q", got, "ExecKind(?)")
	}
}
