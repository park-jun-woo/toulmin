//ff:func feature=cli type=command control=sequence
//ff:what NewCheckCmd — creates the tangl check subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewCheckCmd returns the "check" subcommand: load and validate a .tangeul
// file, printing lint warnings and reporting validation errors.
func NewCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <file.tangeul>",
		Short: "Load and validate a .tangeul file",
		Args:  cobra.ExactArgs(1),
		RunE:  runCheck,
	}
}
