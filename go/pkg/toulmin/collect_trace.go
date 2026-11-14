//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what collectTrace — copies trace entries with short names
package toulmin

// collectTrace creates a copy of trace entries with shortened names.
func collectTrace(trace []TraceEntry) []TraceEntry {
	out := make([]TraceEntry, len(trace))
	for i, te := range trace {
		te.Name = shortName(te.Name)
		out[i] = te
	}
	return out
}
