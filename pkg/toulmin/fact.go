//ff:type feature=engine type=model
//ff:what Fact — optional structured evidence convention (where/expected/actual)
package toulmin

// Fact is an optional convention type for the `any` Evidence return value of a
// rule function. The engine does not require or interpret it — downstream
// consumers (e.g. reins) type-assert TraceEntry.Evidence.(Fact) /
// EvalResult.Evidence.(Fact) themselves.
type Fact struct {
	Where    any
	Expected any
	Actual   any
}
