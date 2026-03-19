//ff:func feature=state type=engine control=sequence
//ff:what Can: 전이 가능 여부 판정
package state

import "fmt"

// Can evaluates whether the transition is allowed. Returns verdict.
func (m *Machine) Can(req *TransitionRequest, ctx *TransitionContext) (float64, error) {
	key := req.From + ":" + req.Event
	t, ok := m.transitions[key]
	if !ok {
		return -1, fmt.Errorf("no transition registered for %s:%s", req.From, req.Event)
	}
	results, err := t.graph.Evaluate(req, ctx)
	if err != nil {
		return -1, err
	}
	if len(results) == 0 {
		return -1, nil
	}
	return results[0].Verdict, nil
}
