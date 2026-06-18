//ff:func feature=engine type=model control=sequence
//ff:what Trace.All — every node's entry in registration order
package toulmin

// All returns every node's entry in registration order.
func (t Trace) All() []TraceEntry { return t.nodes }
