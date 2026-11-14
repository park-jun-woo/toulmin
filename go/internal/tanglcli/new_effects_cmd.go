//ff:func feature=cli type=command control=sequence
//ff:what NewEffectsCmd — creates the tangl effects subcommand
package tanglcli

import "github.com/spf13/cobra"

// NewEffectsCmd returns the "effects" subcommand: print the static
// do/undo effect summary reachable from an endpoint, either as
// "<file.geul> <endpoint>" or as "<subject>.<endpoint> --dir <d>".
func NewEffectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "effects <file.geul> <endpoint> | <subject>.<endpoint>",
		Short: "Print the static do/undo effect summary reachable from an endpoint",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runEffects,
	}
	cmd.Flags().String("dir", ".", "directory to scan for the subject's .geul when a bare <subject>.<endpoint> is given")
	return cmd
}
