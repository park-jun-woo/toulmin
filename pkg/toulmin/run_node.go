//ff:func feature=engine type=engine control=sequence
//ff:what runNode — fires one Active node's RunOn handler then Runs its sub-graph
package toulmin

import "fmt"

// runNode executes the run side-effects for a single Active node: (a) fires the node's
// RunOn handler with the whole Run's Trace, then (b) if the node has a RunGraph, Runs that
// sub-graph with the same ctx at depth+1. Sub-graph verdicts stay isolated; only errors
// propagate, wrapped as `runOn ...` / `run ... → ...`. name is the node's short name.
func runNode(meta *RuleMeta, name string, tr Trace, ctx Context, opt EvalOption, depth int) error {
	if meta.RunOn != nil {
		if herr := safeCallHandler(meta.RunOn, tr); herr != nil {
			return fmt.Errorf("runOn %q: %w", name, herr)
		}
	}
	if meta.RunGraph != nil {
		if _, _, err := meta.RunGraph.runDepth(ctx, opt, depth+1); err != nil {
			return fmt.Errorf("run %q → %q: %w", name, meta.RunGraph.name, err)
		}
	}
	return nil
}
