//ff:func feature=cli type=command control=sequence
//ff:what TestNewEncodeCmd — NewEncodeCmd builds the encode subcommand
package tanglcli

import "testing"

// TestNewEncodeCmd covers the single branch of NewEncodeCmd: it always returns a
// *cobra.Command populated with the "encode <file.md>" use string, short
// description, an exact-one-arg validator, the runEncode handler wired as RunE,
// and an "out"/"o" flag defaulting to "".
func TestNewEncodeCmd(t *testing.T) {
	cmd := NewEncodeCmd()
	if cmd == nil {
		t.Fatal("NewEncodeCmd() returned nil")
	}
	if got, want := cmd.Use, "encode <file.md>"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Encode a text surface (.md) into a .geul stream"; got != want {
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
}
