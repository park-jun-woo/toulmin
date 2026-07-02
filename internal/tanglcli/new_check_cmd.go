//ff:func feature=cli type=command control=sequence
//ff:what NewCheckCmd — creates the tangl check subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewCheckCmd returns the "check" subcommand: parse and validate a TANGL
// v0.3 markdown file, printing lint warnings and reporting validation
// errors.
func NewCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <file.md>",
		Short: "Parse and validate a TANGL v0.3 markdown file",
		Args:  cobra.ExactArgs(1),
		RunE:  runCheck,
	}
}
