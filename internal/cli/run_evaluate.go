//ff:func feature=cli type=command control=sequence
//ff:what runEvaluate — executes example evaluation and prints results
package cli

import (
	"encoding/json"
	"os"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
	"github.com/spf13/cobra"
)

// runEvaluate demonstrates the graph with example warrant + defeater rules.
func runEvaluate(cmd *cobra.Command, args []string) error {
	g := toulmin.NewGraph("evaluate")
	warrant := g.Rule(func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }).
		With(&demoSpec{Value: "Bohm-Jacopini theorem"})
	g.Except(func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }).
		With(&demoSpec{Value: "test files conventionally group multiple test funcs"}).
		Attacks(warrant)
	ctx := toulmin.NewContext()
	results, err := g.Evaluate(ctx)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
