//ff:func feature=cli type=command control=sequence
//ff:what NewDecodeCmd — creates the tangl decode subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewDecodeCmd returns the "decode" subcommand: render a .tangeul stream back to a
// text surface. With no flags it is mode 1 (verbatim, byte-identical Source dump);
// --locale en|ko selects mode 2 (normalized rendering via that locale's codec).
// -o/--out writes to a file instead of stdout.
func NewDecodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode <file.tangeul>",
		Short: "Decode a .tangeul stream back to a text surface",
		Args:  cobra.ExactArgs(1),
		RunE:  runDecode,
	}
	cmd.Flags().String("locale", "", "target locale for mode-2 rendering (en|ko); default is verbatim mode-1")
	cmd.Flags().StringP("out", "o", "", "output path (default: stdout)")
	return cmd
}
