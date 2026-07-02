//ff:func feature=cli type=command control=sequence
//ff:what TestNewCheckCmd — NewCheckCmd builds the check subcommand
package tanglcli

import "testing"

// TestNewCheckCmd covers the single branch of NewCheckCmd: it always returns
// a *cobra.Command populated with the "check <file.md>" use string, short
// description, an exact-one-arg validator, and the runCheck handler wired as
// RunE.
func TestNewCheckCmd(t *testing.T) {
	cmd := NewCheckCmd()
	if cmd == nil {
		t.Fatal("NewCheckCmd() returned nil")
	}
	if got, want := cmd.Use, "check <file.md>"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Parse and validate a TANGL v0.3 markdown file"; got != want {
		t.Errorf("Short = %q, want %q", got, want)
	}
	if cmd.Args == nil {
		t.Fatal("Args must be set")
	}
	if err := cmd.Args(cmd, []string{"one"}); err != nil {
		t.Errorf("Args(1 arg) = %v, want nil", err)
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("Args(0 args) = nil, want error")
	}
	if cmd.RunE == nil {
		t.Fatal("RunE must be set")
	}
}
