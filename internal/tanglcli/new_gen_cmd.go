//ff:func feature=cli type=command control=sequence
//ff:what NewGenCmd — creates the tangl gen subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewGenCmd returns the "gen" subcommand: parse, validate, and generate Go
// source code for a TANGL v0.3 markdown file, printed to stdout or
// written to the -o/--out path.
func NewGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen <file.md>",
		Short: "Generate Go source code for a TANGL v0.3 markdown file",
		Args:  cobra.ExactArgs(1),
		RunE:  runGen,
	}
	cmd.Flags().StringP("out", "o", "", "output file path (default: stdout)")
	return cmd
}
