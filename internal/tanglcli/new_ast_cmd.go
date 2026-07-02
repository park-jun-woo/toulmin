//ff:func feature=cli type=command control=sequence
//ff:what NewAstCmd — creates the tangl ast subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewAstCmd returns the "ast" subcommand: parse a TANGL v0.3 markdown file
// and print its AST as indented JSON.
func NewAstCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ast <file.md>",
		Short: "Parse a TANGL v0.3 markdown file and print its AST as JSON",
		Args:  cobra.ExactArgs(1),
		RunE:  runAst,
	}
}
