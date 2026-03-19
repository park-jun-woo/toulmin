//ff:type feature=state type=model
//ff:what TraceResult: 전이 판정 결과 + 근거
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// TraceResult holds the transition evaluation result with trace.
type TraceResult struct {
	Verdict float64
	From    string
	To      string
	Event   string
	Trace   []toulmin.TraceEntry
}
