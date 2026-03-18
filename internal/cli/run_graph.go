//ff:func feature=cli type=command control=sequence
//ff:what runGraph — reads YAML graph definition, validates, and generates Graph Builder code
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/internal/codegen"
	"github.com/park-jun-woo/toulmin/internal/graphdef"
	"github.com/spf13/cobra"
)

// runGraph orchestrates YAML parsing, validation, and Graph Builder code generation.
func runGraph(cmd *cobra.Command, args []string) error {
	yamlPath := args[0]
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	output, _ := cmd.Flags().GetString("output")
	pkgName, _ := cmd.Flags().GetString("package")
	def, err := graphdef.ParseYAML(yamlPath)
	if err != nil {
		return err
	}
	if err := graphdef.Validate(def); err != nil {
		return err
	}
	if pkgName == "" {
		pkgName = dirToPkg(filepath.Dir(yamlPath))
	}
	code, err := codegen.GenerateGraph(pkgName, def)
	if err != nil {
		return err
	}
	if dryRun {
		fmt.Print(code)
		return nil
	}
	if output == "" {
		output = filepath.Join(filepath.Dir(yamlPath), "graph_gen.go")
	}
	return os.WriteFile(output, []byte(code), 0644)
}
