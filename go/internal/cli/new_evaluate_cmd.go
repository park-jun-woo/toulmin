//ff:func feature=cli type=command control=sequence
//ff:what NewEvaluateCmd — creates cobra evaluate subcommand
package cli

import "github.com/spf13/cobra"

// NewEvaluateCmd returns the evaluate subcommand.
func NewEvaluateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "evaluate",
		Short: "Evaluate example rules and print verdicts",
		RunE:  runEvaluate,
	}
}
