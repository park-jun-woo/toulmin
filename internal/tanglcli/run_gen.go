//ff:func feature=cli type=command control=sequence
//ff:what runGen — loads, validates, and generates Go source for a .tangeul file
package tanglcli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/registry"
	"github.com/park-jun-woo/toulmin/pkg/tangl/gen"
	"github.com/park-jun-woo/toulmin/pkg/tangl/validate"
	"github.com/spf13/cobra"
)

// runGen loads the .tangeul at args[0], validates the Document (aborting on
// error), and generates its Go source. When the document uses alias-qualified
// (cross-document) checking it first loads the whole reachable document graph
// and rejects cross-document checking cycles, so gen only ever emits a
// cross-package call for a verified acyclic reference; gen itself derives each
// alias target's Go import path from d.AST.Path. With no -o/--out flag the
// source is printed to stdout; otherwise it is written to that path (mode 0644).
func runGen(cmd *cobra.Command, args []string) error {
	d, err := loadExecutable(args[0])
	if err != nil {
		return err
	}
	if err := validate.Validate(d.AST); err != nil {
		return err
	}
	if hasAliasChecking(d.AST) {
		reg, err := registry.Load(args[0])
		if err != nil {
			return err
		}
		if err := registry.CheckCrossDocCycle(reg); err != nil {
			return err
		}
	}
	src, err := gen.Generate(d.AST)
	if err != nil {
		return err
	}
	out, err := cmd.Flags().GetString("out")
	if err != nil {
		return err
	}
	if out == "" {
		fmt.Fprint(cmd.OutOrStdout(), src)
		return nil
	}
	return os.WriteFile(out, []byte(src), 0644)
}
