//ff:type feature=engine type=model
//ff:what Trace — read-only view of one Run (all node entries + ctx)
package toulmin

// Trace is the read-only view passed to a RunOn handler and returned by Run:
// every node's entry (registration order) plus this Run's context.
type Trace struct {
	nodes []TraceEntry
	ctx   Context
}
