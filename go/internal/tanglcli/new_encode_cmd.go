//ff:func feature=cli type=command control=sequence
//ff:what NewEncodeCmd — creates the tangl encode subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewEncodeCmd returns the "encode" subcommand: absorb a text surface (.md) into
// a .geul stream. The authoring locale is detected from the first section
// header and its codec drives the encode; -o/--out overrides the default output
// path (the input path with a .geul extension).
func NewEncodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encode <file.md>",
		Short: "Encode a text surface (.md) into a .geul stream",
		Args:  cobra.ExactArgs(1),
		RunE:  runEncode,
	}
	cmd.Flags().StringP("out", "o", "", "output .geul path (default: input path with .geul extension)")
	return cmd
}
