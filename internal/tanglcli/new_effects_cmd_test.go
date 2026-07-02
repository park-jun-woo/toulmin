//ff:func feature=cli type=command control=sequence
//ff:what TestNewEffectsCmd — NewEffectsCmd builds the effects subcommand
package tanglcli

import "testing"

// TestNewEffectsCmd covers the single branch of NewEffectsCmd: it always
// returns a *cobra.Command populated with the effects use string, short
// description, a one-or-two-arg validator, the runEffects handler wired as
// RunE, and a "dir" flag defaulting to ".".
func TestNewEffectsCmd(t *testing.T) {
	cmd := NewEffectsCmd()
	if cmd == nil {
		t.Fatal("NewEffectsCmd() returned nil")
	}
	if got, want := cmd.Use, "effects <file.md> <endpoint> | <subject>.<endpoint>"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Print the static do/undo effect summary reachable from an endpoint"; got != want {
		t.Errorf("Short = %q, want %q", got, want)
	}
	if cmd.Args == nil {
		t.Fatal("Args must be set")
	}
	if err := cmd.Args(cmd, []string{"one"}); err != nil {
		t.Errorf("Args(1 arg) = %v, want nil", err)
	}
	if err := cmd.Args(cmd, []string{"one", "two"}); err != nil {
		t.Errorf("Args(2 args) = %v, want nil", err)
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("Args(0 args) = nil, want error")
	}
	if err := cmd.Args(cmd, []string{"one", "two", "three"}); err == nil {
		t.Error("Args(3 args) = nil, want error")
	}
	if cmd.RunE == nil {
		t.Fatal("RunE must be set")
	}
	dirFlag := cmd.Flags().Lookup("dir")
	if dirFlag == nil {
		t.Fatal("expected a 'dir' flag to be registered")
	}
	if got, want := dirFlag.DefValue, "."; got != want {
		t.Errorf("dir flag default = %q, want %q", got, want)
	}
}
