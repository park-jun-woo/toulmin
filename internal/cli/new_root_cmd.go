//ff:func feature=cli type=command control=sequence
//ff:what NewRootCmd — creates cobra root command
package cli

import "github.com/spf13/cobra"

// NewRootCmd returns the root cobra command for the toulmin CLI.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toulmin",
		Short: "Toulmin argumentation-based rule engine",
	}
	cmd.AddCommand(NewEvaluateCmd())
	cmd.AddCommand(NewGraphCmd())
	return cmd
}
