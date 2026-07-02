//ff:func feature=cli type=command control=sequence
//ff:what runCheck — parses, validates, and lints a TANGL v0.3 file
package tanglcli

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
	"github.com/park-jun-woo/toulmin/pkg/tangl/validate"
	"github.com/spf13/cobra"
)

// runCheck parses args[0], validates the resulting Document, prints every
// Lint warning to stderr, and returns the Validate error (if any) so the
// caller exits non-zero. On success it prints a single "ok" line to stdout.
func runCheck(cmd *cobra.Command, args []string) error {
	doc, err := parser.Parse(args[0])
	if err != nil {
		return err
	}
	verr := validate.Validate(doc)
	printLintWarnings(cmd.ErrOrStderr(), validate.Lint(doc))
	if verr != nil {
		return verr
	}
	fmt.Fprintln(cmd.OutOrStdout(), "ok")
	return nil
}
