//ff:type feature=engine type=model
//ff:what NodeHandler — run handler signature, fired on each Active node
package toulmin

// NodeHandler is invoked by Run on each Active node with the whole Run's Trace.
// The handler's own node is g.Rule(X).RunOn(h) → look it up via t.Get("X").
type NodeHandler func(t Trace) error
