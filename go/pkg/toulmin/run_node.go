//ff:func feature=engine type=engine control=sequence
//ff:what runNode — fires one Active node's RunOn handler (with its entry) then Runs its sub-graph
package toulmin

import "fmt"

// runNode executes the run side-effects for a single Active node: (a) fires the node's
// RunOn handler with its own entry (self) and the whole Run's Trace, then (b) if the node
// has a RunGraph, Runs that sub-graph with the same ctx at depth+1. Sub-graph verdicts stay
// isolated; only errors propagate, wrapped as `runOn ...` / `run ... → ...`. self is the
// node's own TraceEntry.
func runNode(meta *RuleMeta, self TraceEntry, tr Trace, ctx Context, opt EvalOption, depth int) error {
	if meta.RunOn != nil {
		if herr := safeCallHandler(meta.RunOn, self, tr); herr != nil {
			return fmt.Errorf("runOn %q: %w", self.Name, herr)
		}
	}
	if meta.RunGraph != nil {
		if _, _, err := meta.RunGraph.runDepth(ctx, opt, depth+1); err != nil {
			return fmt.Errorf("run %q → %q: %w", self.Name, meta.RunGraph.name, err)
		}
	}
	return nil
}
