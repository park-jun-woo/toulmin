//ff:func feature=approve type=engine control=iteration dimension=1
//ff:what Evaluate: 전체 흐름 판정 (모든 단계 통과해야 승인)
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Evaluate runs all steps sequentially. All steps must pass (verdict > 0).
// Stops at the first failing step.
func (f *Flow) Evaluate(req *ApprovalRequest, ctxBuilder StepContextFunc) (*FlowResult, error) {
	result := &FlowResult{Approved: true}
	for _, step := range f.steps {
		ctx := ctxBuilder(step.Name)
		results, err := step.Graph.EvaluateTrace(req, ctx)
		if err != nil {
			return nil, err
		}
		var verdict float64
		var trace []toulmin.TraceEntry
		if len(results) > 0 {
			verdict = results[0].Verdict
			trace = results[0].Trace
		} else {
			verdict = -1
		}
		result.Steps = append(result.Steps, StepResult{
			Name:    step.Name,
			Verdict: verdict,
			Trace:   trace,
		})
		if verdict <= 0 {
			result.Approved = false
			return result, nil
		}
	}
	return result, nil
}
