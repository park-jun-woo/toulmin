//ff:func feature=cli type=command control=sequence
//ff:what TestNewDecodeCmd — NewDecodeCmd builds the decode subcommand
package tanglcli

import "testing"

// TestNewDecodeCmd covers the single branch of NewDecodeCmd: it always returns a
// *cobra.Command populated with the "decode <file.geul>" use string, short
// description, an exact-one-arg validator, the runDecode handler wired as RunE,
// and the "locale" and "out"/"o" flags defaulting to "".
func TestNewDecodeCmd(t *testing.T) {
	cmd := NewDecodeCmd()
	if cmd == nil {
		t.Fatal("NewDecodeCmd() returned nil")
	}
	if got, want := cmd.Use, "decode <file.geul>"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Decode a .geul stream back to a text surface"; got != want {
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
	if cmd.Flags().Lookup("locale") == nil {
		t.Error("expected a 'locale' flag to be registered")
	}
	if cmd.Flags().Lookup("out") == nil {
		t.Error("expected an 'out' flag to be registered")
	}
}
