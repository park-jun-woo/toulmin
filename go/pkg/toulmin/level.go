//ff:type feature=engine type=model
//ff:what Level — crisp violation level (unset/fail/review) for reins-style aggregation
package toulmin

// Level classifies a rule's crisp reporting level for downstream aggregation
// (e.g. reins: any Fail → FAIL, else any Review → REVIEW, else PASS).
// Unset is the zero value and carries the "no level to report" meaning.
type Level int

const (
	Unset Level = iota
	Fail
	Review
)
