//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestNewRootCmd — NewRootCmd builds the root command with the evaluate subcommand attached
package cli

import "testing"

// TestNewRootCmd covers the single branch of NewRootCmd: it always returns a
// *cobra.Command populated with the "toulmin" use string, short description,
// and the evaluate subcommand registered as a child.
func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd == nil {
		t.Fatal("NewRootCmd() returned nil")
	}
	if got, want := cmd.Use, "toulmin"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Toulmin argumentation-based rule engine"; got != want {
		t.Errorf("Short = %q, want %q", got, want)
	}
	found := false
	for _, c := range cmd.Commands() {
		if c.Use == "evaluate" {
			found = true
		}
	}
	if !found {
		t.Error("expected evaluate subcommand to be registered")
	}
}
