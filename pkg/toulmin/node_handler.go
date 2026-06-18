//ff:type feature=engine type=model
//ff:what NodeHandler — run handler signature, fired on each Active node
package toulmin

// NodeHandler is invoked by Run on each Active node.
// self is the firing node's trace entry; trace is every node's entry (read-only view).
type NodeHandler func(ctx Context, self TraceEntry, trace []TraceEntry) error
