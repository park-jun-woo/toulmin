//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what buildTraceEntries — assembles every node's TraceEntry from the evalContext in registration order
package toulmin

// buildTraceEntries builds one TraceEntry per rule in registration order from the
// filled evalContext. ctx is the shared Ground attached to each entry.
func (g *Graph) buildTraceEntries(ec *evalContext, ctx Context) []TraceEntry {
	entries := make([]TraceEntry, len(g.rules))
	for i := range g.rules {
		entries[i] = g.buildTraceEntry(ec, g.rules[i].Name, ctx)
	}
	return entries
}
