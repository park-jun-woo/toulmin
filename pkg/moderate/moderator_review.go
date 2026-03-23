//ff:func feature=moderate type=engine control=sequence
//ff:what Review: 콘텐츠 모더레이션 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Review evaluates the content and returns a ReviewResult.
// verdict > 0.3 → allow, 0 < verdict <= 0.3 → flag, verdict <= 0 → block.
func (m *Moderator) Review(content *Content, ctx *ContentContext) (*ReviewResult, error) {
	results, err := m.graph.Evaluate(content, ctx, toulmin.EvalOption{Trace: true})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &ReviewResult{Allowed: false, Verdict: -1, Action: ActionBlock}, nil
	}
	v := results[0].Verdict
	var action Action
	switch {
	case v > 0.3:
		action = ActionAllow
	case v > 0:
		action = ActionFlag
	default:
		action = ActionBlock
	}
	return &ReviewResult{
		Allowed: action == ActionAllow,
		Verdict: v,
		Action:  action,
		Trace:   results[0].Trace,
	}, nil
}
