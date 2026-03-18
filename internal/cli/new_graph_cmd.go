//ff:func feature=cli type=command control=sequence
//ff:what NewGraphCmd — creates cobra graph subcommand
package cli

import "github.com/spf13/cobra"

// NewGraphCmd returns the graph subcommand that scans sources
// and generates rule registration code.
func NewGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph <dir>",
		Short: "Scan //rule: annotations and generate RegisterAll code",
		Args:  cobra.ExactArgs(1),
		RunE:  runGraph,
	}
	cmd.Flags().String("output", "", "output file path (default: <dir>/register_gen.go)")
	cmd.Flags().Bool("dry-run", false, "print generated code to stdout")
	return cmd
}
