//ff:func feature=cli type=command control=sequence
//ff:what TestNewAstCmd — NewAstCmd builds the ast subcommand
package tanglcli

import "testing"

// TestNewAstCmd covers the single branch of NewAstCmd: it always returns a
// *cobra.Command populated with the "ast <file.md>" use string, short
// description, an exact-one-arg validator, and the runAst handler wired as
// RunE.
func TestNewAstCmd(t *testing.T) {
	cmd := NewAstCmd()
	if cmd == nil {
		t.Fatal("NewAstCmd() returned nil")
	}
	if got, want := cmd.Use, "ast <file.md>"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Parse a TANGL v0.3 markdown file and print its AST as JSON"; got != want {
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
