//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestNewRootCmd — NewRootCmd builds the tangl root command with all subcommands attached
package tanglcli

import "testing"

// TestNewRootCmd covers the single branch of NewRootCmd: it always returns a
// *cobra.Command populated with the "tangl" use string, short description,
// and the check/ast/effects/gen subcommands registered as children.
func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd == nil {
		t.Fatal("NewRootCmd() returned nil")
	}
	if got, want := cmd.Use, "tangl"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "TANGL v0.3 markdown policy DSL toolchain"; got != want {
		t.Errorf("Short = %q, want %q", got, want)
	}
	wantUses := map[string]bool{
		"check <file.md>": false,
		"ast <file.md>":   false,
		"effects <file.md> <endpoint> | <subject>.<endpoint>": false,
		"gen <file.md>": false,
	}
	for _, c := range cmd.Commands() {
		if _, ok := wantUses[c.Use]; ok {
			wantUses[c.Use] = true
		}
	}
	for use, found := range wantUses {
		if !found {
			t.Errorf("expected subcommand %q to be registered", use)
		}
	}
}
