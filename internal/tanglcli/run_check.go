//ff:func feature=cli type=command control=sequence
//ff:what runCheck — loads, validates, and lints a .tangeul file
package tanglcli

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/registry"
	"github.com/park-jun-woo/toulmin/pkg/tangl/validate"
	"github.com/spf13/cobra"
)

// runCheck loads the .tangeul at args[0] (tiling verified by Load), validates the
// recovered Document, prints every Lint warning to stderr, and returns the
// Validate error (if any) so the caller exits non-zero. When the document uses
// alias-qualified (cross-document) checking it additionally loads the whole
// reachable document graph and rejects cross-document checking cycles; an
// ordinary document skips that path entirely. On success it prints a single
// "ok" line to stdout.
func runCheck(cmd *cobra.Command, args []string) error {
	d, err := loadExecutable(args[0])
	if err != nil {
		return err
	}
	verr := validate.Validate(d.AST)
	printLintWarnings(cmd.ErrOrStderr(), validate.Lint(d.AST))
	if verr != nil {
		return verr
	}
	if hasAliasChecking(d.AST) {
		reg, err := registry.Load(args[0])
		if err != nil {
			return err
		}
		if err := registry.CheckCrossDocCycle(reg); err != nil {
			return err
		}
	}
	fmt.Fprintln(cmd.OutOrStdout(), "ok")
	return nil
}
