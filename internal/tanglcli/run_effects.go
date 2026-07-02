//ff:func feature=cli type=command control=sequence
//ff:what runEffects — resolves an endpoint target and prints its effect closure
package tanglcli

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/effects"
	"github.com/spf13/cobra"
)

// runEffects resolves the (Document, endpoint) target named by args (see
// resolveEffectsTarget), computes its static do/undo effect closure, and
// prints it as an aligned table to stdout.
func runEffects(cmd *cobra.Command, args []string) error {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}
	doc, endpoint, err := resolveEffectsTarget(args, dir)
	if err != nil {
		return err
	}
	entries, err := effects.Closure(doc, endpoint)
	if err != nil {
		return err
	}
	fmt.Fprint(cmd.OutOrStdout(), formatEffects(entries))
	return nil
}
