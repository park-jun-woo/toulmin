//ff:func feature=cli type=command control=sequence
//ff:what NewAstCmd — creates the tangl ast subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewAstCmd returns the "ast" subcommand: load a .tangeul file and print its
// AST as indented JSON.
func NewAstCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ast <file.tangeul>",
		Short: "Load a .tangeul file and print its AST as JSON",
		Args:  cobra.ExactArgs(1),
		RunE:  runAst,
	}
}
