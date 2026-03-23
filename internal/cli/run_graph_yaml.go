//ff:func feature=cli type=command control=sequence
//ff:what runGraphYAML — handles YAML graph validation and code generation
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
	"github.com/spf13/cobra"
)

// runGraphYAML handles the YAML → validate → generate flow.
func runGraphYAML(cmd *cobra.Command, yamlPath string) error {
	check, _ := cmd.Flags().GetBool("check")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	output, _ := cmd.Flags().GetString("output")
	pkgName, _ := cmd.Flags().GetString("package")
	def, err := toulmin.ParseYAML(yamlPath)
	if err != nil {
		return err
	}
	if err := toulmin.ValidateGraphDef(def); err != nil {
		return err
	}
	if check {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: no cycles detected (%d rules, %d defeats)\n",
			def.Graph, len(def.Rules), len(def.Defeats))
		return nil
	}
	if pkgName == "" {
		pkgName = dirToPkg(filepath.Dir(yamlPath))
	}
	code, err := toulmin.GenerateGraph(pkgName, &def)
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
