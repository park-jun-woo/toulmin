//ff:func feature=cli type=command control=sequence
//ff:what TestNewGenCmd — NewGenCmd builds the gen subcommand
package tanglcli

import "testing"

// TestNewGenCmd covers the single branch of NewGenCmd: it always returns a
// *cobra.Command populated with the "gen <file.md>" use string, short
// description, an exact-one-arg validator, the runGen handler wired as RunE,
// and an "out"/"o" flag defaulting to "".
func TestNewGenCmd(t *testing.T) {
	cmd := NewGenCmd()
	if cmd == nil {
		t.Fatal("NewGenCmd() returned nil")
	}
	if got, want := cmd.Use, "gen <file.md>"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Generate Go source code for a TANGL v0.3 markdown file"; got != want {
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
	outFlag := cmd.Flags().Lookup("out")
	if outFlag == nil {
		t.Fatal("expected an 'out' flag to be registered")
	}
	if got, want := outFlag.Shorthand, "o"; got != want {
		t.Errorf("out flag shorthand = %q, want %q", got, want)
	}
	if got, want := outFlag.DefValue, ""; got != want {
		t.Errorf("out flag default = %q, want %q", got, want)
	}
}
