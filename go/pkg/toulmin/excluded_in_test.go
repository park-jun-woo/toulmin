//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what excludedIn — Excluded flag of the entry named short within the warrant's own trace
package toulmin

// excludedIn returns the Excluded flag of the entry named short within result's own trace.
func excludedIn(results []EvalResult, warrant, short string) (bool, bool) {
	for _, r := range results {
		if r.Name != warrant {
			continue
		}
		return traceExcluded(r.Trace, short)
	}
	return false, false
}
