//ff:func feature=cli type=command control=sequence
//ff:what NewRootCmd — creates the tangl CLI root cobra command
package tanglcli

import "github.com/spf13/cobra"

// NewRootCmd returns the root cobra command for the tangl CLI, with the
// check, ast, effects, gen, encode, and decode subcommands registered. Surface
// codecs are wired via registerCodecs before the tree is built.
func NewRootCmd() *cobra.Command {
	registerCodecs()
	cmd := &cobra.Command{
		Use:   "tangl",
		Short: "TANGEUL v0.3 markdown policy DSL toolchain",
	}
	cmd.AddCommand(NewCheckCmd())
	cmd.AddCommand(NewAstCmd())
	cmd.AddCommand(NewEffectsCmd())
	cmd.AddCommand(NewGenCmd())
	cmd.AddCommand(NewEncodeCmd())
	cmd.AddCommand(NewDecodeCmd())
	return cmd
}
