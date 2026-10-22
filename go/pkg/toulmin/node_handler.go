//ff:type feature=engine type=model
//ff:what NodeHandler — run handler signature, fired on each Active node
package toulmin

// NodeHandler is invoked by Run on each Active node. self is that node's entry
// (no name lookup needed); t is the whole Run's Trace (All/Get/Ctx).
type NodeHandler func(self TraceEntry, t Trace) error
