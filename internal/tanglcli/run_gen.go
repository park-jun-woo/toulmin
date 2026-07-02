//ff:func feature=cli type=command control=sequence
//ff:what runGen — parses, validates, and generates Go source for a TANGL v0.3 file
package tanglcli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/toulmin/pkg/tangl/gen"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
	"github.com/park-jun-woo/toulmin/pkg/tangl/validate"
	"github.com/spf13/cobra"
)

// runGen parses args[0], validates the Document (aborting on error), and
// generates its Go source. With no -o/--out flag the source is printed to
// stdout; otherwise it is written to that path (mode 0644).
func runGen(cmd *cobra.Command, args []string) error {
	doc, err := parser.Parse(args[0])
	if err != nil {
		return err
	}
	if err := validate.Validate(doc); err != nil {
		return err
	}
	src, err := gen.Generate(doc)
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
