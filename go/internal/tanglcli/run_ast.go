//ff:func feature=cli type=command control=sequence
//ff:what runAst — loads a .geul file and prints its AST as indented JSON
package tanglcli

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

// runAst loads the .geul at args[0] and writes its ast.Document to stdout as
// indented JSON.
func runAst(cmd *cobra.Command, args []string) error {
	d, err := loadExecutable(args[0])
	if err != nil {
		return err
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(d.AST)
}
