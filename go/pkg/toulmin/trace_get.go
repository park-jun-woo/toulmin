//ff:func feature=engine type=model control=iteration dimension=1
//ff:what Trace.Get — one node's entry by short name
package toulmin

// Get returns the entry whose Name matches name, plus whether it was found.
func (t Trace) Get(name string) (TraceEntry, bool) {
	for i := range t.nodes {
		if t.nodes[i].Name == name {
			return t.nodes[i], true
		}
	}
	return TraceEntry{}, false
}
