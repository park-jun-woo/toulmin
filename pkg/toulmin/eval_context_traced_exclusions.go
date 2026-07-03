//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what tracedExclusions — patch this warrant's own trace with supersede exclusions, return it
package toulmin

// tracedExclusions runs the reporting-only resolveSupersedes pass on the current (per-warrant,
// reset) state and marks every trace entry whose full id is superseded by an Active superseder
// as Excluded, then returns the collected trace. resolveSupersedes may set ec.err via ec.calc;
// the caller checks ec.err after this returns, so patching a possibly-errored trace is harmless.
// Keys are full ids, matching ec.trace[i].Name.
func (ec *evalContext) tracedExclusions(ctx Context) []TraceEntry {
	excluded := ec.resolveSupersedes(ctx)
	for i := range ec.trace {
		if excluded[ec.trace[i].Name] {
			ec.trace[i].Excluded = true
		}
	}
	return collectTrace(ec.trace)
}
