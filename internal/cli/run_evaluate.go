//ff:func feature=cli type=command control=sequence
//ff:what runEvaluate — executes example evaluation and prints results
package cli

import (
	"encoding/json"
	"os"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
	"github.com/spf13/cobra"
)

// runEvaluate demonstrates the engine with example warrant + defeater rules.
func runEvaluate(cmd *cobra.Command, args []string) error {
	eng := toulmin.NewEngine()
	eng.Register(toulmin.RuleMeta{
		Name:      "OneFileOneFunc",
		Qualifier: 1.0,
		Strength:  toulmin.Defeasible,
		Backing:   &demoBacking{Value: "Bohm-Jacopini theorem"},
		Fn:        func(claim any, ground any, backing toulmin.Backing) (bool, any) { return true, nil },
	})
	eng.Register(toulmin.RuleMeta{
		Name:      "TestFileException",
		Qualifier: 1.0,
		Strength:  toulmin.Defeater,
		Defeats:   []string{"OneFileOneFunc"},
		Backing:   &demoBacking{Value: "test files conventionally group multiple test funcs"},
		Fn:        func(claim any, ground any, backing toulmin.Backing) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
