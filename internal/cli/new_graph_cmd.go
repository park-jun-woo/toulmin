//ff:func feature=cli type=command control=sequence
//ff:what NewGraphCmd — creates cobra graph subcommand
package cli

import "github.com/spf13/cobra"

// NewGraphCmd returns the graph subcommand that reads a YAML graph definition
// and generates Graph Builder Go code.
func NewGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph <yaml>",
		Short: "Generate Graph Builder code from YAML graph definition",
		Args:  cobra.ExactArgs(1),
		RunE:  runGraph,
	}
	cmd.Flags().String("output", "", "output file path (default: <yaml dir>/graph_gen.go)")
	cmd.Flags().Bool("dry-run", false, "print generated code to stdout")
	cmd.Flags().String("package", "", "Go package name (default: directory name)")
	return cmd
}
