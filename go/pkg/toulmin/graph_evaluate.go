//ff:func feature=engine type=engine control=sequence
//ff:what Evaluate — evaluates rules by graph traversal and returns verdicts
package toulmin

import "fmt"

// Evaluate traverses the defeats graph from each rule node and returns verdicts.
// Verdict is computed using lazy recursive h-Categoriser. Use EvalOption to enable Trace, Duration, or Recursive (planned) method.
// In Trace mode, each warrant is evaluated independently (rule functions may execute multiple times).
// In non-Trace mode, rule results are shared across warrants (each function executes at most once).
// Returns an error if the defeat graph contains a cycle or a rule panics.
// The maps derived from the graph are built on the first evaluation and reused until the graph changes;
// concurrent evaluations of an unchanging graph are safe as far as its rule functions are.
func (g *Graph) Evaluate(ctx Context, opts ...EvalOption) ([]EvalResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("toulmin: ctx must not be nil")
	}
	opt, err := resolveOption(opts)
	if err != nil {
		return nil, err
	}
	results, ec, err := g.evaluate(ctx, opt, false)
	if ec != nil && !opt.Trace {
		ec.home.release(ec)
	}
	return results, err
}
