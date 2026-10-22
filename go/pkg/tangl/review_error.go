//ff:type feature=tangl type=model
//ff:what ReviewError — a failure that must escalate to human review (cause + failed compensation)
package tangl

// ReviewError signals that an execution failed and its compensation also failed
// (or a compensation-less failure was explicitly escalated). Both the original
// cause and the compensation error are preserved for the human reviewer.
type ReviewError struct {
	Cause      error
	Compensate error
}
