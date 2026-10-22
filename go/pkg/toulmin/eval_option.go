//ff:type feature=engine type=model
//ff:what EvalOption — evaluation options (method, trace, duration)
package toulmin

// EvalOption controls evaluation behavior.
type EvalOption struct {
	// Method selects the verdict computation algorithm. Default is Matrix.
	Method EvalMethod
	// Trace enables per-warrant TraceEntry collection.
	Trace bool
	// Duration enables per-rule execution time measurement in TraceEntry.
	// Automatically enables Trace when true.
	Duration bool
}
