//ff:type feature=moderate type=model
//ff:what ReviewResult: 모더레이션 판정 결과
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ReviewResult holds the moderation evaluation result.
type ReviewResult struct {
	Allowed bool
	Verdict float64
	Action  Action
	Trace   []toulmin.TraceEntry
}
