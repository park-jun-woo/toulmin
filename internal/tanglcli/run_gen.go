//ff:func feature=cli type=command control=sequence
//ff:what runGen — loads, validates, and generates Go source for a .tangeul file
package tanglcli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/toulmin/pkg/tangl/gen"
	"github.com/park-jun-woo/toulmin/pkg/tangl/validate"
	"github.com/spf13/cobra"
)

// runGen loads the .tangeul at args[0], validates the Document (aborting on
// error), and generates its Go source. With no -o/--out flag the source is
// printed to stdout; otherwise it is written to that path (mode 0644).
func runGen(cmd *cobra.Command, args []string) error {
	d, err := loadExecutable(args[0])
	if err != nil {
		return err
	}
	if err := validate.Validate(d.AST); err != nil {
		return err
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
