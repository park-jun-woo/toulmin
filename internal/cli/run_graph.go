//ff:func feature=cli type=command control=sequence
//ff:what runGraph — scans source, validates graph, generates registration code
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/internal/codegen"
	"github.com/park-jun-woo/toulmin/internal/graph"
	"github.com/park-jun-woo/toulmin/internal/scanner"
	"github.com/spf13/cobra"
)

// runGraph orchestrates source scanning, graph validation, and code generation.
func runGraph(cmd *cobra.Command, args []string) error {
	dir := args[0]
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	output, _ := cmd.Flags().GetString("output")
	paths, err := scanner.ScanDir(dir)
	if err != nil {
		return err
	}
	pkgName, metas, err := scanAndParse(paths)
	if err != nil {
		return err
	}
	if len(metas) == 0 {
		fmt.Println("no rules found")
		return nil
	}
	if err := graph.ValidateGraph(metas); err != nil {
		return err
	}
	code, err := codegen.GenerateRegister(pkgName, metas)
	if err != nil {
		return err
	}
	if dryRun {
		fmt.Print(code)
		return nil
	}
	if output == "" {
		output = filepath.Join(dir, "register_gen.go")
	}
	return os.WriteFile(output, []byte(code), 0644)
}
