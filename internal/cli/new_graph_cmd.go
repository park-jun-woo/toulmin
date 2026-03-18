//ff:func feature=cli type=command control=sequence
//ff:what NewGraphCmd — creates cobra graph subcommand
package cli

import "github.com/spf13/cobra"

// NewGraphCmd returns the graph subcommand that validates and generates
// Graph Builder code from YAML or validates Go files for cycles.
func NewGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph <file.yaml|file.go>",
		Short: "Validate graph and generate code from YAML, or check Go files for cycles",
		Args:  cobra.ExactArgs(1),
		RunE:  runGraph,
	}
	cmd.Flags().String("output", "", "output file path (default: <yaml dir>/graph_gen.go)")
	cmd.Flags().Bool("dry-run", false, "print generated code to stdout")
	cmd.Flags().String("package", "", "Go package name (default: directory name)")
	cmd.Flags().Bool("check", false, "validate only, no code generation")
	return cmd
}
