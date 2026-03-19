//ff:type feature=approve type=model
//ff:what StepResult: 단일 승인 단계 판정 결과
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// StepResult holds the result of a single approval step.
type StepResult struct {
	Name    string
	Verdict float64
	Trace   []toulmin.TraceEntry
}
