//ff:func feature=cli type=command control=sequence
//ff:what NewRootCmd — creates the tangl CLI root cobra command
package tanglcli

import "github.com/spf13/cobra"

// NewRootCmd returns the root cobra command for the tangl CLI, with the
// check, ast, effects, and gen subcommands registered.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tangl",
		Short: "TANGL v0.3 markdown policy DSL toolchain",
	}
	cmd.AddCommand(NewCheckCmd())
	cmd.AddCommand(NewAstCmd())
	cmd.AddCommand(NewEffectsCmd())
	cmd.AddCommand(NewGenCmd())
	return cmd
}
