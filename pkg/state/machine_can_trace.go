//ff:func feature=state type=engine control=sequence
//ff:what CanTrace: 전이 가능 여부 + 판정 근거
package state

import "fmt"

// CanTrace evaluates the transition with full trace.
func (m *Machine) CanTrace(req *TransitionRequest, ctx *TransitionContext) (*TraceResult, error) {
	key := req.From + ":" + req.Event
	t, ok := m.transitions[key]
	if !ok {
		return nil, fmt.Errorf("no transition registered for %s:%s", req.From, req.Event)
	}
	results, err := t.graph.EvaluateTrace(req, ctx)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &TraceResult{Verdict: -1, From: t.from, To: t.to, Event: t.event}, nil
	}
	return &TraceResult{
		Verdict: results[0].Verdict,
		From:    t.from,
		To:      t.to,
		Event:   t.event,
		Trace:   results[0].Trace,
	}, nil
}
