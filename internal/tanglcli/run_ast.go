//ff:func feature=cli type=command control=sequence
//ff:what runAst — parses a TANGL v0.3 file and prints its AST as indented JSON
package tanglcli

import (
	"encoding/json"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
	"github.com/spf13/cobra"
)

// runAst parses args[0] and writes its ast.Document to stdout as
// indented JSON.
func runAst(cmd *cobra.Command, args []string) error {
	doc, err := parser.Parse(args[0])
	if err != nil {
		return err
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}
