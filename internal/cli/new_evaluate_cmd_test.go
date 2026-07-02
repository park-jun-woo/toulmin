//ff:func feature=cli type=command control=sequence
//ff:what TestNewEvaluateCmd — NewEvaluateCmd builds the evaluate subcommand
package cli

import "testing"

// TestNewEvaluateCmd covers the single branch of NewEvaluateCmd: it always
// returns a *cobra.Command populated with the "evaluate" use string, short
// description, and the runEvaluate handler wired as RunE.
func TestNewEvaluateCmd(t *testing.T) {
	cmd := NewEvaluateCmd()
	if cmd == nil {
		t.Fatal("NewEvaluateCmd() returned nil")
	}
	if got, want := cmd.Use, "evaluate"; got != want {
		t.Errorf("Use = %q, want %q", got, want)
	}
	if got, want := cmd.Short, "Evaluate example rules and print verdicts"; got != want {
		t.Errorf("Short = %q, want %q", got, want)
	}
	if cmd.RunE == nil {
		t.Fatal("RunE must be set")
	}
}
