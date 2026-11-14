//ff:func feature=cli type=command control=sequence
//ff:what NewCheckCmd — creates the tangl check subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewCheckCmd returns the "check" subcommand: load and validate a .geul
// file, printing lint warnings and reporting validation errors.
func NewCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <file.geul>",
		Short: "Load and validate a .geul file",
		Args:  cobra.ExactArgs(1),
		RunE:  runCheck,
	}
}
