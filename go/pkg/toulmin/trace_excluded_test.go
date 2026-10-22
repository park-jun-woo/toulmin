//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what traceExcluded — Excluded flag of the entry named short within a single trace
package toulmin

// traceExcluded returns the Excluded flag of the entry named short within trace, and whether
// such an entry was found.
func traceExcluded(trace []TraceEntry, short string) (bool, bool) {
	for _, e := range trace {
		if e.Name == short {
			return e.Excluded, true
		}
	}
	return false, false
}
